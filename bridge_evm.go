package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	iCommon "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

// shield deposits an EVM token (ETH/BNB/ERC20/BEP20) into the Incognito chain.
func shield(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

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
	if evmNetwork != "ETH" && evmNetwork != "BSC" {
		return fmt.Errorf("%v is invalid", evmFlag)
	}
	isBSC := evmNetwork == "BSC"

	shieldAmount := c.Float64(shieldAmountFlag)
	tokenAddressStr := c.String(tokenAddressFlag)
	if !isValidTokenAddress(tokenAddressStr) {
		return fmt.Errorf("%v is invalid", tokenAddressFlag)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	incTokenID, err := getIncTokenIDFromEVMTokenID(tokenAddress.String(), isBSC)
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
		if isBSC {
			tokenName = "Binance"
			tokenSymbol = "BNB"
		}
	} else {
		tokenInfo, err := getEVMTokenInfo(tokenAddress.String())
		if err != nil {
			return err
		}
		tokenName = tokenInfo.name
		tokenSymbol = tokenInfo.symbol
		if tokenInfo.network != evmNetwork {
			return fmt.Errorf("expect token to be on `%` network, got `%v`", evmNetwork, tokenInfo.network)
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
	evmTokenBalance, err := acc.checkSufficientBalance(tokenAddress, shieldAmount, isBSC)
	if err != nil {
		return err
	}
	if tokenAddress.String() == nativeToken {
		log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), tokenName, evmTokenBalance)
	} else {
		nativeTokenName := "ETH"
		if isBSC {
			nativeTokenName = "BNB"
		}
		_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), isBSC)
		if err != nil {
			return err
		}
		nativeBalance, _ := tmpNativeBalance.Float64()
		log.Printf("Your %v address: %v, %v: %v, %v: %v\n",
			evmNetwork,
			acc.address.String(), nativeTokenName, nativeBalance, tokenSymbol, evmTokenBalance)
	}
	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] DEPOSIT PUBLIC TOKEN TO SC")
	var evmHash *common.Hash
	if tokenAddress.String() == nativeToken {
		evmHash, err = acc.DepositNative(incAddress, shieldAmount, 0, 0, isBSC)
	} else {
		evmHash, err = acc.DepositToken(incAddress, tokenAddressStr, shieldAmount, 0, 0, isBSC)
	}
	if err != nil {
		return err
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] SHIELD TO INCOGNITO")
	incTxHash, err := Shield(privateKey, incTokenID, evmHash.String(), isBSC)
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
	err := initNetWork()
	if err != nil {
		return err
	}

	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	evmNetwork := c.String(evmFlag)
	if evmNetwork != "ETH" && evmNetwork != "BSC" {
		return fmt.Errorf("%v is invalid", evmFlag)
	}
	isBSC := evmNetwork == "BSC"

	tokenAddressStr := c.String(tokenAddressFlag)
	if !isValidTokenAddress(tokenAddressStr) {
		return fmt.Errorf("%v is invalid", tokenAddressFlag)
	}
	tokenAddress := common.HexToAddress(tokenAddressStr)
	incTokenID, err := getIncTokenIDFromEVMTokenID(tokenAddress.String(), isBSC)
	if err != nil {
		if strings.Contains(err.Error(), "incTokenID not found") {
			log.Printf("IncTokenID not found for %v, perhaps it doesn't exist in the Incognito network.\n", tokenAddress.String())
			incTokenID = fmt.Sprintf("%x", iCommon.RandBytes(32))
			yesNoPrompt(fmt.Sprintf("Newly generated incTokenID: %v. Do you want to continue with this token?", incTokenID))
		} else {
			return err
		}
		return err
	}

	evmTxHashStr := c.String(evmTxHash)
	evmHash := common.HexToHash(evmTxHashStr)

	log.Println("[STEP 1] SHIELD TO INCOGNITO")
	incTxHash, err := Shield(privateKey, incTokenID, evmHash.String(), isBSC)
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
