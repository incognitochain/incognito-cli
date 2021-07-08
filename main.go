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
		Name:        "incognito-cli",
		Version:     "v0.0.1",
		Description: "A simple CLI application for the Incognito network",
	}

	// set app flags
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "network",
			Aliases:     []string{"net"},
			Usage:       "network environment (mainnet, testnet, testnet1, devnet, local, custom)",
			Value:       "mainnet",
			Destination: &network,
		},
		&cli.StringFlag{
			Name:        "host",
			Usage:       "custom full-node host",
			Destination: &host,
		},
	}

	// all account-related commands
	accountCommands := []*cli.Command{
		{
			Name:      "keyinfo",
			Usage:     "print all related-keys of a private key",
			UsageText: "keyinfo --privateKey PRIVATE_KEY",
			ArgsUsage: "--privateKey PRIVATE_KEY",
			Category:  "account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "privateKey",
					Aliases:  []string{"prvKey"},
					Usage:    "a base58-encoded private key",
					Required: true,
				},
			},
			Action: keyInfo,
		},
		{
			Name:      "balance",
			Usage:     "check the balance of an account",
			UsageText: "balance --privateKey PRIVATE_KEY --tokenID TOKEN_ID",
			ArgsUsage: "--privateKey PRIVATE_KEY --tokenID TOKEN_ID",
			Category:  "account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "privateKey",
					Aliases:  []string{"prvKey"},
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "tokenID",
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
			},
			Action: checkBalance,
		},
		{
			Name:      "utxo",
			Usage:     "print the UTXOs of an account",
			UsageText: "utxo --privateKey PRIVATE_KEY --tokenID TOKEN_ID",
			ArgsUsage: "--privateKey PRIVATE_KEY --tokenID TOKEN_ID",
			Category:  "account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "privateKey",
					Aliases:  []string{"prvKey"},
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "tokenID",
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
			},
			Action: checkUTXOs,
		},
		{
			Name:      "consolidate",
			Aliases:   []string{"csl"},
			Usage:     "consolidate UTXOs of an account",
			UsageText: "consolidate --privateKey PRIVATE_KEY --tokenID TOKEN_ID --version VERSION --numThreads NUM_THREADS --enableLog ENABLE_LOG --logFile LOG_FILE",
			ArgsUsage: "--privateKey PRIVATE_KEY --tokenID TOKEN_ID",
			Description: "This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: "account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "privateKey",
					Aliases:  []string{"prvKey"},
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "tokenID",
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
				&cli.IntFlag{
					Name:  "version",
					Usage: "version of the UTXOs being converted (1, 2)",
					Value: 1,
				},
				&cli.IntFlag{
					Name:  "numThreads",
					Usage: "number of threads used in this action",
					Value: 4,
				},
				&cli.BoolFlag{
					Name:  "enableLog",
					Usage: "enable log for this action",
					Value: false,
				},
				&cli.StringFlag{
					Name:  "logFile",
					Usage: "location of the log file",
					Value: "os.Stdout",
				},
			},
			Action: consolidateUTXOs,
		},
		{
			Name:      "history",
			Aliases:   []string{"hst"},
			Usage:     "retrieve the history of an account",
			UsageText: "history --privateKey PRIVATE_KEY --tokenID TOKEN_ID --numThreads NUM_THREADS --enableLog ENABLE_LOG --logFile LOG_FILE --csvFile CSV_FILE",
			Description: "This function helps retrieve the history of an account w.r.t a tokenID. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: "account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "privateKey",
					Aliases:  []string{"prvKey"},
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "tokenID",
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
				&cli.IntFlag{
					Name:  "numThreads",
					Usage: "number of threads used in this action",
					Value: 4,
				},
				&cli.BoolFlag{
					Name:  "enableLog",
					Usage: "enable log for this action",
					Value: false,
				},
				&cli.StringFlag{
					Name:  "logFile",
					Usage: "location of the log file",
					Value: "os.Stdout",
				},
				&cli.StringFlag{
					Name:  "csvFile",
					Usage: "the csv file location to store the history",
				},
			},
			Action: getHistory,
		},
		{
			Name:        "generateaccount",
			Aliases:     []string{"genacc"},
			Usage:       "generate a new account",
			UsageText:   "generateaccount",
			Description: "This function helps generate a new mnemonic phrase and its Incognito account.",
			Category:    "account",
			Action:      genKeySet,
		},
		{
			Name:      "submitkey",
			Aliases:   []string{"sub"},
			Usage:     "submit an ota key to the full-node",
			UsageText: "submitkey --otaKey OTA_KEY --accessToken ACCESS_TOKEN --fromHeight FROM_HEIGHT --isReset IS_RESET",
			Description: "This function submits an otaKey to the full-node to use the full-node's cache. If an access token " +
				"is provided, it will submit the ota key in an authorized manner. See " +
				"https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md " +
				"for more details.",
			Category: "account",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "otaKey",
					Aliases:  []string{"ota"},
					Usage:    "a base58-encoded ota key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "accessToken",
					Usage: "a 64-character long hex-encoded authorized access token",
					Value: "",
				},
				&cli.Uint64Flag{
					Name:  "fromHeight",
					Usage: "the beacon height at which the full-node will sync from",
					Value: 0,
				},
				&cli.BoolFlag{
					Name:  "isReset",
					Usage: "whether the full-node should reset the cache for this ota key",
					Value: false,
				},
			},
			Action: submitKey,
		},
	}

	// all committee-related commands
	committeeCommands := []*cli.Command{
		{
			Name:      "checkrewards",
			Usage:     "get all rewards of a payment address",
			UsageText: "checkRewards --address PAYMENT_ADDRESS",
			ArgsUsage: "--address PAYMENT_ADDRESS",
			Category:  "committee",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "address",
					Aliases:  []string{"addr"},
					Usage:    "a base58-encoded payment address",
					Required: true,
				},
			},
			Action: checkRewards,
		},
		{
			Name:      "withdrawreward",
			Usage:     "withdraw the reward of a privateKey w.r.t to a tokenID.",
			UsageText: "withdrawreward --privateKey PRIVATE_KEY --tokenID TOKEN_ID --version VERSION",
			ArgsUsage: "--privateKey PRIVATE_KEY --tokenID TOKEN_ID",
			Category:  "committee",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "privateKey",
					Aliases:  []string{"prvKey"},
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "address",
					Aliases: []string{"addr"},
					Usage:   "the payment address of a candidate (default: the payment address of the privateKey)",
				},
				&cli.StringFlag{
					Name:  "tokenID",
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
				&cli.IntFlag{
					Name:    "version",
					Aliases: []string{"v"},
					Usage:   "version of the transaction (1 or 2)",
					Value:   2,
				},
			},
			Action: withdrawReward,
		},
	}

	app.Commands = make([]*cli.Command, 0)
	app.Commands = append(app.Commands, accountCommands...)
	app.Commands = append(app.Commands, committeeCommands...)

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
