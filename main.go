package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func main() {
	app := &cli.App{
		Name:    "incognito-cli",
		Usage:   "A simple CLI application for the Incognito network",
		Version: "v1.0.0",
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

	app.Commands = make([]*cli.Command, 0)
	app.Commands = append(app.Commands, accountCommands...)
	app.Commands = append(app.Commands, committeeCommands...)
	app.Commands = append(app.Commands, txCommands...)
	app.Commands = append(app.Commands, pDEXCommands...)
	app.Commands = append(app.Commands, bridgeCommands...)

	for _, command := range app.Commands {
		if len(command.Subcommands) > 0 {
			sort.Sort(cli.CommandsByName(command.Subcommands))
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
