package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/incognito-cli/bridge/evm/vault"
	"log"
	"math/big"
	"time"
)

// DepositETH deposits an amount of ETH to the Incognito contract.
func (acc EVMAccount) DepositETH(incAddress string, depositedAmount, gasLimit, gasPrice uint64) (*common.Hash, error) {
	evmClient := acc.evmConfig.ethClient
	vaultAddress := acc.evmConfig.ethVaultAddress
	prefix := "[DepositETH]"

	balance, err := acc.GetBalance(common.HexToAddress(nativeToken), false)
	if err != nil {
		return nil, err
	}

	v, err := vault.NewVault(vaultAddress, evmClient)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = evmClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	// estimate gasLimit
	amountBigInt := new(big.Int).SetUint64(depositedAmount)
	gasLimit, err = acc.EstimateDepositGas(common.HexToAddress(nativeToken), amountBigInt, incAddress)
	if err != nil {
		return nil, err
	}
	txFee := gasLimit * gasPriceBigInt.Uint64()
	log.Printf("%v gasLimit %v, gasPrice %v, txFee %v\n", prefix, gasLimit, gasPriceBigInt.Uint64(), txFee)

	requiredAmount := depositedAmount + txFee
	if balance <= requiredAmount {
		return nil, fmt.Errorf("[DepositETH] balance insufficient, need %v, got %v", requiredAmount, balance)
	}

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, amountBigInt.Uint64(), nil)
	if err != nil {
		return nil, err
	}

	tx, err := v.Deposit(auth, incAddress)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()
	log.Printf("%v deposited tx: %v\n", prefix, txHash.String())

	if err := acc.wait(txHash); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// DepositERC20 shields an amount of ERC20 to the Incognito network.
func (acc EVMAccount) DepositERC20(incAddress, tokenAddressStr string, depositedAmount, gasLimit, gasPrice uint64) (*common.Hash, error) {
	vaultAddress := acc.evmConfig.ethVaultAddress
	evmClient := acc.evmConfig.ethClient
	isBSC := false
	prefix := "[DepositERC20]"

	// load the vault instance
	v, err := vault.NewVault(vaultAddress, evmClient)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(tokenAddressStr)
	balance, err := acc.GetBalance(common.HexToAddress(nativeToken), isBSC)
	if err != nil {
		return nil, err
	}

	tokenBalance, err := acc.GetBalance(tokenAddress, isBSC)
	if err != nil {
		return nil, err
	}
	log.Printf("%v tokenBalance %v, depositingAmount %v\n", prefix, tokenBalance, depositedAmount)
	if tokenBalance < depositedAmount {
		return nil, fmt.Errorf("%v balance ERC20 insufficient, need %v, got %v", prefix, depositedAmount, balance)
	}

	// get the current ERC20 allowance
	allowance, err := acc.GetAllowance(tokenAddress, vaultAddress, isBSC)
	if err != nil {
		return nil, err
	}
	if allowance < depositedAmount {
		log.Printf("%v insufficient allowance: got %v, need %v\n", prefix, allowance, depositedAmount)
		txHash, err := acc.ApproveERC20(tokenAddress, vaultAddress, depositedAmount, 0, isBSC)
		if err != nil {
			return nil, err
		}
		log.Printf("%v ApprovedERC20 TxHash: %v\n", prefix, txHash.String())
		err = acc.wait(*txHash)
		if err != nil {
			return nil, err
		}
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = acc.evmConfig.ethClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	// estimate gasLimit
	amountBigInt := new(big.Int).SetUint64(depositedAmount)
	gasLimit, err = acc.EstimateDepositGas(tokenAddress, amountBigInt, incAddress)
	if err != nil {
		return nil, err
	}
	txFee := gasLimit * gasPriceBigInt.Uint64()
	log.Printf("%v gasLimit %v, gasPrice %v, txFee %v\n", prefix, gasLimit, gasPriceBigInt.Uint64(), txFee)
	if balance < txFee {
		return nil, fmt.Errorf("%v balance insufficient, need %v, got %v", prefix, txFee, balance)
	}

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, 0, nil)
	if err != nil {
		return nil, err
	}

	// create the depositing transaction
	tx, err := v.DepositERC20(auth, tokenAddress, amountBigInt, incAddress)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()
	log.Printf("%v deposited tx: %v\n", prefix, txHash)

	if err := acc.wait(txHash); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// Shield shields an amount of ETH/ERC20 tokens to the Incognito network.
// This function should be called after the DepositETH or DepositERC20 has finished.
func (acc EVMAccount) Shield(privateKey, pTokenID string, ethTxHashStr string) (string, error) {
	prefix := "[Shield]"

	ethTxHash := common.HexToHash(ethTxHashStr)
	receipt, err := acc.evmConfig.ethClient.TransactionReceipt(context.Background(), ethTxHash)
	if err != nil {
		return "", err
	}
	blockNumber := receipt.BlockNumber.Uint64()
	log.Printf("%v shieldedBlock: %v\n", prefix, blockNumber)
	for {
		header, err := acc.evmConfig.ethClient.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return "", err
		}

		if header.Number.Uint64() > blockNumber+15 {
			log.Println("[Shield] Enough 15 confirmations!!")
			break
		}
		log.Printf("%v currentEVMBlock: %v\n", prefix, header.Number.Uint64())
		time.Sleep(30 * time.Second)
	}

	depositProof, _, err := acc.evmConfig.incClient.GetEVMDepositProof(ethTxHash.String())
	if err != nil {
		return "", err
	}

	encodedTx, incTxHash, err := acc.evmConfig.incClient.CreateIssuingEVMRequestTransaction(privateKey, pTokenID, *depositProof)
	if err != nil {
		return "", err
	}

	tx := new(tx_ver2.Tx)
	rawTxData, _, err := base58.Base58Check{}.Decode(string(encodedTx))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(rawTxData, &tx)
	if err != nil {
		return "", err
	}

	md := tx.GetMetadata().(*metadata.IssuingEVMRequest)
	_, err = acc.verifyProofAndParseReceipt(md)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " Verify proof locally SUCCEEDED!!!")

	err = acc.evmConfig.incClient.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " SendRawTx SUCCEEDED!!")

	return incTxHash, nil
}

// UnShield submits a burn proof of the given incTxHash to the smart contract to obtain back a public token.
func (acc EVMAccount) UnShield(incTxHash string, gasLimit, gasPrice uint64, isBSC bool) (*common.Hash, error) {
	prefix := "[UnShield]"
	evmClient := acc.evmConfig.ethClient
	vaultAddress := acc.evmConfig.ethVaultAddress
	if isBSC {
		evmClient = acc.evmConfig.bscClient
		vaultAddress = acc.evmConfig.bscVaultAddress
	}

	// load the vault instance
	v, err := vault.NewVault(vaultAddress, evmClient)
	if err != nil {
		return nil, err
	}

	balance, err := acc.GetBalance(common.HexToAddress(nativeToken), isBSC)
	if err != nil {
		return nil, err
	}

	// retrieve the burn proof from the incTxHash
	burnProofResult, err := acc.evmConfig.incClient.GetBurnProof(incTxHash)
	if err != nil {
		return nil, err
	}
	burnProof, err := incclient.DecodeBurnProof(burnProofResult)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = evmClient.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}
	if gasLimit == 0 {
		gasLimit, err = acc.EstimateWithdrawalGas(burnProof, isBSC)
		if err != nil {
			return nil, err
		}
	}

	txFee := gasLimit * gasPriceBigInt.Uint64()
	log.Printf("%v gasLimit %v, gasPrice %v, txFee %v\n", prefix, gasLimit, gasPriceBigInt.Uint64(), txFee)
	if balance < txFee {
		return nil, fmt.Errorf("%v balance insufficient, need %v, got %v", prefix, txFee, balance)
	}

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, 0, []byte{})
	tx, err := v.Withdraw(auth,
		burnProof.Instruction,
		burnProof.Heights[0],
		burnProof.InstPaths[0],
		burnProof.InstPathIsLefts[0],
		burnProof.InstRoots[0],
		burnProof.BlkData[0],
		burnProof.SigIndices[0],
		burnProof.SigVs[0],
		burnProof.SigRs[0],
		burnProof.SigSs[0])
	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()
	log.Printf("%v WithdrawTx: %v\n", prefix, txHash.String())

	if err := acc.wait(txHash); err != nil {
		return nil, err
	}
	return &txHash, nil
}