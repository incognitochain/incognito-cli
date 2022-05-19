package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
	"log"
)

func convertUTXOs(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if tokenIDStr == "" {
		return newAppError(InvalidTokenIDError)
	}

	numThreads := c.Int(numThreadsFlag)
	if numThreads == 0 {
		return newAppError(NumThreadsError)
	}

	log.Printf("CONVERTING tokenID %v, numThreads %v\n", tokenIDStr, numThreads)
	utxoList, _, err := cfg.incClient.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return newAppError(GetUnspentOutputCoinsError, err)
	}
	utxoV1Count := 0
	for _, utxo := range utxoList {
		if utxo.GetVersion() == 1 {
			utxoV1Count += 1
		}
	}
	log.Printf("You are currently having %v UTXOs v1\n", utxoV1Count)

	if utxoV1Count == 0 {
		incclient.Logger.Printf("No UTXOs v1 left to be converted")
		return nil
	} else if utxoV1Count <= 30 {
		txHash, err := cfg.incClient.CreateAndSendRawConversionTransaction(privateKey, tokenIDStr)
		if err != nil {
			return newAppError(CreateConversionTransactionError, err)
		}

		incclient.Logger.Println("CONVERSION FINISHED!!")
		incclient.Logger.Println(txHash)

		return nil
	}

	txList, err := cfg.incClient.ConvertAllUTXOs(privateKey, tokenIDStr, numThreads)
	if err != nil {
		return newAppError(CreateConversionTransactionError, err)
	}
	log.Println("CONVERSION FINISHED!!")
	log.Println(txList)

	return nil
}
