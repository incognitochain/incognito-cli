package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/coin"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/common/base58"
	"github.com/incognitochain/go-incognito-sdk-v2/incclient"
	"github.com/incognitochain/go-incognito-sdk-v2/key"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/rpc"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/urfave/cli/v2"
	"strings"
)

func checkBalance(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if privateKey == "" {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if tokenIDStr == "" {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	balance, err := cfg.incClient.GetBalance(privateKey, tokenIDStr)
	if err != nil {
		return err
	}
	fmt.Println(balance)

	return nil
}

func getAllBalanceV2(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if privateKey == "" {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	balances, err := cfg.incClient.GetAllBalancesV2(privateKey)
	if err != nil {
		return err
	}
	jsb, err := json.MarshalIndent(balances, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsb))

	return nil
}

func keyInfo(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if privateKey == "" {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
	if err != nil {
		return err
	}

	jsb, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return fmt.Errorf("marshalling key info error: %v", err)
	}
	fmt.Println(string(jsb))

	return nil
}

func consolidateUTXOs(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if privateKey == "" {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if tokenIDStr == "" {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	version := c.Int(versionFlag)
	if version < 1 || version > 2 {
		return fmt.Errorf("%v is invalid", versionFlag)
	}

	numThreads := c.Int(numThreadsFlag)
	if numThreads == 0 {
		return fmt.Errorf("%v in invalid", numThreadsFlag)
	}

	fmt.Printf("CONSOLIDATING tokenID %v, version %v, numThreads %v\n", tokenIDStr, version, numThreads)

	txList, err := cfg.incClient.Consolidate(privateKey, tokenIDStr, int8(version), numThreads)
	if err != nil {
		return err
	}
	fmt.Println("CONSOLIDATING FINISHED!!")
	fmt.Println(txList)

	return nil
}

func checkUTXOs(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if privateKey == "" {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if tokenIDStr == "" {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	unSpentCoins, idxList, err := cfg.incClient.GetUnspentOutputCoins(privateKey, tokenIDStr, 0)
	if err != nil {
		return err
	}

	numUTXOsV1 := 0
	numUTXOsV2 := 0
	balanceV1 := uint64(0)
	balanceV2 := uint64(0)

	for i, utxo := range unSpentCoins {
		if utxo.GetVersion() == 1 {
			numUTXOsV1++
			balanceV1 += utxo.GetValue()
		} else {
			numUTXOsV2++
			balanceV2 += utxo.GetValue()
		}

		fmt.Printf("idx %v, version %v, pubKey %v, keyImage %v, value %v\n",
			idxList[i].Uint64(), utxo.GetVersion(),
			base58.Base58Check{}.Encode(utxo.GetPublicKey().ToBytesS(), 0),
			base58.Base58Check{}.Encode(utxo.GetKeyImage().ToBytesS(), 0),
			utxo.GetValue())
	}

	fmt.Printf("#numUTXOsV1 %v, #numUTXOsV2 %v\n", numUTXOsV1, numUTXOsV2)
	fmt.Printf("balanceV1 %v, balanceV2 %v, totalBalance %v\n", balanceV1, balanceV2, balanceV1+balanceV2)

	return nil
}

func getOutCoins(c *cli.Context) error {
	address := c.String(addressFlag)
	if !isValidAddress(address) {
		return fmt.Errorf("%v is invalid", addressFlag)
	}

	otaKey := c.String(otaKeyFlag)
	if !isValidOtaKey(otaKey) {
		return fmt.Errorf("%v is invalid", otaKeyFlag)
	}

	readonlyKey := c.String(readonlyKeyFlag)
	if readonlyKey != "" && !isValidReadonlyKey(readonlyKey) {
		return fmt.Errorf("%v is invalid", readonlyKeyFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	outCoinKey := new(rpc.OutCoinKey)
	outCoinKey.SetPaymentAddress(address)
	outCoinKey.SetOTAKey(otaKey)
	outCoinKey.SetReadonlyKey(readonlyKey)

	outCoins, idxList, err := cfg.incClient.GetOutputCoins(outCoinKey, tokenIDStr, 0)
	if err != nil {
		return err
	}

	v1Count := 0
	v2Count := 0
	for i, outCoin := range outCoins {
		if outCoin.GetVersion() == 1 {
			v1Count += 1
		} else {
			v2Count += 1
		}

		fmt.Printf("idx %v, ver %v, encrypted %v, pubKey %v, cmtStr %v\n",
			idxList[i].Int64(),
			outCoin.GetVersion(),
			outCoin.IsEncrypted(),
			base58.Base58Check{}.Encode(outCoin.GetPublicKey().ToBytesS(), 0x00),
			base58.Base58Check{}.Encode(outCoin.GetCommitment().ToBytesS(), 0x00))
	}

	fmt.Printf("#OutCoins: %v, #v1: %v, #v2: %v\n", len(outCoins), v1Count, v2Count)

	return nil
}

func getHistory(c *cli.Context) error {
	privateKey := c.String(privateKeyFlag)
	if !isValidPrivateKey(privateKey) {
		return fmt.Errorf("%v is invalid", privateKeyFlag)
	}

	tokenIDStr := c.String(tokenIDFlag)
	if !isValidTokenID(tokenIDStr) {
		return fmt.Errorf("%v is invalid", tokenIDFlag)
	}

	numThreads := c.Int(numThreadsFlag)
	if numThreads == 0 {
		return fmt.Errorf("%v in invalid", numThreadsFlag)
	}

	csvFile := c.String("csvFile")

	historyProcessor := incclient.NewTxHistoryProcessor(cfg.incClient, numThreads)

	h, err := historyProcessor.GetTokenHistory(privateKey, tokenIDStr)
	if err != nil {
		return err
	}

	if len(csvFile) > 0 {
		err = incclient.SaveTxHistory(h, csvFile)
		if err != nil {
			return err
		}
	} else {
		totalIn := uint64(0)
		fmt.Printf("#TxIns %v\n", len(h.TxInList))
		for _, txIn := range h.TxInList {
			totalIn += txIn.GetAmount()
			fmt.Println(txIn.String())
		}
		fmt.Printf("END TxIns\n\n")

		totalOut := uint64(0)
		fmt.Printf("#TxOuts %v\n", len(h.TxOutList))
		for _, txOut := range h.TxOutList {
			totalOut += txOut.GetAmount()
			fmt.Println(txOut.String())
		}
		fmt.Printf("END TxOuts\n")

		fmt.Printf("TotalIn: %v, TotalOut: %v\n", totalIn, totalOut)
	}

	return nil
}

func genKeySet(c *cli.Context) error {
	w, mnemonic, err := wallet.NewMasterKey()
	if err != nil {
		return err
	}

	numShards := c.Int(numShardsFlag)
	if numShards == 0 {
		return fmt.Errorf("%v is invalid", numShardsFlag)
	}
	common.MaxShardNumber = numShards

	numAccounts := c.Int(numAccountsFlag)

	fmt.Printf("mnemonic: %v\n", mnemonic)
	accounts := make([]*incclient.KeyInfo, 0)
	for i := 0; i < numAccounts; i++ {
		childKey, err := w.DeriveChild(uint32(i))
		if err != nil {
			return err
		}
		privateKey := childKey.Base58CheckSerialize(wallet.PrivateKeyType)
		info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
		if err != nil {
			return err
		}

		accounts = append(accounts, info)
	}
	return jsonPrint(accounts)
}

func importMnemonic(c *cli.Context) error {
	mnemonic := c.String(mnemonicFlag)
	mnemonic = strings.Replace(mnemonic, "-", " ", -1)
	w, err := wallet.NewMasterKeyFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	numShards := c.Int(numShardsFlag)
	if numShards == 0 {
		return fmt.Errorf("%v is invalid", numShardsFlag)
	}
	common.MaxShardNumber = numShards

	numAccounts := c.Int(numAccountsFlag)

	fmt.Printf("mnemonic: %v\n", mnemonic)
	accounts := make([]*incclient.KeyInfo, 0)
	for i := 0; i < numAccounts; i++ {
		childKey, err := w.DeriveChild(uint32(i))
		if err != nil {
			return err
		}
		privateKey := childKey.Base58CheckSerialize(wallet.PrivateKeyType)
		info, err := incclient.GetAccountInfoFromPrivateKey(privateKey)
		if err != nil {
			return err
		}

		accounts = append(accounts, info)
	}
	err = jsonPrint(accounts)

	return nil
}

func submitKey(c *cli.Context) error {
	var err error
	otaKey := c.String(otaKeyFlag)
	if otaKey == "" {
		return fmt.Errorf("%v is invalid", otaKeyFlag)
	}

	accessToken := c.String(accessTokenFlag)
	if accessToken != "" {
		fromHeight := c.Uint64(fromHeightFlag)
		isReset := c.Bool(isResetFlag)
		err = cfg.incClient.AuthorizedSubmitKey(otaKey, accessToken, fromHeight, isReset)
	} else {
		err = cfg.incClient.SubmitKey(otaKey)
	}

	if err != nil {
		return err
	}

	return nil
}

func newOTAReceiver(c *cli.Context) error {
	addr := c.String(addressFlag)
	w, err := wallet.Base58CheckDeserialize(addr)
	if err != nil || len(w.KeySet.PaymentAddress.Pk) == 0 {
		return fmt.Errorf("%v is invalid", addr)
	}

	otaReceiver := new(coin.OTAReceiver)
	err = otaReceiver.FromAddress(w.KeySet.PaymentAddress)
	if err != nil {
		return err
	}

	type Res struct {
		Receiver string
	}

	return jsonPrint(Res{Receiver: otaReceiver.String()})
}

// signDepositOTAReceiver ...
func signDepositOTAReceiver(c *cli.Context) error {
	depositPrivateKeyStr := c.String(depositPrivateKeyFlag)
	if depositPrivateKeyStr == "" {
		return fmt.Errorf("%v is invalid", depositPrivateKeyFlag)
	}
	depositPrivateKey, _, err := base58.Base58Check{}.Decode(depositPrivateKeyStr)
	if err != nil {
		return fmt.Errorf("decode %v error: %v", depositPrivateKeyFlag, err)
	}
	depositKey := key.OTDepositKey{
		PrivateKey: depositPrivateKey,
	}

	otaReceiverStr := c.String(receiverFlag)
	if !isValidOTAReceiver(otaReceiverStr) {
		return fmt.Errorf("%v is invalid", receiverFlag)
	}
	otaReceiver := new(coin.OTAReceiver)
	err = otaReceiver.FromString(otaReceiverStr)
	if err != nil {
		return fmt.Errorf("decode OTAReceiver error %v: %v", otaReceiverStr, err)
	}

	sig, err := incclient.SignDepositData(&depositKey, otaReceiver.Bytes())
	if err != nil {
		return fmt.Errorf("signDepositData error: %v", err)
	}

	type Result struct {
		Signature string
	}

	return jsonPrint(Result{Signature: base58.Base58Check{}.NewEncode(sig, 0)})
}
