package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/transaction/tx_ver2"
	"github.com/incognitochain/incognito-cli/bridge/evm/erc20"
	"github.com/incognitochain/incognito-cli/bridge/evm/prv20"
	"github.com/incognitochain/incognito-cli/bridge/evm/vault"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var (
	nativeToken         = "0x0000000000000000000000000000000000000000"
	nativeTokenDecimals = 18
	EthGasStationAPIKey = "dc93dbe33e5ebd48ae0fbdc0a300d269a722566f98beb11682169319d624"
)

func errEVMNetworkNotSupported(evmNetworkID int) error {
	return fmt.Errorf("evmNetworkID %v not supported", evmNetworkID)
}

// EVMTokenInfo represents the information of an ERC20/BEP20 token.
type EVMTokenInfo struct {
	// network is the name of the network (ETH/BSC) where the token resides in.
	network string

	// address is the address of the token.
	address common.Address

	// name is the name of the token.
	name string

	// symbol is the symbol of the token.
	symbol string

	// totalSupply represents the total supply of the token.
	totalSupply *big.Int
}

// getEVMTokenInfo returns the info of an ERC20/BEP20 token.
func getEVMTokenInfo(tokenAddressStr string, evmNetworkID int) (*EVMTokenInfo, error) {
	tokenAddress := common.HexToAddress(tokenAddressStr)
	if tokenAddress.String() == nativeToken {
		return nil, fmt.Errorf("this is a native token")
	}

	evmClient := cfg.evmClients[evmNetworkID]
	if evmClient == nil {
		return nil, errEVMNetworkNotSupported(evmNetworkID)
	}
	res := new(EVMTokenInfo)
	res.address = tokenAddress
	res.network = "ETH"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		res.network = "BSC"
	case rpc.PLGNetworkID:
		res.network = "PLG"
	case rpc.FTMNetworkID:
		res.network = "FTM"
	}

	erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return nil, err
	}

	res.name, err = erc20Instance.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	res.symbol, err = erc20Instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	res.totalSupply, err = erc20Instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// EVMAccount represents an account on the Ethereum/Binance networks.
type EVMAccount struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    common.Address
}

// NewEVMAccount returns a new EVMAccount given a hex-encoded private key.
func NewEVMAccount(hexPrivateKey string) (*EVMAccount, error) {
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &EVMAccount{
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
		address:    address,
	}, nil

}

// newTransactionOpts creates a new bind.TransactOpts for an EVMAccount.
func (acc EVMAccount) newTransactionOpts(destAddr common.Address, gasPrice, gasLimit, amount uint64, data []byte, evmNetworkID int) (*bind.TransactOpts, error) {
	var err error

	evmClient, _ := getEVMClientAndVaultAddress(evmNetworkID)

	// calculate gas price if needed.
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	//calculate gas limit
	if gasLimit == 0 {
		gasLimit, err = evmClient.EstimateGas(context.Background(), ethereum.CallMsg{From: acc.address, To: &destAddr, Data: data})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}
	}

	nonce, err := evmClient.PendingNonceAt(context.Background(), acc.address)
	if err != nil {
		return nil, fmt.Errorf("get pending nonce error: %v", err)
	}

	chainID, err := evmClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(acc.privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPriceBigInt
	auth.GasLimit = gasLimit
	auth.Nonce = new(big.Int).SetUint64(nonce)
	if amount != 0 {
		value := new(big.Int).SetUint64(amount)
		auth.Value = value
	}

	return auth, nil
}

// estimateDepositGas estimates the gas for depositing a token.
func (acc EVMAccount) estimateDepositGas(tokenAddress common.Address, depositedAmount *big.Int, incAddress string, evmNetworkID int) (uint64, error) {
	evmClient, vaultAddress := getEVMClientAndVaultAddress(evmNetworkID)

	var gasLimit uint64
	vaultABI, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return 0, fmt.Errorf("cannot create vaultABI from file")
	}

	var data []byte
	if tokenAddress.String() == nativeToken {
		data, err = vaultABI.Pack(
			"deposit",
			incAddress,
		)
		if err != nil {
			return 0, err
		}

		gasLimit, err = evmClient.EstimateGas(context.Background(),
			ethereum.CallMsg{From: acc.address, Value: depositedAmount, To: &vaultAddress, Data: data})
		if err != nil {
			return 0, fmt.Errorf("estimateGas for native token error: %v", err)
		}
	} else {
		data, err = vaultABI.Pack(
			"depositERC20",
			tokenAddress,
			depositedAmount,
			incAddress,
		)
		if err != nil {
			return 0, err
		}

		gasLimit, err = evmClient.EstimateGas(context.Background(),
			ethereum.CallMsg{From: acc.address, To: &vaultAddress, Data: data})
		if err != nil {
			return 0, err
		}
	}

	return gasLimit, nil
}

// estimateWithdrawalGas estimates the gas for withdrawing a token.
func (acc EVMAccount) estimateWithdrawalGas(burnProof *incclient.BurnProof, evmNetworkID int) (uint64, error) {
	evmClient, vaultAddress := getEVMClientAndVaultAddress(evmNetworkID)

	vaultABI, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return 0, fmt.Errorf("cannot create vaultABI from file")
	}

	var data []byte
	data, err = vaultABI.Pack(
		"withdraw",
		burnProof.Instruction,
		burnProof.Heights[0],
		burnProof.InstPaths[0],
		burnProof.InstPathIsLefts[0],
		burnProof.InstRoots[0],
		burnProof.BlkData[0],
		burnProof.SigIndices[0],
		burnProof.SigVs[0],
		burnProof.SigRs[0],
		burnProof.SigSs[0],
	)
	if err != nil {
		return 0, err
	}

	gasLimit, err := evmClient.EstimateGas(context.Background(),
		ethereum.CallMsg{From: acc.address, To: &vaultAddress, Data: data})
	if err != nil {
		return 0, fmt.Errorf("estimateGas for withdrawal error: %v", err)
	}

	return gasLimit, nil
}

// getBalance returns the balance of a token.
func (acc EVMAccount) getBalance(tokenAddress common.Address, evmNetworkID int) (*big.Int, *big.Float, error) {
	evmClient, _ := getEVMClientAndVaultAddress(evmNetworkID)

	decimals := uint64(nativeTokenDecimals)
	var balance *big.Int
	var err error
	if tokenAddress.String() == nativeToken {
		balance, err = evmClient.BalanceAt(context.Background(), acc.address, nil)
		if err != nil {
			return nil, nil, err
		}
	} else {
		erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
		if err != nil {
			return nil, nil, err
		}

		balance, err = erc20Instance.BalanceOf(&bind.CallOpts{}, acc.address)
		if err != nil {
			return nil, nil, err
		}

		decimalsBigInt, err := erc20Instance.Decimals(&bind.CallOpts{})
		if err != nil {
			return nil, nil, err
		}

		decimals = decimalsBigInt.Uint64()
	}
	balanceFloat := getSynthesizedAmount(balance, decimals)

	return balance, balanceFloat, nil
}

// getAllowance returns the allowance of an owner to a spender w.r.t to an ERC20 token.
func (acc EVMAccount) getAllowance(tokenAddress, spender common.Address, evmNetworkID int) (uint64, error) {
	evmClient, _ := getEVMClientAndVaultAddress(evmNetworkID)

	erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return 0, err
	}

	allowance, err := erc20Instance.Allowance(&bind.CallOpts{}, acc.address, spender)
	if err != nil {
		return 0, err
	}

	return allowance.Uint64(), nil
}

func (acc EVMAccount) getGasLimitAndPrice(gasLimit, gasPrice uint64, callMsg ethereum.CallMsg, evmNetworkID int) (*big.Int, uint64, error) {
	evmClient, _ := getEVMClientAndVaultAddress(evmNetworkID)
	var err error

	// calculate gas price if needed.
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, 0, err
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	if gasLimit == 0 {
		callMsg.From = acc.address
		gasLimit, err = evmClient.EstimateGas(context.Background(), callMsg)
		if err != nil {
			return nil, 0, err
		}
	}

	return gasPriceBigInt, gasLimit, nil
}

// checkSufficientBalance checks if the balance of the token address is sufficient w.r.t to the requiredAmount.
// It also returns the synthesized balance of the token.
func (acc EVMAccount) checkSufficientBalance(tokenAddress common.Address, requiredAmount float64, evmNetworkID int) (balance float64, err error) {
	_, synthesizedBalance, err := acc.getBalance(tokenAddress, evmNetworkID)
	if err != nil {
		return 0, err
	}

	fBalance, _ := synthesizedBalance.Float64()
	if fBalance < requiredAmount {
		return 0, fmt.Errorf("insufficient balance: required %v, got %v", requiredAmount, fBalance)
	}

	return fBalance, nil
}

// checkAllowance checks if the allowance of the token address is sufficient w.r.t to the requiredAmount.
// It also returns the synthesized allowance of the token.
func (acc EVMAccount) checkAllowance(tokenAddress common.Address, requiredAmount float64, evmNetworkID int) (err error) {
	_, vaultAddress := getEVMClientAndVaultAddress(evmNetworkID)
	prefix := "[CheckAllowanceERC20]"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		prefix = "[CheckAllowanceBEP20]"
	case rpc.PLGNetworkID:
		prefix = "[CheckAllowancePLG20]"
	}

	tokenDecimals, err := getDecimals(tokenAddress, evmNetworkID)
	if err != nil {
		return
	}

	currentAllowance, err := acc.getAllowance(tokenAddress, vaultAddress, evmNetworkID)
	if err != nil {
		return
	}
	allowance, _ := getSynthesizedAmount(new(big.Int).SetUint64(currentAllowance), tokenDecimals).Float64()
	if allowance < requiredAmount {
		approvedAmount := requiredAmount
		if askUser {
			_, err = promptInput(
				fmt.Sprintf("%v insufficient allowance: got %v, need %v. Enter the amount you want to approve (e.g, 0.01, 0.1, 1, 10)", prefix, allowance, requiredAmount),
				&approvedAmount,
			)
			if err != nil {
				return
			}
			if allowance+approvedAmount < requiredAmount {
				err = fmt.Errorf("not enough allowance")
				return
			}
		} else {
			log.Printf("%v insufficient allowance: got %v, need %v\n", prefix, allowance, requiredAmount)
		}
		var txHash *common.Hash
		txHash, err = acc.ApproveERC20(tokenAddress, vaultAddress, approvedAmount, 0, evmNetworkID)
		if err != nil {
			return
		}

		err = wait(*txHash, evmNetworkID)
		if err != nil {
			return
		}
	}

	return
}

func wait(tx common.Hash, evmNetworkID int) error {
	evmClient := cfg.evmClients[evmNetworkID]
	if evmClient == nil {
		return errEVMNetworkNotSupported(evmNetworkID)
	}
	for range time.Tick(10 * time.Second) {
		receipt, err := evmClient.TransactionReceipt(context.Background(), tx)
		if err == nil {
			log.Printf("[EVM Status] TxHash %v: %v\n", tx.String(), receipt.Status)
			if receipt.Status == 0 {
				return fmt.Errorf("tx %v failed", tx.String())
			}
			break
		} else if err == ethereum.NotFound {
			continue
		} else {
			return err
		}
	}
	return nil
}

func verifyProofAndParseReceipt(iReq *metadata.IssuingEVMRequest, evmNetworkID int) (*types.Receipt, error) {
	evmClient := cfg.evmClients[evmNetworkID]
	if evmClient == nil {
		return nil, errEVMNetworkNotSupported(evmNetworkID)
	}

	evmHeader, err := evmClient.HeaderByHash(context.Background(), iReq.BlockHash)
	if err != nil {
		return nil, err
	}
	if evmHeader == nil {
		return nil, fmt.Errorf("WARNING: Could not find out the EVM block header with the hash: %s", iReq.BlockHash.String())
	}

	keyBuf := new(bytes.Buffer)
	keyBuf.Reset()
	err = rlp.Encode(keyBuf, iReq.TxIndex)
	if err != nil {
		return nil, err
	}

	nodeList := new(light.NodeList)
	for _, proofStr := range iReq.ProofStrs {
		proofBytes, err := base64.StdEncoding.DecodeString(proofStr)
		if err != nil {
			return nil, err
		}
		err = nodeList.Put([]byte{}, proofBytes)
		if err != nil {
			return nil, err
		}
	}
	proof := nodeList.NodeSet()
	val, err := trie.VerifyProof(evmHeader.ReceiptHash, keyBuf.Bytes(), proof)
	if err != nil {
		return nil, err
	}

	// Decode value from VerifyProof into Receipt
	constructedReceipt := new(types.Receipt)
	err = rlp.DecodeBytes(val, constructedReceipt)
	if err != nil {
		return nil, err
	}

	if constructedReceipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("the constructedReceipt's status is not success")
	}

	return constructedReceipt, nil
}

func getDecimals(tokenAddress common.Address, evmNetworkID int) (uint64, error) {
	evmClient := cfg.evmClients[evmNetworkID]

	erc20Instance, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return 0, err
	}

	decimals, err := erc20Instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	}

	return decimals.Uint64(), nil
}

// getSynthesizedAmount gets the amount after dividing with 10 to the power of the given decimals.
func getSynthesizedAmount(amount *big.Int, decimals uint64) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		new(big.Float).SetInt(new(big.Int).Exp(new(big.Int).SetUint64(10), new(big.Int).SetUint64(decimals), nil)))
}

func getEVMClientAndVaultAddress(evmNetworkID int) (evmClient *ethclient.Client, vaultAddress common.Address) {
	evmClient = cfg.evmClients[evmNetworkID]
	vaultAddress = cfg.evmVaultAddresses[evmNetworkID]
	return
}

func getAllDecentralizedBridgeTokens() (incToPublic map[string]string, publicToInc map[string]string, err error) {
	incToPublic = make(map[string]string)
	publicToInc = make(map[string]string)

	allTokens, err := cfg.incClient.GetBridgeTokens()
	if err != nil {
		return
	}

	for _, token := range allTokens {
		if token.IsCentralized {
			continue
		}
		incTokenID := token.TokenID.String()
		publicTokenID := fmt.Sprintf("%x", token.ExternalTokenID)

		incToPublic[incTokenID] = publicTokenID
		publicToInc[publicTokenID] = incTokenID
	}

	return
}

func getIncTokenIDFromEVMTokenID(evmTokenID string, evmNetworkID int) (string, error) {
	evmTokenID = strings.Replace(evmTokenID, "0x", "", -1)
	evmTokenID = strings.Replace(evmTokenID, "0X", "", -1)
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		evmTokenID = "425343" + evmTokenID
	case rpc.PLGNetworkID:
		evmTokenID = "504c47" + evmTokenID
	case rpc.FTMNetworkID:
		evmTokenID = "46544d" + evmTokenID
	}

	evmTokenID = strings.ToLower(evmTokenID)

	_, publicToInc, err := getAllDecentralizedBridgeTokens()
	if err != nil {
		return "", err
	}

	if incTokenID, ok := publicToInc[evmTokenID]; ok {
		return incTokenID, nil
	}

	return "", fmt.Errorf("incTokenID not found for evmTokenID %v", evmTokenID)
}

func getEVMTokenIDIncTokenID(incTokenIDStr string) (string, int, error) {
	incTokenIDStr = strings.Replace(incTokenIDStr, "0x", "", -1)
	incTokenIDStr = strings.Replace(incTokenIDStr, "0X", "", -1)

	incTokenIDStr = strings.ToLower(incTokenIDStr)

	incToPublic, _, err := getAllDecentralizedBridgeTokens()
	if err != nil {
		return "", 0, err
	}

	if evmTokenIDStr, ok := incToPublic[incTokenIDStr]; ok {
		if len(evmTokenIDStr) == 46 {
			if strings.Contains(evmTokenIDStr, "425343") {
				return evmTokenIDStr[6:], rpc.BSCNetworkID, nil
			} else if strings.Contains(evmTokenIDStr, "504c47") {
				return evmTokenIDStr[6:], rpc.PLGNetworkID, nil
			} else if strings.Contains(evmTokenIDStr, "46544d") {
				return evmTokenIDStr[6:], rpc.FTMNetworkID, nil
			}
		} else if len(evmTokenIDStr) == 40 {
			return evmTokenIDStr, rpc.ETHNetworkID, nil
		}
	}

	return "", 0, fmt.Errorf("evmToken not found for incTokenIDStr %v", incTokenIDStr)
}

// estimateGasPrice returns the estimated gas price on the EVM network.
func estimateGasPrice(evmNetworkID int) (*big.Int, error) {
	if evmNetworkID == rpc.ETHNetworkID && isMainNet {
		response, err := http.Get(fmt.Sprintf("https://ethgasstation.info/api/ethgasAPI.json?api-key=%v", EthGasStationAPIKey))
		if err == nil {
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			mapRes := make(map[string]interface{})
			err = json.Unmarshal(responseData, &mapRes)
			if err == nil {
				tmpAverageGasPrice, ok := mapRes["average"]
				if ok {
					averageGasPrice, ok := tmpAverageGasPrice.(float64)
					if ok {
						averageGasPrice = averageGasPrice * math.Pow10(9) / 10
						return new(big.Int).SetUint64(uint64(averageGasPrice)), nil
					}
				}
			}
		}
	}
	evmClient, _ := getEVMClientAndVaultAddress(evmNetworkID)
	return evmClient.SuggestGasPrice(context.Background())
}

func getEVMNetworkIDFromName(networkName string) (int, error) {
	evmNetworkID := -1
	switch strings.ToUpper(networkName) {
	case "ETH":
		evmNetworkID = rpc.ETHNetworkID
	case "BSC":
		evmNetworkID = rpc.BSCNetworkID
	case "PLG":
		evmNetworkID = rpc.PLGNetworkID
	case "FTM":
		evmNetworkID = rpc.FTMNetworkID
	default:
		return 0, fmt.Errorf("evmNetwork `%v` not supported", networkName)
	}

	return evmNetworkID, nil
}

func checkAndChangeRPCEndPoint(evmNetworkID int, err error) error {
	if err == nil {
		return nil
	}
	evmNetwork := "ETH"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		evmNetwork = "BSC"
	case rpc.PLGNetworkID:
		evmNetwork = "PLG"
	case rpc.FTMNetworkID:
		evmNetwork = "FTM"
	}

	if strings.Contains(err.Error(), "504 Gateway Timeout") {
		yesNoPrompt(fmt.Sprintf("Gateway time-out. Do you want to change the %v RPC-endpoint?", evmNetwork))
		var newRPCEndPoint string
		var input []byte
		input, err = promptInput(fmt.Sprintf("Enter new %v RPC endpoint", evmNetwork), &newRPCEndPoint, true)
		if err != nil {
			return err
		}
		newRPCEndPoint = string(input)
		cfg.evmClients[evmNetworkID], err = ethclient.Dial(newRPCEndPoint)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

/*
 * These following functions are to interact with the Incognito network to either shield or un-shield assets.
 */

// DepositNative deposits an amount of ETH/BNB to the Incognito contract.
func (acc EVMAccount) DepositNative(incAddress string, depositedAmount float64, gasLimit, gasPrice uint64, evmNetworkID int) (*common.Hash, error) {
	prefix := "[DepositETH]"
	evmClient, vaultAddress := getEVMClientAndVaultAddress(evmNetworkID)
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		prefix = "[DepositBNB]"
	case rpc.PLGNetworkID:
		prefix = "[DepositMATIC]"
	case rpc.FTMNetworkID:
		prefix = "[DepositFTM]"
	}

	v, err := vault.NewVault(vaultAddress, evmClient)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, err
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	// estimate gasLimit
	amountBigInt := new(big.Int).SetUint64(uint64(depositedAmount * math.Pow10(nativeTokenDecimals)))
	gasLimit, err = acc.estimateDepositGas(common.HexToAddress(nativeToken), amountBigInt, incAddress, evmNetworkID)
	if err != nil {
		return nil, err
	}
	txFee, _ := getSynthesizedAmount(new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPriceBigInt), uint64(nativeTokenDecimals)).Float64()
	requiredAmount := txFee + depositedAmount
	_, err = acc.checkSufficientBalance(common.HexToAddress(nativeToken), requiredAmount, evmNetworkID)
	if err != nil {
		return nil, err
	}
	if askUser {
		yesNoPrompt(fmt.Sprintf("%v DepositAmount: %v, GasPrice: %v gWei, DepositFee: %v, TotalAmount: %v. Do you want to continue?",
			prefix, depositedAmount, float64(gasPriceBigInt.Uint64())/math.Pow10(9), txFee, requiredAmount))
	}

	auth, err := acc.newTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, amountBigInt.Uint64(), nil, evmNetworkID)
	if err != nil {
		return nil, err
	}

	tx, err := v.Deposit(auth, incAddress)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()
	log.Printf("%v Deposited Tx: %v\n", prefix, txHash.String())

	if err := wait(txHash, evmNetworkID); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// DepositToken shields an amount of ERC20/BEP20 to the Incognito network.
func (acc EVMAccount) DepositToken(incAddress, tokenAddressStr string, depositedAmount float64, gasLimit, gasPrice uint64, evmNetworkID int) (*common.Hash, error) {
	evmClient, vaultAddress := getEVMClientAndVaultAddress(evmNetworkID)
	prefix := "[DepositERC20]"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		prefix = "[DepositBEP20]"
	case rpc.PLGNetworkID:
		prefix = "[DepositPLG20]"
	case rpc.FTMNetworkID:
		prefix = "[DepositFTM20]"
	}

	// load the vault instance
	v, err := vault.NewVault(vaultAddress, evmClient)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(tokenAddressStr)
	_, err = acc.checkSufficientBalance(tokenAddress, depositedAmount, evmNetworkID)
	if err != nil {
		return nil, err
	}
	err = acc.checkAllowance(tokenAddress, depositedAmount, evmNetworkID)
	if err != nil {
		return nil, err
	}
	tokenDecimals, err := getDecimals(tokenAddress, evmNetworkID)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	// estimate gasLimit
	amountBigInt := new(big.Int).SetUint64(uint64(depositedAmount * math.Pow10(int(tokenDecimals))))
	gasLimit, err = acc.estimateDepositGas(tokenAddress, amountBigInt, incAddress, evmNetworkID)
	if err != nil {
		return nil, err
	}
	txFee, _ := getSynthesizedAmount(new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPriceBigInt), uint64(nativeTokenDecimals)).Float64()
	_, err = acc.checkSufficientBalance(common.HexToAddress(nativeToken), txFee, evmNetworkID)
	if err != nil {
		return nil, err
	}
	if askUser {
		yesNoPrompt(fmt.Sprintf("%v DepositAmount: %v, GasPrice: %v gWei, TxFee: %v. Do you want to continue?",
			prefix, depositedAmount, float64(gasPriceBigInt.Uint64())/math.Pow10(9), txFee))
	}

	auth, err := acc.newTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, 0, nil, evmNetworkID)
	if err != nil {
		return nil, err
	}

	// create the depositing transaction
	tx, err := v.DepositERC20(auth, tokenAddress, amountBigInt, incAddress)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()
	log.Printf("%v Deposited tx: %v\n", prefix, txHash)

	if err := wait(txHash, evmNetworkID); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// ApproveERC20 approves the Incognito Vault to spend an ERC20/BEP20 token of an account.
func (acc EVMAccount) ApproveERC20(tokenAddress, approved common.Address, approvedAmount float64, gasPrice uint64, evmNetworkID int) (*common.Hash, error) {
	prefix := "[ApproveERC20]"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		prefix = "[ApproveBEP20]"
	case rpc.PLGNetworkID:
		prefix = "[ApprovePLG20]"
	case rpc.FTMNetworkID:
		prefix = "[ApproveFTM20]"
	}
	evmClient, _ := getEVMClientAndVaultAddress(evmNetworkID)

	erc20Token, err := erc20.NewErc20(tokenAddress, evmClient)
	if err != nil {
		return nil, err
	}
	tokenDecimals, err := getDecimals(tokenAddress, evmNetworkID)
	if err != nil {
		return nil, err
	}

	// load the ERC20 ABI
	erc20ABI, err := abi.JSON(strings.NewReader(erc20.Erc20ABI))
	if err != nil {
		return nil, err
	}

	// estimate the gas limit
	approvedAmountBigInt := new(big.Int).SetUint64(uint64(approvedAmount * math.Pow10(int(tokenDecimals))))
	data, err := erc20ABI.Pack(
		"approve",
		approved,
		approvedAmountBigInt,
	)
	if err != nil {
		return nil, err
	}

	// estimate the gas limit and gas price
	gasPriceBigInt, gasLimit, err := acc.getGasLimitAndPrice(0, gasPrice, ethereum.CallMsg{To: &tokenAddress, Data: data}, evmNetworkID)
	if err != nil {
		return nil, err
	}
	txFee := getSynthesizedAmount(
		new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPriceBigInt),
		uint64(nativeTokenDecimals),
	)
	if askUser {
		yesNoPrompt(fmt.Sprintf("%v Approve %v to spend %v of token %v. Are you sure?",
			prefix, approved, approvedAmount, tokenAddress.String()))
		yesNoPrompt(fmt.Sprintf("%v GasPrice: %v gWei, TxFee: %v. Do you want to continue?",
			prefix, float64(gasPriceBigInt.Uint64())/math.Pow10(9), txFee.String()))
	} else {
		log.Printf("%v GasPrice: %v, GasLimit %v, TxFee %v\n", prefix, gasPriceBigInt.Uint64(), gasLimit, txFee.String())
	}

	auth, err := acc.newTransactionOpts(tokenAddress, gasPriceBigInt.Uint64(), gasLimit, 0, data, evmNetworkID)
	if err != nil {
		return nil, err
	}

	tx, err := erc20Token.Approve(auth, approved, approvedAmountBigInt)
	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()
	log.Printf("%v TxHash: %v\n", prefix, txHash.String())
	return &txHash, nil
}

// UnShield submits a burn proof of the given incTxHash to the smart contract to obtain back a public token.
func (acc EVMAccount) UnShield(incTxHash string, gasLimit, gasPrice uint64, evmNetworkID int) (*common.Hash, error) {
	prefix := "[UnShield]"

	evmClient, vaultAddress := getEVMClientAndVaultAddress(evmNetworkID)

	// load the vault instance
	v, err := vault.NewVault(vaultAddress, evmClient)
	if err != nil {
		return nil, err
	}

	// retrieve the burn proof from the incTxHash
	burnProofResult, err := cfg.incClient.GetBurnProof(incTxHash, evmNetworkID)
	if err != nil {
		return nil, err
	}
	burnProof, err := incclient.DecodeBurnProof(burnProofResult)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}
	if gasLimit == 0 {
		gasLimit, err = acc.estimateWithdrawalGas(burnProof, evmNetworkID)
		if err != nil {
			return nil, err
		}
	}
	txFee, _ := getSynthesizedAmount(new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPriceBigInt), uint64(nativeTokenDecimals)).Float64()
	_, err = acc.checkSufficientBalance(common.HexToAddress(nativeToken), txFee, evmNetworkID)
	if err != nil {
		return nil, err
	}
	if askUser {
		yesNoPrompt(fmt.Sprintf("%v GasPrice: %v gWei, TxFee: %v. Do you want to continue?",
			prefix, float64(gasPriceBigInt.Uint64())/math.Pow10(9), txFee))
	}

	auth, err := acc.newTransactionOpts(vaultAddress, gasPriceBigInt.Uint64(), gasLimit, 0, []byte{}, evmNetworkID)
	tx, err := v.Withdraw(auth,
		burnProof.Instruction,
		burnProof.Heights[0],
		burnProof.InstPaths[0],
		burnProof.InstPathIsLefts[0],
		burnProof.InstRoots[0],
		burnProof.BlkData[0],
		burnProof.SigIndices[0],
		burnProof.SigVs[0],
		burnProof.SigRs[0],
		burnProof.SigSs[0])
	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()
	log.Printf("%v WithdrawalTx: %v\n", prefix, txHash.String())

	if err := wait(txHash, evmNetworkID); err != nil {
		return nil, err
	}
	return &txHash, nil
}

// Shield shields an amount of EVM tokens to the Incognito network.
// This function should be called after the DepositNative or DepositToken has finished.
func Shield(privateKey, pTokenID string, evmTxHashStr string, evmNetworkID int) (string, error) {
	prefix := "[Shield]"

	evmClient := cfg.evmClients[evmNetworkID]
	if evmClient == nil {
		return "", errEVMNetworkNotSupported(evmNetworkID)
	}

	evmTxHash := common.HexToHash(evmTxHashStr)
	receipt, err := evmClient.TransactionReceipt(context.Background(), evmTxHash)
	if err != nil {
		return "", err
	}
	blockNumber := receipt.BlockNumber.Uint64()
	log.Printf("%v ShieldedBlock: %v\n", prefix, blockNumber)

	numCfms := 15
	if evmNetworkID == rpc.PLGNetworkID {
		numCfms = 35
	} else if evmNetworkID == rpc.FTMNetworkID {
		numCfms = 5
	}
	log.Printf("%v Wait for %v confirmations\n", prefix, numCfms)

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	for {
		header, err := evmClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return "", err
		}

		if header.Number.Uint64() > blockNumber+uint64(numCfms) {
			log.Println(prefix, "Enough confirmations!!")
			break
		}
		log.Printf("%v CurrentEVMBlock: %v\n", prefix, header.Number.Uint64())
		time.Sleep(30 * time.Second)
	}

	depositProof, _, err := cfg.incClient.GetEVMDepositProof(evmTxHash.String(), evmNetworkID)
	if err != nil {
		return "", err
	}

	encodedTx, incTxHash, err := cfg.incClient.CreateIssuingEVMRequestTransaction(privateKey, pTokenID, *depositProof, evmNetworkID)
	if err != nil {
		return "", err
	}

	tx := new(tx_ver2.Tx)
	rawTxData, _, err := base58.Base58Check{}.Decode(string(encodedTx))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(rawTxData, &tx)
	if err != nil {
		return "", err
	}

	md := tx.GetMetadata().(*metadata.IssuingEVMRequest)
	_, err = verifyProofAndParseReceipt(md, evmNetworkID)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " Verify proof locally SUCCEEDED!!!")

	err = cfg.incClient.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " SendRawTx SUCCEEDED!!")
	log.Printf("%v ShieldedTx: %v\n", prefix, incTxHash)

	return incTxHash, nil
}

// ====================== PRV FUNCTIONS ======================

// estimateDepositGas estimates the gas for depositing a token.
func (acc EVMAccount) estimatePRVDepositGas(depositedAmount *big.Int, incAddress string, evmNetworkID int) (uint64, error) {
	if evmNetworkID == rpc.PLGNetworkID || evmNetworkID == rpc.FTMNetworkID {
		return 0, errEVMNetworkNotSupported(evmNetworkID)
	}

	evmClient := cfg.evmClients[evmNetworkID]
	prv20Address := common.HexToAddress(prv20AddressStr)

	var gasLimit uint64
	prvABI, err := abi.JSON(strings.NewReader(prv20.Prv20ABI))
	if err != nil {
		return 0, fmt.Errorf("cannot create prvABI from file")
	}

	var data []byte
	data, err = prvABI.Pack(
		"burn",
		incAddress,
		depositedAmount,
	)
	if err != nil {
		return 0, err
	}

	gasLimit, err = evmClient.EstimateGas(context.Background(),
		ethereum.CallMsg{From: acc.address, To: &prv20Address, Data: data})
	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}

// estimatePRVWithdrawalGas estimates the gas for withdrawing PRV.
func (acc EVMAccount) estimatePRVWithdrawalGas(burnProof *incclient.BurnProof, evmNetworkID int) (uint64, error) {
	if evmNetworkID == rpc.PLGNetworkID || evmNetworkID == rpc.FTMNetworkID {
		return 0, errEVMNetworkNotSupported(evmNetworkID)
	}

	evmClient := cfg.evmClients[evmNetworkID]
	prv20Address := common.HexToAddress(prv20AddressStr)

	prvABI, err := abi.JSON(strings.NewReader(prv20.Prv20ABI))
	if err != nil {
		return 0, fmt.Errorf("cannot create vaultABI from file")
	}

	var data []byte
	data, err = prvABI.Pack(
		"mint",
		burnProof.Instruction,
		burnProof.Heights[0],
		burnProof.InstPaths[0],
		burnProof.InstPathIsLefts[0],
		burnProof.InstRoots[0],
		burnProof.BlkData[0],
		burnProof.SigIndices[0],
		burnProof.SigVs[0],
		burnProof.SigRs[0],
		burnProof.SigSs[0],
	)
	if err != nil {
		return 0, err
	}

	gasLimit, err := evmClient.EstimateGas(context.Background(),
		ethereum.CallMsg{From: acc.address, To: &prv20Address, Data: data})
	if err != nil {
		return 0, fmt.Errorf("estimateGas for withdrawal error: %v", err)
	}

	return gasLimit, nil
}

// BurnPRVOnEVM burns an amount of PRV on an EVM network to shield to the Incognito network.
func (acc EVMAccount) BurnPRVOnEVM(incAddress string, depositedAmount float64, gasLimit, gasPrice uint64, evmNetworkID int) (*common.Hash, error) {
	if evmNetworkID == rpc.PLGNetworkID || evmNetworkID == rpc.FTMNetworkID {
		return nil, errEVMNetworkNotSupported(evmNetworkID)
	}

	evmClient := cfg.evmClients[evmNetworkID]
	prefix := "[DepositPRVERC20]"
	switch evmNetworkID {
	case rpc.BSCNetworkID:
		prefix = "[DepositPRVBEP20]"
	}

	prv20Address := common.HexToAddress(prv20AddressStr)

	// load the vault instance
	prvEVM20, err := prv20.NewPrv20(prv20Address, evmClient)
	if err != nil {
		return nil, err
	}

	_, err = acc.checkSufficientBalance(prv20Address, depositedAmount, evmNetworkID)
	if err != nil {
		return nil, err
	}

	tokenDecimals, err := getDecimals(prv20Address, evmNetworkID)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}

	// estimate gasLimit
	amountBigInt := new(big.Int).SetUint64(uint64(depositedAmount * math.Pow10(int(tokenDecimals))))
	gasLimit, err = acc.estimatePRVDepositGas(amountBigInt, incAddress, evmNetworkID)
	if err != nil {
		return nil, err
	}

	txFee, _ := getSynthesizedAmount(new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPriceBigInt), uint64(nativeTokenDecimals)).Float64()
	_, err = acc.checkSufficientBalance(common.HexToAddress(nativeToken), txFee, evmNetworkID)
	if err != nil {
		return nil, err
	}
	if askUser {
		yesNoPrompt(fmt.Sprintf("%v DepositAmount: %v, GasPrice: %v gWei, TxFee: %v. Do you want to continue?",
			prefix, depositedAmount, float64(gasPriceBigInt.Uint64())/math.Pow10(9), txFee))
	}

	auth, err := acc.newTransactionOpts(prv20Address, gasPriceBigInt.Uint64(), gasLimit, 0, nil, evmNetworkID)
	if err != nil {
		return nil, err
	}

	// create the depositing transaction
	tx, err := prvEVM20.Burn(auth, incAddress, amountBigInt)
	if err != nil {
		return nil, err
	}
	txHash := tx.Hash()
	log.Printf("%v Deposited tx: %v\n", prefix, txHash)

	if err := wait(txHash, evmNetworkID); err != nil {
		return nil, err
	}

	return &txHash, nil
}

// ShieldPRV shields an amount of PRV tokens from EVM chains to the Incognito network.
// This function should be called after the BurnPRVOnEVM has finished.
func ShieldPRV(privateKey, evmTxHashStr string, evmNetworkID int) (string, error) {
	if evmNetworkID == rpc.PLGNetworkID || evmNetworkID == rpc.FTMNetworkID {
		return "", errEVMNetworkNotSupported(evmNetworkID)
	}

	prefix := "[ShieldPRV]"

	evmClient := cfg.evmClients[evmNetworkID]
	if evmClient == nil {
		return "", errEVMNetworkNotSupported(evmNetworkID)
	}

	evmTxHash := common.HexToHash(evmTxHashStr)
	receipt, err := evmClient.TransactionReceipt(context.Background(), evmTxHash)
	if err != nil {
		return "", err
	}
	blockNumber := receipt.BlockNumber.Uint64()
	log.Printf("%v ShieldedBlock: %v\n", prefix, blockNumber)

	numCfms := 15
	log.Printf("%v Wait for %v confirmations\n", prefix, numCfms)

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	for {
		header, err := evmClient.HeaderByNumber(ctx, nil)
		if err != nil {
			return "", err
		}

		if header.Number.Uint64() > blockNumber+uint64(numCfms) {
			log.Println(prefix, "Enough confirmations!!")
			break
		}
		log.Printf("%v CurrentEVMBlock: %v\n", prefix, header.Number.Uint64())
		time.Sleep(30 * time.Second)
	}

	depositProof, _, err := cfg.incClient.GetEVMDepositProof(evmTxHash.String(), evmNetworkID)
	if err != nil {
		return "", err
	}

	encodedTx, incTxHash, err := cfg.incClient.CreateIssuingPRVPeggingRequestTransaction(privateKey, *depositProof, evmNetworkID)
	if err != nil {
		return "", err
	}

	tx := new(tx_ver2.Tx)
	rawTxData, _, err := base58.Base58Check{}.Decode(string(encodedTx))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(rawTxData, &tx)
	if err != nil {
		return "", err
	}

	md := tx.GetMetadata().(*metadata.IssuingEVMRequest)
	_, err = verifyProofAndParseReceipt(md, evmNetworkID)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " Verify proof locally SUCCEEDED!!!")

	err = cfg.incClient.SendRawTx(encodedTx)
	if err != nil {
		return "", err
	}
	log.Println(prefix + " SendRawTx SUCCEEDED!!")
	log.Printf("%v ShieldedTx: %v\n", prefix, incTxHash)

	return incTxHash, nil
}

// UnShieldPRV submits a PRV burn proof of the given incTxHash to the smart contract to obtain PRV on a public EVM chain.
func (acc EVMAccount) UnShieldPRV(incTxHash string, gasLimit, gasPrice uint64, evmNetworkID int) (*common.Hash, error) {
	if evmNetworkID == rpc.PLGNetworkID || evmNetworkID == rpc.FTMNetworkID {
		return nil, errEVMNetworkNotSupported(evmNetworkID)
	}

	prefix := "[UnShieldPRV]"

	evmClient := cfg.evmClients[evmNetworkID]
	if evmClient == nil {
		return nil, errEVMNetworkNotSupported(evmNetworkID)
	}
	prvToken, err := prv20.NewPrv20(common.HexToAddress(prv20AddressStr), evmClient)
	if err != nil {
		return nil, err
	}

	// retrieve the burn proof from the incTxHash
	burnProofResult, err := cfg.incClient.GetBurnPRVPeggingProof(incTxHash, evmNetworkID)
	if err != nil {
		return nil, err
	}
	burnProof, err := incclient.DecodeBurnProof(burnProofResult)
	if err != nil {
		return nil, err
	}

	// calculate gas price
	var gasPriceBigInt *big.Int
	if gasPrice == 0 {
		gasPriceBigInt, err = estimateGasPrice(evmNetworkID)
		if err != nil {
			return nil, fmt.Errorf("cannot get gasPriceBigInt price")
		}
	} else {
		gasPriceBigInt = new(big.Int).SetUint64(gasPrice)
	}
	if gasLimit == 0 {
		gasLimit, err = acc.estimatePRVWithdrawalGas(burnProof, evmNetworkID)
		if err != nil {
			return nil, err
		}
	}
	txFee, _ := getSynthesizedAmount(new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPriceBigInt), uint64(nativeTokenDecimals)).Float64()
	_, err = acc.checkSufficientBalance(common.HexToAddress(nativeToken), txFee, evmNetworkID)
	if err != nil {
		return nil, err
	}
	if askUser {
		yesNoPrompt(fmt.Sprintf("%v GasPrice: %v gWei, TxFee: %v. Do you want to continue?",
			prefix, float64(gasPriceBigInt.Uint64())/math.Pow10(9), txFee))
	}

	auth, err := acc.newTransactionOpts(common.HexToAddress(prv20AddressStr), gasPriceBigInt.Uint64(), gasLimit, 0, []byte{}, evmNetworkID)
	tx, err := prvToken.Mint(auth,
		burnProof.Instruction,
		burnProof.Heights[0],
		burnProof.InstPaths[0],
		burnProof.InstPathIsLefts[0],
		burnProof.InstRoots[0],
		burnProof.BlkData[0],
		burnProof.SigIndices[0],
		burnProof.SigVs[0],
		burnProof.SigRs[0],
		burnProof.SigSs[0])
	if err != nil {
		return nil, err
	}

	txHash := tx.Hash()
	log.Printf("%v WithdrawalTx: %v\n", prefix, txHash.String())

	if err := wait(txHash, evmNetworkID); err != nil {
		return nil, err
	}
	return &txHash, nil
}
