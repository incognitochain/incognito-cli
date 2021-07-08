package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// checkRewards gets all rewards of a payment address.
func checkRewards(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	addr := c.String("address")
	if addr == "" {
		return fmt.Errorf("payment address is invalid")
	}

	rewards, err := client.GetRewardAmount(addr)
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
	err := initNetWork()
	if err != nil {
		return err
	}

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

	txHash, err := client.CreateAndSendWithDrawRewardTransaction(privateKey, addr, tokenIDStr, int8(version))
	if err != nil {
		return err
	}
	fmt.Printf("txHash: %v\n", txHash)

	return nil
}
