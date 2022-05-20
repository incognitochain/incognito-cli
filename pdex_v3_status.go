package main

import (
	"github.com/urfave/cli/v2"
)

// pDEXTradeStatus retrieves the status of a pDEX trade.
func pDEXTradeStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckTradeStatus(txHash)
	if err != nil {
		return newAppError(GetTradeStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXContributionStatus retrieves the status of a pDEX liquidity contribution.
func pDEXContributionStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXLiquidityContributionStatus(txHash)
	if err != nil {
		return newAppError(GetDexContributionStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXOrderAddingStatus retrieves the status of an order-book adding transaction.
func pDEXOrderAddingStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckOrderAddingStatus(txHash)
	if err != nil {
		return newAppError(GetOrderAddingStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXWithdrawalStatus retrieves the status of a pDEX liquidity withdrawal.
func pDEXWithdrawalStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckDEXLiquidityWithdrawalStatus(txHash)
	if err != nil {
		return newAppError(GetDexWithdrawalStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXOrderWithdrawalStatus retrieves the status of an order-book withdrawal.
func pDEXOrderWithdrawalStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckOrderWithdrawalStatus(txHash)
	if err != nil {
		return newAppError(GetOrderWithdrawalStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXStakingStatus retrieves the status of a staking transaction.
func pDEXStakingStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXStakingStatus(txHash)
	if err != nil {
		return newAppError(GetDexStakingStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXUnStakingStatus retrieves the status of a pDEX un-staking transaction.
func pDEXUnStakingStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckDEXUnStakingStatus(txHash)
	if err != nil {
		return newAppError(GetDexUnStakingStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXWithdrawStakingRewardStatus retrieves the status of a pDEX staking reward withdrawal transaction.
func pDEXWithdrawStakingRewardStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	status, err := cfg.incClient.CheckDEXStakingRewardWithdrawalStatus(txHash)
	if err != nil {
		return newAppError(GetDexStakingRewardWithdrawalStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXWithdrawLPFeeStatus retrieves the status of a pDEX LP fee withdrawal transaction.
func pDEXWithdrawLPFeeStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckDEXLPFeeWithdrawalStatus(txHash)
	if err != nil {
		return newAppError(GetLPFeeWithdrawalStatusError, err)
	}

	return jsonPrint(status)
}

// pDEXMintNFTStatus gets the status of a pDEx NFT minting transaction.
func pDEXMintNFTStatus(c *cli.Context) error {
	txHash := c.String(txHashFlag)
	if !isValidIncTxHash(txHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}
	status, err := cfg.incClient.CheckNFTMintingStatus(txHash)
	if err != nil {
		return newAppError(GetNFTMintingStatusError, err)
	}

	return jsonPrint(status)
}
