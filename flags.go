package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
)

var defaultFlags = map[string]cli.Flag{
	networkFlag: &cli.StringFlag{
		Name:        networkFlag,
		Aliases:     aliases[networkFlag],
		Usage:       "network environment (mainnet, testnet, testnet1, devnet, local)",
		Value:       "mainnet",
		Destination: &network,
	},
	hostFlag: &cli.StringFlag{
		Name: hostFlag,
		Usage: "Custom full-node host. This flag is combined with the `network` flag to initialize the environment" +
			" in which the custom host points to.",
		Value:       "",
		Destination: &host,
	},
	clientVersionFlag: &cli.IntFlag{
		Name:        clientVersionFlag,
		Usage:       "version of the incclient",
		Value:       2,
		Destination: &clientVersion,
	},
	debugFlag: &cli.IntFlag{
		Name:        "debug",
		Usage:       "whether to enable the debug mode (0 - disabled, <> 0 - enabled)",
		Value:       1,
		Destination: &debug,
	},

	privateKeyFlag: &cli.StringFlag{
		Name:     privateKeyFlag,
		Aliases:  aliases[privateKeyFlag],
		Usage:    "a base58-encoded Incognito private key",
		Required: true,
	},
	addressFlag: &cli.StringFlag{
		Name:     "address",
		Aliases:  []string{"addr"},
		Usage:    "a base58-encoded payment address",
		Required: true,
	},
	otaKeyFlag: &cli.StringFlag{
		Name:     otaKeyFlag,
		Aliases:  aliases[otaKeyFlag],
		Usage:    "a base58-encoded ota key",
		Required: true,
	},
	readonlyKeyFlag: &cli.StringFlag{
		Name:    readonlyKeyFlag,
		Aliases: aliases[readonlyKeyFlag],
		Usage:   "a base58-encoded read-only key",
		Value:   "",
	},

	tokenIDFlag: &cli.StringFlag{
		Name:    tokenIDFlag,
		Aliases: aliases[tokenIDFlag],
		Usage:   "the Incognito ID of the token",
		Value:   common.PRVIDStr,
	},
	amountFlag: &cli.Uint64Flag{
		Name:     amountFlag,
		Aliases:  aliases[amountFlag],
		Usage:    "the Incognito amount of the action",
		Required: true,
	},
	feeFlag: &cli.Uint64Flag{
		Name:  feeFlag,
		Usage: "the PRV amount for paying the transaction fee",
		Value: incclient.DefaultPRVFee,
	},
	versionFlag: &cli.IntFlag{
		Name:    versionFlag,
		Aliases: aliases[versionFlag],
		Usage:   "version of the transaction (1 or 2)",
		Value:   2,
	},
	numThreadsFlag: &cli.IntFlag{
		Name:  numThreadsFlag,
		Usage: "number of threads used in this action",
		Value: 4,
	},
	enableLogFlag: &cli.BoolFlag{
		Name:  enableLogFlag,
		Usage: "enable log for this action",
		Value: false,
	},
	logFileFlag: &cli.StringFlag{
		Name:  logFileFlag,
		Usage: "location of the log file",
		Value: "os.Stdout",
	},
	csvFileFlag: &cli.StringFlag{
		Name:    csvFileFlag,
		Aliases: aliases[csvFileFlag],
		Usage:   "the csv file location to store the history",
	},
	accessTokenFlag: &cli.StringFlag{
		Name:  accessTokenFlag,
		Usage: "a 64-character long hex-encoded authorized access token",
		Value: "",
	},
	fromHeightFlag: &cli.Uint64Flag{
		Name:  "fromHeight",
		Usage: "the beacon height at which the full-node will sync from",
		Value: 0,
	},
	isResetFlag: &cli.BoolFlag{
		Name:  "isReset",
		Usage: "whether the full-node should reset the cache for this ota key",
		Value: false,
	},
	txHashFlag: &cli.StringFlag{
		Name:     txHashFlag,
		Usage:    "an Incognito transaction hash",
		Required: true,
	},

	tokenIDToSellFlag: &cli.StringFlag{
		Name:     tokenIDToSellFlag,
		Usage:    "ID of the token to sell",
		Required: true,
	},
	tokenIDToBuyFlag: &cli.StringFlag{
		Name:     tokenIDToBuyFlag,
		Usage:    "ID of the token to buy",
		Required: true,
	},
	sellingAmountFlag: &cli.Uint64Flag{
		Name:     sellingAmountFlag,
		Usage:    fmt.Sprintf("the amount of %v wished to sell", tokenIDToSellFlag),
		Required: true,
	},
	minAcceptableAmountFlag: &cli.Uint64Flag{
		Name:  minAcceptableAmountFlag,
		Usage: fmt.Sprintf("the minimum acceptable amount of %v wished to receive", tokenIDToBuyFlag),
		Value: 0,
	},
	tradingFeeFlag: &cli.Uint64Flag{
		Name:  tradingFeeFlag,
		Usage: "the trading fee (measured in nano PRV)",
		Value: 0,
	},
	pairIDFlag: &cli.StringFlag{
		Name:     pairIDFlag,
		Usage:    "the ID of the contributing pair (see https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md)",
		Required: true,
	},
	tokenID1Flag: &cli.StringFlag{
		Name:     tokenID1Flag,
		Usage:    "ID of the first token",
		Required: true,
	},
	tokenID2Flag: &cli.StringFlag{
		Name:  tokenID2Flag,
		Usage: "ID of the second token",
		Value: common.PRVIDStr,
	},
	numShardsFlags: &cli.IntFlag{
		Name:  numShardsFlags,
		Usage: "the number of shard",
		Value: 8,
	},

	evmAddressFlag: &cli.StringFlag{
		Name:  evmAddressFlag,
		Usage: "a hex-encoded address on ETH/BSC networks",
		Value: "",
	},
	tokenAddressFlag: &cli.StringFlag{
		Name:  tokenAddressFlag,
		Usage: "ID of the token on ETH/BSC networks",
		Value: nativeToken,
	},
	shieldAmountFlag: &cli.Float64Flag{
		Name:     shieldAmountFlag,
		Aliases:  aliases[shieldAmountFlag],
		Usage:    "the shielding amount measured in token unit (e.g, 10, 1, 0.1, 0.01)",
		Required: true,
	},
	evmFlag: &cli.StringFlag{
		Name:  evmFlag,
		Usage: "The EVM network (ETH or BSC)",
		Value: "ETH",
	},
	evmTxHash: &cli.StringFlag{
		Name:     evmTxHash,
		Usage:    "the transaction hash on an EVM network (ETH/BSC)",
		Required: true,
	},
}
