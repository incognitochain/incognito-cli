package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	iCommon "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"log"
	"math"
	"testing"
	"time"
)

const (
	testEVMPrivateKey = "d2d4c50537f1c15485463e37cb03d11c444a663eb7b84e8f3230b0db38a4d89c"
	testIncPrivateKey = "112t8rnZDRztVgPjbYQiXS7mJgaTzn66NvHD7Vus2SrhSAY611AzADsPFzKjKQCKWTgbkgYrCPo9atvSMoCf9KT23Sc7Js9RKhzbNJkxpJU6"
	erc20TokenAddress = "4f96fe3b7a6cf9725f59d353f723c1bdb64ca6aa"
	bep20TokenAddress = "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	plg20TokenAddress = "0xfe4F5145f6e09952a5ba9e956ED0C25e3Fa4c7F1"
	pERC20TokenID     = "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0"
	pBEP20TokenID     = "a61df4d870c17a7dc62d7e4c16c6f4f847994403842aaaf21c994d1a0024b032"
	pPLG20TokenID     = "d8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999eb9d18"
	pETH              = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	pBNB              = "e5032c083f0da67ca141331b6005e4a3740c50218f151a5e829e9d03227e33e2"
	pPLG              = "dae027b21d8d57114da11209dce8eeb587d01adf59d4fc356a8be5eedc146859"
)

func init() {
	log.Println("This runs before tests!!")
	askUser = false
}

func TestEstimateGas(t *testing.T) {
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		evmNetworkID := iCommon.RandInt() % 3
		gasPrice, err := estimateGasPrice(evmNetworkID)
		if err != nil {
			panic(err)
		}

		log.Printf("evmNetworkID: %v, gasPrice: %v\n", evmNetworkID, gasPrice.Uint64())
	}
}

func TestEVMAccount_GetBalance(t *testing.T) {
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}
	acc, err := NewEVMAccount(testEVMPrivateKey)

	ethBalance, ethSynthesizedBalance, err := acc.getBalance(common.HexToAddress(nativeToken), rpc.ETHNetworkID)
	if err != nil {
		panic(err)
	}
	fBalance, _ := ethSynthesizedBalance.Float64()
	log.Printf("balanceETH: %v, floatBalance: %v\n", ethBalance.Uint64(), fBalance)

	tokenBalance, tokenSynthesizedBalance, err := acc.getBalance(common.HexToAddress(erc20TokenAddress), rpc.ETHNetworkID)
	if err != nil {
		panic(err)
	}
	log.Printf("balanceToken: %v, floatBalanceToken: %v\n", tokenBalance.Uint64(), tokenSynthesizedBalance.String())
}

func TestEVMAccount_DepositETH(t *testing.T) {
	incclient.Logger.IsEnable = true
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}
	acc, err := NewEVMAccount(testEVMPrivateKey)

	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)
	depositAmount := 0.00001

	// create a deposit transaction.
	txHash, err := acc.DepositNative(incAddress, depositAmount, 0, 0, rpc.ETHNetworkID)
	if err != nil {
		panic(err)
	}

	log.Printf("TxHash: %v\n", txHash.String())
}

func TestEVMAccount_DepositRC20(t *testing.T) {
	incclient.Logger.IsEnable = true
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}
	acc, err := NewEVMAccount(testEVMPrivateKey)

	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)

	for i := 0; i < 10; i++ {
		depositAmount := 0.00001

		// create a deposit transaction.
		txHash, err := acc.DepositToken(incAddress, erc20TokenAddress, depositAmount, 0, 0, rpc.ETHNetworkID)
		if err != nil {
			panic(err)
		}

		log.Printf("TxHash: %v\n", txHash.String())
	}
}

func TestEVMAccount_ShieldNativeToken(t *testing.T) {
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}

	acc, err := NewEVMAccount(testEVMPrivateKey)
	if err != nil {
		panic(err)
	}
	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)
	log.Printf("EVMAddress %v, IncAddress %v\n", acc.address.String(), incAddress)
	for i := 0; i < 10; i++ {
		log.Printf("TEST ATTEMPT %v\n", i)

		evmNetworkID := iCommon.RandInt() % 3
		shieldToken := pETH
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			shieldToken = pBNB
		case rpc.PLGNetworkID:
			shieldToken = pPLG
		}
		log.Printf("ShieldedToken: %v, evmNetworkID:%v\n", shieldToken, evmNetworkID)

		oldIncBalance, err := cfg.incClient.GetBalance(testIncPrivateKey, shieldToken)
		if err != nil {
			panic(err)
		}
		log.Printf("oldIncBalance %v\n", oldIncBalance)

		depositAmount := float64(1+iCommon.RandUint64()%10000) / 1e9
		log.Printf("DepositAmount: %v\n", depositAmount)

		ethTxHash, err := acc.DepositNative(incAddress, depositAmount, 0, 0, evmNetworkID)
		if err != nil {
			panic(err)
		}

		ethTxHashStr := ethTxHash.String()
		incTxHash, err := Shield(testIncPrivateKey, shieldToken, ethTxHashStr, evmNetworkID)
		if err != nil {
			panic(err)
		}
		log.Printf("IncognitoShieldedTx: %v\n", incTxHash)

		for {
			status, err := cfg.incClient.CheckShieldStatus(incTxHash)
			if err != nil || status <= 1 {
				log.Printf("ShieldingStatus: %v\n", status)
				log.Println("Sleep 10 seconds!!")
				time.Sleep(10 * time.Second)
				continue
			} else if status == 2 {
				log.Println("Shielding SUCCEEDED!!")
				break
			} else {
				panic("Shielding FAILED!!")
			}
		}

		expectedReceivedAmount := uint64(depositAmount * math.Pow10(9))
		for {
			newIncBalance, err := cfg.incClient.GetBalance(testIncPrivateKey, shieldToken)
			if err != nil {
				panic(err)
			}
			if newIncBalance != oldIncBalance {
				if newIncBalance-oldIncBalance != expectedReceivedAmount {
					panic(fmt.Sprintf("expectedReceived %v, got %v\n", expectedReceivedAmount, newIncBalance-oldIncBalance))
				}
				log.Printf("newIncBalance: %v\n", newIncBalance)
				break
			} else {
				log.Println("balance not updated!!")
				time.Sleep(10 * time.Second)
			}
		}

		log.Printf("FINISHED ATTEMTP %v\n\n", i)
	}
}

func TestEVMAccount_ShieldToken(t *testing.T) {
	//incclient.Logger.IsEnable = true
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}

	acc, err := NewEVMAccount(testEVMPrivateKey)
	if err != nil {
		panic(err)
	}
	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)

	for i := 0; i < 10; i++ {
		log.Printf("TEST ATTEMPT %v\n", i)

		evmNetworkID := iCommon.RandInt() % 3
		shieldToken := pERC20TokenID
		publicTokenAddress := erc20TokenAddress
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			shieldToken = pBEP20TokenID
			publicTokenAddress = bep20TokenAddress
		case rpc.PLGNetworkID:
			shieldToken = pPLG20TokenID
			publicTokenAddress = plg20TokenAddress
		}
		log.Printf("ShieldedToken: %v, evmNetworkID:%v\n", shieldToken, evmNetworkID)

		tokenDecimals, err := getDecimals(common.HexToAddress(publicTokenAddress), evmNetworkID)
		if err != nil {
			panic(err)
		}
		depositAmount := float64(1+iCommon.RandUint64()%uint64(math.Pow10(int(tokenDecimals-3)))) / math.Pow10(int(tokenDecimals))

		oldIncBalance, err := cfg.incClient.GetBalance(testIncPrivateKey, shieldToken)
		if err != nil {
			panic(err)
		}
		log.Printf("oldIncBalance %v\n", oldIncBalance)

		evmTxHash, err := acc.DepositToken(incAddress, publicTokenAddress, depositAmount, 0, 0, evmNetworkID)
		if err != nil {
			panic(err)
		}

		ethTxHashStr := evmTxHash.String()
		incTxHash, err := Shield(testIncPrivateKey, shieldToken, ethTxHashStr, evmNetworkID)
		if err != nil {
			panic(err)
		}

		for {
			status, err := cfg.incClient.CheckShieldStatus(incTxHash)
			if err != nil || status <= 1 {
				log.Printf("ShieldingStatus: %v\n", status)
				log.Println("Sleep 10 seconds!!")
				time.Sleep(10 * time.Second)
				continue
			} else if status == 2 {
				log.Println("Shielding SUCCEEDED!!")
				break
			} else {
				panic("Shielding FAILED!!")
			}
		}

		for {
			newIncBalance, err := cfg.incClient.GetBalance(testIncPrivateKey, shieldToken)
			if err != nil {
				panic(err)
			}
			if newIncBalance != oldIncBalance {
				log.Printf("newIncBalance: %v\n", newIncBalance)
				break
			} else {
				log.Println("balance not updated!!")
				time.Sleep(10 * time.Second)
			}
		}

		log.Printf("FINISHED ATTEMTP %v\n\n", i)
	}
}

func TestEVMAccount_UnShieldNativeToken(t *testing.T) {
	//incclient.Logger.IsEnable = true
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}

	acc, err := NewEVMAccount(testEVMPrivateKey)
	if err != nil {
		panic(err)
	}
	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)
	log.Printf("EVMAddress %v, IncAddress %v\n", acc.address.String(), incAddress)

	for i := 0; i < 10; i++ {
		log.Printf("TEST ATTEMPT %v\n", i)

		evmNetworkID := iCommon.RandInt() % 3
		unShieldToken := pETH
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			unShieldToken = pBNB
		case rpc.PLGNetworkID:
			unShieldToken = pPLG
		}

		log.Printf("UnShieldedToken: %v, evmNetworkID:%v\n", unShieldToken, evmNetworkID)

		oldEVMBalance, _, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
		if err != nil {
			panic(err)
		}
		log.Printf("oldEVMBalance %v\n", oldEVMBalance)

		withdrawalAmount := 1 + iCommon.RandUint64()%1000
		log.Printf("WithdrawalAmount: %v\n", withdrawalAmount)

		incTxHash, err := cfg.incClient.CreateAndSendBurningRequestTransaction(
			testIncPrivateKey,
			acc.address.String(),
			unShieldToken,
			withdrawalAmount,
			evmNetworkID,
		)
		if err != nil {
			panic(err)
		}
		log.Printf("incTxHash: %v\n", incTxHash)
		for {
			burnProof, err := cfg.incClient.GetBurnProof(incTxHash, evmNetworkID)
			if burnProof == nil || err != nil {
				log.Println("Sleep 20 seconds for the burnProof!!!")
				time.Sleep(20 * time.Second)
			} else {
				log.Println("Had a burn proof!!!")
				break
			}
		}

		ethTxHash, err := acc.UnShield(incTxHash, 0, 0, evmNetworkID)
		if err != nil {
			panic(err)
		}
		log.Printf("ethWithdrawalTxHash: %v\n", ethTxHash)
		time.Sleep(30 * time.Second)

		newIncBalance, _, err := acc.getBalance(common.HexToAddress(nativeToken), evmNetworkID)
		if err != nil {
			panic(err)
		}
		log.Printf("newBalace: %v\n", newIncBalance)

		log.Printf("FINISHED ATTEMTP %v\n\n", i)
	}
}

func TestEVMAccount_UnshieldToken(t *testing.T) {
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}

	acc, err := NewEVMAccount(testEVMPrivateKey)
	if err != nil {
		panic(err)
	}
	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)
	log.Printf("EVMAddress %v, IncAddress %v\n", acc.address.String(), incAddress)

	for i := 0; i < 10; i++ {
		log.Printf("TEST ATTEMPT %v\n", i)

		evmNetworkID := iCommon.RandInt() % 3
		unShieldToken := pERC20TokenID
		publicTokenAddress := erc20TokenAddress
		switch evmNetworkID {
		case rpc.BSCNetworkID:
			unShieldToken = pBEP20TokenID
			publicTokenAddress = bep20TokenAddress
		case rpc.PLGNetworkID:
			unShieldToken = pPLG20TokenID
			publicTokenAddress = plg20TokenAddress
		}
		log.Printf("UnShieldedToken: %v, evmNetworkID:%v\n", unShieldToken, evmNetworkID)

		oldEVMBalance, _, err := acc.getBalance(common.HexToAddress(publicTokenAddress), evmNetworkID)
		if err != nil {
			panic(err)
		}
		log.Printf("oldEVMBalance %v\n", oldEVMBalance)

		withdrawalAmount := 1 + iCommon.RandUint64()%1000
		log.Printf("WithdrawalAmount: %v\n", withdrawalAmount)

		incTxHash, err := cfg.incClient.CreateAndSendBurningRequestTransaction(
			testIncPrivateKey,
			acc.address.String(),
			unShieldToken,
			withdrawalAmount,
			evmNetworkID,
		)
		if err != nil {
			panic(err)
		}
		log.Printf("incTxHash: %v\n", incTxHash)
		for {
			burnProof, err := cfg.incClient.GetBurnProof(incTxHash, evmNetworkID)
			if burnProof == nil || err != nil {
				log.Println("Sleep 10 seconds for the burnedProof!!!")
				time.Sleep(10 * time.Second)
			} else {
				log.Println("Had a burn proof!!!")
				break
			}
		}

		ethTxHash, err := acc.UnShield(incTxHash, 0, 0, evmNetworkID)
		if err != nil {
			panic(err)
		}
		log.Printf("ethWithdrawalTxHash: %v\n", ethTxHash)
		time.Sleep(30 * time.Second)

		newEVMBalance, _, err := acc.getBalance(common.HexToAddress(publicTokenAddress), evmNetworkID)
		if err != nil {
			panic(err)
		}
		diff := newEVMBalance.Uint64() - oldEVMBalance.Uint64()
		log.Printf("newBalace: %v, diff: %v\n", newEVMBalance, diff)

		log.Printf("FINISHED ATTEMTP %v\n\n", i)
	}
}
