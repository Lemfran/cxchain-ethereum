package block

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
)

type Header struct {
    Root       common.Hash
	ParentHash common.Hash
	Height     uint64
	Coinbase   common.Address
	Timestamp  uint64
	Nonce uint64
}

type Body struct {
    Transactions []*common.Transaction
	Receipts []*common.Receipt
}

type Block struct {
    CurrentHeader Header
	StateDB       *statdb.StatDB
}


