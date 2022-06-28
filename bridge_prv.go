package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	iCommon "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

var shieldPRVMessage = "This function helps shield an EVM-pegged PRV token into the Incognito network. " +
	"It will ask for users' EVM PRIVATE KEY to proceed. " +
	"The shielding process consists of the following operations.\n" +
	"\t 1. Deposit the EVM asset into the corresponding smart contract.\n" +
	"\t\t 1.1. A PRV-approval transaction is performed (if needed) the before the " +
	"actual deposit. For this operation, a prompt will be displayed to ask for user's approval.\n" +
	"\t 2. Get the deposited EVM transaction, parse the depositing proof and submit it to the Incognito network. " +
	"This step requires an Incognito private key with a sufficient amount of PRV to create an issuing transaction.\n\n" +
	"Note that EVM shielding is a complicated process, users MUST understand how the process works before using this function. " +
	"We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.\n" +
	"DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE SHIELDING PROCESS."

var unShieldPRVMessage = "This function helps withdraw PRV to an EVM network. " +
	"The un-shielding process consists the following operations.\n" +
	"\t 1. Users burn the token inside the Incognito chain.\n" +
	"\t 2. After the burning is successful, wait for 1-2 Incognito blocks and retrieve the corresponding burn proof from " +
	"the Incognito chain.\n" +
	"\t 3. After successfully retrieving the burn proof, users submit the burn proof to the smart contract to get back the " +
	"corresponding public token. This step will ask for users' EVM PRIVATE KEY to proceed. Note that ONLY UNTIL this step, " +
	"it is feasible to estimate the actual un-shielding fee (mainly is the fee interacting with the smart contract).\n\n" +
	"Please be aware that EVM un-shielding is a complicated process; and once burned, there is NO WAY to recover the asset inside the " +
	"Incognito network. Therefore, use this function IF ADN ONLY IF you understand the way un-shielding works. " +
	"Otherwise, use the un-shielding function from the Incognito app. " +
	"We RECOMMEND users test the function with test networks BEFORE performing it on the live networks.\n" +
	"DO NOT USE THIS FUNCTION UNLESS YOU UNDERSTAND THE UN-SHIELDING PROCESS."

var prv20AddressStr string

func prvInitFunc(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}
	if evmNetworkID == rpc.PLGNetworkID || evmNetworkID == rpc.FTMNetworkID {
		return errEVMNetworkNotSupported(evmNetworkID)
	}

	if evmNetworkID == rpc.ETHNetworkID {
		prv20AddressStr = incclient.MainNetPRVERC20ContractAddressStr
		switch network {
		case "testnet":
			prv20AddressStr = incclient.TestNetPRVERC20ContractAddressStr
		case "testnet1":
			prv20AddressStr = incclient.TestNet1PRVERC20ContractAddressStr
		}
	}
	if evmNetworkID == rpc.BSCNetworkID {
		prv20AddressStr = incclient.MainNetPRVBEP20ContractAddressStr
		switch network {
		case "testnet":
			prv20AddressStr = incclient.TestNetPRVBEP20ContractAddressStr
		case "testnet1":
			prv20AddressStr = incclient.TestNet1PRVBEP20ContractAddressStr
		}
	}

	if !isValidEVMAddress(prv20AddressStr) {
		return fmt.Errorf("PRV20 address is invalid")
	}

	return nil
}

// shieldPRV deposits PRV tokens (on ETH/BSC) into the Incognito chain.
func shieldPRV(c *cli.Context) error {
	fmt.Println(shieldPRVMessage)
	yesNoPrompt("Do you want to continue?")
	fmt.Println()

	log.Println("[STEP 0] PREPARE DATA")
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	incAddress := c.String(addressFlag)
	if incAddress == "" {
		incAddress = incclient.PrivateKeyToPaymentAddress(privateKey, -1)
	}
	if !isValidAddress(incAddress) {
		return newAppError(InvalidPaymentAddressError)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return newAppError(GetEVMNetworkError, err)
	}
	prvTokenAddress := common.HexToAddress(prv20AddressStr)

	shieldAmount := c.Float64(shieldAmountFlag)

	log.Printf("Network: %v, Token: PRV, TokenAddress: %v, ShieldAmount: %v",
		evmNetwork, prv20AddressStr, shieldAmount)
	yesNoPrompt("Do you want to continue?")
	log.Printf("[STEP 0] FINISHED!\n\n")

	log.Println("[STEP 1] CHECK INCOGNITO BALANCE")
	prvBalance, err := checkSufficientIncBalance(privateKey, iCommon.PRVIDStr, incclient.DefaultPRVFee)
	if err != nil {
		return newAppError(InsufficientBalanceError, err)
	}
	log.Printf("Current PRV balance: %v\n", prvBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Printf("[STEP 2] IMPORT %v ACCOUNT\n", evmNetwork)

	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return newAppError(UserInputError, err)
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return newAppError(NewEVMAccountError, err)
	}

	for {
		evmTokenBalance, err := acc.checkSufficientBalance(prvTokenAddress, shieldAmount, evmNetworkID)
		err = checkAndChangeRPCEndPoint(evmNetworkID, err)
		if err != nil {
			return err
		}

		nativeTokenName := "ETH"
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			nativeTokenName = "BNB"
		case rpc.PLGNetworkID:
			nativeTokenName = "MATIC"
		case rpc.FTMNetworkID:
			nativeTokenName = "FTM"
		}
		_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
		if err != nil {
			return newAppError(GetEVMBalanceError, err)
		}
		nativeBalance, _ := tmpNativeBalance.Float64()

		log.Printf("Your %v address: %v, %v: %v, PRV: %v\n",
			evmNetwork,
			acc.address.String(), nativeTokenName, nativeBalance, evmTokenBalance)
		break
	}

	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] DEPOSIT PUBLIC TOKEN TO SC")
	evmHash, err := acc.BurnPRVOnEVM(incAddress, shieldAmount, 0, 0, evmNetworkID)
	if err != nil {
		return newAppError(EVMBurnPRVError, err)
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] SHIELD TO INCOGNITO")
	incTxHash, err := ShieldPRV(privateKey, evmHash.String(), evmNetworkID)
	if err != nil {
		return newAppError(CreatePRVShieldingTransactionError, err)
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

// retryShieldPRV retries to shield PRV with an already-deposited evm TxHash.
func retryShieldPRV(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return newAppError(GetEVMNetworkError, err)
	}

	evmTxHashStr := c.String(externalTxIDFlag)

	log.Println("[STEP 1] SHIELD TO INCOGNITO")
	incTxHash, err := ShieldPRV(privateKey, evmTxHashStr, evmNetworkID)
	if err != nil {
		return newAppError(CreatePRVShieldingTransactionError, err)
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

// unShieldPRV withdraws an amount of PRV on the Incognito network and mint to an EVM network.
func unShieldPRV(c *cli.Context) error {
	fmt.Println(unShieldPRVMessage)
	yesNoPrompt("Do you want to continue?")
	fmt.Println()

	log.Println("[STEP 0] PREPARE DATA")
	// get the private key
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return newAppError(InvalidPrivateKeyError)
	}

	// get the un-shield amount
	unShieldAmount := c.Uint64(amountFlag)
	if unShieldAmount == 0 {
		return newAppError(InvalidAmountError)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return newAppError(GetEVMNetworkError, err)
	}

	log.Printf("Network: %v, Token: PRV, TokenAddress: %v, UnShieldAmount: %v",
		evmNetwork, prv20AddressStr, unShieldAmount)
	yesNoPrompt("Do you want to continue?")
	log.Printf("[STEP 0] FINISHED!\n\n")

	log.Println("[STEP 1] CHECK INCOGNITO BALANCE")
	prvBalance, err := checkSufficientIncBalance(privateKey, iCommon.PRVIDStr, incclient.DefaultPRVFee+unShieldAmount)
	if err != nil {
		return newAppError(InsufficientBalanceError, err)
	}

	log.Printf("Current PRVBalance: %v\n", prvBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Printf("[STEP 2] IMPORT %v ACCOUNT\n", evmNetwork)

	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return newAppError(UserInputError, err)
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return newAppError(NewEVMAccountError, err)
	}

	nativeTokenName := "ETH"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		nativeTokenName = "BNB"
	case rpc.PLGNetworkID:
		nativeTokenName = "MATIC"
	case rpc.FTMNetworkID:
		nativeTokenName = "FTM"
	}

	_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
	if err != nil {
		return newAppError(GetEVMBalanceError, err)
	}
	nativeBalance, _ := tmpNativeBalance.Float64()
	log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), nativeTokenName, nativeBalance)
	evmAddress := acc.address
	var res string
	resInBytes, err := promptInput(
		fmt.Sprintf("Un-shield to the following address: %v. Continue? (y/n)", evmAddress.String()),
		&res)
	if err != nil {
		return newAppError(UnexpectedError, err)
	}
	res = string(resInBytes)
	if !strings.Contains(res, "y") && !strings.Contains(res, "Y") {
		resInBytes, err = promptInput(
			fmt.Sprintf("Enter the address you want to un-shield to"),
			&res)
		if err != nil {
			return newAppError(UserInputError, err)
		}
		res = string(resInBytes)
		if !isValidEVMAddress(res) {
			return newAppError(InvalidExternalAddressError)
		}
		evmAddress = common.HexToAddress(res)
	}
	log.Printf("[STEP 2] FINISHED!\n\n")

	log.Println("[STEP 3] BURN INCOGNITO TOKEN")
	incTxHash, err := cfg.incClient.CreateAndSendBurningPRVPeggingRequestTransaction(privateKey, evmAddress.String(), unShieldAmount, evmNetworkID)
	if err != nil {
		return newAppError(CreatePRVUnShieldingTransactionError, err)
	}
	log.Printf("incTxHash: %v\n", incTxHash)
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] RETRIEVE THE BURN PROOF")
	for {
		burnProof, err := cfg.incClient.GetBurnPRVPeggingProof(incTxHash, evmNetworkID)
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
	_, err = acc.UnShieldPRV(incTxHash, 0, 0, evmNetworkID)
	if err != nil {
		return newAppError(EVMMintPRVError, err)
	}
	log.Printf("[STEP 5] FINISHED!\n\n")

	return nil
}

// retryUnShieldPRV retries to un-shield PRV with an already-burned Incognito TxHash.
func retryUnShieldPRV(c *cli.Context) error {
	yesNoPrompt("Do you want to continue?")

	incTxHash := c.String(txHashFlag)
	if !isValidIncTxHash(incTxHash) {
		return newAppError(InvalidIncognitoTxHashError)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return newAppError(GetEVMNetworkError, err)
	}

	nativeTokenName := "ETH"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		nativeTokenName = "BNB"
	case rpc.PLGNetworkID:
		nativeTokenName = "MATIC"
	case rpc.FTMNetworkID:
		nativeTokenName = "FTM"
	}

	log.Printf("[STEP 1] IMPORT %v ACCOUNT\n", evmNetwork)
	// Get EVM account
	var privateEVMKey string
	input, err := promptInput(fmt.Sprintf("Enter your %v private key", evmNetwork), &privateEVMKey, true)
	if err != nil {
		return newAppError(UserInputError, err)
	}
	privateEVMKey = string(input)
	acc, err := NewEVMAccount(privateEVMKey)
	if err != nil {
		return newAppError(NewEVMAccountError, err)
	}
	_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
	if err != nil {
		return newAppError(GetEVMBalanceError, err)
	}
	nativeBalance, _ := tmpNativeBalance.Float64()
	log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), nativeTokenName, nativeBalance)
	log.Printf("[STEP 1] FINISHED!\n\n")

	log.Println("[STEP 2] RETRIEVE THE BURN PROOF")
	for {
		burnProof, err := cfg.incClient.GetBurnPRVPeggingProof(incTxHash, evmNetworkID)
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
	_, err = acc.UnShieldPRV(incTxHash, 0, 0, evmNetworkID)
	if err != nil {
		return newAppError(EVMMintPRVError, err)
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	return nil
}
