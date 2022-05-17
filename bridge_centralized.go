package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func shieldCentralized(c *cli.Context) error {
	adminPrivateKey := c.String(adminPrivateKeyFlag)
	if !isValidPrivateKey(adminPrivateKey) {
		return fmt.Errorf("%v is invalid", adminPrivateKey)
	}

	receiver := c.String(addressFlag)
	if !isValidAddress(receiver) {
		return fmt.Errorf("%v is invalid", addressFlag)
	}

	amt := c.Uint64(amountFlag)
	if amt == 0 {
		return fmt.Errorf("%v must not be 0", amountFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	tokenName := c.String(tokenNameFlag)

	txHash, err := cfg.incClient.CreateAndSendIssuingRequestTransaction(adminPrivateKey,
		receiver, tokenIDStr, tokenName, amt)
	if err != nil {
		return err
	}

	return jsonPrint(map[string]interface{}{"TxHash": txHash})
}
