package block

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/txpool"
	"github.com/ethereum/go-ethereum/rlp"
)

type Header struct {
	Root       common.Hash
	ParentHash common.Hash
	Height     uint64
	Coinbase   common.Address
	Timestamp  uint64
	Nonce      uint64
}

type Body struct {
	Transactions []*common.Transaction
	Receipts     []*common.Receipt
}

type Block struct {
	CurrentHeader *Header
	CurrentBody   *Body
	StateDB       *statdb.StatDB
}

type Blockchain struct {
	CurrentHeader Header
	Statedb       *statdb.StatDB
	Txpool        *txpool.TxPool
}

func NewBlockchain(txpool *txpool.TxPool, statedb *statdb.StatDB) *Blockchain {
	return &Blockchain{
		Statedb:       statedb,
		Txpool:        txpool,
	}
}

func NewBody() *Body {
	return &Body{
		Transactions: make([]*common.Transaction, 0),
		Receipts:     make([]*common.Receipt, 0),
	}
}

func NewHeader(parent *Header) *Header {
	header := &Header{
		Root:       parent.Root,
		ParentHash: parent.Hash(),
		Height:     parent.Height + 1,
	}
	return header
}

func (header *Header) Hash() common.Hash {
	// 序列化header
	headerBytes, err := header.SerializeHeader()
	if err != nil {
		return common.Hash{}
	}
	return common.NewHash(headerBytes)
}

func (header *Header) SerializeHeader() ([]byte, error) {
	// 序列化header
	headerBytes, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, err
	}
	return headerBytes, nil
}

func (header *Header) DeserializeHeader(data []byte) Header {
	// 反序列化header
	err := rlp.DecodeBytes(data, header)
	if err != nil {
		return Header{}
	}
	return *header
}

func (body *Body) DeserializeBody(data []byte) Body {
	// 反序列化body
	err := rlp.DecodeBytes(data, body)
	if err != nil {
		return Body{}
	}
	return *body
}

func (body *Body) SerializeBody() ([]byte, error) {
	// 序列化body
	bodyBytes, err := rlp.EncodeToBytes(body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func (block *Block) SerializeBlock() ([]byte, error) {
	// 序列化header
	headerBytes, err := block.CurrentHeader.SerializeHeader()
	if err != nil {
		return nil, err
	}
	// 序列化body
	bodyBytes, err := block.CurrentBody.SerializeBody()
	if err != nil {
		return nil, err
	}
	// 合并header和body
	blockBytes := append(headerBytes, bodyBytes...)
	return blockBytes, nil
}


