package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
)

// send creates and sends a transaction from one wallet to another w.r.t a tokenID.
func send(c *cli.Context) error {
	incclient.Logger.IsEnable = true
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String("privateKey")
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("private key is invalid")
	}

	address := c.String("address")
	if !isValidAddress(address) {
		return fmt.Errorf("receiver address is not valid")
	}

	tokenIDStr := c.String("tokenID")
	if !isValidTokenID(tokenIDStr) {
		return fmt.Errorf("tokenID is invalid")
	}

	amount := c.Uint64("amount")
	if amount == 0 {
		return fmt.Errorf("amount cannot be zero")
	}

	fee := c.Uint64("fee")
	if fee == 0 {
		return fmt.Errorf("fee cannot be zero")
	}

	version := c.Int("version")
	if !isSupportedVersion(int8(version)) {
		return fmt.Errorf("version is not supported")
	}

	fmt.Printf("Send %v of token %v from %v to %v with version %v\n", amount, tokenIDStr, privateKey, address, version)

	var txHash string
	if tokenIDStr == common.PRVIDStr {
		txHash, err = client.CreateAndSendRawTransaction(privateKey,
			[]string{address},
			[]uint64{amount},
			int8(version), nil)
	} else {
		txHash, err = client.CreateAndSendRawTokenTransaction(privateKey,
			[]string{address},
			[]uint64{amount},
			tokenIDStr,
			int8(version), nil)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Success!! TxHash %v\n", txHash)

	return nil
}
