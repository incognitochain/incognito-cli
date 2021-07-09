package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "incognito-cli",
		Usage:   "A simple CLI application for the Incognito network",
		Version: "v0.0.1",
		Description: "A simple CLI application for the Incognito network. With this tool, you can run some basic functions" +
			" on your computer to interact with the Incognito network such as checking balances, transferring PRV or tokens," +
			" consolidating and converting your UTXOs, etc.",
		Authors: []*cli.Author{
			{
				Name:  "Incognito Devs Team",
				Email: "support@incognito.org",
			},
		},
		Copyright: "This tool is created and managed by the Incognito Devs Team. It is free for anyone. However, any " +
			"commercial usages should be acknowledged by the Incognito Devs Team.",
	}
	app.EnableBashCompletion = true

	// set app defaultFlags
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        networkFlag,
			Aliases:     aliases[networkFlag],
			Usage:       "network environment (mainnet, testnet, testnet1, devnet, local, custom)",
			Value:       "mainnet",
			Destination: &network,
		},
		&cli.StringFlag{
			Name:        hostFlag,
			Usage:       "custom full-node host",
			Value:       "",
			Destination: &host,
		},
	}

	// all account-related commands
	accountCommands := []*cli.Command{
		{
			Name:     "keyinfo",
			Usage:    "print all related-keys of a private key",
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
			Usage:    "check the balance of an account",
			Category: accountCat,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     privateKeyFlag,
					Aliases:  aliases[privateKeyFlag],
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  tokenIDFlag,
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
			},
			Action: checkBalance,
		},
		{
			Name:     "utxo",
			Usage:    "print the UTXOs of an account",
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
			Usage:   "consolidate UTXOs of an account",
			Description: "This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDFlag],
				defaultFlags[versionFlag],
				defaultFlags[numThreadsFlag],
				defaultFlags[enableLogFlag],
				defaultFlags[logFileFlag],
			},
			Action: consolidateUTXOs,
		},
		{
			Name:    "history",
			Aliases: []string{"hst"},
			Usage:   "retrieve the history of an account",
			Description: "This function helps retrieve the history of an account w.r.t a tokenID. " +
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
				defaultFlags[enableLogFlag],
				defaultFlags[logFileFlag],
				defaultFlags[csvFileFlag],
			},
			Action: getHistory,
		},
		{
			Name:        "generateaccount",
			Aliases:     []string{"genacc"},
			Usage:       "generate a new account",
			Description: "This function helps generate a new mnemonic phrase and its Incognito account.",
			Category:    accountCat,
			Action:      genKeySet,
		},
		{
			Name:    "submitkey",
			Aliases: []string{"sub"},
			Usage:   "submit an ota key to the full-node",
			Description: "This function submits an otaKey to the full-node to use the full-node's cache. If an access token " +
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
			Name:     "checkrewards",
			Usage:    "get all rewards of a payment address",
			Category: committeeCat,
			Flags: []cli.Flag{
				defaultFlags[addressFlag],
			},
			Action: checkRewards,
		},
		{
			Name:     "withdrawreward",
			Usage:    "withdraw the reward of a privateKey w.r.t to a tokenID.",
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
			Usage: "send an amount of PRV or token from one wallet to another wallet",
			Description: "This function sends an amount of PRV or token from one wallet to another wallet. By default, " +
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
			Usage: "convert UTXOs of an account w.r.t a tokenID",
			Description: "This function helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: transactionCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDFlag],
				defaultFlags[numThreadsFlag],
				defaultFlags[enableLogFlag],
				defaultFlags[logFileFlag],
			},
			Action: convertUTXOs,
		},
	}

	app.Commands = make([]*cli.Command, 0)
	app.Commands = append(app.Commands, accountCommands...)
	app.Commands = append(app.Commands, committeeCommands...)
	app.Commands = append(app.Commands, txCommands...)

	for _, command := range app.Commands {
		buildUsageTextFromCommand(command)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
