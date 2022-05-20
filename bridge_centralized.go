package main

import (
	"github.com/urfave/cli/v2"
)

func shieldCentralized(c *cli.Context) error {
	adminPrivateKey := c.String(adminPrivateKeyFlag)
	if !isValidPrivateKey(adminPrivateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	receiver := c.String(addressFlag)
	if !isValidAddress(receiver) {
		return newAppError(InvalidPaymentAddressError)
	}

	amt := c.Uint64(amountFlag)
	if amt == 0 {
		return newAppError(InvalidAmountError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	tokenName := c.String(tokenNameFlag)

	txHash, err := cfg.incClient.CreateAndSendIssuingRequestTransaction(adminPrivateKey,
		receiver, tokenIDStr, tokenName, amt)
	if err != nil {
		return newAppError(CentralizedShieldError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}
