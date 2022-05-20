package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/urfave/cli/v2"
)

// send creates and sends a transaction from one wallet to another w.r.t a tokenID.
func send(c *cli.Context) error {
	var err error

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	address := c.String(addressFlag)
	if !isValidAddress(address) {
		return newAppError(InvalidPaymentAddressError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return newAppError(InvalidAmountError)
	}

	version := c.Int(versionFlag)
	if !isSupportedVersion(int8(version)) {
		return newAppError(VersionError)
	}

	fmt.Printf("Send %v of token %v from %v to %v with version %v\n", amount, tokenIDStr, privateKey, address, version)

	var txHash string
	if tokenIDStr == common.PRVIDStr {
		txHash, err = cfg.incClient.CreateAndSendRawTransaction(privateKey,
			[]string{address},
			[]uint64{amount},
			int8(version), nil)
	} else {
		txHash, err = cfg.incClient.CreateAndSendRawTokenTransaction(privateKey,
			[]string{address},
			[]uint64{amount},
			tokenIDStr,
			int8(version), nil)
	}
	if err != nil {
		return newAppError(CreateTransferTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// checkReceiver if a user is a receiver of a transaction.
func checkReceiver(c *cli.Context) error {
	var err error

	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}

	otaKey := c.String(otaKeyFlag)
	if !isValidOtaKey(otaKey) {
		return newAppError(InvalidOTAKeyError)
	}

	readonlyKey := c.String(readonlyKeyFlag)
	if readonlyKey != "" && !isValidReadonlyKey(readonlyKey) {
		return newAppError(InvalidReadonlyKeyError)
	}

	var received bool
	var res map[string]uint64
	if readonlyKey == "" {
		received, res, err = cfg.incClient.GetReceivingInfo(txHash, otaKey)
	} else {
		received, res, err = cfg.incClient.GetReceivingInfo(txHash, otaKey, readonlyKey)
	}

	if err != nil {
		return newAppError(GetReceivingInfoError, err)
	}

	type receivingInfo struct {
		Received      bool
		ReceivingInfo map[string]uint64 `json:"ReceivingInfo"`
	}

	return jsonPrint(receivingInfo{Received: received, ReceivingInfo: res})
}
