package common

import "math/big"

type Transaction struct {
	txdata
	signature
}

type txdata struct {
	To       []byte
	Value    uint64
	Nonce    uint64
	GasLimit uint64
	GasPrice uint64
	Input    []byte
}

type signature struct {
	R, S *big.Int
	V uint8
}