package txpool

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
)

// pool 接口定义了交易池的基本操作
// Pop: 从池中取出一个交易
// NewTX: 向池中添加新交易
// NewTxPool: 创建新的交易池
type pool interface {
	Pop() *common.Transaction
	NewTX(tx *common.Transaction) error
	NewTxPool(statDB *statdb.StatDB) *TxPool
}

// txbox 结构体表示一个交易盒子
// txs: 包含的交易列表
// GasPrice: 该盒子中所有交易的最低GasPrice
type txbox struct {
	txs []*common.Transaction
	GasPrice uint64
}

// txboxer 接口定义了交易盒子的基本操作
type txboxer interface{
	GetGasPrice() uint64
	Push(tx *common.Transaction)
	Pop() *common.Transaction
	GetNonce() uint64
}

// boxes 是txbox的切片类型
type boxes []txbox

// Len 返回boxes中交易盒子的数量
func (b boxes) Len() int {
	return len(b)
}

// Less 比较两个交易盒子的GasPrice
func (b boxes) Less(i, j int) bool {
	return b[i].GasPrice > b[j].GasPrice
}

// Swap 交换两个交易盒子的位置
func (b boxes) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// GetFirstNonce 返回交易盒子中第一个交易的Nonce值
func (b txbox) GetFirstNonce() uint64 {
	return b.txs[0].Nonce
}

// GetLastNonce 返回交易盒子中最后一个交易的Nonce值
func (b txbox) GetLastNonce() uint64 {
	return b.txs[len(b.txs)-1].Nonce
}

// GetGasPrice 返回交易盒子的GasPrice
func (b txbox) GetGasPrice() uint64 {
	return b.GasPrice
}

// push 向交易盒子中添加一个新交易
func (b *txbox) push(tx *common.Transaction) {
	b.txs = append(b.txs, tx)
}

// pop 从交易盒子中取出并返回第一个交易
// 如果盒子中只有一个交易，取出后盒子置空
// 否则取出第一个交易并更新盒子
func (b *txbox) pop() *common.Transaction {
	if len(b.txs) == 1 {
		tx := b.txs[0]
		b.txs = nil
		return tx
	}
	tx := b.txs[0]
	b.txs = b.txs[1:]
	return tx
}

// replace 替换交易盒子中指定Nonce的交易
func (b *txbox) replace(tx *common.Transaction) {
	for i, Tx := range b.txs {
		if tx.Nonce == Tx.Nonce {
			b.txs[i] = tx
			break
		}
	}
}

// GetAddress 返回交易盒子中第一个交易的发送方地址
// 如果盒子为空，返回空地址
func (b txbox) GetAddress() common.Address {
	if len(b.txs) == 0 {
		return common.Address{}
	}
	return b.txs[0].From()
}



