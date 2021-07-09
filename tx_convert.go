package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
)

func convertUTXOs(c *cli.Context) error {
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
	if !isValidPrivateKey(privateKey) {
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

	fmt.Printf("CONVERTING tokenID %v, numThreads %v, enableLog %v\n", tokenIDStr, numThreads, enableLog)
	utxoList, _, err := client.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return err
	}
	utxoV1Count := 0
	for _, utxo := range utxoList {
		if utxo.GetVersion() == 1 {
			utxoV1Count += 1
		}
	}
	fmt.Printf("You are currently having %v UTXOs v1\n", utxoV1Count)

	if utxoV1Count == 0 {
		fmt.Printf("No UTXOs v1 left to be converted")
		return nil
	} else if utxoV1Count <= 30 {
		txHash, err := client.CreateAndSendRawConversionTransaction(privateKey, tokenIDStr)
		if err != nil {
			return err
		}

		fmt.Println("CONVERSION FINISHED!!")
		fmt.Println(txHash)

		return nil
	}

	txList, err := client.ConvertAllUTXOs(privateKey, tokenIDStr, numThreads)
	if err != nil {
		return err
	}
	fmt.Println("CONVERSION FINISHED!!")
	fmt.Println(txList)

	return nil
}
