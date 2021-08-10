package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

var network string
var host string
var clientVersion int

func initNetWork() error {
	if host != "" {
		fmt.Printf("host: %v\n", host)
		return initClient(host, "", clientVersion)
	}
	switch network {
	case "mainnet":
		return initMainNet()
	case "testnet":
		return initTestNet()
	case "testnet1":
		return initTestNet1()
	case "devnet":
		return initDevNet()
	case "local":
		return initLocal("")
	}

	return fmt.Errorf("network not found")
}
func initMainNet() error {
	var err error
	client, err = incclient.NewMainNetClient()

	return err
}
func initTestNet() error {
	var err error
	client, err = incclient.NewTestNetClient()

	return err
}
func initTestNet1() error {
	var err error
	client, err = incclient.NewTestNet1Client()

	return err
}
func initDevNet() error {
	var err error
	client, err = incclient.NewDevNetClient()
	if err != nil {
		return err
	}

	return nil
}
func initLocal(port string) error {
	var err error
	client, err = incclient.NewLocalClient(port)
	if err != nil {
		return err
	}

	return nil
}
func initClient(rpcHost, ethHost string, version int) error {
	var err error
	client, err = incclient.NewIncClient(rpcHost, ethHost, version)

	return err
}
