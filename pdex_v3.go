package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// pDEXCheckPrice checks the price of two tokenIds.
func pDEXCheckPrice(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return fmt.Errorf("%v is invalid", tokenIDToSellFlag)
	}

	tokenIdToBuy := c.String(tokenIDToBuyFlag)
	if !isValidTokenID(tokenIdToBuy) {
		return fmt.Errorf("%v is invalid", tokenIDToBuyFlag)
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return fmt.Errorf("%v cannot be zero", sellingAmountFlag)
	}

	pairID := c.String(pairIDFlag)
	bestExpectedReceive := uint64(0)
	if pairID != "" {
		pairs, err := cfg.incClient.GetPdexPoolPair(0, tokenIdToSell, tokenIdToBuy)
		if err != nil {
			return err
		}
		for path, _ := range pairs {
			expectedPrice, err := cfg.incClient.CheckPrice(path, tokenIdToSell, sellingAmount)
			if err != nil {
				fmt.Println(path, err)
				continue
			}
			if expectedPrice > bestExpectedReceive {
				bestExpectedReceive = expectedPrice
				pairID = path
			}
		}
	} else {
		bestExpectedReceive, err = cfg.incClient.CheckPrice(pairID, tokenIdToSell, sellingAmount)
		if err != nil {
			return err
		}
	}

	if bestExpectedReceive == 0 {
		return fmt.Errorf("cannot find a proper path")
	}

	fmt.Printf("bestPairID %v: %v\n", pairID, bestExpectedReceive)
	return nil
}

// pDEXTrade creates and sends a trade to the pDEX.
func pDEXTrade(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return fmt.Errorf("%v is invalid", tokenIDToSellFlag)
	}

	tokenIdToBuy := c.String(tokenIDToBuyFlag)
	if !isValidTokenID(tokenIdToBuy) {
		return fmt.Errorf("%v is invalid", tokenIDToBuyFlag)
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return fmt.Errorf("%v cannot be zero", sellingAmountFlag)
	}

	minAcceptableAmount := c.Uint64(minAcceptableAmountFlag)
	tradingFee := c.Uint64(tradingFeeFlag)

	prvFee := c.Int(prvFeeFlag)

	txHash, err := cfg.incClient.CreateAndSendPdexv3TradeTransaction(
		privateKey,
		[]string{},
		tokenIdToSell,
		tokenIdToBuy,
		sellingAmount,
		minAcceptableAmount,
		tradingFee,
		prvFee != 0,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)

	return nil
}