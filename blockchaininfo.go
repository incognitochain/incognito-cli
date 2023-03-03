package main

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
)

// stake creates a staking transaction.
func getBeaconStaker(c *cli.Context) (err error) {
	cPK := c.String(committeePubKeyFlag)
	if cPK == "" {
		privateKey := c.String(privateKeyFlag)
		if !isValidPrivateKey(privateKey) {
			return newAppError(InvalidPrivateKeyError)
		}
		miningKey := c.String(miningKeyFlag)
		if miningKey == "" {
			miningKey = incclient.PrivateKeyToMiningKey(privateKey)
		}
		cPK, err = incclient.CommitteePublicKeyFromPrivateKey(privateKey, miningKey)
		if err != nil {
			return err
		}
	}

	res, err := cfg.incClient.GetBeaconStaker(6466590, cPK)
	if err != nil {
		return newAppError(CreateStakingTransactionError, err)
	}
	fmt.Printf("%s", res)

	return nil
}

func getShardStaker(c *cli.Context) (err error) {
	cPK := c.String(committeePubKeyFlag)
	if cPK == "" {
		privateKey := c.String(privateKeyFlag)
		if !isValidPrivateKey(privateKey) {
			return newAppError(InvalidPrivateKeyError)
		}
		miningKey := c.String(miningKeyFlag)
		if miningKey == "" {
			miningKey = incclient.PrivateKeyToMiningKey(privateKey)
		}
		cPK, err = incclient.CommitteePublicKeyFromPrivateKey(privateKey, miningKey)
		if err != nil {
			return err
		}
	}

	res, err := cfg.incClient.GetShardStaker(6466590, cPK)
	if err != nil {
		return newAppError(CreateStakingTransactionError, err)
	}
	fmt.Printf("%s", res)

	return nil
}

func getBeaconCommitteeState(c *cli.Context) (err error) {
	beaconHeight := c.Uint64("beaconheight")

	res, err := cfg.incClient.GetBeaconCommitteeState(beaconHeight)
	if err != nil {
		return newAppError(CreateStakingTransactionError, err)
	}
	fmt.Printf("%s", res)

	return nil
}
