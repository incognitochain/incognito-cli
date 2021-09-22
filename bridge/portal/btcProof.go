package portal

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// MerkleProof represents a Merkle path.
type MerkleProof struct {
	ProofHash *chainhash.Hash
	IsLeft    bool
}

// BTCProof represents a Merkle proof for a BTC transaction.
type BTCProof struct {
	MerkleProofs []*MerkleProof
	BTCTx        *wire.MsgTx
	BlockHash    *chainhash.Hash
}

func buildMerkleTreeStoreFromTxHashes(txHashes []*chainhash.Hash) []*chainhash.Hash {
	nextPoT := nextPowerOfTwo(len(txHashes))
	arraySize := nextPoT*2 - 1
	merkelTrees := make([]*chainhash.Hash, arraySize)

	for i, txHash := range txHashes {
		merkelTrees[i] = txHash
	}

	offset := nextPoT
	for i := 0; i < arraySize-1; i += 2 {
		switch {
		case merkelTrees[i] == nil:
			merkelTrees[offset] = nil

		case merkelTrees[i+1] == nil:
			newHash := HashMerkleBranches(merkelTrees[i], merkelTrees[i])
			merkelTrees[offset] = newHash

		default:
			newHash := HashMerkleBranches(merkelTrees[i], merkelTrees[i+1])
			merkelTrees[offset] = newHash
		}
		offset++
	}

	return merkelTrees
}

// verifyMerkleProof checks if a MerkleProof has been properly constructed.
func verifyMerkleProof(
	merkleRoot *chainhash.Hash,
	merkleProofs []*MerkleProof,
	txHash *chainhash.Hash,
) bool {
	curHash := txHash
	for _, mklProof := range merkleProofs {
		if mklProof.IsLeft {
			curHash = HashMerkleBranches(mklProof.ProofHash, curHash)
		} else {
			curHash = HashMerkleBranches(curHash, mklProof.ProofHash)
		}
	}
	return curHash.String() == merkleRoot.String()
}

// BuildMerkleProof returns a list of MerkleProof of a transaction hash.
func BuildMerkleProof(txHashes []*chainhash.Hash, targetedTxHash *chainhash.Hash) []*MerkleProof {
	merkleTree := buildMerkleTreeStoreFromTxHashes(txHashes)
	nextPoT := nextPowerOfTwo(len(txHashes))
	layers := make([][]*chainhash.Hash, 0)
	left := 0
	right := nextPoT
	for left < right {
		layers = append(layers, merkleTree[left:right])
		curLen := len(merkleTree[left:right])
		left = right
		right = right + curLen/2
	}

	merkleProofs := make([]*MerkleProof, 0)
	curHash := targetedTxHash
	for _, layer := range layers {
		if len(layer) == 1 {
			break
		}

		for i := 0; i < len(layer); i++ {
			if layer[i] == nil || layer[i].String() != curHash.String() {
				continue
			}
			if i%2 == 0 {
				if layer[i+1] == nil {
					curHash = HashMerkleBranches(layer[i], layer[i])
					merkleProofs = append(
						merkleProofs,
						&MerkleProof{
							ProofHash: layer[i],
							IsLeft:    false,
						},
					)
				} else {
					curHash = HashMerkleBranches(layer[i], layer[i+1])
					merkleProofs = append(
						merkleProofs,
						&MerkleProof{
							ProofHash: layer[i+1],
							IsLeft:    false,
						},
					)
				}
			} else {
				if layer[i-1] == nil {
					curHash = HashMerkleBranches(layer[i], layer[i])
					merkleProofs = append(
						merkleProofs,
						&MerkleProof{
							ProofHash: layer[i],
							IsLeft:    true,
						},
					)
				} else {
					curHash = HashMerkleBranches(layer[i-1], layer[i])
					merkleProofs = append(
						merkleProofs,
						&MerkleProof{
							ProofHash: layer[i-1],
							IsLeft:    true,
						},
					)
				}
			}
			break // process next layer
		}
	}
	return merkleProofs
}

// BuildProof constructs the proof for a BTC transaction hash.
func (b *BTCClient) BuildProof(txHashStr string, blkHeight uint64) (string, error) {
	if b.isNil() {
		return "", ErrBTCClientNotInitialized
	}
	var blkHash *chainhash.Hash
	var merkleRoot *chainhash.Hash
	var msgBlk *wire.MsgBlock
	var msgTx *wire.MsgTx
	var err error
	txHashes := make([]*chainhash.Hash, 0)

	txHash, err := chainhash.NewHashFromStr(txHashStr)
	if err != nil {
		return "", err
	}

	// get the block and all transactions in the block.
	if b.rpcClient != nil {
		blkHash, err = b.rpcClient.GetBlockHash(int64(blkHeight))
		if err != nil {
			return "", err
		}

		msgBlk, err = b.rpcClient.GetBlock(blkHash)
		if err != nil {
			return "", err
		}

		// get all txs in the blocks.
		for _, tx := range msgBlk.Transactions {
			tmpTxHash := tx.TxHash()
			txHashes = append(txHashes, &tmpTxHash)
		}

		tx, err := b.rpcClient.GetRawTransaction(txHash)
		if err != nil {
			return "", err
		}
		msgTx = tx.MsgTx()
	} else if b.cypherBlockClient != nil {
		var block gobcy.Block

		cur := 0
		for {
			block, err = b.cypherBlockClient.GetBlock(
				int(blkHeight),
				"",
				map[string]string{
					"txstart": fmt.Sprintf("%d", cur),
					"limit":   fmt.Sprintf("%d", cur+500),
				},
			)

			txIDs := block.TXids
			for i := 0; i < len(txIDs); i++ {
				tmpTxHash, err := chainhash.NewHashFromStr(txIDs[i])
				if err != nil {
					return "", err
				}
				txHashes = append(txHashes, tmpTxHash)
			}

			if len(txHashes) == block.NumTX {
				fmt.Println("numTxs", len(txHashes))
				break
			} else {
				cur += 500
			}
		}

		if err != nil {
			return "", err
		}

		blkHash, err = chainhash.NewHashFromStr(block.Hash)
		if err != nil {
			return "", err
		}

		msgTx, err = b.BuildMsgTxFromCypher(txHashStr)
		if err != nil {
			return "", err
		}

		merkleRoot, err = chainhash.NewHashFromStr(block.MerkleRoot)
		if err != nil {
			return "", err
		}
	}

	// build the Merkle proof for the transaction.
	merkleProofs := BuildMerkleProof(txHashes, txHash)
	if merkleRoot != nil {
		valid := verifyMerkleProof(merkleRoot, merkleProofs, txHash)
		if !valid {
			return "", fmt.Errorf("invalid merkleProofs")
		}
	}
	btcProof := BTCProof{
		MerkleProofs: merkleProofs,
		BTCTx:        msgTx,
		BlockHash:    blkHash,
	}
	btcProofBytes, _ := json.Marshal(btcProof)
	btcProofStr := base64.StdEncoding.EncodeToString(btcProofBytes)

	return btcProofStr, nil
}

// BuildMsgTxFromCypher returns a wire.MsgTx given a BTC transaction hash.
func (b *BTCClient) BuildMsgTxFromCypher(txHashStr string) (*wire.MsgTx, error) {
	if b.cypherBlockClient == nil {
		return nil, fmt.Errorf("cypherBlock client is not initialized")
	}
	bc := b.cypherBlockClient
	cypherTx, _ := bc.GetTX(txHashStr, nil)

	txIns := make([]*wire.TxIn, 0)
	for _, cypherTxIn := range cypherTx.Inputs {
		prevHash, _ := chainhash.NewHashFromStr(cypherTxIn.PrevHash)

		signatureScript, err := hex.DecodeString(cypherTxIn.Script)
		if err != nil {
			return nil, err
		}

		in := &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  *prevHash,
				Index: uint32(cypherTxIn.OutputIndex),
			},
			SignatureScript: signatureScript,
			Sequence:        uint32(cypherTxIn.Sequence),
		}
		txIns = append(txIns, in)
	}

	txOuts := make([]*wire.TxOut, 0)
	for _, cypherTxOut := range cypherTx.Outputs {
		pkScript, err := hex.DecodeString(cypherTxOut.Script)
		if err != nil {
			return nil, err
		}
		out := &wire.TxOut{
			Value:    cypherTxOut.Value.Int64(),
			PkScript: pkScript,
		}
		txOuts = append(txOuts, out)
	}

	msgTx := wire.MsgTx{
		Version:  int32(cypherTx.Ver),
		TxIn:     txIns,
		TxOut:    txOuts,
		LockTime: uint32(cypherTx.LockTime),
	}
	return &msgTx, nil
}
