package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/incognito-cli/bridge/evm/erc20"
	"github.com/incognitochain/incognito-cli/bridge/evm/vault"
	"log"
	"math/big"
	"strings"
	"time"
)

var (
	nativeToken         = "0x0000000000000000000000000000000000000000"
	nativeTokenDecimals = 18
)

// EVMAccount represents an account on the Ethereum/Binance networks.
type EVMAccount struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
}

// NewEVMAccount returns a new EVMAccount given a hex-encoded private key.
func NewEVMAccount(hexPrivateKey string) (*EVMAccount, error) {
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &EVMAccount{
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
		address:    address,
	}, nil

}

// NewTransactionOpts creates a new bind.TransactOpts for an EVMAccount.
func (acc EVMAccount) NewTransactionOpts(destAddr common.Address, gasPrice, gasLimit, amount uint64, data []byte, isBSC bool) (*bind.TransactOpts, error) {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}

	var err error

	// calculate gas price if needed.
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = evmClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	//calculate gas limit
	if gasLimit == 0 {
		gasLimit, err = evmClient.EstimateGas(context.Background(), ethereum.CallMsg{From: acc.address, To: &destAddr, Data: data})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}
	}

	nonce, err := evmClient.PendingNonceAt(context.Background(), acc.address)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce error: %v", err)
	}

	chainID, err := evmClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(acc.privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPriceBigInt
	auth.GasLimit = gasLimit
	auth.Nonce = new(big.Int).SetUint64(nonce)
	if amount != 0 {
		value := new(big.Int).SetUint64(amount)
		auth.Value = value
	}

	return auth, nil
}

// EstimateDepositGas estimates the gas for depositing a token.
func (acc EVMAccount) EstimateDepositGas(tokenAddress common.Address, depositedAmount *big.Int, incAddress string, isBSC bool) (uint64, error) {
	evmClient := cfg.ethClient
	vaultAddress := cfg.ethVaultAddress
	if isBSC {
		evmClient = cfg.bscClient
		vaultAddress = cfg.bscVaultAddress
	}

	var gasLimit uint64
	vaultABI, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return 0, fmt.Errorf("cannot create vaultABI from file")
	}

	var data []byte
	if tokenAddress.String() == nativeToken {
		data, err = vaultABI.Pack(
			"deposit",
			incAddress,
		)
		if err != nil {
			return 0, err
		}

		gasLimit, err = evmClient.EstimateGas(context.Background(),
			ethereum.CallMsg{From: acc.address, Value: depositedAmount, To: &vaultAddress, Data: data})
		if err != nil {
			return 0, fmt.Errorf("estimateGas for native token error: %v", err)
		}
	} else {
		data, err = vaultABI.Pack(
			"depositERC20",
			tokenAddress,
			depositedAmount,
			incAddress,
		)
		if err != nil {
			return 0, err
		}

		gasLimit, err = evmClient.EstimateGas(context.Background(),
			ethereum.CallMsg{From: acc.address, To: &vaultAddress, Data: data})
		if err != nil {
			return 0, err
		}
	}

	return gasLimit, nil
}

// EstimateWithdrawalGas estimates the gas for withdrawing a token.
func (acc EVMAccount) EstimateWithdrawalGas(burnProof *incclient.BurnProof, isBSC bool) (uint64, error) {
	evmClient := cfg.ethClient
	vaultAddress := cfg.ethVaultAddress
	if isBSC {
		evmClient = cfg.bscClient
		vaultAddress = cfg.bscVaultAddress
	}

	vaultABI, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return 0, fmt.Errorf("cannot create vaultABI from file")
	}

	var data []byte
	data, err = vaultABI.Pack(
		"withdraw",
		burnProof.Instruction,
		burnProof.Heights[0],
		burnProof.InstPaths[0],
		burnProof.InstPathIsLefts[0],
		burnProof.InstRoots[0],
		burnProof.BlkData[0],
		burnProof.SigIndices[0],
		burnProof.SigVs[0],
		burnProof.SigRs[0],
		burnProof.SigSs[0],
	)
	if err != nil {
		return 0, err
	}

	gasLimit, err := evmClient.EstimateGas(context.Background(),
		ethereum.CallMsg{From: acc.address, To: &vaultAddress, Data: data})
	if err != nil {
		return 0, fmt.Errorf("estimateGas for withdrawal error: %v", err)
	}

	return gasLimit, nil
}

// GetBalance returns the balance of a token.
func (acc EVMAccount) GetBalance(tokenAddress common.Address, isBSC bool) (uint64, error) {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}
	if tokenAddress.String() == nativeToken {
		balance, err := evmClient.BalanceAt(context.Background(), acc.address, nil)
		if err != nil {
			return 0, nil
		}

		return balance.Uint64(), nil
	} else {
		erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
		if err != nil {
			return 0, err
		}

		balance, err := erc20Instance.BalanceOf(&bind.CallOpts{}, acc.address)
		if err != nil {
			return 0, err
		}

		return balance.Uint64(), nil
	}
}

// GetAllowance returns the allowance of an owner to a spender w.r.t to an ERC20 token.
func (acc EVMAccount) GetAllowance(tokenAddress, spender common.Address, isBSC bool) (uint64, error) {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}

	erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return 0, err
	}

	allowance, err := erc20Instance.Allowance(&bind.CallOpts{}, acc.address, spender)
	if err != nil {
		return 0, err
	}

	return allowance.Uint64(), nil
}

// ApproveERC20 approves the Incognito Vault to spend an ERC20/BEP20 token of an account.
func (acc EVMAccount) ApproveERC20(tokenAddress, approved common.Address, approvedAmount, gasPrice uint64, isBSC bool) (*common.Hash, error) {
	prefix := "[ApproveERC20]"
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
		prefix = "[ApproveBEP20]"
	}

	erc20Token, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return nil, err
	}

	// load the ERC20 ABI
	erc20ABI, err := abi.JSON(strings.NewReader(erc20.Erc20ABI))
	if err != nil {
		return nil, err
	}

	// estimate the gas limit
	approvedAmountBigInt := new(big.Int).SetUint64(approvedAmount)
	data, err := erc20ABI.Pack(
		"approve",
		approved,
		approvedAmountBigInt,
	)
	if err != nil {
		return nil, err
	}

	gasPriceBigInt, gasLimit, err := acc.getGasLimitAndPrice(0, gasPrice, ethereum.CallMsg{To: &tokenAddress, Data: data}, isBSC)
	if err != nil {
		return nil, err
	}
	txFee := gasLimit * gasPriceBigInt.Uint64()
	log.Printf("%v gasPrice: %v, gasLimit %v, txFee %v\n", prefix, gasPriceBigInt.Uint64(), gasLimit, txFee)

	auth, err := acc.NewTransactionOpts(tokenAddress, gasPriceBigInt.Uint64(), gasLimit, 0, data, isBSC)
	if err != nil {
		return nil, err
	}
	amountBigInt := new(big.Int).SetUint64(approvedAmount)

	tx, err := erc20Token.Approve(auth, approved, amountBigInt)
	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()
	log.Printf("%v TxHash: %v\n", prefix, txHash.String())
	return &txHash, nil
}

func (acc EVMAccount) wait(tx common.Hash, isBSC bool) error {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}
	for range time.Tick(10 * time.Second) {
		receipt, err := evmClient.TransactionReceipt(context.Background(), tx)
		if err == nil {
			log.Printf("[EVM Status] TxHash %v: %v\n", tx.String(), receipt.Status)
			if receipt.Status == 0 {
				return fmt.Errorf("tx %v failed", tx.String())
			}
			break
		} else if err == ethereum.NotFound {
			continue
		} else {
			return err
		}
	}
	return nil
}

func (acc EVMAccount) getGasLimitAndPrice(gasLimit, gasPrice uint64, callMsg ethereum.CallMsg, isBSC bool) (*big.Int, uint64, error) {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}

	var err error

	// calculate gas price if needed.
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = evmClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, 0, err
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	if gasLimit == 0 {
		callMsg.From = acc.address
		gasLimit, err = evmClient.EstimateGas(context.Background(), callMsg)
		if err != nil {
			return nil, 0, err
		}
	}

	return gasPriceBigInt, gasLimit, nil
}

func (acc EVMAccount) verifyProofAndParseReceipt(iReq *metadata.IssuingEVMRequest, isBSC bool) (*types.Receipt, error) {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}

	evmHeader, err := evmClient.HeaderByHash(context.Background(), iReq.BlockHash)
	if err != nil {
		return nil, err
	}
	if evmHeader == nil {
		return nil, fmt.Errorf("WARNING: Could not find out the EVM block header with the hash: %s", iReq.BlockHash.String())
	}

	keyBuf := new(bytes.Buffer)
	keyBuf.Reset()
	err = rlp.Encode(keyBuf, iReq.TxIndex)
	if err != nil {
		return nil, err
	}

	nodeList := new(light.NodeList)
	for _, proofStr := range iReq.ProofStrs {
		proofBytes, err := base64.StdEncoding.DecodeString(proofStr)
		if err != nil {
			return nil, err
		}
		err = nodeList.Put([]byte{}, proofBytes)
		if err != nil {
			return nil, err
		}
	}
	proof := nodeList.NodeSet()
	val, err := trie.VerifyProof(evmHeader.ReceiptHash, keyBuf.Bytes(), proof)
	if err != nil {
		return nil, err
	}

	// Decode value from VerifyProof into Receipt
	constructedReceipt := new(types.Receipt)
	err = rlp.DecodeBytes(val, constructedReceipt)
	if err != nil {
		return nil, err
	}

	if constructedReceipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("the constructedReceipt's status is not success")
	}

	return constructedReceipt, nil
}

func getDecimals(tokenAddress common.Address, isBSC bool) (uint64, error) {
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}

	erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return 0, err
	}

	decimals, err := erc20Instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}

	return decimals.Uint64(), nil
}

func getAllDecentralizedBridgeTokens() (incToPublic map[string]string, publicToInc map[string]string, err error) {
	incToPublic = make(map[string]string)
	publicToInc = make(map[string]string)

	allTokens, err := cfg.incClient.GetBridgeTokens()
	if err != nil {
		return
	}

	for _, token := range allTokens {
		if token.IsCentralized {
			continue
		}
		incTokenID := token.TokenID.String()
		publicTokenID := fmt.Sprintf("%x", token.ExternalTokenID)

		incToPublic[incTokenID] = publicTokenID
		publicToInc[publicTokenID] = incTokenID
	}

	return
}