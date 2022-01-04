package pdex_v3

import (
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	v2 "github.com/incognitochain/incognito-cli/pdex_v3/v2utils"
	"log"
	"math/big"
	"sort"
)

type Node struct {
	TokenIDStr     string
	TokenPoolValue *big.Int
}

type SimplePoolNodeData struct {
	Token0ID  string
	Token1ID  string
	Token0Liq *big.Int
	Token1Liq *big.Int
}

type PriceCalculator struct {
	Graph map[string][]Node
}

func (pc *PriceCalculator) findPaths(
	maxPathLen uint,
	simplePools []*SimplePoolNodeData,
	tokenIDStrSource string,
	tokenIDStrDest string,
) [][]string {
	pc.buildGraph(simplePools)

	visited := make(map[string]bool)
	for tokenIDStr := range pc.Graph {
		visited[tokenIDStr] = false
	}

	var path []string
	var allPaths [][]string
	pc.findPath(
		maxPathLen,
		tokenIDStrSource,
		tokenIDStrDest,
		visited,
		path,
		&allPaths,
	)

	return allPaths
}

func (pc *PriceCalculator) findPath(
	maxPathLen uint,
	tokenIDStrSource string,
	tokenIDStrDest string,
	visited map[string]bool,
	path []string,
	allPaths *[][]string,
) {
	if len(*allPaths) == MaxPaths {
		log.Println("MaxPaths exceeded")
		return
	}
	path = append(path, tokenIDStrSource)
	visited[tokenIDStrSource] = true

	if tokenIDStrSource == tokenIDStrDest {
		// Beware, we need to make deep copy of path array
		newPath := make([]string, len(path))
		copy(newPath, path)
		*allPaths = append(*allPaths, newPath)
	} else if len(path) < int(maxPathLen) {
		nodes, found := pc.Graph[tokenIDStrSource]
		if found {
			for _, node := range nodes {
				if visited[node.TokenIDStr] {
					continue
				}
				pc.findPath(maxPathLen, node.TokenIDStr, tokenIDStrDest, visited, path, allPaths)
			}
		}
	}
	path = path[:len(path)-1]
	visited[tokenIDStrSource] = false
}

// NOTEs: the built graph would be undirected graph
func (pc *PriceCalculator) buildGraph(
	simplePools []*SimplePoolNodeData,
) {
	pc.Graph = make(map[string][]Node)
	for _, pool := range simplePools {
		addEdge(
			pool.Token0ID,
			pool.Token1ID,
			pool.Token0Liq,
			pc.Graph,
		)
		addEdge(
			pool.Token1ID,
			pool.Token0ID,
			pool.Token1Liq,
			pc.Graph,
		)
	}

	// sort by descending order of liquidity
	for _, nodes := range pc.Graph {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].TokenPoolValue.Cmp(nodes[j].TokenPoolValue) > 0
		})
	}
}

func addEdge(
	tokenIDStrSource string,
	tokenIDStrDest string,
	tokenLiqSource *big.Int,
	graph map[string][]Node,
) {
	dest := Node{
		TokenIDStr:     tokenIDStrDest,
		TokenPoolValue: tokenLiqSource,
	}
	_, found := graph[tokenIDStrSource]
	if !found {
		graph[tokenIDStrSource] = []Node{dest}
	} else {
		isExisted := false
		for _, existedDest := range graph[tokenIDStrSource] {
			if existedDest.TokenIDStr == dest.TokenIDStr {
				if existedDest.TokenPoolValue.Cmp(dest.TokenPoolValue) < 0 {
					*existedDest.TokenPoolValue = *dest.TokenPoolValue
				}
				isExisted = true
				break
			}
		}
		if !isExisted {
			graph[tokenIDStrSource] = append(graph[tokenIDStrSource], dest)
		}
	}
}

func clonePoolPairState(p *jsonresult.Pdexv3PoolPairState) *jsonresult.Pdexv3PoolPairState {
	res := &jsonresult.Pdexv3PoolPairState{}

	tmpPoolPair := jsonresult.Pdexv3PoolPair{
		ShareAmount:      p.State.ShareAmount,
		Token0ID:         p.State.Token0ID,
		Token1ID:         p.State.Token1ID,
		Token0RealAmount: p.State.Token0RealAmount,
		Token1RealAmount: p.State.Token1RealAmount,
		Amplifier:        p.State.Amplifier,
	}
	tmpPoolPair.Token0VirtualAmount = new(big.Int).Set(p.State.Token0VirtualAmount)
	tmpPoolPair.Token1VirtualAmount = new(big.Int).Set(p.State.Token1VirtualAmount)
	res.State = tmpPoolPair
	res.Shares = make(map[string]*jsonresult.Pdexv3Share)
	for k, v := range p.Shares {
		tmpRes := &jsonresult.Pdexv3Share{}
		tmpRes.Amount = v.Amount
		tmpRes.TradingFees = map[common.Hash]uint64{}
		for key, value := range v.TradingFees {
			tmpRes.TradingFees[key] = value
		}
		tmpRes.LastLPFeesPerShare = map[common.Hash]*big.Int{}
		for key, value := range v.LastLPFeesPerShare {
			tmpRes.LastLPFeesPerShare[key] = new(big.Int).Set(value)
		}
		res.Shares[k] = tmpRes
	}

	res.LpFeesPerShare = make(map[common.Hash]*big.Int)
	for k, v := range p.LpFeesPerShare {
		res.LpFeesPerShare[k] = big.NewInt(0).Set(v)
	}

	res.ProtocolFees = make(map[common.Hash]uint64)
	for k, v := range p.ProtocolFees {
		res.ProtocolFees[k] = v
	}

	res.StakingPoolFees = make(map[common.Hash]uint64)
	for k, v := range p.StakingPoolFees {
		res.StakingPoolFees[k] = v
	}

	res.Orderbook = jsonresult.Pdexv3Orderbook{}
	res.Orderbook.Orders = make([]*jsonresult.Pdexv3Order, len(p.Orderbook.Orders))
	for i, order := range p.Orderbook.Orders {
		tmp := *order
		res.Orderbook.Orders[i] = &tmp
	}

	return res
}

func TradePathFromState(
	sellToken common.Hash,
	tradePath []string,
	pairs map[string]*jsonresult.Pdexv3PoolPairState,
) (
	[]*jsonresult.Pdexv3PoolPair, []v2.OrderBookIterator, []byte, error,
) {
	var results []*jsonresult.Pdexv3PoolPair
	var orderBookList []v2.OrderBookIterator
	var tradeDirections []byte

	nextTokenToSell := sellToken
	for _, pairID := range tradePath {
		if pair, exists := pairs[pairID]; exists {
			pair = clonePoolPairState(pair) // work on cloned state in case trade is rejected
			results = append(results, &pair.State)

			ob := NewOrderBook(pair.Orderbook)
			orderBookList = append(orderBookList, &ob)
			var td byte
			switch nextTokenToSell {
			case pair.State.Token0ID:
				td = v2.TradeDirectionSell0
				// set token to sell for next iteration. If this is the last iteration, it's THE token to buy
				nextTokenToSell = pair.State.Token1ID
			case pair.State.Token1ID:
				td = v2.TradeDirectionSell1
				nextTokenToSell = pair.State.Token0ID
			default:
				return nil, nil, nil, fmt.Errorf("incompatible selling token %s vs next pair %s", nextTokenToSell.String(), pairID)
			}
			tradeDirections = append(tradeDirections, td)
		} else {
			return nil, nil, nil, fmt.Errorf("path contains nonexistent pair %s", pairID)
		}
	}
	return results, orderBookList, tradeDirections, nil
}
