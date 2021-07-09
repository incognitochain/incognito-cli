package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
)

var defaultFlags = map[string]cli.Flag{
	privateKeyFlag: &cli.StringFlag{
		Name:     privateKeyFlag,
		Aliases:  aliases[privateKeyFlag],
		Usage:    "a base58-encoded private key",
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
	tokenIDFlag: &cli.StringFlag{
		Name:  tokenIDFlag,
		Usage: "ID of the token",
		Value: common.PRVIDStr,
	},
	amountFlag: &cli.Uint64Flag{
		Name:     amountFlag,
		Aliases:  aliases[amountFlag],
		Usage:    "the amount being transferred",
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
		Value:   1,
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
}
