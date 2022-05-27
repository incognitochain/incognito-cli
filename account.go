package main

import (
	"encoding/csv"
	"fmt"
	"github.com/incognitochain/bridge-eth/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

func checkBalance(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	balance, err := cfg.incClient.GetBalance(privateKey, tokenIDStr)
	if err != nil {
		return newAppError(GetBalanceError, err)
	}

	return jsonPrintWithKey("Balance", balance)
}

func getAllBalanceV2(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	balances, err := cfg.incClient.GetAllBalancesV2(privateKey)
	if err != nil {
		return newAppError(GetAllBalancesError, err)
	}

	return jsonPrint(balances)
}

func keyInfo(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		return newAppError(GetAccountInfoError, err)
	}

	return jsonPrint(info)
}

func consolidateUTXOs(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	version := c.Int(versionFlag)
	if version < 1 || version > 2 {
		return newAppError(VersionError)
	}

	numThreads := c.Int(numThreadsFlag)
	if numThreads == 0 {
		return newAppError(NumThreadsError)
	}

	fmt.Printf("CONSOLIDATING tokenID %v, version %v, numThreads %v\n", tokenIDStr, version, numThreads)

	txList, err := cfg.incClient.Consolidate(privateKey, tokenIDStr, int8(version), numThreads)
	if err != nil {
		return newAppError(ConsolidateAccountError, err)
	}
	fmt.Println("CONSOLIDATING FINISHED!!")

	return jsonPrintWithKey("TxList", txList)
}

func checkUTXOs(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	unSpentCoins, idxList, err := cfg.incClient.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return newAppError(GetUnspentOutputCoinsError, err)
	}

	numUTXOsV1 := 0
	numUTXOsV2 := 0
	balanceV1 := uint64(0)
	balanceV2 := uint64(0)

	for i, utxo := range unSpentCoins {
		if utxo.GetVersion() == 1 {
			numUTXOsV1++
			balanceV1 += utxo.GetValue()
		} else {
			numUTXOsV2++
			balanceV2 += utxo.GetValue()
		}

		fmt.Printf("idx %v, version %v, pubKey %v, keyImage %v, value %v\n",
			idxList[i].Uint64(), utxo.GetVersion(),
			base58.Base58Check{}.Encode(utxo.GetPublicKey().ToBytesS(), 0),
			base58.Base58Check{}.Encode(utxo.GetKeyImage().ToBytesS(), 0),
			utxo.GetValue())
	}

	fmt.Printf("#numUTXOsV1 %v, #numUTXOsV2 %v\n", numUTXOsV1, numUTXOsV2)
	fmt.Printf("balanceV1 %v, balanceV2 %v, totalBalance %v\n", balanceV1, balanceV2, balanceV1+balanceV2)

	return nil
}

func getOutCoins(c *cli.Context) error {
	address := c.String(addressFlag)
	if !isValidAddress(address) {
		return newAppError(InvalidPaymentAddressError)
	}

	otaKey := c.String(otaKeyFlag)
	if !isValidOtaKey(otaKey) {
		return newAppError(InvalidOTAKeyError)
	}

	readonlyKey := c.String(readonlyKeyFlag)
	if readonlyKey != "" && !isValidReadonlyKey(readonlyKey) {
		return newAppError(InvalidReadonlyKeyError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	outCoinKey := new(rpc.OutCoinKey)
	outCoinKey.SetPaymentAddress(address)
	outCoinKey.SetOTAKey(otaKey)
	outCoinKey.SetReadonlyKey(readonlyKey)

	outCoins, idxList, err := cfg.incClient.GetOutputCoins(outCoinKey, tokenIDStr, 0)
	if err != nil {
		return newAppError(GetOutputCoinsError, err)
	}

	v1Count := 0
	v2Count := 0
	for i, outCoin := range outCoins {
		if outCoin.GetVersion() == 1 {
			v1Count += 1
		} else {
			v2Count += 1
		}

		fmt.Printf("idx %v, ver %v, encrypted %v, pubKey %v, cmtStr %v\n",
			idxList[i].Int64(),
			outCoin.GetVersion(),
			outCoin.IsEncrypted(),
			base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0x00),
			base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), 0x00))
	}

	fmt.Printf("#OutCoins: %v, #v1: %v, #v2: %v\n", len(outCoins), v1Count, v2Count)

	return nil
}

func getHistory(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	numThreads := c.Int(numThreadsFlag)
	if numThreads == 0 {
		return newAppError(NumThreadsError)
	}

	csvFile := c.String("csvFile")

	historyProcessor := incclient.NewTxHistoryProcessor(cfg.incClient, numThreads)

	h, err := historyProcessor.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		return newAppError(GetHistoryError, err)
	}

	if len(csvFile) > 0 {
		err = incclient.SaveTxHistory(h, csvFile)
		if err != nil {
			return newAppError(SaveHistoryError, err)
		}
	} else {
		totalIn := uint64(0)
		fmt.Printf("#TxIns %v\n", len(h.TxInList))
		for _, txIn := range h.TxInList {
			totalIn += txIn.GetAmount()
			fmt.Println(txIn.String())
		}
		fmt.Printf("END TxIns\n\n")

		totalOut := uint64(0)
		fmt.Printf("#TxOuts %v\n", len(h.TxOutList))
		for _, txOut := range h.TxOutList {
			totalOut += txOut.GetAmount()
			fmt.Println(txOut.String())
		}
		fmt.Printf("END TxOuts\n")

		fmt.Printf("TotalIn: %v, TotalOut: %v\n", totalIn, totalOut)
	}

	return nil
}

func financialExport(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	numThreads := c.Int(numThreadsFlag)
	if numThreads == 0 {
		return newAppError(NumThreadsError)
	}

	csvFile := c.String(csvFileFlag)
	if len(csvFile) == 0 {
		csvFile = incclient.DefaultTxHistory
	}

	historyProcessor := incclient.NewTxHistoryProcessor(cfg.incClient, numThreads)

	historyMap, err := historyProcessor.GetAllHistory(privateKey)
	if err != nil {
		return newAppError(GetHistoryError, err)
	}

	history := new(incclient.TxHistory)
	history.TxInList = make([]incclient.TxIn, 0)
	history.TxOutList = make([]incclient.TxOut, 0)
	for tokenID, h := range historyMap {
		if tokenID == common.ConfidentialAssetID.String() || tokenID == common.PRVIDStr {
			continue
		}

		history.TxInList = append(history.TxInList, h.TxInList...)
		history.TxOutList = append(history.TxOutList, h.TxOutList...)
	}

	if historyMap[common.PRVIDStr] != nil {
		history.TxInList = append(history.TxInList, historyMap[common.PRVIDStr].TxInList...)

	}

	for _, txOut := range historyMap[common.PRVIDStr].TxOutList {
		if txOut.Amount == 0 {
			continue
		}
		history.TxOutList = append(history.TxOutList, txOut)
	}

	//fmt.Println(historyMap[common.PRVIDStr].TxOutList)

	sort.Slice(history.TxInList, func(i, j int) bool {
		return history.TxInList[i].LockTime > history.TxInList[j].LockTime
	})
	sort.Slice(history.TxOutList, func(i, j int) bool {
		return history.TxOutList[i].LockTime > history.TxOutList[j].LockTime
	})

	f, err := os.OpenFile(csvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return newAppError(UnexpectedError, fmt.Errorf("cannot open file %v: %v", csvFile, err))
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	w := csv.NewWriter(f)
	defer w.Flush()

	_ = f.Truncate(0)

	log.Println("Building up the history file...")

	var dateTimeFormat = "2006/01/02 15:04:05"
	historyPattern := []string{"Date", "TxHash", "Received Quantity", "Received Currency", "Sent Quantity", "Sent Currency", "Fee Amount", "Fee Currency", "Tag"}
	err = w.Write(historyPattern)
	if err != nil {
		return newAppError(SaveHistoryError, err)
	}

	writtenData := make(map[string]bool)
	for _, txIn := range history.TxInList {
		if writtenData[common.HashH([]byte(txIn.String())).String()] {
			continue
		} else {
			writtenData[common.HashH([]byte(txIn.String())).String()] = true
		}
		toBeWritten := make([]string, 0)
		toBeWritten = append(toBeWritten, time.Unix(txIn.GetLockTime(), 0).Format(dateTimeFormat))
		toBeWritten = append(toBeWritten, txIn.TxHash)
		toBeWritten = append(toBeWritten, fmt.Sprintf("%v", float64(txIn.Amount)/math.Pow10(getTokenDecimals(txIn.TokenID))))
		toBeWritten = append(toBeWritten, getTokenName(txIn.TokenID))
		toBeWritten = append(toBeWritten, "")
		toBeWritten = append(toBeWritten, "")
		toBeWritten = append(toBeWritten, "")
		toBeWritten = append(toBeWritten, "")
		toBeWritten = append(toBeWritten, txIn.Note)

		err = w.Write(toBeWritten)
		if err != nil {
			return newAppError(SaveHistoryError, fmt.Errorf("write txHash %v error: %v", txIn.GetTxHash(), err))
		}
	}

	for _, txOut := range history.TxOutList {
		if writtenData[common.HashH([]byte(txOut.String())).String()] {
			continue
		} else {
			writtenData[common.HashH([]byte(txOut.String())).String()] = true
		}
		fee := txOut.PRVFee
		tokenFee := common.PRVIDStr
		if fee == 0 {
			fee = txOut.TokenFee
			tokenFee = txOut.TokenID
		}
		if fee == 0 {
			tokenFee = ""
		}

		toBeWritten := make([]string, 0)
		toBeWritten = append(toBeWritten, time.Unix(txOut.GetLockTime(), 0).Format(dateTimeFormat))
		toBeWritten = append(toBeWritten, txOut.TxHash)
		toBeWritten = append(toBeWritten, "")
		toBeWritten = append(toBeWritten, "")
		toBeWritten = append(toBeWritten, fmt.Sprintf("%v", float64(txOut.Amount)/math.Pow10(getTokenDecimals(txOut.TokenID))))
		toBeWritten = append(toBeWritten, getTokenName(txOut.TokenID))
		toBeWritten = append(toBeWritten, fmt.Sprintf("%v", float64(fee)/math.Pow10(getTokenDecimals(tokenFee))))
		toBeWritten = append(toBeWritten, getTokenName(tokenFee))
		toBeWritten = append(toBeWritten, txOut.Note)

		err = w.Write(toBeWritten)
		if err != nil {
			return newAppError(SaveHistoryError, fmt.Errorf("write txHash %v error: %v", txOut.GetTxHash(), err))
		}
	}
	log.Printf("Report written to file `%v`\n", csvFile)

	return nil
}

type accountInfo struct {
	Index int
	*incclient.KeyInfo
}

type masterKeyInfo struct {
	Mnemonic string `json:"Mnemonic,omitempty"`
	Accounts []*accountInfo
}

func genKeySet(c *cli.Context) error {
	w, mnemonic, err := wallet.NewMasterKey()
	if err != nil {
		return newAppError(GenerateMasterKeyError, err)
	}

	numShards := c.Int(numShardsFlag)
	if numShards == 0 {
		return newAppError(InvalidNumberShardsError)
	}
	common.MaxShardNumber = numShards

	shardID := c.Int(shardIDFlag)
	if shardID < -2 || shardID >= common.MaxShardNumber {
		return newAppError(InvalidShardError, fmt.Errorf("expected shardID from -2 to %v", common.MaxShardNumber-1))
	}
	supportedShards := make(map[byte]bool)
	if shardID == -1 {
		for i := 0; i < common.MaxShardNumber; i++ {
			supportedShards[byte(i)] = true
		}
	} else if shardID >= 0 {
		supportedShards[byte(shardID)] = true
	}

	numAccounts := c.Int(numAccountsFlag)

	accounts := make([]*accountInfo, 0)
	genCount := 0
	index := 1
	for {
		if genCount == numAccounts {
			break
		}
		childKey, err := w.DeriveChild(uint32(index))
		if err != nil {
			return newAppError(DeriveChildError, err)
		}
		privateKey := childKey.Base58CheckSerialize(wallet.PrivateKeyType)
		info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
		if err != nil {
			return newAppError(GetAccountInfoError, err)
		}
		if index == 1 && shardID == -2 {
			supportedShards[info.ShardID] = true
		}
		if supportedShards[info.ShardID] {
			accounts = append(accounts, &accountInfo{Index: index, KeyInfo: info})
			genCount++
		}

		index++
	}
	return jsonPrint(masterKeyInfo{Mnemonic: mnemonic, Accounts: accounts})
}

func importMnemonic(c *cli.Context) error {
	mnemonic := c.String(mnemonicFlag)
	mnemonic = strings.Replace(mnemonic, "-", " ", -1)
	w, err := wallet.NewMasterKeyFromMnemonic(mnemonic)
	if err != nil {
		return newAppError(ImportMnemonicError)
	}

	numShards := c.Int(numShardsFlag)
	if numShards == 0 {
		return newAppError(InvalidNumberShardsError)
	}
	common.MaxShardNumber = numShards

	shardID := c.Int(shardIDFlag)
	if shardID < -2 || shardID >= common.MaxShardNumber {
		return newAppError(InvalidShardError, fmt.Errorf("expected shardID from -2 to %v", common.MaxShardNumber-1))
	}
	supportedShards := make(map[byte]bool)
	if shardID == -1 {
		for i := 0; i < common.MaxShardNumber; i++ {
			supportedShards[byte(i)] = true
		}
	} else if shardID >= 0 {
		supportedShards[byte(shardID)] = true
	}

	numAccounts := c.Int(numAccountsFlag)

	accounts := make([]*accountInfo, 0)
	genCount := 0
	index := 1
	for {
		if genCount == numAccounts {
			break
		}
		childKey, err := w.DeriveChild(uint32(index))
		if err != nil {
			return newAppError(DeriveChildError, err)
		}
		privateKey := childKey.Base58CheckSerialize(wallet.PrivateKeyType)
		info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
		if err != nil {
			return newAppError(GetAccountInfoError, err)
		}
		if index == 1 && shardID == -2 {
			supportedShards[info.ShardID] = true
		}
		if supportedShards[info.ShardID] {
			accounts = append(accounts, &accountInfo{Index: index, KeyInfo: info})
			genCount++
		}

		index++
	}
	return jsonPrint(accounts)
}

func submitKey(c *cli.Context) error {
	var err error
	otaKey := c.String(otaKeyFlag)
	if !isValidOtaKey(otaKey) {
		return newAppError(InvalidOTAKeyError)
	}

	accessToken := c.String(accessTokenFlag)
	if accessToken != "" {
		fromHeight := c.Uint64(fromHeightFlag)
		isReset := c.Bool(isResetFlag)
		err = cfg.incClient.AuthorizedSubmitKey(otaKey, accessToken, fromHeight, isReset)
	} else {
		err = cfg.incClient.SubmitKey(otaKey)
	}

	if err != nil {
		return newAppError(SubmitKeyError, err)
	}

	return nil
}
