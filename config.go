package main

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/incognito-cli/bridge/portal"
	"github.com/spf13/viper"
)

type FConfig struct {
	Network           string   `mapstructure:"network"`
	Host              string   `mapstructure:"host"`
	PrivateKeys       string   `mapstructure:"privatekeys"`
	MiningKeys        string   `mapstructure:"miningkeys"`
	SelfCache         bool     `mapstructure:"selfcache"`
	Debug             bool     `mapstructure:"debug"`
	MaxGetCoinThreads int      `mapstructure:"maxthreads"`
	WantedCoin        []string `mapstructure:"wantedcoins"`
}

var defaultFConfig = map[string]interface{}{
	"network":     "local",
	"host":        "127.0.0.1:9334",
	"privatekeys": []string{},
	"miningkeys":  []string{},
	"selfcache":   true,
	"debug":       true,
	"maxthreads":  runtime.NumCPU(),
	"wantedcoins": []string{
		"0000000000000000000000000000000000000000000000000000000000000004", //PRV
	},
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func LoadConfig(filePath string) (config *FConfig, err error) {
	for k, v := range defaultFConfig {
		viper.SetDefault(k, v)
	}
	dir, fName := filepath.Split(filePath)
	if dir == "" {
		dir = "."
	}
	fExt := filepath.Ext(fName)
	if fExt == "" {
		fExt = "yml"
	} else {
		fExt = fExt[1:]
	}

	viper.AddConfigPath(dir)
	viper.SetConfigName(strings.TrimSuffix(fName, filepath.Ext(fName)))
	viper.SetConfigType(fExt)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	config = &FConfig{}

	err = viper.Unmarshal(config)
	return
}

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

func NewConfigFromFile(fConf *FConfig) (*incclient.IncClient, error) {
	var (
		ethHost   string
		bscHost   string
		plgHost   string
		fmtHost   string
		ethSCAddr string
		bscSCAddr string
		plgSCAddr string
		fmtSCAddr string
		btcClient *portal.BTCClient
		err       error
		incClient *incclient.IncClient
	)

	switch fConf.Network {
	case "mainnet":
		ethHost = incclient.MainNetETHHost
		bscHost = incclient.MainNetBSCHost
		plgHost = incclient.MainNetPLGHost
		fmtHost = incclient.MainNetFTMHost
		ethSCAddr = incclient.MainNetETHContractAddressStr
		bscSCAddr = incclient.MainNetBSCContractAddressStr
		plgSCAddr = incclient.MainNetPLGContractAddressStr
		fmtSCAddr = incclient.MainNetFTMContractAddressStr
		btcClient, err = portal.NewBTCMainNetClient()
	case "testnet":
		ethHost = incclient.TestNetETHHost
		bscHost = incclient.TestNetBSCHost
		plgHost = incclient.TestNetPLGHost
		fmtHost = incclient.TestNetFTMHost
		ethSCAddr = incclient.TestNetETHContractAddressStr
		bscSCAddr = incclient.TestNetBSCContractAddressStr
		plgSCAddr = incclient.TestNetPLGContractAddressStr
		fmtSCAddr = incclient.TestNetFTMContractAddressStr
		btcClient, err = portal.NewBTCTestNetClient()
	case "testnet1":
		ethHost = incclient.TestNet1ETHHost
		bscHost = incclient.TestNet1BSCHost
		plgHost = incclient.TestNet1PLGHost
		fmtHost = incclient.TestNet1FTMHost
		ethSCAddr = incclient.TestNet1ETHContractAddressStr
		bscSCAddr = incclient.TestNet1BSCContractAddressStr
		plgSCAddr = incclient.TestNet1PLGContractAddressStr
		fmtSCAddr = incclient.TestNet1FTMContractAddressStr
		btcClient, err = portal.NewBTCTestNetClient()
	case "local":
		ethHost = incclient.LocalETHHost
		bscHost = incclient.LocalETHHost
		plgHost = incclient.LocalETHHost
		fmtHost = incclient.LocalETHHost
		ethSCAddr = incclient.LocalETHContractAddressStr
		bscSCAddr = incclient.LocalETHContractAddressStr
		plgSCAddr = incclient.LocalETHContractAddressStr
		fmtSCAddr = incclient.LocalETHContractAddressStr
		btcClient, err = portal.NewBTCTestNetClient()
	}

	ethClient, err := ethclient.Dial(ethHost)
	if err != nil {
		return nil, err
	}

	bscClient, err := ethclient.Dial(bscHost)
	if err != nil {
		return nil, err
	}

	plgClient, err := ethclient.Dial(plgHost)
	if err != nil {
		return nil, err
	}

	fmtClient, err := ethclient.Dial(fmtHost)
	if err != nil {
		return nil, err
	}

	evmClients := map[int]*ethclient.Client{
		rpc.ETHNetworkID: ethClient,
		rpc.BSCNetworkID: bscClient,
		rpc.PLGNetworkID: plgClient,
		rpc.FTMNetworkID: fmtClient,
	}

	evmVaultAddresses := map[int]common.Address{
		rpc.ETHNetworkID: common.HexToAddress(ethSCAddr),
		rpc.BSCNetworkID: common.HexToAddress(bscSCAddr),
		rpc.PLGNetworkID: common.HexToAddress(plgSCAddr),
		rpc.FTMNetworkID: common.HexToAddress(fmtSCAddr),
	}
	if fConf.SelfCache {
		incClient, err = incclient.NewIncClientWithCache(fConf.Host, ethHost, 2, fConf.Network)
	} else {
		incClient, err = incclient.NewIncClient(fConf.Host, ethHost, 2, fConf.Network)
	}
	incclient.Logger.IsEnable = fConf.Debug

	cfg = NewConfig(incClient, evmClients, btcClient, evmVaultAddresses)
	return incClient, err
}
