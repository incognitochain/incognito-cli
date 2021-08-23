package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
	"time"
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

	incclient.Logger.Printf("CONVERTING tokenID %v, numThreads %v, enableLog %v\n", tokenIDStr, numThreads, enableLog)
	utxoList, _, err := cfg.incClient.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return err
	}
	utxoV1Count := 0
	for _, utxo := range utxoList {
		if utxo.GetVersion() == 1 {
			utxoV1Count += 1
		}
	}
	incclient.Logger.Printf("You are currently having %v UTXOs v1\n", utxoV1Count)

	if utxoV1Count == 0 {
		incclient.Logger.Printf("No UTXOs v1 left to be converted")
		return nil
	} else if utxoV1Count <= 30 {
		txHash, err := cfg.incClient.CreateAndSendRawConversionTransaction(privateKey, tokenIDStr)
		if err != nil {
			return err
		}

		incclient.Logger.Println("CONVERSION FINISHED!!")
		incclient.Logger.Println(txHash)

		return nil
	}

	txList, err := cfg.incClient.ConvertAllUTXOs(privateKey, tokenIDStr, numThreads)
	if err != nil {
		return err
	}
	incclient.Logger.Println("CONVERSION FINISHED!!")
	incclient.Logger.Println(txList)

	return nil
}

func convertAll(c *cli.Context) error {
	logFile := c.String("logFile")
	if logFile == "" || logFile == "os.Stdout" {
		incclient.Logger = incclient.NewLogger(true)
	} else {
		incclient.Logger = incclient.NewLogger(true, logFile)
	}

	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String("privateKey")
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("private key is invalid")
	}

	numThreads := c.Int("numThreads")
	if numThreads == 0 {
		return fmt.Errorf("numThreads in invalid")
	}

	start := time.Now()

	//Convert PRV first (if have)
	incclient.Logger.Println("CONVERTING PRV")
	txList, err := cfg.incClient.ConvertAllUTXOs(privateKey, common.PRVIDStr, numThreads)
	if err != nil {
		if err.Error() == "no UTXOs to convert" {
			incclient.Logger.Printf("no UTXOs to convert\n\n")
		} else {
			return err
		}
	} else {
		incclient.Logger.Printf("txList: %v\n\n", txList)
	}
	incclient.Logger.Printf("[CONVERTALL] timeElapsed: %v\n\n", time.Since(start).Seconds())

	incclient.Logger.Printf("CHECKING TOKENS V1...\n")
	listTokens, err := cfg.incClient.GetListToken()
	if err != nil {
		return fmt.Errorf("cannot get list tokens: %v", err)
	}
	incclient.Logger.Printf("There are currently %v tokens on the network\n", len(listTokens))
	incclient.Logger.Println("Checking v1 UTXOs of these token...")
	listTokensV1 := make(map[string]int)
	numChecked := 0
	for tokenIDStr := range listTokens {
		if numChecked%50 == 0 {
			incclient.Logger.Printf("[CONVERTALL] numChecked: %v, timeElapsed: %v\n", numChecked, time.Since(start).Seconds())
		}
		utxoList, _, err := cfg.incClient.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
		if err != nil {
			return err
		}
		utxoV1Count := 0
		for _, utxo := range utxoList {
			if utxo.GetVersion() == 1 {
				utxoV1Count += 1
			}
		}
		if utxoV1Count > 0 {
			listTokensV1[tokenIDStr] = utxoV1Count
		}
		numChecked++
	}

	incclient.Logger.Printf("There are %v tokens with UTXOs v1\n", len(listTokensV1))
	numConverted := 0
	for tokenIDStr := range listTokensV1 {
		incclient.Logger.Printf("[%v] CONVERTING TOKEN %v\n", numConverted, tokenIDStr)
		txList, err = cfg.incClient.ConvertAllUTXOs(privateKey, tokenIDStr, numThreads)
		if err != nil {
			if err.Error() == "no UTXOs to convert" {
				incclient.Logger.Printf("no UTXOs to convert\n")
			} else {
				return err
			}
		} else {
			incclient.Logger.Printf("txList: %v\n", txList)
		}
		numConverted++
		incclient.Logger.Printf("timeElapsed: %v\n\n", time.Since(start))
	}

	incclient.Logger.Printf("SUCCESS!! TIME: %v\n", time.Since(start).Seconds())

	return nil
}
