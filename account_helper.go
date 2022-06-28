package main

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
)

var listTokenInfo map[string]TokenInfo

type TokenInfo struct {
	TokenID   string `json:"TokenID"`
	Name      string `json:"Name"`
	Symbol    string `json:"Symbol"`
	PSymbol   string `json:"PSymbol"`
	IsBridge  bool   `json:"IsBridge"`
	Verified  bool   `json:"Verified"`
	PDecimals int    `json:"PDecimals"`
	Network   string `json:"Network"`
}

func initForFinancialReport(c *cli.Context) error {
	err := initNetWork()
	if err != nil {
		return err
	}

	if network != "mainnet" {
		return nil
	}

	listTokenInfo, err = getAllTokenInfo()
	if err != nil {
		listTokenInfo = nil
	}

	return nil
}

func getAllTokenInfo() (map[string]TokenInfo, error) {
	url := "https://api-coinservice.incognito.org/coins/tokenlist?all=true"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type TmpRes struct {
		Result []TokenInfo `json:"Result"`
		Error  string      `json:"Error"`
	}
	var tmpRes TmpRes
	err = json.Unmarshal(body, &tmpRes)
	if err != nil {
		return nil, err
	}
	if tmpRes.Error != "" {
		return nil, fmt.Errorf("retrieving token info encountered an error: %v", tmpRes.Error)
	}

	res := make(map[string]TokenInfo)
	for _, tokenInfo := range tmpRes.Result {
		res[tokenInfo.TokenID] = tokenInfo
	}

	return res, nil
}

func getTokenName(tokenID string) string {
	if tokenID == common.PRVIDStr {
		return "PRV"
	}

	if listTokenInfo == nil {
		return tokenID
	}

	if tokenInfo, ok := listTokenInfo[tokenID]; ok {
		if !tokenInfo.Verified {
			return tokenID
		}
		res := tokenInfo.Symbol
		if tokenInfo.Network != "" {
			res = fmt.Sprintf("%v (%v)", res, tokenInfo.Network)
		}

		return res
	}

	return tokenID
}

func getTokenDecimals(tokenID string) int {
	if tokenID == common.PRVIDStr {
		return 9
	}

	if listTokenInfo == nil {
		return 0
	}

	if tokenInfo, ok := listTokenInfo[tokenID]; ok {
		if !tokenInfo.Verified {
			return 0
		}

		return tokenInfo.PDecimals
	}

	return 0
}
