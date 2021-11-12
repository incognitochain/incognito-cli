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

// pDEXStakingStatus retrieves the status of a staking transaction.
func pDEXStakingStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXStakingStatus(txHash)
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

// pDEXUnStakingStatus retrieves the status of a pDEX un-staking transaction.
func pDEXUnStakingStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXUnStakingStatus(txHash)
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

// pDEXWithdrawStakingRewardStatus retrieves the status of a pDEX staking reward withdrawal transaction.
func pDEXWithdrawStakingRewardStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXStakingRewardWithdrawalStatus(txHash)
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

// pDEXWithdrawLPFeeStatus retrieves the status of a pDEX LP fee withdrawal transaction.
func pDEXWithdrawLPFeeStatus(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXLPFeeWithdrawalStatus(txHash)
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

// pDEXMintNFTStatus gets the status of a pDEx NFT minting transaction.
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
