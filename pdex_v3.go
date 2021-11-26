package main

import (
	"fmt"
	"github.com/incognitochain/incognito-cli/pdex_v3"
	"github.com/urfave/cli/v2"
	"strings"
)

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
	if tradingFee == 0 {
		return fmt.Errorf("%v cannot be zero", tradingFeeFlag)
	}

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
	if len(tradingPath) > int(maxPaths) {
		return fmt.Errorf("maximum trading path length %v, got %v", maxPaths, len(tradingPath))
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

// pDEXAddOrder places an order to the pDEX.
func pDEXAddOrder(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return fmt.Errorf("%v is invalid", pairHashFlag)
	}
	tokenIDs := strings.Split(pairID, "-")[:2]

	nftID := c.String(nftIDFlag)
	myNFTs, err := cfg.incClient.GetMyNFTs(privateKey)
	if err != nil {
		return err
	}
	nftExist := false
	for _, nft := range myNFTs {
		if nft == nftID {
			nftExist = true
			break
		}
	}
	if !nftExist {
		return fmt.Errorf("nftID %v does not belong to the private key %v", nftID, privateKey)
	}

	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return fmt.Errorf("%v is invalid", tokenIDToSellFlag)
	}
	if tokenIdToSell != tokenIDs[0] && tokenIdToSell != tokenIDs[1] {
		return fmt.Errorf("tokenToSell %v not belong to pool pair %v", tokenIdToSell, pairID)
	}
	tokenIdToBuy := tokenIDs[1]
	if tokenIdToSell == tokenIDs[1] {
		tokenIdToBuy = tokenIDs[0]
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return fmt.Errorf("%v cannot be zero", sellingAmountFlag)
	}

	minAcceptableAmount := c.Uint64(minAcceptableAmountFlag)
	if minAcceptableAmount == 0 {
		return fmt.Errorf("%v cannot be zero", minAcceptableAmount)
	}
	txHash, err := cfg.incClient.CreateAndSendPdexv3AddOrderTransaction(
		privateKey,
		pairID,
		tokenIdToSell,
		tokenIdToBuy,
		nftID,
		sellingAmount,
		minAcceptableAmount,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)

	return nil
}

// pDEXWithdrawOrder withdraws an order from the pDEX.
func pDEXWithdrawOrder(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return fmt.Errorf("%v is invalid", pairHashFlag)
	}
	tmpTokenIDs := strings.Split(pairID, "-")[:2]
	nftID := c.String(nftIDFlag)
	orderID := c.String(orderIDFlag)

	tokenId1 := c.String(tokenID1Flag)
	if !isValidTokenID(tokenId1) && tokenId1 != tmpTokenIDs[0] && tokenId1 != tmpTokenIDs[1] {
		return fmt.Errorf("%v is invalid", tokenID1Flag)
	}

	tokenId2 := c.String(tokenID2Flag)
	if tokenId2 != "" && !isValidTokenID(tokenId2) && tokenId2 != tmpTokenIDs[0] && tokenId2 != tmpTokenIDs[1] {
		return fmt.Errorf("%v is invalid", tokenID2Flag)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return fmt.Errorf("%v cannot be zero", amountFlag)
	}

	tokenIDs := []string{tokenId1}
	if tokenId2 != "" {
		tokenIDs = append(tokenIDs, tokenId2)
	}
	txHash, err := cfg.incClient.CreateAndSendPdexv3WithdrawOrderTransaction(
		privateKey,
		pairID,
		orderID,
		nftID,
		amount,
		tokenIDs...,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)

	return nil
}

// pDEXStake creates a pDEX staking transaction.
func pDEXStake(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	nftID := c.String(nftIDFlag)

	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return fmt.Errorf("%v cannot be zero", amountFlag)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3StakingTransaction(
		privateKey,
		tokenID,
		nftID,
		amount,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)

	return nil
}

// pDEXUnStake creates a pDEX un-staking transaction.
func pDEXUnStake(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	nftID := c.String(nftIDFlag)

	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return fmt.Errorf("%v cannot be zero", amountFlag)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3UnstakingTransaction(
		privateKey,
		tokenID,
		nftID,
		amount,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)

	return nil
}

// CheckDEXStakingReward returns the estimated pDEX staking rewards.
func CheckDEXStakingReward(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	nftID := c.String(nftIDFlag)
	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	res, err := cfg.incClient.GetEstimatedDEXStakingReward(0, tokenID, nftID)
	if err != nil {
		return err
	}
	return jsonPrint(res)
}

// pDEXWithdrawStakingReward creates a transaction withdrawing the staking rewards from the pDEX.
func pDEXWithdrawStakingReward(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	nftID := c.String(nftIDFlag)

	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3WithdrawStakeRewardTransaction(
		privateKey,
		tokenID,
		nftID,
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

// pDEXWithdrawLPFee creates a transaction withdrawing the LP fees for an nftID from the pDEX.
func pDEXWithdrawLPFee(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	nftID := c.String(nftIDFlag)

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return fmt.Errorf("%v is invalid", pairIDFlag)
	}

	lpValue, err := cfg.incClient.GetEstimatedLPValue(0, pairID, nftID)
	if err != nil {
		return err
	}
	if len(lpValue.TradingFee) == 0 {
		return fmt.Errorf("not enough LP fee to withdraw")
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3WithdrawLPFeeTransaction(
		privateKey,
		pairID,
		nftID,
	)
	if err != nil {
		return err
	}

	fmt.Printf("TxHash: %v\n", txHash)

	return nil
}

// pDEXGetEstimatedLPValue returns the estimated LP values of an LP in a given pool.
func pDEXGetEstimatedLPValue(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	poolPairID := c.String(pairIDFlag)
	if !isValidDEXPairID(poolPairID) {
		return fmt.Errorf("%v is invalid", pairIDFlag)
	}
	nftID := c.String(nftIDFlag)

	res, err := cfg.incClient.GetEstimatedLPValue(0, poolPairID, nftID)
	if err != nil {
		return err
	}
	err = jsonPrint(res)

	return err
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
		for path := range pairs {
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

// pDEXGetAllNFTs returns the list of NFTs for a given private key.
func pDEXGetAllNFTs(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	allNFTs, err := cfg.incClient.GetMyNFTs(privateKey)
	if err != nil {
		return err
	}
	err = jsonPrint(allNFTs)

	return err
}
