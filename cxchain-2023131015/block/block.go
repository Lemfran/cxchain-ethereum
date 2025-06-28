package block

import (
	"cxchain-2023131015/common"  
	"cxchain-2023131015/statdb"  
	"cxchain-2023131015/txpool"   
	"github.com/ethereum/go-ethereum/rlp" 
)

// Header 结构体表示区块头信息
// Root: 状态树根哈希
// ParentHash: 父区块哈希
// Height: 区块高度(区块号)
// Coinbase: 矿工地址(打包该区块的矿工)
// Timestamp: 区块创建时间戳
// Nonce: 随机数(用于工作量证明)
type Header struct {
	Root       common.Hash
	ParentHash common.Hash
	Height     uint64
	Coinbase   common.Address
	Timestamp  uint64
	Nonce      uint64
}

// Body 结构体表示区块体信息
// Transactions: 区块包含的交易列表
// Receipts: 交易执行后的收据列表
type Body struct {
	Transactions []*common.Transaction
	Receipts     []*common.Receipt
}

// Block 结构体表示完整的区块
// CurrentHeader: 当前区块头指针
// CurrentBody: 当前区块体指针
// StateDB: 状态数据库指针(区块执行后的状态)
type Block struct {
	CurrentHeader *Header
	CurrentBody   *Body
	StateDB       *statdb.StatDB
}

// Blockchain 结构体表示区块链
// CurrentHeader: 当前链上最新的区块头
// Statedb: 状态数据库实例
// Txpool: 交易池实例
type Blockchain struct {
	CurrentHeader Header
	Statedb       *statdb.StatDB
	Txpool        *txpool.TxPool
}

// NewBlockchain 创建并初始化一个新的区块链实例
// txpool: 交易池实例
// statedb: 状态数据库实例
// 返回: 初始化后的区块链指针
func NewBlockchain(txpool *txpool.TxPool, statedb *statdb.StatDB) *Blockchain {
	return &Blockchain{
		Statedb:       statedb,  // 初始化状态数据库
		Txpool:        txpool,   // 初始化交易池
	}
}

// NewBody 创建并初始化一个新的区块体
// 返回: 初始化后的区块体指针
func NewBody() *Body {
	return &Body{
		Transactions: make([]*common.Transaction, 0), // 初始化空交易列表
		Receipts:     make([]*common.Receipt, 0),    // 初始化空收据列表
	}
}

// NewHeader 根据父区块头创建新的区块头
// parent: 父区块头指针
// 返回: 初始化后的区块头指针
func NewHeader(parent *Header) *Header {
	header := &Header{
		Root:       parent.Root,       // 继承父区块的状态根
		ParentHash: parent.Hash(),     // 设置父区块哈希
		Height:     parent.Height + 1, // 高度递增
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



