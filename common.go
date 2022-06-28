package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/camelcase"
	iCommon "github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/urfave/cli/v2"
	"regexp"
	"strings"
)

// flag constants
const (
	networkFlag       = "network"
	hostFlag          = "host"
	clientVersionFlag = "clientVersion"
	debugFlag         = "debug"
	cacheFlag         = "utxoCache"
	privateKeyFlag    = "privateKey"
	addressFlag       = "address"
	otaKeyFlag        = "otaKey"
	readonlyKeyFlag   = "readonlyKey"
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

	tokenIDToSellFlag        = "sellTokenID"
	tokenIDToBuyFlag         = "buyTokenID"
	sellingAmountFlag        = "sellingAmount"
	minAcceptableAmountFlag  = "minAcceptAmount"
	tradingFeeFlag           = "tradingFee"
	pairIDFlag               = "pairID"
	tokenID1Flag             = "tokenID1"
	tokenID2Flag             = "tokenID2"
	prvFeeFlag               = "prvFee"
	tradingPathFlag          = "tradingPath"
	maxTradingPathLengthFlag = "maxPaths"
	nftIDFlag                = "nftID"
	orderIDFlag              = "orderID"
	pairHashFlag             = "pairHash"
	amplifierFlag            = "amplifier"

	mnemonicFlag    = "mnemonic"
	numShardsFlag   = "numShards"
	numAccountsFlag = "numAccounts"
	shardIDFlag     = "shardID"

	evmAddressFlag      = "evmAddress"
	tokenAddressFlag    = "externalTokenAddress"
	shieldAmountFlag    = "shieldAmount"
	evmFlag             = "evm"
	externalTxIDFlag    = "externalTxHash"
	externalAddressFlag = "externalAddress"

	miningKeyFlag        = "miningKey"
	candidateAddressFlag = "candidateAddress"
	rewardReceiverFlag   = "rewardAddress"
	autoReStakeFlag      = "autoReStake"

	adminPrivateKeyFlag = "adminPrivateKey"
	tokenNameFlag       = "tokenName"
)

// aliases for defaultFlags
var aliases = map[string][]string{
	networkFlag:          {"net"},
	debugFlag:            {"d"},
	privateKeyFlag:       {"p", "prvKey"},
	otaKeyFlag:           {"ota"},
	readonlyKeyFlag:      {"ro"},
	addressFlag:          {"addr"},
	tokenIDFlag:          {"id", "ID"},
	tokenID1Flag:         {"id1", "ID1"},
	tokenID2Flag:         {"id2", "ID2"},
	amountFlag:           {"amt"},
	versionFlag:          {"v"},
	csvFileFlag:          {"csv"},
	shieldAmountFlag:     {"amt"},
	externalAddressFlag:  {"eAddr"},
	tokenAddressFlag:     {"evmTokenAddr"},
	txHashFlag:           {"iTxID"},
	externalTxIDFlag:     {"eTxID"},
	miningKeyFlag:        {"mKey", "vKey"},
	candidateAddressFlag: {"canAddr"},
	rewardReceiverFlag:   {"rwdAddr"},
	autoReStakeFlag:      {"reStake"},

	tokenIDToSellFlag:       {"sellID", "sellId"},
	tokenIDToBuyFlag:        {"buyID", "buyId"},
	sellingAmountFlag:       {"sellAmt"},
	minAcceptableAmountFlag: {"minAmt"},
	pairIDFlag:              {"pairId"},
	nftIDFlag:               {"nftId"},
	orderIDFlag:             {"orderId"},
	amplifierFlag:           {"amp"},
}

// category constants
const (
	accountCat     = "ACCOUNTS"
	committeeCat   = "COMMITTEES"
	transactionCat = "TRANSACTIONS"
	pDEXCat        = "DEX"
	cenBridgeCat   = "CENTRALIZED BRIDGE"
	evmBridgeCat   = "BRIDGE"
	portalCat      = "BRIDGE"
)

var cfg *Config

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

// isValidOtaKey checks if a base58-encoded ota key is valid or not.
func isValidOtaKey(otaKeyStr string) bool {
	if otaKeyStr == "" {
		return false
	}

	kWallet, err := wallet.Base58CheckDeserialize(otaKeyStr)
	if err != nil {
		return false
	}

	otaKey := kWallet.KeySet.OTAKey

	if otaKey.GetPublicSpend() == nil || otaKey.GetOTASecretKey() == nil {
		return false
	}

	return true
}

// isValidReadonlyKey checks if a base58-encoded read-only key is valid or not.
func isValidReadonlyKey(readonlyKeyStr string) bool {
	if readonlyKeyStr == "" {
		return false
	}

	kWallet, err := wallet.Base58CheckDeserialize(readonlyKeyStr)
	if err != nil {
		return false
	}

	readonlyKey := kWallet.KeySet.ReadonlyKey

	if readonlyKey.GetPublicSpend() == nil || readonlyKey.GetPrivateView() == nil {
		return false
	}

	return true
}

// isValidMiningKey checks if a base58-encoded mining key is valid or not.
func isValidMiningKey(miningKeyStr string) bool {
	if miningKeyStr == "" {
		return false
	}

	_, _, err := base58.Base58Check{}.Decode(miningKeyStr)
	if err != nil {
		return false
	}

	return true
}

// isValidTokenID checks if a string tokenIDStr is valid or not.
func isValidTokenID(tokenIDStr string) bool {
	if tokenIDStr == "" {
		return false
	}

	_, err := iCommon.Hash{}.NewHashFromStr(tokenIDStr)
	if err != nil {
		return false
	}

	return true
}

// isValidDEXPairID checks if a string pairIDStr is valid or not.
func isValidDEXPairID(pairIDStr string) bool {
	if pairIDStr == "" {
		return false
	}
	tmpStrings := strings.Split(pairIDStr, "-")
	if len(tmpStrings) != 3 {
		return false
	}
	for _, tmp := range tmpStrings {
		_, err := iCommon.Hash{}.NewHashFromStr(tmp)
		if err != nil {
			return false
		}
	}

	return true
}

// isValidEVMAddress checks if a string tokenAddress is valid or not.
func isValidEVMAddress(tokenAddress string) bool {
	if tokenAddress == "" {
		return false
	}
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(tokenAddress) {
		return false
	}

	if tokenAddress == nativeToken {
		return true
	}

	tmpTokenAddress := common.HexToAddress(tokenAddress)
	if tmpTokenAddress.String() == nativeToken {
		return false
	}

	return true
}

// isValidIncTxHash checks if an Incognito string txHash is valid or not.
func isValidIncTxHash(txHash string) bool {
	if txHash == "" {
		return false
	}

	_, err := iCommon.Hash{}.NewHashFromStr(txHash)
	if err != nil {
		return false
	}

	return true
}

// isSupportedVersion checks if the given version of transaction is supported or not.
func isSupportedVersion(version int8) bool {
	return version == 1 || version == 2
}

func jsonPrintWithKey(key string, val interface{}) error {
	return jsonPrint(map[string]interface{}{key: val})
}

func jsonPrint(val interface{}) error {
	jsb, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))
	return nil
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

func buildUsageTextFromCommand(command *cli.Command, parents ...string) {
	parent := ""
	if len(parents) > 0 {
		parent = parents[0]
	}
	res := command.Name
	hasOptionalFlags := false
	for _, f := range command.Flags {
		flagString := fmt.Sprintf(" --%v %v", f.Names()[0], flagToVariable(f.Names()[0]))
		if requiredFlag, ok := f.(cli.RequiredFlag); ok {
			if !requiredFlag.IsRequired() {
				// optional flag is put inside a [] symbol.
				flagString = fmt.Sprintf(" [--%v %v]", f.Names()[0], flagToVariable(f.Names()[0]))
				hasOptionalFlags = true
			}
		}
		res += flagString
	}
	if parent != "" {
		res = fmt.Sprintf("%v %v", parent, res)
	}

	command.UsageText = res
	if hasOptionalFlags {
		command.UsageText += "\n\n\t OPTIONAL flags are denoted by a [] bracket."
	}
}
