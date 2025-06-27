package block

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/txpool"
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

type Blockchain struct {
	CurrentHeader Header
	Statedb       statdb.StatDB
	Txpool        txpool.TxPool
}

func NewBlock() *Body {
	return &Body{
		Transactions: make([]*common.Transaction, 0),
		Receipts:  make([]*common.Receipt, 0),
	}
}

func NewHeader(parent *Header) *Header {
	return &Header{
		Root:       parent.Root,
		ParentHash: parent.ParentHash,
		Height:     parent.Height + 1,
		Coinbase:   parent.Coinbase,
		Timestamp:  parent.Timestamp,
		Nonce:      parent.Nonce,
	}
}


