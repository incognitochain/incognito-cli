package main

import (
	"fmt"
	"github.com/incognitochain/incognito-cli/pdex_v3"
	"github.com/urfave/cli/v2"
	"strings"
)

// pDEXTrade creates and sends a trade to the pDEX.
func pDEXTrade(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return newAppError(InvalidSellTokenIDError)
	}

	tokenIdToBuy := c.String(tokenIDToBuyFlag)
	if !isValidTokenID(tokenIdToBuy) {
		return newAppError(InvalidBuyTokenIDError)
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return newAppError(InvalidSellAmountError)
	}

	minAcceptableAmount := c.Uint64(minAcceptableAmountFlag)
	tradingFee := c.Uint64(tradingFeeFlag)
	if tradingFee == 0 {
		return newAppError(InvalidTradingFeeError)
	}

	maxPaths := c.Uint(maxTradingPathLengthFlag)
	if maxPaths > pdex_v3.MaxPaths {
		return newAppError(InvalidMaxTradingPathError, fmt.Errorf("maximum trading path length allowed %v, got %v", pdex_v3.MaxPaths, maxPaths))
	}

	allPoolPairs, err := cfg.incClient.GetAllPdexPoolPairs(0)
	if err != nil {
		return newAppError(GetAllDexPoolPairsError, err)
	}
	tmpTradingPath := c.String(tradingPathFlag)
	tradingPath := make([]string, 0)
	if tmpTradingPath != "" {
		tradingPath = strings.Split(tmpTradingPath, ",")
		for _, poolID := range tradingPath {
			if _, ok := allPoolPairs[poolID]; !ok {
				return newAppError(UnexpectedError, fmt.Errorf("poolID %v not existed", poolID))
			}
		}
	} else {
		_, tradingPath, _ = pdex_v3.FindGoodTradePath(maxPaths, allPoolPairs, tokenIdToSell, tokenIdToBuy, sellingAmount)
	}
	if len(tradingPath) == 0 {
		return newAppError(InvalidTradingPathError, fmt.Errorf("no trading path is found for the pair %v-%v with maxPaths = %v", tokenIdToSell, tokenIdToBuy, maxPaths))
	}
	if len(tradingPath) > int(maxPaths) {
		return newAppError(InvalidTradingPathError, fmt.Errorf("maximum trading path length %v, got %v", maxPaths, len(tradingPath)))
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
		return newAppError(CreateDexTradeTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXMintNFT creates and sends a transaction that mints a new C-NFT for a given user.
func pDEXMintNFT(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	encodedTx, txHash, err := cfg.incClient.CreatePdexv3MintNFT(privateKey)
	if err != nil {
		return newAppError(CreateMintNFTTransactionError, err)
	}
	err = cfg.incClient.SendRawTx(encodedTx)
	if err != nil {
		return newAppError(SendRawTxError)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXContribute contributes a token to the pDEX.
func pDEXContribute(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	nftID := c.String(nftIDFlag)

	pairHash := c.String(pairHashFlag)
	if pairHash == "" {
		return newAppError(InvalidTokenIDError)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return newAppError(InvalidAmountError)
	}

	amplifier := c.Uint64(amplifierFlag)
	if amplifier == 0 {
		return newAppError(InvalidAmplifierError)
	}

	pairID := c.String(pairIDFlag)

	tokenId := c.String(tokenIDFlag)
	if !isValidTokenID(tokenId) {
		return newAppError(InvalidTokenIDError)
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
		return newAppError(CreateDexContributionTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXWithdraw withdraws a pair of tokens from the pDEX.
func pDEXWithdraw(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return fmt.Errorf("%v is invalid", pairHashFlag)
	}
	tmpTokenIDs := strings.Split(pairID, "-")[:2]
	nftID := c.String(nftIDFlag)

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
		tmpTokenIDs[0],
		tmpTokenIDs[1],
		nftID,
		shareAmount,
	)
	if err != nil {
		return err
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXAddOrder places an order to the pDEX.
func pDEXAddOrder(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return newAppError(InvalidPoolPairIDError)
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
		return newAppError(InvalidNFTError, fmt.Errorf("nftID %v does not belong to the private key %v", nftID, privateKey))
	}

	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return newAppError(InvalidSellTokenIDError)
	}
	if tokenIdToSell != tokenIDs[0] && tokenIdToSell != tokenIDs[1] {
		return newAppError(InvalidSellTokenIDError, fmt.Errorf("tokenToSell %v not belong to pool pair %v", tokenIdToSell, pairID))
	}
	tokenIdToBuy := tokenIDs[1]
	if tokenIdToSell == tokenIDs[1] {
		tokenIdToBuy = tokenIDs[0]
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return newAppError(InvalidSellAmountError)
	}

	minAcceptableAmount := c.Uint64(minAcceptableAmountFlag)
	if minAcceptableAmount == 0 {
		return newAppError(InvalidMinAcceptableAmountError)
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
		return newAppError(CreateAddOrderTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXWithdrawOrder withdraws an order from the pDEX.
func pDEXWithdrawOrder(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return newAppError(InvalidPoolPairIDError)
	}
	tmpTokenIDs := strings.Split(pairID, "-")[:2]
	nftID := c.String(nftIDFlag)
	orderID := c.String(orderIDFlag)

	tokenId1 := c.String(tokenID1Flag)
	if !isValidTokenID(tokenId1) && tokenId1 != tmpTokenIDs[0] && tokenId1 != tmpTokenIDs[1] {
		return newAppError(InvalidTokenIDError, fmt.Errorf("%v is invalid", tokenID1Flag))
	}

	tokenId2 := c.String(tokenID2Flag)
	if tokenId2 != "" && !isValidTokenID(tokenId2) && tokenId2 != tmpTokenIDs[0] && tokenId2 != tmpTokenIDs[1] {
		return newAppError(InvalidTokenIDError, fmt.Errorf("%v is invalid", tokenID2Flag))
	}

	amount := c.Uint64(amountFlag)

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
		return newAppError(CreateWithdrawOrderTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXStake creates a pDEX staking transaction.
func pDEXStake(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	nftID := c.String(nftIDFlag)

	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return newAppError(InvalidTokenIDError)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return newAppError(InvalidAmountError)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3StakingTransaction(
		privateKey,
		tokenID,
		nftID,
		amount,
	)
	if err != nil {
		return newAppError(CreateDexStakingTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXUnStake creates a pDEX un-staking transaction.
func pDEXUnStake(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	nftID := c.String(nftIDFlag)

	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return newAppError(InvalidTokenIDError)
	}

	amount := c.Uint64(amountFlag)
	if amount == 0 {
		return newAppError(InvalidAmountError)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3UnstakingTransaction(
		privateKey,
		tokenID,
		nftID,
		amount,
	)
	if err != nil {
		return newAppError(CreateDexUnStakingTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// CheckDEXStakingReward returns the estimated pDEX staking rewards.
func CheckDEXStakingReward(c *cli.Context) error {
	nftID := c.String(nftIDFlag)
	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return newAppError(InvalidTokenIDError)
	}

	res, err := cfg.incClient.GetEstimatedDEXStakingReward(0, tokenID, nftID)
	if err != nil {
		return newAppError(EstimateDEXStakingRewardError, err)
	}
	return jsonPrint(res)
}

// pDEXWithdrawStakingReward creates a transaction withdrawing the staking rewards from the pDEX.
func pDEXWithdrawStakingReward(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	nftID := c.String(nftIDFlag)

	tokenID := c.String(tokenIDFlag)
	if !isValidTokenID(tokenID) {
		return newAppError(InvalidTokenIDError)
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3WithdrawStakeRewardTransaction(
		privateKey,
		tokenID,
		nftID,
	)
	if err != nil {
		return newAppError(CreateDexStakingRewardWithdrawalTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXGetShare returns the share amount of a pDEX nftID with-in a given poolID.
func pDEXGetShare(c *cli.Context) error {
	pairID := c.String(pairIDFlag)
	nftID := c.String(nftIDFlag)

	share, err := cfg.incClient.GetPoolShareAmount(pairID, nftID)
	if err != nil {
		return newAppError(GetPoolShareError, err)
	}

	return jsonPrintWithKey("Share", share)
}

// pDEXWithdrawLPFee creates a transaction withdrawing the LP fees for an nftID from the pDEX.
func pDEXWithdrawLPFee(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	nftID := c.String(nftIDFlag)

	pairID := c.String(pairIDFlag)
	if !isValidDEXPairID(pairID) {
		return newAppError(InvalidPoolPairIDError)
	}

	lpValue, err := cfg.incClient.GetEstimatedLPValue(0, pairID, nftID)
	if err != nil {
		return newAppError(GetEstimatedLPValueError, err)
	}
	if len(lpValue.PoolReward) == 0 {
		return newAppError(CreateLPFeeWithdrawalTransactionError, fmt.Errorf("not enough reward to withdraw"))
	}

	txHash, err := cfg.incClient.CreateAndSendPdexv3WithdrawLPFeeTransaction(
		privateKey,
		pairID,
		nftID,
	)
	if err != nil {
		return newAppError(CreateLPFeeWithdrawalTransactionError, err)
	}

	return jsonPrintWithKey("TxHash", txHash)
}

// pDEXGetEstimatedLPValue returns the estimated LP values of an LP in a given pool.
func pDEXGetEstimatedLPValue(c *cli.Context) error {
	poolPairID := c.String(pairIDFlag)
	if !isValidDEXPairID(poolPairID) {
		return newAppError(InvalidPoolPairIDError)
	}
	nftID := c.String(nftIDFlag)

	res, err := cfg.incClient.GetEstimatedLPValue(0, poolPairID, nftID)
	if err != nil {
		return newAppError(GetEstimatedLPValueError, err)
	}

	return jsonPrint(res)
}

// pDEXFindPath finds a proper trading path.
func pDEXFindPath(c *cli.Context) error {
	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return newAppError(InvalidSellTokenIDError)
	}

	tokenIdToBuy := c.String(tokenIDToBuyFlag)
	if !isValidTokenID(tokenIdToBuy) {
		return newAppError(InvalidBuyTokenIDError)
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return newAppError(InvalidSellAmountError)
	}

	maxPaths := c.Uint(maxTradingPathLengthFlag)
	if maxPaths > pdex_v3.MaxPaths {
		return newAppError(InvalidMaxTradingPathError, fmt.Errorf("maximum trading path length allowed %v, got %v", pdex_v3.MaxPaths, maxPaths))
	}

	allPoolPairs, err := cfg.incClient.GetAllPdexPoolPairs(0)
	if err != nil {
		return newAppError(GetAllDexPoolPairsError, err)
	}
	_, tradingPath, maxReceived := pdex_v3.FindGoodTradePath(maxPaths, allPoolPairs, tokenIdToSell, tokenIdToBuy, sellingAmount)
	if len(tradingPath) == 0 {
		return newAppError(FindTradingPathError,
			fmt.Errorf("no trading path is found for the pair %v-%v with maxPaths = %v", tokenIdToSell, tokenIdToBuy, maxPaths))
	}

	return jsonPrint(map[string]interface{}{"MaxReceived": maxReceived, "TradingPath": tradingPath})
}

// pDEXCheckPrice checks the price of two tokenIds.
func pDEXCheckPrice(c *cli.Context) error {
	var err error
	tokenIdToSell := c.String(tokenIDToSellFlag)
	if !isValidTokenID(tokenIdToSell) {
		return newAppError(InvalidSellTokenIDError)
	}

	tokenIdToBuy := c.String(tokenIDToBuyFlag)
	if !isValidTokenID(tokenIdToBuy) {
		return newAppError(InvalidBuyTokenIDError)
	}

	sellingAmount := c.Uint64(sellingAmountFlag)
	if sellingAmount == 0 {
		return newAppError(InvalidSellAmountError)
	}

	pairID := c.String(pairIDFlag)
	bestExpectedReceive := uint64(0)
	if pairID != "" {
		pairs, err := cfg.incClient.GetPdexPoolPair(0, tokenIdToSell, tokenIdToBuy)
		if err != nil {
			return newAppError(GetDexPoolPairError, err)
		}
		for path := range pairs {
			expectedPrice, err := cfg.incClient.CheckPrice(path, tokenIdToSell, sellingAmount)
			if err != nil {
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
			return newAppError(DexPriceCheckingError, err)
		}
	}

	if bestExpectedReceive == 0 {
		return newAppError(DexPriceCheckingError, fmt.Errorf("cannot find a proper path"))
	}

	return jsonPrint(map[string]interface{}{"BestPairID": pairID, "BestReceived": bestExpectedReceive})
}

// pDEXGetAllNFTs returns the list of NFTs for a given private key.
func pDEXGetAllNFTs(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	allNFTs, err := cfg.incClient.GetMyNFTs(privateKey)
	if err != nil {
		return newAppError(GetAllDexNFTsError, err)
	}

	return jsonPrint(allNFTs)
}

// pDEXGetOrderByID returns the detail of an order given its id.
func pDEXGetOrderByID(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	orderID := c.String(orderIDFlag)
	if orderID == "" {
		return newAppError(InvalidOrderIDError)
	}

	order, err := cfg.incClient.GetOrderByID(0, orderID)
	if err != nil {
		return newAppError(GetOrderByIDError, err)
	}

	return jsonPrint(order)
}
