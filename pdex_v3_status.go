package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
)

// pDEXTradeStatus retrieves the status of a pDEX trade.
func pDEXTradeStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckTradeStatus(txHash)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}

// pDEXContributionStatus retrieves the status of a pDEX liquidity contribution.
func pDEXContributionStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXLiquidityContributionStatus(txHash)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}

// pDEXOrderAddingStatus retrieves the status of an order-book adding transaction.
func pDEXOrderAddingStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckOrderAddingStatus(txHash)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}

// pDEXWithdrawalStatus retrieves the status of a pDEX liquidity withdrawal.
func pDEXWithdrawalStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXLiquidityWithdrawalStatus(txHash)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}

// pDEXOrderWithdrawalStatus retrieves the status of an order-book withdrawal.
func pDEXOrderWithdrawalStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckOrderWithdrawalStatus(txHash)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}

// pDEXMintNFTStatus gets the status of a (c)NFT minting transaction.
func pDEXMintNFTStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckNFTMintingStatus(txHash)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(status, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}
