package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

var errMsg = "If an error is thrown, it is mainly because the transaction has not yet reached the beacon chain or the txHash is invalid."

func main() {
	app := &cli.App{
		Name:    "incognito-cli",
		Usage:   "A simple CLI application for the Incognito network",
		Version: "v0.0.3",
		Description: "A simple CLI application for the Incognito network. With this tool, you can run some basic functions" +
			" on your computer to interact with the Incognito network such as checking balances, transferring PRV or tokens," +
			" consolidating and converting your UTXOs, transferring tokens, manipulating with the pDEX, shielding or un-shielding " +
			"ETH/BNB/ERC20/BEP20, etc.",
		Authors: []*cli.Author{
			{
				Name: "Incognito Devs Team",
			},
		},
		Copyright: "This tool is developed and maintained by the Incognito Devs Team. It is free for anyone. However, any " +
			"commercial usages should be acknowledged by the Incognito Devs Team.",
	}
	app.EnableBashCompletion = true

	// set app defaultFlags
	app.Flags = []cli.Flag{
		defaultFlags[networkFlag],
		defaultFlags[hostFlag],
		defaultFlags[debugFlag],
		defaultFlags[cacheFlag],
	}

	// all account-related commands
	accountCommands := []*cli.Command{
		{
			Name:     "keyinfo",
			Usage:    "Print all related-keys of a private key.",
			Category: accountCat,
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
			Name:     "balance",
			Usage:    "Check the balance of an account for a tokenID.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDFlag],
			},
			Action: checkBalance,
		},
		//{
		//	Name:  "balanceall",
		//	Usage: "Return the non-zero balances of an account for all tokenIDs.",
		//	Description: "This command returns the non-zero balances of an account for all tokenIDs. Due to the large number of " +
		//		"tokens on the network, This command requires a long amount of time to proceed.",
		//	Category: accountCat,
		//	Flags: []cli.Flag{
		//		defaultFlags[privateKeyFlag],
		//	},
		//	Action: checkBalanceAll,
		//},
		{
			Name:     "outcoin",
			Usage:    "Print the output coins of an account.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[addressFlag],
				defaultFlags[otaKeyFlag],
				defaultFlags[readonlyKeyFlag],
				defaultFlags[tokenIDFlag],
			},
			Action: getOutCoins,
		},
		{
			Name:     "utxo",
			Usage:    "Print the UTXOs of an account.",
			Category: accountCat,
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
			Category: accountCat,
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
			Category: accountCat,
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
			Name:        "generateaccount",
			Aliases:     []string{"genacc"},
			Usage:       "Generate a new Incognito account.",
			Description: "This command helps generate a new mnemonic phrase and its Incognito accounts.",
			Category:    accountCat,
			Flags: []cli.Flag{
				defaultFlags[numShardsFlag],
			},
			Action: genKeySet,
		},
		{
			Name:        "importeaccount",
			Aliases:     []string{"import"},
			Usage:       "Import a mnemonic of 12 words.",
			Description: "This command helps generate Incognito accounts given a mnemonic.",
			Category:    accountCat,
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
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[otaKeyFlag],
				defaultFlags[accessTokenFlag],
				defaultFlags[fromHeightFlag],
				defaultFlags[isResetFlag],
			},
			Action: submitKey,
		},
	}

	// all committee-related commands
	committeeCommands := []*cli.Command{
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

	// all tx-related commands
	txCommands := []*cli.Command{
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
		{
			Name:  "convertall",
			Usage: "Convert UTXOs of an account for all assets.",
			Description: "This command helps convert UTXOs v1 of a user to UTXO v2 for all assets. " +
				"It will automatically check for all UTXOs v1 of all tokens and convert them. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: transactionCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[numThreadsFlag],
			},
			Action: convertAll,
		},
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

	// pDEX commands
	pDEXCommands := []*cli.Command{
		{
			Name:  "pdecheckprice",
			Usage: "Check the price between two tokenIDs.",
			Description: "This command checks the price of a pair of tokenIds. It must be supplied with the selling amount " +
				"since the pDEX uses the AMM algorithm.",
			Category: pDEXCat,
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
			Name:        "pdefindpath",
			Usage:       "Find a `good` trading path for a trade.",
			Description: "This command helps find a good trading path for a trade.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
				defaultFlags[maxTradingPathLengthFlag],
			},
			Action: pDEXFindPath,
		},
		{
			Name:        "pdetrade",
			Usage:       "Create a trade transaction.",
			Description: "This command creates a trade transaction on the pDEX.",
			Category:    pDEXCat,
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
			Name:  "pdetradestatus",
			Usage: "Check the status of a pDEX trade.",
			Description: "This command retrieves the status of a pDEX trade. If an error is thrown, it is mainly " +
				"because the transaction has not yet reached the beacon chain or the txHash is invalid.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXTradeStatus,
		},
		{
			Name:        "pdemintnft",
			Usage:       "Create a (pDEX) NFT minting transaction.",
			Description: "This command creates and broadcasts a transaction that mints a new (pDEX) NFT for the pDEX.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
			},
			Action: pDEXMintNFT,
		},
		{
			Name:        "pdemintnftstatus",
			Usage:       "Check the status of a (pDEX) NFT minting transaction.",
			Description: "This command retrieves the status of a (pDEX) NFT minting transaction.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXMintNFTStatus,
		},
		{
			Name:        "pdecontribute",
			Usage:       "Create a pDEX liquidity-contributing transaction.",
			Description: "This command creates a pDEX liquidity-contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md",
			Category:    pDEXCat,
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
			Name:  "pdecontributionstatus",
			Usage: "Check the status of a pDEX liquidity contribution.",
			Description: "This command retrieves the status of a pDEX liquidity contribution. If an error is thrown, it is mainly " +
				"because the transaction has not yet reached the beacon chain or the txHash is invalid.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXContributionStatus,
		},
		{
			Name:        "pdewithdraw",
			Usage:       "Create a pDEX liquidity-withdrawal transaction.",
			Description: "This command creates a transaction withdrawing an amount of `share` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:     pairIDFlag,
					Usage:    "The ID of the contributed pool pair",
					Required: true,
				},
				defaultFlags[nftIDFlag],
				defaultFlags[tokenID1Flag],
				defaultFlags[tokenID2Flag],
				&cli.Uint64Flag{
					Name:    amountFlag,
					Aliases: aliases[amountFlag],
					Usage:   "The amount of share wished to withdraw. If set to 0, it will withdraw all of the share.",
				},
			},
			Action: pDEXWithdraw,
		},
		{
			Name:  "pdewithdrawalstatus",
			Usage: "Check the status of a pDEX liquidity withdrawal.",
			Description: "This command retrieves the status of a pDEX liquidity withdrawal. If an error is thrown, it is mainly " +
				"because the transaction has not yet reached the beacon chain or the txHash is invalid.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXWithdrawalStatus,
		},
		{
			Name:        "pdeaddorder",
			Usage:       "Add an order book to the pDEX.",
			Description: "This command creates a transaction adding an order to the pDEX.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
				defaultFlags[minAcceptableAmountFlag],
			},
			Action: pDEXAddOrder,
		},
		{
			Name:  "pdeorderaddstatus",
			Usage: "Check the status of a pDEX order-adding withdrawal.",
			Description: "This command retrieves the status of a pDEX order-adding withdrawal. If an error is thrown, it is mainly " +
				"because the transaction has not yet reached the beacon chain or the txHash is invalid.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXOrderAddingStatus,
		},
		{
			Name:        "pdewithdraworder",
			Usage:       "Withdraw an order from the pDEX.",
			Description: "This command creates a transaction withdraing an order to the pDEX.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[orderIDFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
				defaultFlags[amountFlag],
				defaultFlags[tokenID1Flag],
				&cli.StringFlag{
					Name:    tokenID2Flag,
					Aliases: aliases[tokenID2Flag],
					Usage:   "ID of the second token (if have). In the case of withdrawing a single token, leave it empty.",
					Value:   "",
				},
			},
			Action: pDEXWithdrawOrder,
		},
		{
			Name:  "pdeorderwithdrawalstatus",
			Usage: "Check the status of a pDEX order-withdrawal transaction.",
			Description: "This command retrieves the status of a pDEX order-withdrawal transaction. If an error is thrown, it is mainly " +
				"because the transaction has not yet reached the beacon chain or the txHash is invalid.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXOrderWithdrawalStatus,
		},
		{
			Name:        "pdestake",
			Usage:       "Stake a token to the pDEX.",
			Description: "This command creates a transaction staking a token to the pDEX.",
			Category:    pDEXCat,
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
			Name:        "pdeunstake",
			Usage:       "Un-stake a token from the pDEX.",
			Description: "This command creates a transaction un-staking a token from the pDEX.",
			Category:    pDEXCat,
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
			Name:        "pdewithdrawstakereward",
			Usage:       "Withdraw staking rewards from the pDEX.",
			Description: "This command creates a transaction withdrawing staking rewards from the pDEX.",
			Category:    pDEXCat,
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
			Name:        "pdelpvalue",
			Usage:       "Check the estimated LP value in a given pool.",
			Description: "This command retrieves the information about the value of an LP in a given pool.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
			},
			Action: pDEXGetEstimatedLPValue,
		},
		{
			Name:        "pdewithdrawlpfee",
			Usage:       "Withdraw LP fees from the pDEX.",
			Description: "This command creates a transaction withdrawing LP fees from the pDEX.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[nftIDFlag],
			},
			Action: pDEXWithdrawLPFee,
		},
		{
			Name:        "pdeshare",
			Usage:       "Retrieve the share amount of a pDEX poolID given an nftID.",
			Description: "This command returns the share amount of an nftID within a pDEX poolID.",
			Category:    pDEXCat,
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
			Name:        "pdemynft",
			Usage:       "Retrieve the list of NFTs for a given private key.",
			Description: "This command returns the list of NFTs for a given private key.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
			},
			Action: pDEXGetAllNFTs,
		},
		{
			Name:        "pdestakereward",
			Usage:       "Retrieve the estimated pDEX staking rewards.",
			Description: "This command returns the estimated pDEX staking rewards of an nftID within a pDEX staking pool.",
			Category:    pDEXCat,
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
	}
	pDEXStatusCommands := &cli.Command{
		Name:        "pdestatus",
		Usage:       "Retrieve the status of a pDEX action.",
		Description: fmt.Sprintf("This command helps retrieve the status of a pDEX action given its hash. %v", errMsg),
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
	pDEXCommands = append(pDEXCommands, pDEXStatusCommands)

	// Bridge commands
	evmBridgeCommands := []*cli.Command{
		{
			Name:        "evmshield",
			Usage:       "Shield an EVM (ETH/BNB/ERC20/BEP20) token into the Incognito network.",
			Description: shieldMessage,
			Category:    evmBridgeCat,
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
			Name:        "evmretryshield",
			Usage:       "Retry a shield from the given already-been-deposited-to-sc EVM transaction.",
			Description: "This command re-shields an already-been-deposited-to-sc transaction in case of prior failure.",
			Category:    evmBridgeCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[externalTxIDFlag],
				defaultFlags[evmFlag],
				defaultFlags[tokenAddressFlag],
			},
			Action: retryShield,
		},
		{
			Name:        "evmunshield",
			Usage:       "Withdraw an EVM (ETH/BNB/ERC20/BEP20) token from the Incognito network.",
			Description: unShieldMessage,
			Category:    evmBridgeCat,
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
			Name:        "evmretryunshield",
			Usage:       "Retry an un-shielding request from the given already-been-burned Incognito transaction.",
			Description: "This command tries to un-shield an asset from an already-been-burned Incognito transaction in case of prior failure.",
			Category:    evmBridgeCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
				defaultFlags[evmFlag],
			},
			Action: retryUnShield,
		},
	}

	// portal commands
	portalCommands := []*cli.Command{
		{
			Name:        "portalshieldaddress",
			Usage:       "Generate a portal shielding address.",
			Description: "This command helps generate the portal shielding address for a payment address and a tokenID.",
			Category:    portalCat,
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
			Name:  "portalshield",
			Usage: "Shield a portal token (e.g, BTC) into the Incognito network.",
			Description: "This command helps shield a portal token into the Incognito network after the fund has been " +
				"transferred to the depositing address (generated by `portalshieldaddress`).",
			Category: portalCat,
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
			Name:  "portalshieldstatus",
			Usage: "Get the status of a portal shielding request.",
			Description: "This command helps retrieve the status of a portal shielding request.\n" +
				"Status should be understood as: " +
				"0 - rejected; 1 - accepted.\n" +
				"If you encounter an error, it might be because the request hasn't reached the " +
				"beacon chain yet. Please try again a few minutes later.",
			Category: portalCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: getPortalShieldStatus,
		},
		{
			Name:        "portalunshield",
			Usage:       "Withdraw portal tokens (BTC) from the Incognito network.",
			Description: "This command helps withdraw portal tokens (BTC) out of the Incognito network.",
			Category:    portalCat,
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
			Name:  "portalunshieldstatus",
			Usage: "Get the status of a portal un-shielding request.",
			Description: "This command helps retrieve the status of a portal un-shielding request.\n" +
				"Status should be understood as: " +
				"0 - waiting; 1 - processed but not completed; 2 - completed; 3 - rejected.\n" +
				"If you encounter an error saying \"unexpected end of JSON input\", it might be because the request hasn't reached the " +
				"beacon chain yet. Please try again a few minutes later.",
			Category: portalCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: getPortalUnShieldStatus,
		},
	}

	app.Commands = make([]*cli.Command, 0)
	app.Commands = append(app.Commands, accountCommands...)
	app.Commands = append(app.Commands, committeeCommands...)
	app.Commands = append(app.Commands, txCommands...)
	app.Commands = append(app.Commands, pDEXCommands...)
	app.Commands = append(app.Commands, evmBridgeCommands...)
	app.Commands = append(app.Commands, portalCommands...)

	for _, command := range app.Commands {
		if len(command.Subcommands) > 0 {
			for _, subCommand := range command.Subcommands {
				buildUsageTextFromCommand(subCommand, command.Name)
			}
		}
		buildUsageTextFromCommand(command)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	//_ = generateDocsToFile(app, "commands.md") // un-comment this line to generate docs for the app's commands.

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
