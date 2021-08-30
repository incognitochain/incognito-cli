package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/urfave/cli/v2"
	"log"
)

// portalUnShield creates and sends a port un-shielding transaction.
func portalUnShield(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if tokenIDStr == common.PRVIDStr {
		tokenIDStr = "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696" // main-net tokenID of BTC
	}
	if !isValidTokenID(tokenIDStr) {
		return fmt.Errorf("%v is invalid", tokenIDStr)
	}

	unShieldAmount := c.Uint64(amountFlag)
	if unShieldAmount == 0 {
		return fmt.Errorf("%v cannot be zero", amountFlag)
	}

	remoteAddress := c.String(remoteAddressFlag)
	if remoteAddress == "" {
		return fmt.Errorf("%v is invalid", remoteAddress)
	}

	// create a transaction to burn the Incognito token.
	txHash, err := cfg.incClient.CreateAndSendPortalUnShieldTransaction(
		privateKey,
		tokenIDStr,
		remoteAddress,
		unShieldAmount,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)
	fmt.Println("Please wait for ~30 minutes for the fund to be released!!")

	return nil
}

// getPortalUnShieldStatus returns the status of a portal un-shielding request.
func getPortalUnShieldStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	if txHash == "" {
		return fmt.Errorf("%v is invalid", txHashFlag)
	}

	status, err := cfg.incClient.GetPortalUnShieldingRequestStatus(txHash)
	if err != nil {
		return err
	}
	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	log.Println(string(jsb))

	return nil
}