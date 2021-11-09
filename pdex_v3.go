package main

import (
	"fmt"
	"github.com/incognitochain/incognito-cli/pdex_v3"
	"github.com/urfave/cli/v2"
	"strings"
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

	maxPaths := c.Uint(maxTradingPathLengthFlag)
	if maxPaths > pdex_v3.MaxPaths {
		return fmt.Errorf("maximum trading path length allowed %v, got %v", pdex_v3.MaxPaths, maxPaths)
	}

	allPoolPairs, err := cfg.incClient.GetAllPdexPoolPairs(0)
	if err != nil {
		return err
	}
	tmpTradingPath := c.String(tradingPathFlag)
	tradingPath := make([]string, 0)
	if tmpTradingPath != "" {
		tradingPath = strings.Split(tmpTradingPath, ",")
		for _, poolID := range tradingPath {
			if _, ok := allPoolPairs[poolID]; !ok {
				return fmt.Errorf("poolID %v not existed", poolID)
			}
		}
	} else {
		_, tradingPath, _ = pdex_v3.FindGoodTradePath(maxPaths, allPoolPairs, tokenIdToSell, tokenIdToBuy, sellingAmount)
	}
	if len(tradingPath) == 0 {
		return fmt.Errorf("no trading path is found for the pair %v-%v with maxPaths = %v", tokenIdToSell, tokenIdToBuy, maxPaths)
	}

	prvFee := c.Int(prvFeeFlag)

	txHash, err := cfg.incClient.CreateAndSendPdexv3TradeTransaction(
		privateKey,
		tradingPath,
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
	fmt.Printf("Status: %v\n", status)

	return nil
}

// pDEXFindPath finds a proper trading path.
func pDEXFindPath(c *cli.Context) error {
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

	maxPaths := c.Uint(maxTradingPathLengthFlag)
	if maxPaths > pdex_v3.MaxPaths {
		return fmt.Errorf("maximum trading path length allowed %v, got %v", pdex_v3.MaxPaths, maxPaths)
	}

	allPoolPairs, err := cfg.incClient.GetAllPdexPoolPairs(0)
	if err != nil {
		return err
	}
	_, tradingPath, maxReceived := pdex_v3.FindGoodTradePath(maxPaths, allPoolPairs, tokenIdToSell, tokenIdToBuy, sellingAmount)
	if len(tradingPath) == 0 {
		return fmt.Errorf("no trading path is found for the pair %v-%v with maxPaths = %v", tokenIdToSell, tokenIdToBuy, maxPaths)
	}

	fmt.Printf("MaxReceived: %v\n", maxReceived)
	fmt.Printf("TradingPath: %v\n", tradingPath)

	return nil
}

// pDEXMintNFT creates and sends a transaction that mints a new C-NFT for a given user.
func pDEXMintNFT(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	encodedTx, txHash, err := cfg.incClient.CreatePdexv3MintNFT(privateKey)
	if err != nil {
		return err
	}
	err = cfg.incClient.SendRawTx(encodedTx)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)
	return nil
}

// pDEXCheckMintNFT gets the status of a (c)NFT minting transaction.
func pDEXCheckMintNFT(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	txHash := c.String(txHashFlag)
	status, nftID, err := cfg.incClient.CheckNFTMintingStatus(txHash)
	if err != nil {
		return err
	}

	if !status {
		fmt.Printf("Minting FAILED\n")
	} else {
		fmt.Printf("Minting SUCCEEDED with new nftID: %v\n", nftID)
	}

	return nil
}

// pDEXContribute contributes a token to the pDEX.
func pDEXContribute(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	nftID := c.String(nftIDFlag)

	pairHash := c.String(pairHashFlag)
	if pairHash == "" {
		return fmt.Errorf("%v is invalid", pairHashFlag)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return fmt.Errorf("%v cannot be zero", amountFlag)
	}

	amplifier := c.Uint64(amplifierFlag)
	if amplifier == 0 {
		return fmt.Errorf("%v cannot be zero", amplifierFlag)
	}

	pairID := c.String(pairIDFlag)

	tokenId := c.String(tokenIDFlag)
	if !isValidTokenID(tokenId) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3ContributeTransaction(
		privateKey,
		pairID,
		pairHash,
		tokenId,
		nftID,
		amount,
		amplifier,
		)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)
	return nil
}

// pDEXWithdraw withdraws a pair of tokens from the pDEX.
func pDEXWithdraw(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	pairID := c.String(pairIDFlag)
	nftID := c.String(nftIDFlag)

	tokenID1 := c.String(tokenID1Flag)
	if !isValidTokenID(tokenID1) {
		return fmt.Errorf("%v is invalid", tokenID1Flag)
	}
	tokenID2 := c.String(tokenID2Flag)
	if !isValidTokenID(tokenID2) {
		return fmt.Errorf("%v is invalid", tokenID2Flag)
	}

	shareAmount := c.Uint64(amountFlag)
	myShare, err := cfg.incClient.GetPoolShareAmount(pairID, nftID)
	if err != nil {
		return err
	}
	if shareAmount == 0 {
		shareAmount = myShare
	}
	if shareAmount > myShare {
		return fmt.Errorf("maximum share allowed to withdraw: %v", myShare)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3WithdrawLiquidityTransaction(
		privateKey,
		pairID,
		tokenID1,
		tokenID2,
		nftID,
		shareAmount,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)
	return nil
}

// pDEXGetShare returns the share amount of a pDEX nftID with-in a given poolID.
func pDEXGetShare(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	pairID := c.String(pairIDFlag)
	nftID := c.String(nftIDFlag)

	share, err := cfg.incClient.GetPoolShareAmount(pairID, nftID)
	if err != nil {
		return err
	}

	fmt.Printf("Share: %v\n", share)
	return nil
}
