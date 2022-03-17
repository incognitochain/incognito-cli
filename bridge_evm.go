package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	iCommon "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/crypto"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

var shieldMessage = "This function helps shield an EVM (ETH/BNB/ERC20/BEP20, etc.) token into the Incognito network. " +
	"It will ask for users' EVM PRIVATE KEY to proceed. " +
	"The shielding process consists of the following operations.\n" +
	"\t 1. Deposit the EVM asset into the corresponding smart contract.\n" +
	"\t\t 1.1. In case the asset is an ERC20/BEP20 token, an approval transaction is performed (if needed) the before the " +
	"actual deposit. For this operation, a prompt will be displayed to ask for user's approval.\n" +
	"\t 2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. " +
	"This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.\n\n" +
	"Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. " +
	"We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.\n" +
	"DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS."

var unShieldMessage = "This function helps withdraw an EVM (ETH/BNB/ERC20/BEP20, etc.) token out of the Incognito network. " +
	"The un-shielding process consists the following operations.\n" +
	"\t 1. Users burn the token inside the Incognito chain.\n" +
	"\t 2. After the burning is success, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from " +
	"the Incognito chain.\n" +
	"\t 3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the " +
	"corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, " +
	"it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).\n\n" +
	"Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the " +
	"Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. " +
	"Otherwise, use the un-shielding function from the Incognito app. " +
	"We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.\n" +
	"DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS."

// shield deposits an EVM token (ETH/BNB/ERC20/BEP20) into the Incognito chain.
func shield(c *cli.Context) error {
	fmt.Println(shieldMessage)
	yesNoPrompt("Do you want to continue?")
	fmt.Println()

	log.Println("[STEP 0] PREPARE DATA")
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	incAddress := c.String(addressFlag)
	if incAddress == "" {
		incAddress = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}

	shieldAmount := c.Float64(shieldAmountFlag)
	tokenAddressStr := c.String(tokenAddressFlag)
	if !isValidEVMAddress(tokenAddressStr) {
		return fmt.Errorf("%v is invalid", tokenAddressFlag)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	incTokenID, err := getIncTokenIDFromEVMTokenID(tokenAddress.String(), evmNetworkID)
	if err != nil {
		if strings.Contains(err.Error(), "incTokenID not found") {
			log.Printf("IncTokenID not found for %v, perhaps it doesn't exist in the Incognito network.\n", tokenAddress.String())
			incTokenID = fmt.Sprintf("%x", iCommon.RandBytes(32))
			yesNoPrompt(fmt.Sprintf("Newly generated incTokenID: %v. Do you want to continue with this token?", incTokenID))
		} else {
			return err
		}
	}

	var tokenName, tokenSymbol string
	if tokenAddress.String() == nativeToken {
		tokenName = "Ethereum"
		tokenSymbol = "ETH"
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			tokenName = "Binance"
			tokenSymbol = "BNB"
		case rpc.PLGNetworkID:
			tokenName = "Matic"
			tokenSymbol = "MATIC"
		}
	} else {
		tokenInfo, err := getEVMTokenInfo(tokenAddress.String(), evmNetworkID)
		if err != nil {
			return err
		}
		tokenName = tokenInfo.name
		tokenSymbol = tokenInfo.symbol
		if tokenInfo.network != evmNetwork {
			return fmt.Errorf("expect token to be on `%v` network, got `%v`", evmNetwork, tokenInfo.network)
		}
	}
	log.Printf("Network: %v, TokenName: %v, TokenSymbol: %v, TokenAddress: %v, ShieldAmount: %v",
		evmNetwork, tokenName, tokenSymbol, tokenAddress.String(), shieldAmount)
	yesNoPrompt("Do you want to continue?")
	log.Printf("[STEP 0] FINISHED!\n\n")

	log.Println("[STEP 1] CHECK INCOGNITO BALANCE")
	prvBalance, err := checkSufficientIncBalance(privateKey, iCommon.PRVIDStr, incclient.DefaultPRVFee)
	if err != nil {
		return err
	}
	log.Printf("Current PRV balance: %v\n", prvBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Printf("[STEP 2] IMPORT %v ACCOUNT\n", evmNetwork)

	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return err
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return err
	}

	for {
		evmTokenBalance, err := acc.checkSufficientBalance(tokenAddress, shieldAmount, evmNetworkID)
		err = checkAndChangeRPCEndPoint(evmNetworkID, err)
		if err != nil {
			return err
		}
		if tokenAddress.String() == nativeToken {
			log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), tokenName, evmTokenBalance)
		} else {
			nativeTokenName := "ETH"
			switch evmNetworkID {
			case rpc.BSCNetworkID:
				nativeTokenName = "BNB"
			case rpc.PLGNetworkID:
				nativeTokenName = "MATIC"
			}
			_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
			if err != nil {
				return err
			}
			nativeBalance, _ := tmpNativeBalance.Float64()
			log.Printf("Your %v address: %v, %v: %v, %v: %v\n",
				evmNetwork,
				acc.address.String(), nativeTokenName, nativeBalance, tokenSymbol, evmTokenBalance)
		}
		break
	}

	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] DEPOSIT PUBLIC TOKEN TO SC")
	var evmHash *common.Hash
	if tokenAddress.String() == nativeToken {
		evmHash, err = acc.DepositNative(incAddress, shieldAmount, 0, 0, evmNetworkID)
	} else {
		evmHash, err = acc.DepositToken(incAddress, tokenAddressStr, shieldAmount, 0, 0, evmNetworkID)
	}
	if err != nil {
		return err
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] SHIELD TO INCOGNITO")
	incTxHash, err := Shield(privateKey, incTokenID, evmHash.String(), evmNetworkID)
	if err != nil {
		return err
	}
	log.Printf("[STEP 4] FINISHED!\n\n")

	log.Println("[STEP 5] CHECK SHIELD STATUS")
	for {
		status, err := cfg.incClient.CheckShieldStatus(incTxHash)
		if err != nil || status <= 1 {
			log.Printf("ShieldingStatus: %v\n", status)
			time.Sleep(40 * time.Second)
			continue
		} else if status == 2 {
			log.Println("Shielding SUCCEEDED!!")
			break
		} else {
			panic("Shielding FAILED!!")
		}
	}
	log.Printf("[STEP 5] FINISHED!\n\n")
	return nil
}

// shieldWithDepositKey deposits an EVM token (ETH/BNB/ERC20/BEP20) into the Incognito chain using deposit keys.
func shieldWithDepositKey(c *cli.Context) error {
	fmt.Println(shieldMessage)
	yesNoPrompt("Do you want to continue?")
	fmt.Println()

	log.Println("[STEP 0] PREPARE DATA")
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	shieldAmount := c.Float64(shieldAmountFlag)

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}
	tokenAddressStr := c.String(tokenAddressFlag)
	if !isValidEVMAddress(tokenAddressStr) {
		return fmt.Errorf("%v is invalid", tokenAddressFlag)
	}

	tokenAddress := common.HexToAddress(tokenAddressStr)
	incTokenID, err := getIncTokenIDFromEVMTokenID(tokenAddress.String(), evmNetworkID)
	if err != nil {
		if strings.Contains(err.Error(), "incTokenID not found") {
			log.Printf("IncTokenID not found for %v, perhaps it doesn't exist in the Incognito network.\n", tokenAddress.String())
			incTokenID = fmt.Sprintf("%x", iCommon.RandBytes(32))
			yesNoPrompt(fmt.Sprintf("Newly generated incTokenID: %v. Do you want to continue with this token?", incTokenID))
		} else {
			return err
		}
	}

	var depositPrivateKey []byte
	depositPrivateKeyStr := c.String(depositPrivateKeyFlag)
	if depositPrivateKeyStr == "" {
		depositKeyIndex := c.Uint64(depositIndexFlag)
		depositKey, err := cfg.incClient.GenerateDepositKeyFromPrivateKey(privateKey, incTokenID, depositKeyIndex)
		if err != nil {
			return err
		}
		depositPrivateKey = depositKey.PrivateKey
		depositPrivateKeyStr = base58.Base58Check{}.NewEncode(depositKey.PrivateKey, 0)
	} else {
		depositPrivateKey, _, err = base58.Base58Check{}.Decode(depositPrivateKeyStr)
		if err != nil {
			return fmt.Errorf("invalid deposit private key")
		}
	}
	depositPubKey := new(crypto.Point).ScalarMultBase(new(crypto.Scalar).FromBytesS(depositPrivateKey)).ToBytesS()
	depositPubKeyStr := base58.Base58Check{}.NewEncode(depositPubKey, 0)

	signature := c.String(signatureFlag)
	receiver := c.String(receiverFlag)
	if receiver != "" && !isValidOTAReceiver(receiver) {
		return fmt.Errorf("%v is invalid", receiverFlag)
	}

	var tokenName, tokenSymbol string
	if tokenAddress.String() == nativeToken {
		tokenName = "Ethereum"
		tokenSymbol = "ETH"
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			tokenName = "Binance"
			tokenSymbol = "BNB"
		case rpc.PLGNetworkID:
			tokenName = "Matic"
			tokenSymbol = "MATIC"
		}
	} else {
		tokenInfo, err := getEVMTokenInfo(tokenAddress.String(), evmNetworkID)
		if err != nil {
			return err
		}
		tokenName = tokenInfo.name
		tokenSymbol = tokenInfo.symbol
		if tokenInfo.network != evmNetwork {
			return fmt.Errorf("expect token to be on `%v` network, got `%v`", evmNetwork, tokenInfo.network)
		}
	}
	log.Printf("Network: %v, TokenName: %v, TokenSymbol: %v, TokenAddress: %v, ShieldAmount: %v",
		evmNetwork, tokenName, tokenSymbol, tokenAddress.String(), shieldAmount)
	yesNoPrompt("Do you want to continue?")
	log.Printf("[STEP 0] FINISHED!\n\n")

	log.Println("[STEP 1] CHECK INCOGNITO BALANCE")
	prvBalance, err := checkSufficientIncBalance(privateKey, iCommon.PRVIDStr, incclient.DefaultPRVFee)
	if err != nil {
		return err
	}
	log.Printf("Current PRV balance: %v\n", prvBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Printf("[STEP 2] IMPORT %v ACCOUNT\n", evmNetwork)
	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return err
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return err
	}

	for {
		evmTokenBalance, err := acc.checkSufficientBalance(tokenAddress, shieldAmount, evmNetworkID)
		err = checkAndChangeRPCEndPoint(evmNetworkID, err)
		if err != nil {
			return err
		}
		if tokenAddress.String() == nativeToken {
			log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), tokenName, evmTokenBalance)
		} else {
			nativeTokenName := "ETH"
			switch evmNetworkID {
			case rpc.BSCNetworkID:
				nativeTokenName = "BNB"
			case rpc.PLGNetworkID:
				nativeTokenName = "MATIC"
			}
			_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
			if err != nil {
				return err
			}
			nativeBalance, _ := tmpNativeBalance.Float64()
			log.Printf("Your %v address: %v, %v: %v, %v: %v\n",
				evmNetwork,
				acc.address.String(), nativeTokenName, nativeBalance, tokenSymbol, evmTokenBalance)
		}
		break
	}
	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] DEPOSIT PUBLIC TOKEN TO SC")
	var evmHash *common.Hash
	if tokenAddress.String() == nativeToken {
		evmHash, err = acc.DepositNative(depositPubKeyStr, shieldAmount, 0, 0, evmNetworkID)
	} else {
		evmHash, err = acc.DepositToken(depositPubKeyStr, tokenAddressStr, shieldAmount, 0, 0, evmNetworkID)
	}
	if err != nil {
		return err
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] SHIELD TO INCOGNITO")
	dp := incclient.EVMDepositParams{
		MetadataType:      metadata.IssuingETHRequestMeta,
		TokenID:           incTokenID,
		DepositPrivateKey: depositPrivateKeyStr,
		Receiver:          receiver,
		Signature:         signature,
	}
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		dp.MetadataType = metadata.IssuingBSCRequestMeta
	case rpc.PLGNetworkID:
		dp.MetadataType = metadata.IssuingPLGRequestMeta
	}
	incTxHash, err := ShieldWithDepositKey(privateKey, evmHash.String(), evmNetworkID, dp)
	if err != nil {
		return err
	}
	log.Printf("[STEP 4] FINISHED!\n\n")

	log.Println("[STEP 5] CHECK SHIELD STATUS")
	for {
		status, err := cfg.incClient.CheckShieldStatus(incTxHash)
		if err != nil || status <= 1 {
			log.Printf("ShieldingStatus: %v\n", status)
			time.Sleep(40 * time.Second)
			continue
		} else if status == 2 {
			log.Println("Shielding SUCCEEDED!!")
			break
		} else {
			panic("Shielding FAILED!!")
		}
	}
	log.Printf("[STEP 5] FINISHED!\n\n")
	return nil
}

// retryShield retries to shield a token with an already-deposited evm TxHash.
func retryShield(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}

	tokenAddressStr := c.String(tokenAddressFlag)
	if !isValidEVMAddress(tokenAddressStr) {
		return fmt.Errorf("%v is invalid", tokenAddressFlag)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	incTokenID, err := getIncTokenIDFromEVMTokenID(tokenAddress.String(), evmNetworkID)
	if err != nil {
		if strings.Contains(err.Error(), "incTokenID not found") {
			log.Printf("IncTokenID not found for %v, perhaps it doesn't exist in the Incognito network.\n", tokenAddress.String())
			incTokenID = fmt.Sprintf("%x", iCommon.RandBytes(32))
			yesNoPrompt(fmt.Sprintf("Newly generated incTokenID: %v. Do you want to continue with this token?", incTokenID))
		} else {
			return err
		}
	}

	evmTxHashStr := c.String(externalTxIDFlag)
	evmHash := common.HexToHash(evmTxHashStr)

	log.Println("[STEP 1] SHIELD TO INCOGNITO")
	incTxHash, err := Shield(privateKey, incTokenID, evmHash.String(), evmNetworkID)
	if err != nil {
		return err
	}
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Println("[STEP 2] CHECK SHIELD STATUS")
	for {
		status, err := cfg.incClient.CheckShieldStatus(incTxHash)
		if err != nil || status <= 1 {
			log.Printf("ShieldingStatus: %v\n", status)
			time.Sleep(40 * time.Second)
			continue
		} else if status == 2 {
			log.Println("Shielding SUCCEEDED!!")
			break
		} else {
			panic("Shielding FAILED!!")
		}
	}
	log.Printf("[STEP 2] FINISHED!\n\n")
	return nil
}

// unShield withdraws an EVM token (ETH/BNB/ERC20/BEP20) from the Incognito chain.
func unShield(c *cli.Context) error {
	fmt.Println(unShieldMessage)
	yesNoPrompt("Do you want to continue?")
	fmt.Println()

	log.Println("[STEP 0] PREPARE DATA")
	// get the private key
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	// get the un-shield amount
	unShieldAmount := c.Uint64(amountFlag)
	if unShieldAmount == 0 {
		return fmt.Errorf("%v is invalid", amountFlag)
	}

	// get the Incognito tokenID, evmTokenID, name and symbol.
	incTokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(incTokenIDStr) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}
	evmTokenIDStr, evmNetworkID, err := getEVMTokenIDIncTokenID(incTokenIDStr)
	if err != nil {
		return err
	}
	evmTokenAddress := common.HexToAddress(evmTokenIDStr)
	evmNetwork := "ETH"
	nativeTokenName := "ETH"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		evmNetwork = "BSC"
		nativeTokenName = "BNB"
	case rpc.PLGNetworkID:
		evmNetwork = "PLG"
		nativeTokenName = "MATIC"
	}
	var tokenName, tokenSymbol string
	if evmTokenAddress.String() == nativeToken {
		tokenName = "Ethereum"
		tokenSymbol = "ETH"
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			tokenName = "Binance"
			tokenSymbol = "BNB"
		case rpc.PLGNetworkID:
			tokenName = "Matic"
			tokenSymbol = "MATIC"
		}
	} else {
		tokenInfo, err := getEVMTokenInfo(evmTokenAddress.String(), evmNetworkID)
		if err != nil {
			return err
		}
		tokenName = tokenInfo.name
		tokenSymbol = tokenInfo.symbol
	}
	log.Printf("Network: %v, TokenName: %v, TokenSymbol: %v, TokenAddress: %v, UnShieldAmount: %v",
		evmNetwork, tokenName, tokenSymbol, evmTokenAddress.String(), unShieldAmount)
	yesNoPrompt("Do you want to continue?")
	log.Printf("[STEP 0] FINISHED!\n\n")

	log.Println("[STEP 1] CHECK INCOGNITO BALANCE")
	prvBalance, err := checkSufficientIncBalance(privateKey, iCommon.PRVIDStr, incclient.DefaultPRVFee)
	if err != nil {
		return err
	}
	incTokenBalance, err := checkSufficientIncBalance(privateKey, incTokenIDStr, unShieldAmount)
	if err != nil {
		return err
	}
	log.Printf("Current PRVBalance: %v, TokenBalance: %v\n", prvBalance, incTokenBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Printf("[STEP 2] IMPORT %v ACCOUNT\n", evmNetwork)

	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return err
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return err
	}
	_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
	if err != nil {
		return err
	}
	nativeBalance, _ := tmpNativeBalance.Float64()
	log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), nativeTokenName, nativeBalance)
	evmAddress := acc.address
	var res string
	resInBytes, err := promptInput(
		fmt.Sprintf("Un-shield to the following address: %v. Continue? (y/n)", evmAddress.String()),
		&res)
	if err != nil {
		return err
	}
	res = string(resInBytes)
	if !strings.Contains(res, "y") && !strings.Contains(res, "Y") {
		resInBytes, err = promptInput(
			fmt.Sprintf("Enter the address you want to un-shield to"),
			&res)
		if err != nil {
			return err
		}
		res = string(resInBytes)
		if !isValidEVMAddress(res) {
			return fmt.Errorf("%v is not a valid EVM address", res)
		}
		evmAddress = common.HexToAddress(res)
	}
	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] BURN INCOGNITO TOKEN")
	incTxHash, err := cfg.incClient.CreateAndSendBurningRequestTransaction(privateKey, evmAddress.String(), incTokenIDStr, unShieldAmount, evmNetworkID)
	if err != nil {
		return err
	}
	log.Printf("incTxHash: %v\n", incTxHash)
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] RETRIEVE THE BURN PROOF")
	for {
		burnProof, err := cfg.incClient.GetBurnProof(incTxHash, evmNetworkID)
		if burnProof == nil || err != nil {
			time.Sleep(40 * time.Second)
			log.Println("Wait for the burn proof!")
		} else {
			log.Println("Had the burn proof!!!")
			break
		}
	}
	log.Printf("[STEP 4] FINISHED!\n\n")

	log.Println("[STEP 5] SUBMIT THE BURN PROOF TO THE SC")
	_, err = acc.UnShield(incTxHash, 0, 0, evmNetworkID)
	if err != nil {
		panic(err)
	}
	log.Printf("[STEP 5] FINISHED!\n\n")

	return nil
}

// retryUnShield retries to un-shield a token with an already-burned Incognito TxHash.
func retryUnShield(c *cli.Context) error {
	yesNoPrompt("Do you want to continue?")

	incTxHash := c.String(txHashFlag)

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}

	nativeTokenName := "ETH"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		nativeTokenName = "BNB"
	case rpc.PLGNetworkID:
		nativeTokenName = "MATIC"
	}

	log.Printf("[STEP 1] IMPORT %v ACCOUNT\n", evmNetwork)
	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return err
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return err
	}
	_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
	if err != nil {
		return err
	}
	nativeBalance, _ := tmpNativeBalance.Float64()
	log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), nativeTokenName, nativeBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Println("[STEP 2] RETRIEVE THE BURN PROOF")
	for {
		burnProof, err := cfg.incClient.GetBurnProof(incTxHash, evmNetworkID)
		if burnProof == nil || err != nil {
			time.Sleep(40 * time.Second)
			log.Println("Wait for the burn proof!")
		} else {
			log.Println("Had the burn proof!!!")
			break
		}
	}
	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] SUBMIT THE BURN PROOF TO THE SC")
	_, err = acc.UnShield(incTxHash, 0, 0, evmNetworkID)
	if err != nil {
		panic(err)
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	return nil
}

// bb6755af237f30a00d0b1eb72524f7a7ada41a369fec5fce237380fc0059a70d
