package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	network       string
	host          string
	debug         int
	cache         int
	askUser       = true
	isMainNet     = false
	clientVersion = 2
)

func defaultBeforeFunc(_ *cli.Context) error {
	return initNetWork()
}

func initNetWork() error {
	if cache != 0 {
		incclient.MaxGetCoinThreads = 20
	}
	if debug != 0 {
		incclient.Logger.IsEnable = true
	}
	if host != "" {
		fmt.Printf("host: %v, version: %v\n", host, clientVersion)
		return initClient(host, clientVersion)
	}
	switch network {
	case "mainnet":
		return NewMainNetConfig(nil)
	case "testnet":
		return NewTestNetConfig(nil)
	case "testnet1":
		return NewTestNet1Config(nil)
	case "local":
		return NewLocalConfig(nil)
	}

	return fmt.Errorf("network not found")
}
func initClient(rpcHost string, version int) error {
	ethNode := incclient.MainNetETHHost
	var err error
	switch network {
	case "testnet":
		ethNode = incclient.TestNetETHHost
		err = NewTestNetConfig(nil)
	case "testnet1":
		ethNode = incclient.TestNet1ETHHost
		err = NewTestNet1Config(nil)
	case "local":
		ethNode = incclient.LocalETHHost
		err = NewLocalConfig(nil)
	default:
		err = NewMainNetConfig(nil)
	}
	if err != nil {
		return err
	}

	incClient, err := incclient.NewIncClient(rpcHost, ethNode, version, network)
	if cache != 0 {
		incClient, err = incclient.NewIncClientWithCache(rpcHost, ethNode, version, network)
	}
	if err != nil {
		return err
	}

	cfg.incClient = incClient
	return nil
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
