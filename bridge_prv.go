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
	if evmNetworkID == rpc.PLGNetworkID {
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
	prvTokenAddress := common.HexToAddress(prv20AddressStr)

	shieldAmount := c.Float64(shieldAmountFlag)

	log.Printf("Network: %v, Token: PRV, TokenAddress: %v, ShieldAmount: %v",
		evmNetwork, prv20AddressStr, shieldAmount)
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
		}
		_, tmpNativeBalance, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
		if err != nil {
			return err
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
		return err
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	log.Println("[STEP 4] SHIELD TO INCOGNITO")
	incTxHash, err := ShieldPRV(privateKey, evmHash.String(), evmNetworkID)
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

// retryShieldPRV retries to shield PRV with an already-deposited evm TxHash.
func retryShieldPRV(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}

	evmTxHashStr := c.String(externalTxIDFlag)

	log.Println("[STEP 1] SHIELD TO INCOGNITO")
	incTxHash, err := ShieldPRV(privateKey, evmTxHashStr, evmNetworkID)
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

// unShieldPRV withdraws an amount of PRV on the Incognito network and mint to an EVM network.
func unShieldPRV(c *cli.Context) error {
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

	evmNetwork := c.String(evmFlag)
	evmNetworkID, err := getEVMNetworkIDFromName(evmNetwork)
	if err != nil {
		return err
	}

	log.Printf("Network: %v, Token: PRV, TokenAddress: %v, UnShieldAmount: %v",
		evmNetwork, prv20AddressStr, unShieldAmount)
	yesNoPrompt("Do you want to continue?")
	log.Printf("[STEP 0] FINISHED!\n\n")

	log.Println("[STEP 1] CHECK INCOGNITO BALANCE")
	prvBalance, err := checkSufficientIncBalance(privateKey, iCommon.PRVIDStr, incclient.DefaultPRVFee+unShieldAmount)
	if err != nil {
		return err
	}

	log.Printf("Current PRVBalance: %v\n", prvBalance)
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
	log.Printf("Your %v address: %v, %v: %v\n", evmNetwork, acc.address.String(), nativeTokenName, nativeBalance)
	evmAddress := acc.address
	var res string
	resInBytes, err := promptInput(
		fmt.Sprintf("Un-shield to the following address: %v? Continue? (y/n)", evmAddress.String()),
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
	incTxHash, err := cfg.incClient.CreateAndSendBurningPRVPeggingRequestTransaction(privateKey, evmAddress.String(), unShieldAmount, evmNetworkID)
	if err != nil {
		return err
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
		panic(err)
	}
	log.Printf("[STEP 5] FINISHED!\n\n")

	return nil
}

// retryUnShieldPRV retries to un-shield PRV with an already-burned Incognito TxHash.
func retryUnShieldPRV(c *cli.Context) error {
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
		panic(err)
	}
	log.Printf("[STEP 3] FINISHED!\n\n")

	return nil
}
