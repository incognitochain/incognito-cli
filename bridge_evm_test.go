package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	iCommon "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"log"
	"testing"
	"time"
)

const (
	testEVMPrivateKey = "d2d4c50537f1c15485463e37cb03d11c444a663eb7b84e8f3230b0db38a4d89c"
	testIncPrivateKey = "112t8rnZDRztVgPjbYQiXS7mJgaTzn66NvHD7Vus2SrhSAY611AzADsPFzKjKQCKWTgbkgYrCPo9atvSMoCf9KT23Sc7Js9RKhzbNJkxpJU6"
	erc20TokenAddress = "4f96fe3b7a6cf9725f59d353f723c1bdb64ca6aa"
	bep20TokenAddress = "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	pERC20TokenID     = "c7545459764224a000a9b323850648acf271186238210ce474b505cd17cc93a0"
	pBEP20TokenID     = "a61df4d870c17a7dc62d7e4c16c6f4f847994403842aaaf21c994d1a0024b032"
	pETH              = "ffd8d42dc40a8d166ea4848baf8b5f6e9fe0e9c30d60062eb7d44a8df9e00854"
	pBNB              = "e5032c083f0da67ca141331b6005e4a3740c50218f151a5e829e9d03227e33e2"
)

func TestEVMAccount_GetBalance(t *testing.T) {
	incclient.Logger.IsEnable = true
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}
	acc, err := NewEVMAccount(testEVMPrivateKey)

	ethBalance, err := acc.GetBalance(common.HexToAddress(nativeToken), false)
	if err != nil {
		panic(err)
	}
	log.Printf("balanceETH: %v\n", ethBalance)

	tokenBalance, err := acc.GetBalance(common.HexToAddress(erc20TokenAddress), false)
	if err != nil {
		panic(err)
	}
	log.Printf("balanceToken: %v\n", tokenBalance)
}

func TestEVMAccount_DepositETH(t *testing.T) {
	incclient.Logger.IsEnable = true
	err := NewTestNetConfig(nil)
	if err != nil {
		panic(err)
	}
	acc, err := NewEVMAccount(testEVMPrivateKey)

	incAddress := incclient.PrivateKeyToPaymentAddress(testIncPrivateKey, -1)
	depositAmount := uint64(1000000)

	// create a deposit transaction.
	txHash, err := acc.DepositNative(incAddress, depositAmount, 0, 0)
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
		depositAmount := uint64(1000000)

		// create a deposit transaction.
		txHash, err := acc.DepositToken(incAddress, erc20TokenAddress, depositAmount, 0, 0)
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

		isBSC := (iCommon.RandInt() % 2) == 1
		shieldToken := pETH
		if isBSC {
			shieldToken = pBNB
		}
		log.Printf("ShieldedToken: %v, isBSC:%v\n", shieldToken, isBSC)

		oldIncBalance, err := cfg.incClient.GetBalance(testIncPrivateKey, shieldToken)
		if err != nil {
			panic(err)
		}
		log.Printf("oldIncBalance %v\n", oldIncBalance)

		depositAmount := (1 + iCommon.RandUint64()%10000) * 1e9
		log.Printf("DepositAmount: %v\n", depositAmount)

		ethTxHash, err := acc.DepositNative(incAddress, depositAmount, 0, 0, isBSC)
		if err != nil {
			panic(err)
		}

		ethTxHashStr := ethTxHash.String()
		incTxHash, err := acc.Shield(testIncPrivateKey, shieldToken, ethTxHashStr, isBSC)
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

		expectedReceivedAmount := depositAmount / 1e9
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

		isBSC := (iCommon.RandInt() % 2) == 1
		shieldToken := pERC20TokenID
		publicTokenAddress := erc20TokenAddress
		if isBSC {
			shieldToken = pBEP20TokenID
			publicTokenAddress = bep20TokenAddress
		}
		log.Printf("ShieldedToken: %v, isBSC:%v\n", shieldToken, isBSC)

		depositAmount := (1 + iCommon.RandUint64()%1000) * 1e12
		log.Printf("Address: %v\n", acc.address.String())

		oldIncBalance, err := cfg.incClient.GetBalance(testIncPrivateKey, shieldToken)
		if err != nil {
			panic(err)
		}
		log.Printf("oldIncBalance %v\n", oldIncBalance)

		evmTxHash, err := acc.DepositToken(incAddress, publicTokenAddress, depositAmount, 0, 0, isBSC)
		if err != nil {
			panic(err)
		}

		ethTxHashStr := evmTxHash.String()
		incTxHash, err := acc.Shield(testIncPrivateKey, shieldToken, ethTxHashStr, isBSC)
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

		isBSC := (iCommon.RandInt() % 2) == 0
		unShieldToken := pETH
		if isBSC {
			unShieldToken = pBNB
		}
		log.Printf("UnShieldedToken: %v, isBSC:%v\n", unShieldToken, isBSC)

		oldEVMBalance, err := acc.GetBalance(common.HexToAddress(nativeToken), isBSC)
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
			isBSC,
		)
		if err != nil {
			panic(err)
		}
		log.Printf("incTxHash: %v\n", incTxHash)
		for {
			burnProof, err := cfg.incClient.GetBurnProof(incTxHash, isBSC)
			if burnProof == nil || err != nil {
				log.Println("Sleep 20 seconds for the burnProof!!!")
				time.Sleep(20 * time.Second)
			} else {
				log.Println("Had a burn proof!!!")
				break
			}
		}

		ethTxHash, err := acc.UnShield(incTxHash, 0, 0, isBSC)
		if err != nil {
			panic(err)
		}
		log.Printf("ethWithdrawalTxHash: %v\n", ethTxHash)
		time.Sleep(30 * time.Second)

		newIncBalance, err := acc.GetBalance(common.HexToAddress(nativeToken), isBSC)
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

		isBSC := (iCommon.RandInt() % 2) == 1
		unShieldToken := pERC20TokenID
		publicTokenAddress := erc20TokenAddress
		if isBSC {
			unShieldToken = pBEP20TokenID
			publicTokenAddress = bep20TokenAddress
		}
		log.Printf("UnShieldedToken: %v, isBSC:%v\n", unShieldToken, isBSC)

		oldEVMBalance, err := acc.GetBalance(common.HexToAddress(publicTokenAddress), isBSC)
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
			isBSC,
		)
		if err != nil {
			panic(err)
		}
		log.Printf("incTxHash: %v\n", incTxHash)
		for {
			burnProof, err := cfg.incClient.GetBurnProof(incTxHash, isBSC)
			if burnProof == nil || err != nil {
				log.Println("Sleep 10 seconds for the burnedProof!!!")
				time.Sleep(10 * time.Second)
			} else {
				log.Println("Had a burn proof!!!")
				break
			}
		}

		ethTxHash, err := acc.UnShield(incTxHash, 0, 0, isBSC)
		if err != nil {
			panic(err)
		}
		log.Printf("ethWithdrawalTxHash: %v\n", ethTxHash)
		time.Sleep(30 * time.Second)

		newEVMBalance, err := acc.GetBalance(common.HexToAddress(publicTokenAddress), isBSC)
		if err != nil {
			panic(err)
		}
		diff := newEVMBalance - oldEVMBalance
		log.Printf("newBalace: %v, diff: %v\n", newEVMBalance, diff)

		log.Printf("FINISHED ATTEMTP %v\n\n", i)
	}
}