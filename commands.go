package main

import (
	"fmt"

	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/urfave/cli/v2"
)

var dexStatusErrMsg = "If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid."

// accountCommands consists of all account-related commands
var accountCommands = []*cli.Command{
	{
		Name:        "account",
		Aliases:     []string{"acc"},
		Usage:       "Manage an Incognito account.",
		Description: fmt.Sprintf("This command helps perform an account-related action."),
		Category:    accountCat,
		Subcommands: []*cli.Command{
			{
				Name:  "keyinfo",
				Usage: "Print all related-keys of a private key.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     privateKeyFlag,
						Aliases:  aliases[privateKeyFlag],
						Usage:    "a base58-encoded private key",
						Required: true,
					},
				},
				Action: keyInfo,
			},
			{
				Name:  "balance",
				Usage: "Check the balance of an account for a tokenID.",
				Flags: []cli.Flag{
					defaultFlags[privateKeyFlag],
					defaultFlags[tokenIDFlag],
				},
				Action: checkBalance,
			},
			{
				Name: "balanceall",
				Usage: "Check all non-zero balances (calculated based on v2 UTXOs only) of a private key. In case you have v1 UTXOs left, try using " +
					"regular `balance` command with each token for the best result.",
				Flags: []cli.Flag{
					defaultFlags[privateKeyFlag],
				},
				Action: getAllBalanceV2,
			},
			{
				Name:  "outcoin",
				Usage: "Print the output coins of an account.",
				Flags: []cli.Flag{
					defaultFlags[addressFlag],
					defaultFlags[otaKeyFlag],
					defaultFlags[readonlyKeyFlag],
					defaultFlags[tokenIDFlag],
				},
				Action: getOutCoins,
			},
			{
				Name:  "utxo",
				Usage: "Print the UTXOs of an account.",
				Flags: []cli.Flag{
					defaultFlags[privateKeyFlag],
					defaultFlags[tokenIDFlag],
				},
				Action: checkUTXOs,
			},
			{
				Name:    "consolidate",
				Aliases: []string{"csl"},
				Usage:   "Consolidate UTXOs of an account.",
				Description: "This command helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. " +
					"Please note that this process is time-consuming and requires a considerable amount of CPU.",
				Flags: []cli.Flag{
					defaultFlags[privateKeyFlag],
					defaultFlags[tokenIDFlag],
					defaultFlags[versionFlag],
					defaultFlags[numThreadsFlag],
				},
				Action: consolidateUTXOs,
			},
			{
				Name:    "history",
				Aliases: []string{"hst"},
				Usage:   "Retrieve the history of an account.",
				Description: "This command helps retrieve the history of an account w.r.t a tokenID. " +
					"Please note that this process is time-consuming and requires a considerable amount of CPU.",
				Flags: []cli.Flag{
					defaultFlags[privateKeyFlag],
					&cli.StringFlag{
						Name:  "tokenID",
						Usage: "ID of the token",
						Value: common.PRVIDStr,
					},
					defaultFlags[numThreadsFlag],
					defaultFlags[csvFileFlag],
				},
				Action: getHistory,
			},
			{
				Name:        "generate",
				Aliases:     []string{"gen"},
				Usage:       "Generate a new Incognito account.",
				Description: "This command helps generate a new mnemonic phrase and its Incognito accounts.",
				Flags: []cli.Flag{
					defaultFlags[numShardsFlag],
				},
				Action: genKeySet,
			},
			{
				Name:        "importaccount",
				Aliases:     []string{"import"},
				Usage:       "Import a mnemonic of 12 words.",
				Description: "This command helps generate Incognito accounts given a mnemonic.",
				Flags: []cli.Flag{
					defaultFlags[mnemonicFlag],
					defaultFlags[numShardsFlag],
				},
				Action: importMnemonic,
			},
			{
				Name:    "submitkey",
				Aliases: []string{"sub"},
				Usage:   "Submit an ota key to the full-node.",
				Description: "This command submits an otaKey to the full-node to use the full-node's cache. If an access token " +
					"is provided, it will submit the ota key in an authorized manner. See " +
					"https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md " +
					"for more details.",
				Flags: []cli.Flag{
					defaultFlags[otaKeyFlag],
					defaultFlags[accessTokenFlag],
					defaultFlags[fromHeightFlag],
					defaultFlags[isResetFlag],
				},
				Action: submitKey,
			},
		},
	},
}

// committeeCommands consists of all committee-related commands
var committeeCommands = []*cli.Command{
	{
		Name:     "stake",
		Usage:    "Create a staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/stake.md).",
		Category: committeeCat,
		Flags: []cli.Flag{
			defaultFlags[privateKeyFlag],
			defaultFlags[miningKeyFlag],
			defaultFlags[candidateAddressFlag],
			defaultFlags[rewardReceiverFlag],
			defaultFlags[autoReStakeFlag],
		},
		Action: stake,
	},
	{
		Name:     "unstake",
		Usage:    "Create an un-staking transaction (https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/staking/unstake.md).",
		Category: committeeCat,
		Flags: []cli.Flag{
			defaultFlags[privateKeyFlag],
			defaultFlags[miningKeyFlag],
			defaultFlags[candidateAddressFlag],
		},
		Action: unStake,
	},
	{
		Name:     "checkrewards",
		Usage:    "Get all rewards of a payment address.",
		Category: committeeCat,
		Flags: []cli.Flag{
			defaultFlags[addressFlag],
		},
		Action: checkRewards,
	},
	{
		Name:     "withdrawreward",
		Usage:    "Withdraw the reward of a privateKey w.r.t to a tokenID.",
		Category: committeeCat,
		Flags: []cli.Flag{
			defaultFlags[privateKeyFlag],
			&cli.StringFlag{
				Name:    addressFlag,
				Aliases: aliases[addressFlag],
				Usage:   "the payment address of a candidate (default: the payment address of the privateKey)",
			},
			defaultFlags[tokenIDFlag],
			defaultFlags[versionFlag],
		},
		Action: withdrawReward,
	},
}

// txCommands consists of all (normal) tx-related commands
var txCommands = []*cli.Command{
	{
		Name:  "send",
		Usage: "Send an amount of PRV or token from one wallet to another wallet.",
		Description: "This command sends an amount of PRV or token from one wallet to another wallet. By default, " +
			"it used 100 nano PRVs to pay the transaction fee.",
		Category: transactionCat,
		Flags: []cli.Flag{
			defaultFlags[privateKeyFlag],
			defaultFlags[addressFlag],
			defaultFlags[amountFlag],
			defaultFlags[tokenIDFlag],
			defaultFlags[feeFlag],
			defaultFlags[versionFlag],
		},
		Action: send,
	},
	{
		Name:  "convert",
		Usage: "Convert UTXOs of an account w.r.t a tokenID.",
		Description: "This command helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. " +
			"Please note that this process is time-consuming and requires a considerable amount of CPU.",
		Category: transactionCat,
		Flags: []cli.Flag{
			defaultFlags[privateKeyFlag],
			defaultFlags[tokenIDFlag],
			defaultFlags[numThreadsFlag],
		},
		Action: convertUTXOs,
	},
	//{
	//	Name:  "convertall",
	//	Usage: "Convert UTXOs of an account for all assets.",
	//	Description: "This command helps convert UTXOs v1 of a user to UTXO v2 for all assets. " +
	//		"It will automatically check for all UTXOs v1 of all tokens and convert them. " +
	//		"Please note that this process is time-consuming and requires a considerable amount of CPU.",
	//	Category: transactionCat,
	//	Flags: []cli.Flag{
	//		defaultFlags[privateKeyFlag],
	//		defaultFlags[numThreadsFlag],
	//	},
	//	Action: convertAll,
	//},
	{
		Name:  "checkreceiver",
		Usage: "Check if an OTA key is a receiver of a transaction.",
		Description: "This command checks if an OTA key is a receiver of a transaction. If so, it will try to decrypt " +
			"the received outputs and return the receiving info.",
		Category: transactionCat,
		Flags: []cli.Flag{
			defaultFlags[txHashFlag],
			defaultFlags[otaKeyFlag],
			defaultFlags[readonlyKeyFlag],
		},
		Action: checkReceiver,
	},
}

// bridgeCommands consists of all bridge-related commands
var bridgeCommands = []*cli.Command{
	evmBridgeCommands, portalCommands,
}

// evmBridgeCommands consists of all EVM-related commands
var evmBridgeCommands = &cli.Command{
	Name:        "evm",
	Usage:       "Perform an EVM action (e.g, shield, unshield, etc.).",
	Description: fmt.Sprintf("This command helps perform an EVM action (e.g, shield, unshield, etc.)."),
	Category:    evmBridgeCat,
	Subcommands: []*cli.Command{
		{
			Name:        "shield",
			Usage:       "Shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network.",
			Description: shieldMessage,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[shieldAmountFlag],
				defaultFlags[evmFlag],
				defaultFlags[tokenAddressFlag],
				&cli.StringFlag{
					Name:    addressFlag,
					Aliases: aliases[addressFlag],
					Usage:   "The Incognito payment address to receive the shielding asset (default: the payment address of the privateKey)",
				},
			},
			Action: shield,
		},
		{
			Name:        "retryshield",
			Usage:       "Retry a shield from the given already-been-deposited-to-sc EVM transaction.",
			Description: "This command re-shields an already-been-deposited-to-sc transaction in case of prior failure.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[externalTxIDFlag],
				defaultFlags[evmFlag],
				defaultFlags[tokenAddressFlag],
			},
			Action: retryShield,
		},
		{
			Name:        "unshield",
			Usage:       "Withdraw an EVM (ETH/BNB/ERC20/BEP20) token from the Incognito network.",
			Description: unShieldMessage,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:     tokenIDFlag,
					Aliases:  aliases[tokenIDFlag],
					Usage:    "The Incognito tokenID of the un-shielding asset",
					Required: true,
				},
				defaultFlags[amountFlag],
			},
			Action: unShield,
		},
		{
			Name:        "retryunshield",
			Usage:       "Retry an un-shielding request from the given already-been-burned Incognito transaction.",
			Description: "This command tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
				defaultFlags[evmFlag],
			},
			Action: retryUnShield,
		},
	},
}

// portalCommands consists of all Portal-related commands
var portalCommands = &cli.Command{
	Name:        "portal",
	Usage:       "Perform a portal action (e.g, shield, unshield, etc.).",
	Description: fmt.Sprintf("This command helps perform a portal action (e.g, shield, unshield, etc.)."),
	Category:    portalCat,
	Subcommands: []*cli.Command{
		{
			Name:        "shieldaddress",
			Usage:       "Generate a portal shielding address.",
			Description: "This command helps generate the portal shielding address for a payment address and a tokenID.",
			Flags: []cli.Flag{
				defaultFlags[addressFlag],
				&cli.StringFlag{
					Name:    tokenIDFlag,
					Aliases: aliases[tokenIDFlag],
					Usage:   "The Incognito tokenID of the shielding asset",
					Value:   "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696",
				},
			},
			Action: getPortalDepositAddress,
		},
		{
			Name:  "shield",
			Usage: "Shield a portal token (e.g, BTC) into the Incognito network.",
			Description: "This command helps shield a portal token into the Incognito network after the fund has been " +
				"transferred to the depositing address (generated by `portalshieldaddress`).",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[externalTxIDFlag],
				&cli.StringFlag{
					Name:    tokenIDFlag,
					Aliases: aliases[tokenIDFlag],
					Usage:   "The Incognito tokenID of the shielding asset",
					Value:   "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696",
				},
				&cli.StringFlag{
					Name:    addressFlag,
					Aliases: aliases[addressFlag],
					Usage:   "The Incognito payment address to receive the shielding asset (default: the payment address of the privateKey)",
				},
			},
			Action: portalShield,
		},
		{
			Name:  "shieldstatus",
			Usage: "Get the status of a portal shielding request.",
			Description: "This command helps retrieve the status of a portal shielding request.\n" +
				"Status should be understood as: " +
				"0 - rejected; 1 - accepted.\n" +
				"If you encounter an error, it might be because the request hasn't reached the " +
				"beacon chain yet. Please try again a few minutes later.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: getPortalShieldStatus,
		},
		{
			Name:        "unshield",
			Usage:       "Withdraw portal tokens (BTC) from the Incognito network.",
			Description: "This command helps withdraw portal tokens (BTC) out of the Incognito network.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:     externalAddressFlag,
					Aliases:  aliases[externalAddressFlag],
					Usage:    "A valid remote address for the currently-processed tokenID. User MUST make sure this address is valid to avoid the loss of money.",
					Required: true,
				},
				defaultFlags[amountFlag],
				&cli.StringFlag{
					Name:    tokenIDFlag,
					Aliases: aliases[tokenIDFlag],
					Usage:   "The Incognito tokenID of the un-shielding asset",
					Value:   "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696",
				},
			},
			Action: portalUnShield,
		},
		{
			Name:  "unshieldstatus",
			Usage: "Get the status of a portal un-shielding request.",
			Description: "This command helps retrieve the status of a portal un-shielding request.\n" +
				"Status should be understood as: " +
				"0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.\n" +
				"If you encounter an error saying \"unexpected end of JSON input\", it might be because the request hasn't reached the " +
				"beacon chain yet. Please try again a few minutes later.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: getPortalUnShieldStatus,
		},
	},
}

// pDEXCommands consists of all pDEX-related commands
var pDEXCommands = []*cli.Command{
	pDEXActionCommands,
	pDEXInfoCommands,
	pDEXStatusCommands,
}

// pDEXActionCommands consists of all pDEX-related commands for a pDEX action (e.g, trade, add orders, etc.).
var pDEXActionCommands = &cli.Command{
	Name:        "pdeaction",
	Usage:       "Perform a pDEX action.",
	Description: fmt.Sprintf("This command helps perform a pDEX action. Most of the terms here are based on the SDK tutorial series (https://github.com/incognitochain/go-incognito-sdk-v2/blob/dev/pdex-v3/tutorials/docs/pdex/intro.md)."),
	Category:    pDEXCat,
	Subcommands: []*cli.Command{
		{
			Name:        "modifyparams",
			Usage:       "Create a modify params transaction.",
			Description: "This command creates a modify params transaction on the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:     paramsFlag,
					Usage:    "The parameters for the new desired pDEX params",
					Required: true,
				},
			},
			Action: pDEXModifyParams,
		},
		{
			Name:        "trade",
			Usage:       "Create a trade transaction.",
			Description: "This command creates a trade transaction on the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
				defaultFlags[tradingFeeFlag],
				defaultFlags[minAcceptableAmountFlag],
				defaultFlags[tradingPathFlag],
				defaultFlags[prvFeeFlag],
				defaultFlags[maxTradingPathLengthFlag],
			},
			Action: pDEXTrade,
		},
		{
			Name:        "mintnft",
			Usage:       "Create a (pDEX) NFT minting transaction.",
			Description: "This command creates and broadcasts a transaction that mints a new (pDEX) NFT for the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
			},
			Action: pDEXMintNFT,
		},
		{
			Name:        "contribute",
			Usage:       "Create a pDEX liquidity-contributing transaction.",
			Description: "This command creates a pDEX liquidity-contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[pairHashFlag],
				defaultFlags[amountFlag],
				defaultFlags[amplifierFlag],
				defaultFlags[tokenIDFlag],
				&cli.StringFlag{
					Name:  pairIDFlag,
					Usage: "The ID of the contributing pool pair. For pool-initializing transactions (e.g, first contribution in the pool), it should be left empty.",
					Value: "",
				},
			},
			Action: pDEXContribute,
		},
		{
			Name:        "withdraw",
			Usage:       "Create a pDEX liquidity-withdrawal transaction.",
			Description: "This command creates a transaction withdrawing an amount of `share` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:     pairIDFlag,
					Usage:    "The ID of the contributed pool pair",
					Required: true,
				},
				defaultFlags[nftIDFlag],
				&cli.Uint64Flag{
					Name:    amountFlag,
					Aliases: aliases[amountFlag],
					Usage:   "The amount of share wished to withdraw. If set to 0, it will withdraw all of the share.",
				},
			},
			Action: pDEXWithdraw,
		},
		{
			Name:        "addorder",
			Usage:       "Add an order book to the pDEX.",
			Description: "This command creates a transaction adding an order to the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[sellingAmountFlag],
				&cli.Uint64Flag{
					Name:     minAcceptableAmountFlag,
					Aliases:  aliases[minAcceptableAmountFlag],
					Usage:    fmt.Sprintf("The minimum acceptable amount of %v wished to receive", tokenIDToBuyFlag),
					Required: true,
				},
			},
			Action: pDEXAddOrder,
		},
		{
			Name:        "withdraworder",
			Usage:       "Withdraw an order from the pDEX.",
			Description: "This command creates a transaction withdrawing an order to the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[orderIDFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[tokenID1Flag],
				&cli.StringFlag{
					Name:    tokenID2Flag,
					Aliases: aliases[tokenID2Flag],
					Usage:   "ID of the second token (if have). In the case of withdrawing a single token, leave it empty",
					Value:   "",
				},
				&cli.Uint64Flag{
					Name:    amountFlag,
					Aliases: aliases[amountFlag],
					Usage:   "Amount to withdraw (0 for all)",
					Value:   0,
				},
			},
			Action: pDEXWithdrawOrder,
		},
		{
			Name:        "stake",
			Usage:       "Stake a token to the pDEX.",
			Description: "This command creates a transaction staking a token to the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[amountFlag],
				&cli.StringFlag{
					Name:  tokenIDFlag,
					Usage: "The ID of the target staking pool ID (or token ID)",
					Value: common.PRVIDStr,
				},
			},
			Action: pDEXStake,
		},
		{
			Name:        "unstake",
			Usage:       "Un-stake a token from the pDEX.",
			Description: "This command creates a transaction un-staking a token from the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[amountFlag],
				&cli.StringFlag{
					Name:  tokenIDFlag,
					Usage: "The ID of the target staking pool ID (or token ID)",
					Value: common.PRVIDStr,
				},
			},
			Action: pDEXUnStake,
		},
		{
			Name:        "withdrawstakereward",
			Usage:       "Withdraw staking rewards from the pDEX.",
			Description: "This command creates a transaction withdrawing staking rewards from the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[nftIDFlag],
				&cli.StringFlag{
					Name:  tokenIDFlag,
					Usage: "The ID of the target staking pool ID (or token ID)",
					Value: common.PRVIDStr,
				},
			},
			Action: pDEXWithdrawStakingReward,
		},
		{
			Name:        "withdrawlpfee",
			Usage:       "Withdraw LP fees from the pDEX.",
			Description: "This command creates a transaction withdrawing LP fees from the pDEX.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
			},
			Action: pDEXWithdrawLPFee,
		},
	},
}

// pDEXInfoCommands consists of all pDEX-related commands for retrieving information (sub-commands of pDEXCommands).
var pDEXInfoCommands = &cli.Command{
	Name:        "pdeinfo",
	Usage:       "Retrieve pDEX information.",
	Description: fmt.Sprintf("This command helps retrieve some information of the pDEX. Most of the terms here are based on the SDK tutorial series (https://github.com/incognitochain/go-incognito-sdk-v2/blob/dev/pdex-v3/tutorials/docs/pdex/intro.md)."),
	Category:    pDEXCat,
	Subcommands: []*cli.Command{
		{
			Name:        "mynft",
			Usage:       "Retrieve the list of NFTs for a given private key.",
			Description: "This command returns the list of NFTs for a given private key.",
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
			},
			Action: pDEXGetAllNFTs,
		},
		{
			Name:        "getorder",
			Usage:       "Retrieve the detail of an order given its id.",
			Description: "This command returns the detail of an order given its id.",
			Flags: []cli.Flag{
				defaultFlags[orderIDFlag],
			},
			Action: pDEXGetOrderByID,
		},
		{
			Name:        "share",
			Usage:       "Retrieve the share amount of a pDEX poolID given an nftID.",
			Description: "This command returns the share amount of an nftID within a pDEX poolID.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     pairIDFlag,
					Usage:    "The ID of the target pool pair",
					Required: true,
				},
				defaultFlags[nftIDFlag],
			},
			Action: pDEXGetShare,
		},
		{
			Name:        "stakereward",
			Usage:       "Retrieve the estimated pDEX staking rewards.",
			Description: "This command returns the estimated pDEX staking rewards of an nftID within a pDEX staking pool.",
			Flags: []cli.Flag{
				defaultFlags[nftIDFlag],
				&cli.StringFlag{
					Name:  tokenIDFlag,
					Usage: "The ID of the target staking pool ID (or token ID)",
					Value: common.PRVIDStr,
				},
			},
			Action: CheckDEXStakingReward,
		},
		{
			Name:        "lpvalue",
			Usage:       "Check the estimated LP value in a given pool.",
			Description: "This command retrieves the information about the value of an LP in a given pool.",
			Flags: []cli.Flag{
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
			},
			Action: pDEXGetEstimatedLPValue,
		},
		{
			Name:  "checkprice",
			Usage: "Check the price between two tokenIDs.",
			Description: "This command checks the price of a pair of tokenIds. It must be supplied with the selling amount " +
				"since the pDEX uses the AMM algorithm.",
			Flags: []cli.Flag{
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
				&cli.StringFlag{
					Name:     pairIDFlag,
					Usage:    "The ID of the target pool pair",
					Required: true,
				},
			},
			Action: pDEXCheckPrice,
		},
		{
			Name:        "findpath",
			Usage:       "Find a `good` trading path for a trade.",
			Description: "This command helps find a good trading path for a trade.",
			Flags: []cli.Flag{
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
				defaultFlags[maxTradingPathLengthFlag],
			},
			Action: pDEXFindPath,
		},
	},
}

// pDEXStatusCommands consists of all pDEX-related commands for getting statuses (sub-commands of pDEXCommands).
var pDEXStatusCommands = &cli.Command{
	Name:        "pdestatus",
	Usage:       "Retrieve the status of a pDEX action.",
	Description: fmt.Sprintf("This command helps retrieve the status of a pDEX action given its hash. %v", dexStatusErrMsg),
	Category:    pDEXCat,
	Subcommands: []*cli.Command{
		{
			Name:  "trade",
			Usage: "Check the status of a pDEX trade.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXTradeStatus,
		},
		{
			Name:  "mintnft",
			Usage: "Check the status of a (pDEX) NFT minting transaction.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXMintNFTStatus,
		},
		{
			Name:  "contribute",
			Usage: "Check the status of a pDEX liquidity contribution.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXContributionStatus,
		},
		{
			Name:  "withdraw",
			Usage: "Check the status of a pDEX liquidity withdrawal.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXWithdrawalStatus,
		},
		{
			Name:  "addorder",
			Usage: "Check the status of a pDEX order-adding withdrawal.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXOrderAddingStatus,
		},
		{
			Name:  "withdraworder",
			Usage: "Check the status of a pDEX order-withdrawal transaction.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXOrderWithdrawalStatus,
		},
		{
			Name:  "stake",
			Usage: "Check the status of a pDEX staking transaction.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXStakingStatus,
		},
		{
			Name:  "unstake",
			Usage: "Check the status of a pDEX un-staking transaction.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXUnStakingStatus,
		},
		{
			Name:  "withdrawstakereward",
			Usage: "Check the status of a pDEX staking reward withdrawal transaction.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXWithdrawStakingRewardStatus,
		},
		{
			Name:  "withdrawlpfee",
			Usage: "Check the status of a pDEX LP fee withdrawal transaction.",
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXWithdrawLPFeeStatus,
		},
	},
}
