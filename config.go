package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/incognito-cli/bridge/portal"
)

// Config represents the config of an environment of the CLI tool.
type Config struct {
	incClient *incclient.IncClient

	evmClients map[int]*ethclient.Client

	btcClient *portal.BTCClient

	evmVaultAddresses map[int]common.Address
}

// NewConfig returns a new Config from given parameters.
func NewConfig(
	incClient *incclient.IncClient,
	evmClients map[int]*ethclient.Client,
	btcClient *portal.BTCClient,
	evmVaultAddresses map[int]common.Address,
) *Config {
	return &Config{
		incClient:         incClient,
		evmClients:        evmClients,
		evmVaultAddresses: evmVaultAddresses,
		btcClient:         btcClient,
	}
}

// NewTestNetConfig creates a new testnet Config.
func NewTestNetConfig(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewTestNetClient()
		} else {
			incClient, err = incclient.NewTestNetClientWithCache()
		}
		if err != nil {
			return err
		}
	}

	ethClient, err := ethclient.Dial(incclient.TestNetETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.TestNetBSCHost)
	if err != nil {
		return err
	}

	plgClient, err := ethclient.Dial(incclient.TestNetPLGHost)
	if err != nil {
		return err
	}

	ftmClient, err := ethclient.Dial(incclient.TestNetFTMHost)
	if err != nil {
		return err
	}

	evmClients := map[int]*ethclient.Client{
		rpc.ETHNetworkID: ethClient,
		rpc.BSCNetworkID: bscClient,
		rpc.PLGNetworkID: plgClient,
		rpc.FTMNetworkID: ftmClient,
	}

	evmVaultAddresses := map[int]common.Address{
		rpc.ETHNetworkID: common.HexToAddress(incclient.TestNetETHContractAddressStr),
		rpc.BSCNetworkID: common.HexToAddress(incclient.TestNetBSCContractAddressStr),
		rpc.PLGNetworkID: common.HexToAddress(incclient.TestNetPLGContractAddressStr),
		rpc.FTMNetworkID: common.HexToAddress(incclient.TestNetFTMContractAddressStr),
	}

	btcClient, err := portal.NewBTCTestNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, evmClients, btcClient, evmVaultAddresses)

	return nil
}

// NewTestNet1Config creates a new testnet1 Config.
func NewTestNet1Config(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewTestNet1Client()
		} else {
			incClient, err = incclient.NewTestNet1ClientWithCache()
		}
		if err != nil {
			return err
		}
	}

	ethClient, err := ethclient.Dial(incclient.TestNet1ETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.TestNet1BSCHost)
	if err != nil {
		return err
	}

	plgClient, err := ethclient.Dial(incclient.TestNet1PLGHost)
	if err != nil {
		return err
	}

	ftmClient, err := ethclient.Dial(incclient.TestNet1FTMHost)
	if err != nil {
		return err
	}

	evmClients := map[int]*ethclient.Client{
		rpc.ETHNetworkID: ethClient,
		rpc.BSCNetworkID: bscClient,
		rpc.PLGNetworkID: plgClient,
		rpc.FTMNetworkID: ftmClient,
	}

	evmVaultAddresses := map[int]common.Address{
		rpc.ETHNetworkID: common.HexToAddress(incclient.TestNet1ETHContractAddressStr),
		rpc.BSCNetworkID: common.HexToAddress(incclient.TestNet1BSCContractAddressStr),
		rpc.PLGNetworkID: common.HexToAddress(incclient.TestNet1PLGContractAddressStr),
		rpc.FTMNetworkID: common.HexToAddress(incclient.TestNet1FTMContractAddressStr),
	}

	btcClient, err := portal.NewBTCTestNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, evmClients, btcClient, evmVaultAddresses)
	return nil
}

// NewMainNetConfig creates a new main-net Config.
func NewMainNetConfig(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewMainNetClient()
		} else {
			incClient, err = incclient.NewMainNetClientWithCache()
		}
		if err != nil {
			return err
		}
	}
	isMainNet = true

	ethClient, err := ethclient.Dial(incclient.MainNetETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.MainNetBSCHost)
	if err != nil {
		return err
	}

	plgClient, err := ethclient.Dial(incclient.MainNetPLGHost)
	if err != nil {
		return err
	}

	ftmClient, err := ethclient.Dial(incclient.MainNetFTMHost)
	if err != nil {
		return err
	}

	evmClients := map[int]*ethclient.Client{
		rpc.ETHNetworkID: ethClient,
		rpc.BSCNetworkID: bscClient,
		rpc.PLGNetworkID: plgClient,
		rpc.FTMNetworkID: ftmClient,
	}

	evmVaultAddresses := map[int]common.Address{
		rpc.ETHNetworkID: common.HexToAddress(incclient.MainNetETHContractAddressStr),
		rpc.BSCNetworkID: common.HexToAddress(incclient.MainNetBSCContractAddressStr),
		rpc.PLGNetworkID: common.HexToAddress(incclient.MainNetPLGContractAddressStr),
		rpc.FTMNetworkID: common.HexToAddress(incclient.MainNetFTMContractAddressStr),
	}

	btcClient, err := portal.NewBTCMainNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, evmClients, btcClient, evmVaultAddresses)
	return nil
}

// NewLocalConfig creates a new local Config.
func NewLocalConfig(incClient *incclient.IncClient) error {
	var err error
	if incClient == nil {
		if cache == 0 {
			incClient, err = incclient.NewLocalClient("")
		} else {
			incClient, err = incclient.NewLocalClientWithCache()
		}
		if err != nil {
			return err
		}
	}

	ethClient, err := ethclient.Dial(incclient.LocalETHHost)
	if err != nil {
		return err
	}

	bscClient, err := ethclient.Dial(incclient.LocalETHHost)
	if err != nil {
		return err
	}

	plgClient, err := ethclient.Dial(incclient.LocalETHHost)
	if err != nil {
		return err
	}

	fmtClient, err := ethclient.Dial(incclient.LocalETHHost)
	if err != nil {
		return err
	}

	evmClients := map[int]*ethclient.Client{
		rpc.ETHNetworkID: ethClient,
		rpc.BSCNetworkID: bscClient,
		rpc.PLGNetworkID: plgClient,
		rpc.FTMNetworkID: fmtClient,
	}

	evmVaultAddresses := map[int]common.Address{
		rpc.ETHNetworkID: common.HexToAddress(incclient.LocalETHContractAddressStr),
		rpc.BSCNetworkID: common.HexToAddress(incclient.LocalETHContractAddressStr),
		rpc.PLGNetworkID: common.HexToAddress(incclient.LocalETHContractAddressStr),
		rpc.FTMNetworkID: common.HexToAddress(incclient.LocalETHContractAddressStr),
	}

	btcClient, err := portal.NewBTCTestNetClient()
	if err != nil {
		return err
	}

	cfg = NewConfig(incClient, evmClients, btcClient, evmVaultAddresses)
	return nil
}
