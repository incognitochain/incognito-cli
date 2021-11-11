package pdex_v3

import (
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	"github.com/incognitochain/incognito-cli/pdex_v3/v2utils"
	"math/big"
)

func chooseBestPoolFromAPair(
	poolPairStates map[string]*jsonresult.Pdexv3PoolPairState,
	tokenIDStrNodeSource string,
	tokenIDStrNodeDest string,
	sellAmt uint64,
) (*jsonresult.Pdexv3PoolPair, string, uint64) {
	maxReceive := uint64(0)
	var chosenPool *jsonresult.Pdexv3PoolPair
	var chosenPoolID string
	for poolID, poolState := range poolPairStates {
		pool := poolState.State
		if (tokenIDStrNodeSource == pool.Token0ID.String() && tokenIDStrNodeDest == pool.Token1ID.String()) ||
			(tokenIDStrNodeSource == pool.Token1ID.String() && tokenIDStrNodeDest == pool.Token0ID.String()) {
			receive := trade(
				poolID,
				poolPairStates,
				tokenIDStrNodeDest,
				sellAmt,
			)
			if receive > maxReceive {
				maxReceive = receive
				chosenPool = &pool
				chosenPoolID = poolID
			}
		}
	}
	return chosenPool, chosenPoolID, maxReceive
}

func trade(
	poolID string,
	poolPairStates map[string]*jsonresult.Pdexv3PoolPairState,
	tokenIDToBuyStr string,
	sellAmount uint64,
) uint64 {
	pool := poolPairStates[poolID].State
	tokenIDToBuy := pool.Token0ID
	tokenIDToSell := pool.Token1ID
	if tokenIDToBuyStr != tokenIDToBuy.String() {
		tokenIDToBuy = pool.Token1ID
		tokenIDToSell = pool.Token0ID
	}

	var tradePath []string
	tradePath = append(tradePath, poolID)

	// get relevant, cloned data from state for the trade path
	reserves, orderBookList, tradeDirections, err :=
		TradePathFromState(tokenIDToSell, tradePath, poolPairStates)

	expectedReceived, err := v2utils.EstimateReceivingAmount(
		sellAmount,
		0,
		reserves,
		tradeDirections,
		0,
		orderBookList,
	)

	if err != nil {
		return 0
	}
	return expectedReceived
}

// FindGoodTradePath attempts to find a good enough trading path for the given trading pair, selling amount, and pool pairs.
func FindGoodTradePath(
	maxPathLen uint,
	poolPairStates map[string]*jsonresult.Pdexv3PoolPairState,
	tokenIDStrSource string,
	tokenIDStrDest string,
	originalSellAmount uint64,
) ([]*jsonresult.Pdexv3PoolPair, []string, uint64) {

	pc := &PriceCalculator{
		Graph: make(map[string][]Node),
	}

	simplePools := make([]*SimplePoolNodeData, 0)

	for _, poolState := range poolPairStates {
		pool := poolState.State
		token0Liq := new(big.Int).Mul(pool.Token0VirtualAmount, big.NewInt(int64(BaseAmplifier)))
		token0Liq.Div(token0Liq, new(big.Int).SetUint64(uint64(pool.Amplifier)))
		token1Liq := new(big.Int).Mul(pool.Token1VirtualAmount, big.NewInt(int64(BaseAmplifier)))
		token1Liq.Div(token1Liq, new(big.Int).SetUint64(uint64(pool.Amplifier)))

		simplePools = append(simplePools, &SimplePoolNodeData{
			Token0ID:  pool.Token0ID.String(),
			Token1ID:  pool.Token1ID.String(),
			Token0Liq: token0Liq,
			Token1Liq: token1Liq,
		})
	}

	allPaths := pc.findPaths(maxPathLen+1, simplePools, tokenIDStrSource, tokenIDStrDest)

	if len(allPaths) == 0 {
		return nil, nil, 0
	}

	maxReceive := uint64(0)
	var chosenPairs []*jsonresult.Pdexv3PoolPair
	var chosenPath []string
	for _, path := range allPaths {
		sellAmt := originalSellAmount

		var poolsByPath []*jsonresult.Pdexv3PoolPair
		var poolIDsByPath []string

		for i := 0; i < len(path)-1; i++ {
			tokenIDStrNodeSource := path[i]
			tokenIDStrNodeDest := path[i+1]

			pool, poolID, receive := chooseBestPoolFromAPair(poolPairStates, tokenIDStrNodeSource, tokenIDStrNodeDest, sellAmt)
			sellAmt = receive
			poolsByPath = append(poolsByPath, pool)
			poolIDsByPath = append(poolIDsByPath, poolID)
		}

		if len(poolsByPath) == 0 || sellAmt > maxReceive {
			maxReceive = sellAmt
			chosenPairs = poolsByPath
			chosenPath = poolIDsByPath
		}
	}

	return chosenPairs, chosenPath, maxReceive
}
