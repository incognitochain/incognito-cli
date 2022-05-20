package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
)

// stake creates a staking transaction.
func stake(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}
	canAddr := c.String(candidateAddressFlag)
	if canAddr == "" {
		canAddr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(canAddr) {
		return newAppError(InvalidPaymentAddressError, fmt.Errorf("%v", canAddr))
	}
	rewardAddr := c.String(rewardReceiverFlag)
	if rewardAddr == "" {
		rewardAddr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(rewardAddr) {
		return newAppError(InvalidPaymentAddressError, fmt.Errorf("%v", rewardAddr))
	}
	miningKey := c.String(miningKeyFlag)
	if miningKey == "" {
		miningKey = incclient.PrivateKeyToMiningKey(privateKey)
	}
	if !isValidMiningKey(miningKey) {
		return newAppError(InvalidMiningKeyError)
	}
	reStake := c.Int(autoReStakeFlag)
	autoReStake := reStake != 0

	txHash, err := cfg.incClient.CreateAndSendShardStakingTransaction(privateKey, miningKey, canAddr, rewardAddr, autoReStake)
	if err != nil {
		return newAppError(CreateStakingTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// unStake creates an un-staking transaction.
func unStake(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	canAddr := c.String(candidateAddressFlag)
	if canAddr == "" {
		canAddr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(canAddr) {
		return newAppError(InvalidPaymentAddressError)
	}

	miningKey := c.String(miningKeyFlag)
	if miningKey == "" {
		miningKey = incclient.PrivateKeyToMiningKey(privateKey)
	}
	if !isValidMiningKey(miningKey) {
		return newAppError(InvalidMiningKeyError)
	}

	txHash, err := cfg.incClient.CreateAndSendUnStakingTransaction(privateKey, miningKey, canAddr)
	if err != nil {
		return newAppError(CreateUnStakingTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// checkRewards gets all rewards of a payment address.
func checkRewards(c *cli.Context) error {
	addr := c.String(addressFlag)
	if !isValidAddress(addr) {
		return newAppError(InvalidPaymentAddressError)
	}

	rewards, err := cfg.incClient.GetRewardAmount(addr)
	if err != nil {
		return newAppError(GetRewardAmountError, err)
	}

	if len(rewards) == 0 {
		fmt.Printf("There is not rewards found for the address %v\n", addr)
	} else {
		return jsonPrint(rewards)
	}

	return nil
}

// withdrawReward withdraws the reward of a privateKey w.r.t to a tokenID.
func withdrawReward(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	addr := c.String(addressFlag)
	if addr == "" {
		addr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(addr) {
		return newAppError(InvalidPaymentAddressError)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return newAppError(InvalidTokenIDError)
	}

	version := c.Int(versionFlag)
	if !isSupportedVersion(int8(version)) {
		return newAppError(VersionError)
	}

	fmt.Printf("Withdrawing the reward for tokenID %v, using tx version %v\n", tokenIDStr, version)

	txHash, err := cfg.incClient.CreateAndSendWithDrawRewardTransaction(privateKey, addr, tokenIDStr, int8(version))
	if err != nil {
		return newAppError(CreateWithdrawRewardTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}
