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
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}
	canAddr := c.String(candidateAddressFlag)
	if canAddr == "" {
		canAddr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(canAddr) {
		return fmt.Errorf("%v is invalid", candidateAddressFlag)
	}
	rewardAddr := c.String(rewardReceiverFlag)
	if rewardAddr == "" {
		rewardAddr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(rewardAddr) {
		return fmt.Errorf("%v is invalid", rewardReceiverFlag)
	}
	miningKey := c.String(miningKeyFlag)
	if miningKey == "" {
		miningKey = incclient.PrivateKeyToMiningKey(privateKey)
	}
	if !isValidMiningKey(miningKey) {
		return fmt.Errorf("%v is invalid", miningKeyFlag)
	}
	reStake := c.Int(autoReStakeFlag)
	autoReStake := reStake != 0

	txHash, err := cfg.incClient.CreateAndSendShardStakingTransaction(privateKey, miningKey, canAddr, rewardAddr, autoReStake)
	if err != nil {
		return err
	}
	fmt.Printf("txHash: %v\n", txHash)

	return nil
}

// unStake creates an un-staking transaction.
func unStake(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}
	canAddr := c.String(candidateAddressFlag)
	if canAddr == "" {
		canAddr = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(canAddr) {
		return fmt.Errorf("%v is invalid", candidateAddressFlag)
	}
	miningKey := c.String(miningKeyFlag)
	if miningKey == "" {
		miningKey = incclient.PrivateKeyToMiningKey(privateKey)
	}
	if !isValidMiningKey(miningKey) {
		return fmt.Errorf("%v is invalid", miningKeyFlag)
	}

	txHash, err := cfg.incClient.CreateAndSendUnStakingTransaction(privateKey, miningKey, canAddr)
	if err != nil {
		return err
	}
	fmt.Printf("txHash: %v\n", txHash)

	return nil
}

// checkRewards gets all rewards of a payment address.
func checkRewards(c *cli.Context) error {
	addr := c.String("address")
	if addr == "" {
		return fmt.Errorf("payment address is invalid")
	}

	rewards, err := cfg.incClient.GetRewardAmount(addr)
	if err != nil {
		return err
	}

	if len(rewards) == 0 {
		fmt.Printf("There is not rewards found for the address %v\n", addr)
	} else {
		fmt.Printf("Rewards of address %v:\n", addr)
		for tokenID, amount := range rewards {
			fmt.Printf("%v: %v\n", tokenID, amount)
		}
	}

	return nil
}

// withdrawReward withdraws the reward of a privateKey w.r.t to a tokenID.
func withdrawReward(c *cli.Context) error {
	privateKey := c.String("privateKey")
	if privateKey == "" {
		return fmt.Errorf("private key is invalid")
	}

	addr := c.String("address")

	tokenIDStr := c.String("tokenID")
	if tokenIDStr == "" {
		return fmt.Errorf("tokenID is invalid")
	}

	version := c.Int("version")
	if version != 1 && version != 2 {
		return fmt.Errorf("version must be 1 or 2")
	}

	fmt.Printf("Withdrawing the reward for tokenID %v, using tx version %v\n", tokenIDStr, version)

	txHash, err := cfg.incClient.CreateAndSendWithDrawRewardTransaction(privateKey, addr, tokenIDStr, int8(version))
	if err != nil {
		return err
	}
	fmt.Printf("txHash: %v\n", txHash)

	return nil
}
