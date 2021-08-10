package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/bridge-eth/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/urfave/cli/v2"
)

func checkBalance(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String("privateKey")
	if privateKey == "" {
		return fmt.Errorf("private key is invalid")
	}

	tokenIDStr := c.String("tokenID")
	if tokenIDStr == "" {
		return fmt.Errorf("tokenID is invalid")
	}

	balance, err := client.GetBalance(privateKey, tokenIDStr)
	if err != nil {
		return err
	}
	fmt.Println(balance)

	return nil
}

func keyInfo(c *cli.Context) error {
	privateKey := c.String("privateKey")
	if privateKey == "" {
		return fmt.Errorf("private key is invalid")
	}

	info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return fmt.Errorf("marshalling key info error: %v", err)
	}
	fmt.Println(string(jsb))

	return nil
}

func consolidateUTXOs(c *cli.Context) error {
	enableLog := c.Bool("enableLog")
	if enableLog {
		logFile := c.String("logFile")
		if logFile == "" || logFile == "os.Stdout" {
			incclient.Logger = incclient.NewLogger(true)
		} else {
			incclient.Logger = incclient.NewLogger(true, logFile)
		}
	}

	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String("privateKey")
	if privateKey == "" {
		return fmt.Errorf("private key is invalid")
	}

	tokenIDStr := c.String("tokenID")
	if tokenIDStr == "" {
		return fmt.Errorf("tokenID is invalid")
	}

	version := c.Int("version")
	if version < 1 || version > 2 {
		return fmt.Errorf("version is invalid")
	}

	numThreads := c.Int("numThreads")
	if numThreads == 0 {
		return fmt.Errorf("numThreads in invalid")
	}

	fmt.Printf("CONSOLIDATING tokenID %v, version %v, numThreads %v, enableLog %v\n", tokenIDStr, version, numThreads, enableLog)

	txList, err := client.Consolidate(privateKey, tokenIDStr, int8(version), numThreads)
	if err != nil {
		return err
	}
	fmt.Println("CONSOLIDATING FINISHED!!")
	fmt.Println(txList)

	return nil
}

func checkUTXOs(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String("privateKey")
	if privateKey == "" {
		return fmt.Errorf("private key is invalid")
	}

	tokenIDStr := c.String("tokenID")
	if tokenIDStr == "" {
		return fmt.Errorf("tokenID is invalid")
	}

	unSpentCoins, idxList, err := client.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return err
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

func getHistory(c *cli.Context) error {
	enableLog := c.Bool("enableLog")
	if enableLog {
		logFile := c.String("logFile")
		if logFile == "" || logFile == "os.Stdout" {
			incclient.Logger = incclient.NewLogger(true)
		} else {
			incclient.Logger = incclient.NewLogger(true, logFile)
		}
	}

	err := initClient("https://beta-fullnode.incognito.org/fullnode", incclient.MainNetETHHost, 1)
	if err != nil {
		return err
	}

	privateKey := c.String("privateKey")
	if privateKey == "" {
		return fmt.Errorf("private key is invalid")
	}

	tokenIDStr := c.String("tokenID")
	if tokenIDStr == "" {
		return fmt.Errorf("tokenID is invalid")
	}

	numThreads := c.Int("numThreads")
	if numThreads == 0 {
		return fmt.Errorf("numThreads in invalid")
	}

	csvFile := c.String("csvFile")

	historyProcessor := incclient.NewTxHistoryProcessor(client, numThreads)

	h, err := historyProcessor.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		return err
	}

	if len(csvFile) > 0 {
		err = incclient.SaveTxHistory(h, csvFile)
		if err != nil {
			return err
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

func genKeySet(c *cli.Context) error {
	w, mnemonic, err := wallet.NewMasterKey()
	if err != nil {
		return err
	}

	numShards := c.Int(numShardsFlags)
	if numShards == 0 {
		return fmt.Errorf("%v is invalid", numShardsFlags)
	}
	common.MaxShardNumber = numShards

	privateKey := w.Base58CheckSerialize(wallet.PrivateKeyType)
	info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		return err
	}

	fmt.Printf("mnemonic: %v\n", mnemonic)
	jsb, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return fmt.Errorf("marshalling key info error: %v", err)
	}
	fmt.Println(string(jsb))

	return nil
}

func submitKey(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	otaKey := c.String("otaKey")
	if otaKey == "" {
		return fmt.Errorf("ota key is invalid")
	}

	accessToken := c.String("accessToken")
	if accessToken != "" {
		fromHeight := c.Uint64("fromHeight")
		isReset := c.Bool("isReset")
		err = client.AuthorizedSubmitKey(otaKey, accessToken, fromHeight, isReset)
	} else {
		err = client.SubmitKey(otaKey)
	}

	if err != nil {
		return err
	}

	return nil
}
