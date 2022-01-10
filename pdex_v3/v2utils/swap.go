package v2utils

import (
	"fmt"
	"math/big"
)

func calculateBuyAmount(amountIn uint64, reserveIn uint64, reserveOut uint64, virtualReserveIn *big.Int, virtualReserveOut *big.Int) (uint64, error) {
	if amountIn <= 0 {
		return 0, fmt.Errorf("insufficient input amount")
	}
	if reserveIn <= 0 || reserveOut <= 0 {
		return 0, fmt.Errorf("insufficient liquidity")
	}
	amount := big.NewInt(0).SetUint64(amountIn)
	num := big.NewInt(0).Mul(amount, virtualReserveOut)
	den := big.NewInt(0).Add(amount, virtualReserveIn)
	result := num.Div(num, den)
	if !result.IsUint64() {
		return 0, fmt.Errorf("number out of range uint64")
	}
	return result.Uint64(), nil
}

func calculateAmountToSell(amountOut uint64, reserveIn uint64, reserveOut uint64, virtualReserveIn *big.Int, virtualReserveOut *big.Int) (uint64, error) {
	if amountOut <= 0 {
		return 0, fmt.Errorf("insufficient input amount")
	}
	if reserveIn <= 0 || reserveOut <= 0 {
		return 0, fmt.Errorf("insufficient liquidity")
	}
	num := big.NewInt(0).Mul(virtualReserveIn, big.NewInt(0).SetUint64(amountOut))
	den := big.NewInt(0).Sub(virtualReserveOut, big.NewInt(0).SetUint64(amountOut))
	result := num.Div(num, den)
	result.Add(result, big.NewInt(1))
	if !result.IsUint64() {
		return 0, fmt.Errorf("number out of range uint64")
	}
	return result.Uint64(), nil
}
