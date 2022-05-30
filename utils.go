package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	network              string
	host                 string
	debug                int
	cache                int
	askUser              = true
	isMainNet            = false
	clientVersion        = 2
	clientCacheDirectory = "incognito-cli"
)

func defaultBeforeFunc(_ *cli.Context) error {
	return initNetWork()
}

func initNetWork() error {
	clientConfig := *config.MainNetConfig
	switch strings.ToLower(network) {
	case "mainnet", "main-net":
	case "testnet", "test-net":
		clientConfig = *config.TestNetConfig
	case "testnet1", "test-net-1", "test-net1":
		clientConfig = *config.TestNet1Config
	case "local":
		clientConfig = *config.LocalConfig
	default:
		return fmt.Errorf("network not found")
	}

	// if UTXO is enabled, use the CLI cache folder instead of the default.
	if cache != 0 {
		clientConfig.UTXOCache.Enable = true
		clientConfig.UTXOCache.MaxGetCoinThreads = 20
		homeDirectory := os.Getenv("HOME")
		if homeDirectory != "" {
			clientConfig.UTXOCache.CacheLocation = fmt.Sprintf("%v/.cache/%v", homeDirectory, clientCacheDirectory)
		}
	}
	if debug != 0 {
		clientConfig.LogConfig.Enable = true
	}
	if host != "" {
		clientConfig.RPCHost = host
		clientConfig.Version = clientVersion
	}

	return initConfig(&clientConfig)
}

// checkSufficientIncBalance checks if the Incognito balance is not less than the requiredAmount.
func checkSufficientIncBalance(privateKey, tokenIDStr string, requiredAmount uint64) (balance uint64, err error) {
	balance, err = cfg.incClient.GetBalance(privateKey, tokenIDStr)
	if err != nil {
		return
	}
	if balance < requiredAmount {
		err = fmt.Errorf("need at least %v of token %v to continue", requiredAmount, tokenIDStr)
	}

	return
}

// promptInput asks for input from the user and saves input to `response`.
// If isSecret is `true`, it will not echo user's input on the terminal.
func promptInput(message string, response interface{}, isSecret ...bool) ([]byte, error) {
	fmt.Printf("%v %v: ", time.Now().Format("2006/01/02 15:04:05"), message)

	var input []byte
	var err error
	if len(isSecret) > 0 && isSecret[0] {
		input, err = terminal.ReadPassword(0)
		if err != nil {
			return nil, err
		}
		fmt.Println()
	} else {
		reader := bufio.NewReader(os.Stdin)
		tmpInput, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		tmpInput = parseInput(tmpInput)
		input = []byte(tmpInput)
	}

	switch reflect.TypeOf(response).String() {
	case "*string", "string":
		response = string(input)
	default:
		err = json.Unmarshal(input, response)
		if err != nil {
			return nil, err
		}
	}

	return input, nil
}

// yesNoPrompt asks for a yes/no decision from the user.
func yesNoPrompt(message string) {
	fmt.Printf("%v %v (y/n): ", time.Now().Format("2006/01/02 15:04:05"), message)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = parseInput(input)

	if !strings.Contains(input, "y") && !strings.Contains(input, "Y") {
		log.Fatal("Abort!!")
	}
}

func parseInput(text string) string {
	if len(text) == 0 {
		return text
	}
	if text[len(text)-1] == 13 || text[len(text)-1] == 10 {
		text = text[:len(text)-1]
	}

	return text
}
