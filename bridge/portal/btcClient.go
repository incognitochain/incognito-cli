package portal

import (
	"fmt"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

var (
	// ErrBTCClientNotInitialized is thrown when trying to call a non-initialized BTCClient.
	ErrBTCClientNotInitialized = fmt.Errorf("the btcClient has not been initialized")
)

// BTCClient implements a wrapped Go client for retrieving information of the BTC network.
type BTCClient struct {
	rpcClient         *rpcclient.Client
	cypherBlockClient *gobcy.API
}

// NewBTCMainNetClient returns a BTC main-net client.
func NewBTCMainNetClient() (*BTCClient, error) {
	b := new(BTCClient)
	b.cypherBlockClient = &gobcy.API{
		Coin:  "btc",
		Chain: "main",
	}

	return b, nil
}

// NewBTCTestNetClient returns a BTC test-net client.
func NewBTCTestNetClient() (*BTCClient, error) {
	b := new(BTCClient)
	b.cypherBlockClient = &gobcy.API{
		Coin:  "btc",
		Chain: "test3",
	}

	return b, nil
}

func (b *BTCClient) isNil() bool {
	return b.rpcClient == nil && b.cypherBlockClient == nil
}

// IsConfirmedTx returns a boolean indicator of whether a BTC txHashStr has been confirmed, and the block height at which
// the transaction was included.
func (b *BTCClient) IsConfirmedTx(txHashStr string) (bool, uint64, error) {
	if b.isNil() {
		return false, 0, ErrBTCClientNotInitialized
	}

	if b.rpcClient != nil {
		txID, _ := chainhash.NewHashFromStr(txHashStr)
		tx, err := b.rpcClient.GetRawTransactionVerbose(txID)
		if err != nil {
			return false, 0, err
		}
		if tx.Confirmations >= 6 {
			blkHash, _ := chainhash.NewHashFromStr(tx.BlockHash)
			blk, err := b.rpcClient.GetBlockHeaderVerbose(blkHash)
			if err != nil {
				return false, 0, err
			}
			return true, uint64(blk.Height), nil
		}

		return false, 0, nil
	}

	tx, err := b.cypherBlockClient.GetTX(txHashStr, nil)
	if err != nil {
		return false, 0, err
	}
	return tx.Confirmations >= 6, uint64(tx.BlockHeight), nil
}
