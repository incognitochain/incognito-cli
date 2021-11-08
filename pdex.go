package main

//// pDEXCheckPrice checks the price of two tokenIds.
//func pDEXCheckPrice(c *cli.Context) error {
//	err := initNetWork()
//	if err != nil {
//		return err
//	}
//
//	tokenIdToSell := c.String(tokenIDToSellFlag)
//	if !isValidTokenID(tokenIdToSell) {
//		return fmt.Errorf("%v is invalid", tokenIDToSellFlag)
//	}
//
//	tokenIdToBuy := c.String(tokenIDToBuyFlag)
//	if !isValidTokenID(tokenIdToBuy) {
//		return fmt.Errorf("%v is invalid", tokenIDToBuyFlag)
//	}
//
//	sellingAmount := c.Uint64(sellingAmountFlag)
//	if sellingAmount == 0 {
//		return fmt.Errorf("%v cannot be zero", sellingAmountFlag)
//	}
//
//	expectedPrice, err := cfg.incClient.CheckXPrice(tokenIdToSell, tokenIdToBuy, sellingAmount)
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("%v\n", expectedPrice)
//
//	return nil
//}

//// pDEXTrade creates and sends a trade to the pDEX.
//func pDEXTrade(c *cli.Context) error {
//	err := initNetWork()
//	if err != nil {
//		return err
//	}
//
//	privateKey := c.String(privateKeyFlag)
//	if !isValidPrivateKey(privateKey) {
//		return fmt.Errorf("%v is invalid", privateKeyFlag)
//	}
//
//	tokenIdToSell := c.String(tokenIDToSellFlag)
//	if !isValidTokenID(tokenIdToSell) {
//		return fmt.Errorf("%v is invalid", tokenIDToSellFlag)
//	}
//
//	tokenIdToBuy := c.String(tokenIDToBuyFlag)
//	if !isValidTokenID(tokenIdToBuy) {
//		return fmt.Errorf("%v is invalid", tokenIDToBuyFlag)
//	}
//
//	sellingAmount := c.Uint64(sellingAmountFlag)
//	if sellingAmount == 0 {
//		return fmt.Errorf("%v cannot be zero", sellingAmountFlag)
//	}
//
//	minAcceptableAmount := c.Uint64(minAcceptableAmountFlag)
//	tradingFee := c.Uint64(tradingFeeFlag)
//
//	txHash, err := cfg.incClient.CreateAndSendPDETradeTransaction(
//		privateKey,
//		tokenIdToSell,
//		tokenIdToBuy,
//		sellingAmount,
//		minAcceptableAmount,
//		tradingFee,
//	)
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("TxHash: %v\n", txHash)
//
//	return nil
//}

//// pDEXContribute contributes a token to the pDEX.
//func pDEXContribute(c *cli.Context) error {
//	err := initNetWork()
//	if err != nil {
//		return err
//	}
//
//	privateKey := c.String(privateKeyFlag)
//	if !isValidPrivateKey(privateKey) {
//		return fmt.Errorf("%v is invalid", privateKeyFlag)
//	}
//
//	pairID := c.String(pairIDFlag)
//
//	tokenId := c.String(tokenIDFlag)
//	if !isValidTokenID(tokenId) {
//		return fmt.Errorf("%v is invalid", tokenIDFlag)
//	}
//
//	amount := c.Uint64(amountFlag)
//	if amount == 0 {
//		return fmt.Errorf("%v cannot be zero", amountFlag)
//	}
//
//	version := c.Int(versionFlag)
//	if !isSupportedVersion(int8(version)) {
//		return fmt.Errorf("%v is not supported", versionFlag)
//	}
//
//	txHash, err := cfg.incClient.CreateAndSendPDEContributeTransaction(
//		privateKey,
//		pairID,
//		tokenId,
//		amount,
//		int8(version),
//	)
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("TxHash: %v\n", txHash)
//
//	return nil
//}
//
//// pDEXWithdraw withdraws assets from the pDEX.
//func pDEXWithdraw(c *cli.Context) error {
//	err := initNetWork()
//	if err != nil {
//		return err
//	}
//
//	privateKey := c.String(privateKeyFlag)
//	if !isValidPrivateKey(privateKey) {
//		return fmt.Errorf("%v is invalid", privateKeyFlag)
//	}
//
//	tokenId1 := c.String(tokenID1Flag)
//	if !isValidTokenID(tokenId1) {
//		return fmt.Errorf("%v is invalid", tokenID1Flag)
//	}
//
//	tokenId2 := c.String(tokenID2Flag)
//	if !isValidTokenID(tokenId2) {
//		return fmt.Errorf("%v is invalid", tokenID2Flag)
//	}
//
//	amount := c.Uint64(amountFlag)
//	if amount == 0 {
//		return fmt.Errorf("%v cannot be zero", amountFlag)
//	}
//
//	version := c.Int(versionFlag)
//	if !isSupportedVersion(int8(version)) {
//		return fmt.Errorf("%v is not supported", versionFlag)
//	}
//
//	txHash, err := cfg.incClient.CreateAndSendPDEWithdrawalTransaction(
//		privateKey,
//		tokenId1,
//		tokenId2,
//		amount,
//		int8(version),
//	)
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("TxHash: %v\n", txHash)
//
//	return nil
//}
//
//// pDEXGetShare returns the share amount of a user w.r.t a pDEX pair.
//func pDEXGetShare(c *cli.Context) error {
//	err := initNetWork()
//	if err != nil {
//		return err
//	}
//
//	address := c.String(addressFlag)
//	if !isValidAddress(address) {
//		return fmt.Errorf("%v is invalid", addressFlag)
//	}
//
//	tokenId1 := c.String(tokenID1Flag)
//	if !isValidTokenID(tokenId1) {
//		return fmt.Errorf("%v is invalid", tokenID1Flag)
//	}
//
//	tokenId2 := c.String(tokenID2Flag)
//	if !isValidTokenID(tokenId2) {
//		return fmt.Errorf("%v is invalid", tokenID2Flag)
//	}
//
//	shareAmount, err := cfg.incClient.GetShareAmount(0, tokenId1, tokenId2, address)
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("shareAmount: %v\n", shareAmount)
//
//	return nil
//}
//
//// pDEXTradeStatus returns the status of a trade.
//func pDEXTradeStatus(c *cli.Context) error {
//	err := initNetWork()
//	if err != nil {
//		return err
//	}
//
//	txHash := c.String(txHashFlag)
//	tradeStatus, err := cfg.incClient.CheckTradeStatus(txHash)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("tradeStatus: %v\n", tradeStatus)
//
//	return nil
//}
