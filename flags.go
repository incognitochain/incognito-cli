package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/incognito-cli/pdex_v3"
	"github.com/urfave/cli/v2"
)

var defaultFlags = map[string]cli.Flag{
	networkFlag: &cli.StringFlag{
		Name:        networkFlag,
		Aliases:     aliases[networkFlag],
		Usage:       "Network environment (mainnet, testnet, testnet1, local)",
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
		Usage:       "Version of the incclient",
		Value:       2,
		Destination: &clientVersion,
	},
	debugFlag: &cli.IntFlag{
		Name:        debugFlag,
		Aliases:     []string{"d"},
		Usage:       "Whether to enable the debug mode (0 - disabled, <> 0 - enabled)",
		Value:       0,
		Destination: &debug,
	},
	cacheFlag: &cli.IntFlag{
		Name:        cacheFlag,
		Aliases:     []string{"c", "cache"},
		Usage:       "Whether to use the UTXO cache (0 - disabled, <> 0 - enabled). See https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/utxo_cache.md for more information.",
		Value:       0,
		Destination: &cache,
	},
	privateKeyFlag: &cli.StringFlag{
		Name:     privateKeyFlag,
		Aliases:  aliases[privateKeyFlag],
		Usage:    "A base58-encoded Incognito private key",
		Required: true,
	},
	addressFlag: &cli.StringFlag{
		Name:     addressFlag,
		Aliases:  []string{"addr"},
		Usage:    "A base58-encoded payment address",
		Required: true,
	},
	otaKeyFlag: &cli.StringFlag{
		Name:     otaKeyFlag,
		Aliases:  aliases[otaKeyFlag],
		Usage:    "A base58-encoded ota key",
		Required: true,
	},
	readonlyKeyFlag: &cli.StringFlag{
		Name:    readonlyKeyFlag,
		Aliases: aliases[readonlyKeyFlag],
		Usage:   "A base58-encoded read-only key",
		Value:   "",
	},

	tokenIDFlag: &cli.StringFlag{
		Name:    tokenIDFlag,
		Aliases: aliases[tokenIDFlag],
		Usage:   "The Incognito ID of the token",
		Value:   common.PRVIDStr,
	},
	amountFlag: &cli.Uint64Flag{
		Name:     amountFlag,
		Aliases:  aliases[amountFlag],
		Usage:    "The Incognito (uint64) amount of the action (e.g, 1000, 1000000, 1000000000)",
		Required: true,
	},
	feeFlag: &cli.Uint64Flag{
		Name:  feeFlag,
		Usage: "The PRV amount for paying the transaction fee",
		Value: incclient.DefaultPRVFee,
	},
	versionFlag: &cli.IntFlag{
		Name:    versionFlag,
		Aliases: aliases[versionFlag],
		Usage:   "Version of the transaction (1 or 2)",
		Value:   2,
	},
	numThreadsFlag: &cli.IntFlag{
		Name:  numThreadsFlag,
		Usage: "Number of threads used in this action",
		Value: 4,
	},
	enableLogFlag: &cli.BoolFlag{
		Name:  enableLogFlag,
		Usage: "Enable log for this action",
		Value: false,
	},
	logFileFlag: &cli.StringFlag{
		Name:  logFileFlag,
		Usage: "Location of the log file",
		Value: "os.Stdout",
	},
	csvFileFlag: &cli.StringFlag{
		Name:    csvFileFlag,
		Aliases: aliases[csvFileFlag],
		Usage:   "The csv file location to store the history",
	},
	accessTokenFlag: &cli.StringFlag{
		Name:  accessTokenFlag,
		Usage: "A 64-character long hex-encoded authorized access token",
		Value: "",
	},
	fromHeightFlag: &cli.Uint64Flag{
		Name:  fromHeightFlag,
		Usage: "The beacon height at which the full-node will sync from",
		Value: 0,
	},
	isResetFlag: &cli.BoolFlag{
		Name:  isResetFlag,
		Usage: "Whether the full-node should reset the cache for this ota key",
		Value: false,
	},
	txHashFlag: &cli.StringFlag{
		Name:     txHashFlag,
		Aliases:  aliases[txHashFlag],
		Usage:    "An Incognito transaction hash",
		Required: true,
	},

	tokenIDToSellFlag: &cli.StringFlag{
		Name:     tokenIDToSellFlag,
		Aliases:  aliases[tokenIDToSellFlag],
		Usage:    "ID of the token to sell",
		Required: true,
	},
	tokenIDToBuyFlag: &cli.StringFlag{
		Name:     tokenIDToBuyFlag,
		Aliases:  aliases[tokenIDToBuyFlag],
		Usage:    "ID of the token to buy",
		Required: true,
	},
	sellingAmountFlag: &cli.Uint64Flag{
		Name:     sellingAmountFlag,
		Aliases:  aliases[sellingAmountFlag],
		Usage:    fmt.Sprintf("The amount of %v wished to sell", tokenIDToSellFlag),
		Required: true,
	},
	minAcceptableAmountFlag: &cli.Uint64Flag{
		Name:    minAcceptableAmountFlag,
		Aliases: aliases[minAcceptableAmountFlag],
		Usage:   fmt.Sprintf("The minimum acceptable amount of %v wished to receive", tokenIDToBuyFlag),
		Value:   0,
	},
	tradingFeeFlag: &cli.Uint64Flag{
		Name:     tradingFeeFlag,
		Usage:    "The trading fee",
		Required: true,
	},
	tokenID1Flag: &cli.StringFlag{
		Name:     tokenID1Flag,
		Aliases:  aliases[tokenID1Flag],
		Usage:    "ID of the first token",
		Required: true,
	},
	tokenID2Flag: &cli.StringFlag{
		Name:    tokenID2Flag,
		Aliases: aliases[tokenID2Flag],
		Usage:   "ID of the second token",
		Value:   common.PRVIDStr,
	},
	prvFeeFlag: &cli.IntFlag{
		Name:  prvFeeFlag,
		Usage: "Whether or not to pay fee in PRV (0 - no, <> 0 - yes)",
		Value: 1,
	},
	tradingPathFlag: &cli.StringFlag{
		Name: tradingPathFlag,
		Usage: "A list of trading pair IDs seperated by a comma (Example: `pairID1,pairID2`). If none is given, the tool will automatically find " +
			"a suitable path.",
		Value: "",
	},
	maxTradingPathLengthFlag: &cli.UintFlag{
		Name:  maxTradingPathLengthFlag,
		Usage: "The maximum length of the trading path.",
		Value: pdex_v3.MaxPaths,
	},
	nftIDFlag: &cli.StringFlag{
		Name:     nftIDFlag,
		Aliases:  aliases[nftIDFlag],
		Usage:    "A pDEX NFT generated by the nft minting command",
		Required: true,
	},
	orderIDFlag: &cli.StringFlag{
		Name:     orderIDFlag,
		Aliases:  aliases[orderIDFlag],
		Usage:    "The ID of the order.",
		Required: true,
	},
	pairHashFlag: &cli.StringFlag{
		Name:     pairHashFlag,
		Usage:    "A unique string representing the contributing pair",
		Required: true,
	},
	pairIDFlag: &cli.StringFlag{
		Name:     pairIDFlag,
		Aliases:  aliases[pairIDFlag],
		Usage:    "The ID of the target pool pair",
		Required: true,
	},
	amplifierFlag: &cli.Uint64Flag{
		Name:     amplifierFlag,
		Aliases:  aliases[amplifierFlag],
		Usage:    "The amplifier for the target contributing pool",
		Required: true,
	},

	mnemonicFlag: &cli.StringFlag{
		Name:     mnemonicFlag,
		Aliases:  []string{"m"},
		Usage:    "A 12-word mnemonic phrase, words are separated by a \"-\", or put in \"\" (Examples: artist-decline-pepper-spend-good-enemy-caught-sister-sure-opinion-hundred-lake, \"artist decline pepper spend good enemy caught sister sure opinion hundred lake\").",
		Required: true,
	},
	numShardsFlag: &cli.IntFlag{
		Name:  numShardsFlag,
		Usage: "The number of shards",
		Value: 8,
	},
	numAccountsFlag: &cli.IntFlag{
		Name:  numAccountsFlag,
		Usage: "The number of accounts",
		Value: 1,
	},
	shardIDFlag: &cli.IntFlag{
		Name:  shardIDFlag,
		Usage: fmt.Sprintf("A specific shardID (-2: same shard as the first account (i.e, `Anon`); -1: any shard)"),
		Value: -2,
	},

	evmAddressFlag: &cli.StringFlag{
		Name:  evmAddressFlag,
		Usage: "A hex-encoded address on ETH/BSC networks",
		Value: "",
	},
	tokenAddressFlag: &cli.StringFlag{
		Name:    tokenAddressFlag,
		Aliases: aliases[tokenAddressFlag],
		Usage:   "ID of the token on ETH/BSC networks",
		Value:   nativeToken,
	},
	shieldAmountFlag: &cli.Float64Flag{
		Name:     shieldAmountFlag,
		Aliases:  aliases[shieldAmountFlag],
		Usage:    "The shielding amount measured in token unit (e.g, 10, 1, 0.1, 0.01)",
		Required: true,
	},
	evmFlag: &cli.StringFlag{
		Name:  evmFlag,
		Usage: "The EVM network (ETH, BSC, PLG or FTM)",
		Value: "ETH",
	},
	externalTxIDFlag: &cli.StringFlag{
		Name:     externalTxIDFlag,
		Aliases:  aliases[externalTxIDFlag],
		Usage:    "The external transaction hash",
		Required: true,
	},
	externalAddressFlag: &cli.StringFlag{
		Name:    externalAddressFlag,
		Aliases: aliases[externalAddressFlag],
		Usage:   "A valid remote address for the currently-processed tokenID. User MUST make sure this address is valid to avoid the loss of money.",
		Value:   "",
	},

	miningKeyFlag: &cli.StringFlag{
		Name:     miningKeyFlag,
		Aliases:  aliases[miningKeyFlag],
		Usage:    "An Incognito mining key of the committee candidate (default: the mining key associated with the privateKey)",
		Required: false,
	},
	candidateAddressFlag: &cli.StringFlag{
		Name:     candidateAddressFlag,
		Aliases:  aliases[candidateAddressFlag],
		Usage:    "The Incognito payment address of the committee candidate (default: the payment address of the privateKey)",
		Required: false,
	},
	rewardReceiverFlag: &cli.StringFlag{
		Name:     rewardReceiverFlag,
		Aliases:  aliases[rewardReceiverFlag],
		Usage:    "The Incognito payment address of the reward receiver (default: the payment address of the privateKey)",
		Required: false,
	},
	autoReStakeFlag: &cli.IntFlag{
		Name:     autoReStakeFlag,
		Aliases:  aliases[autoReStakeFlag],
		Usage:    "Whether or not to automatically re-stake (0 - false, <> 0 - true)",
		Value:    1,
		Required: false,
	},

	adminPrivateKeyFlag: &cli.StringFlag{
		Name:     adminPrivateKeyFlag,
		Aliases:  aliases[adminPrivateKeyFlag],
		Usage:    "A base58-encoded Incognito private key of the admin account",
		Required: true,
	},
	tokenNameFlag: &cli.StringFlag{
		Name:     tokenNameFlag,
		Aliases:  aliases[tokenNameFlag],
		Usage:    "The name of the shielding token",
		Required: true,
	},
}
