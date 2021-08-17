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
		return initClientWithCache(host, "", clientVersion)
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
	client, err = incclient.NewMainNetClientWithCache()

	return err
}
func initTestNet() error {
	var err error
	client, err = incclient.NewTestNetClientWithCache()

	return err
}
func initTestNet1() error {
	var err error
	client, err = incclient.NewTestNet1ClientWithCache()

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
	client, err = incclient.NewLocalClientWithCache()
	if err != nil {
		return err
	}

	return nil
}
func initClientWithCache(rpcHost, ethHost string, version int) error {
	var err error
	client, err = incclient.NewIncClientWithCache(rpcHost, ethHost, version)

	return err
}
