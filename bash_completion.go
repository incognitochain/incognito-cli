package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func defaultSubCommandCompletion(c *cli.Context) {
	for _, f := range c.Command.Flags {
		flagName := f.Names()[0]
		fmt.Printf("--%v\n", flagName)
	}
}
