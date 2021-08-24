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

// DepositNative deposits an amount of ETH/BNB to the Incognito contract.
func (acc EVMAccount) DepositNative(incAddress string, depositedAmount, gasLimit, gasPrice uint64, isOnBSC ...bool) (*common.Hash, error) {
	prefix := "[DepositETH]"
	evmClient := cfg.ethClient
	vaultAddress := cfg.ethVaultAddress

	isBSC := false
	if len(isOnBSC) != 0 && isOnBSC[0] {
		isBSC = true
		prefix = "[DepositBNB]"
		evmClient = cfg.bscClient
		vaultAddress = cfg.bscVaultAddress
	}

	balance, err := acc.GetBalance(common.HexToAddress(nativeToken), isBSC)
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
	gasLimit, err = acc.EstimateDepositGas(common.HexToAddress(nativeToken), amountBigInt, incAddress, isBSC)
	if err != nil {
		return nil, err
	}
	txFee := gasLimit * gasPriceBigInt.Uint64()
	log.Printf("%v gasLimit %v, gasPrice %v, txFee %v\n", prefix, gasLimit, gasPriceBigInt.Uint64(), txFee)

	requiredAmount := depositedAmount + txFee
	if balance <= requiredAmount {
		return nil, fmt.Errorf("[DepositNative] balance insufficient, need %v, got %v", requiredAmount, balance)
	}

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, amountBigInt.Uint64(), nil, isBSC)
	if err != nil {
		return nil, err
	}

	tx, err := v.Deposit(auth, incAddress)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()
	log.Printf("%v deposited tx: %v\n", prefix, txHash.String())

	if err := acc.wait(txHash, isBSC); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// DepositToken shields an amount of ERC20/BEP20 to the Incognito network.
func (acc EVMAccount) DepositToken(incAddress, tokenAddressStr string, depositedAmount, gasLimit, gasPrice uint64, isOnBSC ...bool) (*common.Hash, error) {
	prefix := "[DepositERC20]"
	isBSC := false
	if len(isOnBSC) != 0 && isOnBSC[0] {
		isBSC = true
		prefix = "[DepositBEP20]"
	}

	// load the vault address
	vaultAddress := cfg.ethVaultAddress
	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
		vaultAddress = cfg.bscVaultAddress
	}

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

		err = acc.wait(*txHash, isBSC)
		if err != nil {
			return nil, err
		}
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

	// estimate gasLimit
	amountBigInt := new(big.Int).SetUint64(depositedAmount)
	gasLimit, err = acc.EstimateDepositGas(tokenAddress, amountBigInt, incAddress, isBSC)
	if err != nil {
		return nil, err
	}
	txFee := gasLimit * gasPriceBigInt.Uint64()
	log.Printf("%v gasLimit %v, gasPrice %v, txFee %v\n", prefix, gasLimit, gasPriceBigInt.Uint64(), txFee)
	if balance < txFee {
		return nil, fmt.Errorf("%v balance insufficient, need %v, got %v", prefix, txFee, balance)
	}

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, 0, nil, isBSC)
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

	if err := acc.wait(txHash, isBSC); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// Shield shields an amount of ETH/ERC20 tokens to the Incognito network.
// This function should be called after the DepositNative or DepositToken has finished.
func (acc EVMAccount) Shield(privateKey, pTokenID string, ethTxHashStr string, isOnBSC ...bool) (string, error) {
	prefix := "[Shield]"

	isBSC := false
	if len(isOnBSC) != 0 && isOnBSC[0] {
		isBSC = true
	}

	evmClient := cfg.ethClient
	if isBSC {
		evmClient = cfg.bscClient
	}

	ethTxHash := common.HexToHash(ethTxHashStr)
	receipt, err := evmClient.TransactionReceipt(context.Background(), ethTxHash)
	if err != nil {
		return "", err
	}
	blockNumber := receipt.BlockNumber.Uint64()
	log.Printf("%v shieldedBlock: %v\n", prefix, blockNumber)
	log.Printf("%v Wait for 15 confirmations\n", prefix)

	ctx, cancel := context.WithTimeout(context.Background(), 120 * time.Second)
	defer cancel()
	for {
		header, err := evmClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return "", err
		}

		if header.Number.Uint64() > blockNumber+15 {
			log.Println(prefix, "Enough confirmations!!")
			break
		}
		log.Printf("%v currentEVMBlock: %v\n", prefix, header.Number.Uint64())
		time.Sleep(30 * time.Second)
	}

	depositProof, _, err := cfg.incClient.GetEVMDepositProof(ethTxHash.String(), isBSC)
	if err != nil {
		return "", err
	}

	encodedTx, incTxHash, err := cfg.incClient.CreateIssuingEVMRequestTransaction(privateKey, pTokenID, *depositProof, isBSC)
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
	_, err = acc.verifyProofAndParseReceipt(md, isBSC)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " Verify proof locally SUCCEEDED!!!")

	err = cfg.incClient.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " SendRawTx SUCCEEDED!!")

	return incTxHash, nil
}

// UnShield submits a burn proof of the given incTxHash to the smart contract to obtain back a public token.
func (acc EVMAccount) UnShield(incTxHash string, gasLimit, gasPrice uint64, isOnBSC ...bool) (*common.Hash, error) {
	prefix := "[UnShield]"

	isBSC := false
	if len(isOnBSC) != 0 && isOnBSC[0] {
		isBSC = true
	}

	evmClient := cfg.ethClient
	vaultAddress := cfg.ethVaultAddress
	if isBSC {
		evmClient = cfg.bscClient
		vaultAddress = cfg.bscVaultAddress
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
	burnProofResult, err := cfg.incClient.GetBurnProof(incTxHash, isBSC)
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

	auth, err := acc.NewTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, 0, []byte{}, isBSC)
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

	if err := acc.wait(txHash, isBSC); err != nil {
		return nil, err
	}
	return &txHash, nil
}