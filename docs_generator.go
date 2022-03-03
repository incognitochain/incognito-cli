package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
)

// generateDocsToFile automatically creates a doc file for the application.
func generateDocsToFile(app *cli.App, file string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	write(f, "There are two options for you to run the Incognito CLI by:\n1. Downloading the pre-compiled executable binary file, you can find it in the [releases](https://github.com/incognitochain/incognito-cli/releases).\n2. Compiling your own executable binary file from source as in the Installation instruction above.\n\n")
	write(f, "Then execute the binary file with the following commands.\n\n")
	write(f, "```shell\n")
	write(f, fmt.Sprintf("$ %v help\n", app.Name))

	app.Writer = f
	err = app.Run([]string{"help"})
	if err != nil {
		return err
	}
	write(f, "```\n")
	write(f, "# Commands\n")
	write(f, "<!-- commands -->\n")

	categories, commandsInCategory := groupCommandByCategory(app.Commands)
	for i := range categories {
		cat := categories[i]
		write(f, fmt.Sprintf("* [`%v`](#%v)\n", cat, strings.ToLower(cat)))
		for j := range commandsInCategory[cat] {
			cmdName := commandsInCategory[cat][j].Name
			write(f, fmt.Sprintf("\t* [`%v`](#%v)\n", cmdName, strings.ToLower(cmdName)))
			for _, subCommand := range commandsInCategory[cat][j].Subcommands {
				subCmdName := fmt.Sprintf("%v_%v", cmdName, subCommand.Name)
				write(f, fmt.Sprintf("\t\t* [`%v %v`](#%v)\n", cmdName, subCommand.Name, strings.ToLower(subCmdName)))
			}
		}

	}

	for i := range categories {
		cat := categories[i]
		write(f, fmt.Sprintf("## %v\n", cat))
		for j := range commandsInCategory[cat] {
			cmd := commandsInCategory[cat][j]
			err = createDocsForCommand(app, cmd, f)
			if err != nil {
				return err
			}
			if len(cmd.Subcommands) > 0 {
				for _, subCmd := range cmd.Subcommands {
					err = createDocsForCommand(app, subCmd, f, cmd.Name)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	write(f, "<!-- commandsstop -->\n")

	return nil
}

// createDocsForCommand automatically creates a doc file for the application.
func createDocsForCommand(app *cli.App, command *cli.Command, f *os.File, parents ...string) error {
	parent := ""
	if len(parents) > 0 {
		parent = parents[0]
	}

	if parent == "" {
		write(f, fmt.Sprintf("### %v\n", command.Name))
		usageString := command.Description
		if usageString == "" {
			usageString = command.Usage
		}
		write(f, fmt.Sprintf("%v\n", usageString))
		write(f, "```shell\n")

		// Write command name
		write(f, fmt.Sprintf("$ %v help %v\n", app.Name, command.Name))
		err := app.Run([]string{"help", "help", command.Name})
		if err != nil {
		}
		write(f, "```\n\n")
	} else {
		write(f, fmt.Sprintf("#### %v_%v\n", parent, command.Name))
		usageString := command.Description
		if usageString == "" {
			usageString = command.Usage
		}
		write(f, fmt.Sprintf("%v\n", usageString))
		write(f, "```shell\n")

		// Write command name
		write(f, fmt.Sprintf("$ %v %v help %v\n", app.Name, parent, command.Name))
		err := app.Run([]string{"help", parent, "help", command.Name})
		if err != nil {
		}
		write(f, "```\n\n")
	}

	return nil
}

func write(f *os.File, content string) {
	_, err := f.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func groupCommandByCategory(commands []*cli.Command) (category []string, commandsByCategory map[string][]*cli.Command) {
	category = make([]string, 0)
	commandsByCategory = make(map[string][]*cli.Command)

	for _, cmd := range commands {
		if cmd.Category == "" {
			continue
		}
		if cmdInCat, ok := commandsByCategory[cmd.Category]; ok {
			cmdInCat = append(cmdInCat, cmd)
			commandsByCategory[cmd.Category] = cmdInCat
		} else {
			cmdInCat = make([]*cli.Command, 0)
			cmdInCat = append(cmdInCat, cmd)
			category = append(category, cmd.Category)
			commandsByCategory[cmd.Category] = cmdInCat
		}
	}

	sort.Strings(category)
	for _, tmpCommands := range commandsByCategory {
		sort.Sort(cli.CommandsByName(tmpCommands))
	}

	return
}
