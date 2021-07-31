package main

import (
	"fmt"
	"github.com/fatih/camelcase"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/urfave/cli/v2"
	"strings"
)

// flag constants
const (
	networkFlag       = "network"
	hostFlag          = "host"
	clientVersionFlag = "clientVersion"
	privateKeyFlag    = "privateKey"
	addressFlag       = "address"
	otaKeyFlag        = "otaKey"
	tokenIDFlag       = "tokenID"
	amountFlag        = "amount"
	feeFlag           = "fee"
	versionFlag       = "version"
	numThreadsFlag    = "numThreads"
	logFileFlag       = "logFile"
	enableLogFlag     = "enableLog"
	csvFileFlag       = "csvFile"
	accessTokenFlag   = "accessToken"
	fromHeightFlag    = "fromHeight"
	isResetFlag       = "isReset"
	txHashFlag        = "txHash"

	tokenIDToSellFlag       = "sellTokenID"
	tokenIDToBuyFlag        = "buyTokenID"
	sellingAmountFlag       = "sellingAmount"
	minAcceptableAmountFlag = "minAcceptAmount"
	tradingFeeFlag          = "tradingFee"
	pairIDFlag              = "pairId"
	tokenID1Flag            = "tokenID1"
	tokenID2Flag            = "tokenID2"
)

// aliases for defaultFlags
var aliases = map[string][]string{
	networkFlag:    {"net"},
	privateKeyFlag: {"prvKey"},
	otaKeyFlag:     {"ota"},
	addressFlag:    {"addr"},
	amountFlag:     {"amt"},
	versionFlag:    {"v"},
	csvFileFlag:    {"csv"},
}

// category constants
const (
	accountCat     = "ACCOUNTS"
	committeeCat   = "COMMITTEES"
	transactionCat = "TRANSACTIONS"
	pDEXCat        = "PDEX"
)

var client *incclient.IncClient

// isValidPrivateKey checks if a base58-encoded private key is valid or not.
func isValidPrivateKey(privateKey string) bool {
	if privateKey == "" {
		return false
	}

	kWallet, err := wallet.Base58CheckDeserialize(privateKey)
	if err != nil {
		return false
	}

	if kWallet.KeySet.PrivateKey == nil {
		return false
	}

	return true
}

// isValidAddress checks if a base58-encoded payment address is valid or not.
func isValidAddress(address string) bool {
	if address == "" {
		return false
	}

	kWallet, err := wallet.Base58CheckDeserialize(address)
	if err != nil {
		return false
	}

	if kWallet.KeySet.PaymentAddress.Pk == nil {
		return false
	}
	if kWallet.KeySet.PaymentAddress.Tk == nil {
		return false
	}

	return true
}

// isValidTokenID checks if a string tokenIDStr is valid or not.
func isValidTokenID(tokenIDStr string) bool {
	if tokenIDStr == "" {
		return false
	}

	_, err := common.Hash{}.NewHashFromStr(tokenIDStr)
	if err != nil {
		return false
	}

	return true
}

// isSupportedVersion checks if the given version of transaction is supported or not.
func isSupportedVersion(version int8) bool {
	return version == 1 || version == 2
}

// flagToVariable gets the variable representation for a flag.
// The variable representation of a flag is a ALL_UPPER_CASE form of a flag.
//
// For example, the variable resp of the flag `privateKey` is `PRIVATE_KEY`.
func flagToVariable(f string) string {
	f = strings.Replace(f, "Flag", "", 1)

	words := camelcase.Split(f)
	res := ""
	for _, word := range words {
		if res == "" {
			res += strings.ToUpper(word)
		} else {
			res += "_" + strings.ToUpper(word)
		}
	}

	return res
}

func buildUsageTextFromCommand(command *cli.Command) {
	res := command.Name
	for _, f := range command.Flags {
		flagString := fmt.Sprintf(" --%v %v", f.Names()[0], flagToVariable(f.Names()[0]))
		if requiredFlag, ok := f.(cli.RequiredFlag); ok {
			if !requiredFlag.IsRequired() {
				// optional flag is put inside a [] symbol.
				flagString = fmt.Sprintf(" [--%v %v]", f.Names()[0], flagToVariable(f.Names()[0]))
			}
		}
		res += flagString
	}

	command.UsageText = res + "\n\n\t OPTIONAL flags are denoted by a [] bracket."
}
