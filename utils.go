package main

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
)

var network string
var host string
var clientVersion int
var debug int

func initNetWork() error {
	if debug != 0 {
		incclient.Logger.IsEnable = true
	}
	if host != "" {
		fmt.Printf("host: %v\n", host)
		return initClient(host, clientVersion)
	}
	switch network {
	case "mainnet":
		return NewMainNetConfig(nil)
	case "testnet":
		return NewTestNetConfig(nil)
	case "testnet1":
		return NewTestNet1Config(nil)
	case "devnet":
		return NewDevNetConfig(nil)
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
	case "devnet":
		ethNode = incclient.DevNetETHHost
		err = NewDevNetConfig(nil)
	case "local":
		ethNode = incclient.LocalETHHost
		err = NewLocalConfig(nil)
	default:
		err = NewMainNetConfig(nil)
	}
	if err != nil {
		return err
	}

	incClient, err := incclient.NewIncClientWithCache(rpcHost, ethNode, version)
	if err != nil {
		return err
	}

	cfg.incClient = incClient
	return nil
}
