package pdex_v3

import (
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler/jsonresult"
	v2 "github.com/incognitochain/incognito-cli/pdex_v3/v2utils"
	"math/big"
	"sort"
	"strings"
)

type Order = jsonresult.Pdexv3Order

type OrderBook struct {
	orders []*Order
}

func NewOrderBook(tmpOrder jsonresult.Pdexv3Orderbook) OrderBook {
	res := OrderBook{}
	orders := make([]*Order, len(tmpOrder.Orders))
	for i, order := range tmpOrder.Orders {
		orders[i] = order
	}
	res.orders = orders
	return res
}

func (ob OrderBook) MarshalJSON() ([]byte, error) {
	temp := struct {
		Orders []*Order `json:"orders,omitempty"`
	}{ob.orders}
	return json.Marshal(temp)
}

func (ob *OrderBook) UnmarshalJSON(data []byte) error {
	var temp struct {
		Orders []*Order `json:"orders"`
	}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	ob.orders = temp.Orders
	return nil
}

// InsertOrder appends a new order while keeping the list sorted (ascending by Token1Rate / Token0Rate)
func (ob *OrderBook) InsertOrder(ord *Order) {
	insertAt := func(lst []*Order, i int, newItem *Order) []*Order {
		if i == len(lst) {
			return append(lst, newItem)
		}
		lst = append(lst[:i+1], lst[i:]...)
		lst[i] = newItem
		return lst
	}
	index := sort.Search(len(ob.orders), func(i int) bool {
		ordRate := big.NewInt(0).SetUint64(ob.orders[i].Token0Rate)
		ordRate.Mul(ordRate, big.NewInt(0).SetUint64(ord.Token1Rate))
		myRate := big.NewInt(0).SetUint64(ob.orders[i].Token1Rate)
		myRate.Mul(myRate, big.NewInt(0).SetUint64(ord.Token0Rate))

		// compare Token1Rate / Token0Rate of current order in the list to ord
		rateCmp := ordRate.Cmp(myRate)
		// break equality of rate by comparing ID
		if rateCmp == 0 {
			// sell0 orders must precede sell1 orders of the same rate
			if ord.TradeDirection != ob.orders[i].TradeDirection {
				return ord.TradeDirection == TradeDirectionSell0
			}
			// no equality for ID since duplicate ID was handled in addOrder flow
			idCmp := strings.Compare(ord.Id, ob.orders[i].Id)
			// best rate for sell0 is at start of list, so we put smaller ID first to match. The opposite is true for sell1
			if ord.TradeDirection == TradeDirectionSell0 {
				return idCmp < 0
			} else {
				return idCmp > 0
			}
		}

		return ordRate.Cmp(myRate) < 0
	})
	ob.orders = insertAt(ob.orders, index, ord)
}

// NextOrder returns the matchable order with the best rate that has any outstanding balance to sell
func (ob *OrderBook) NextOrder(tradeDirection byte) (*v2.MatchingOrder, string, error) {
	lstLen := len(ob.orders)
	switch tradeDirection {
	case v2.TradeDirectionSell0:
		for i := lstLen - 1; i >= 0; i-- {
			currentOrder := &v2.MatchingOrder{Pdexv3Order: ob.orders[i]}
			if check, err := currentOrder.CanMatch(tradeDirection); check && err == nil {
				return currentOrder, ob.orders[i].Id, nil
			}
		}
		// no active order
		return nil, "", nil
	case v2.TradeDirectionSell1:
		for i := 0; i < lstLen; i++ {
			currentOrder := &v2.MatchingOrder{Pdexv3Order: ob.orders[i]}
			if check, err := currentOrder.CanMatch(tradeDirection); check && err == nil {
				return currentOrder, ob.orders[i].Id, nil
			}
		}
		// no active order
		return nil, "", nil
	default:
		return nil, "", fmt.Errorf("Invalid trade direction %d", tradeDirection)
	}
}

// RemoveOrder removes one order by its index
func (ob *OrderBook) RemoveOrder(index int) error {
	if index < 0 || index >= len(ob.orders) {
		return fmt.Errorf("Invalid order index %d for orderbook length %d", index, len(ob.orders))
	}
	ob.orders = append(ob.orders[:index], ob.orders[index+1:]...)
	return nil
}
