package txpool

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
)

type pool interface {
	Pop() *common.Transaction
	NewTX(tx *common.Transaction) error
	NewTxPool(statDB *statdb.StatDB) *TxPool
}

type txbox struct {
	txs []*common.Transaction
	GasPrice uint64
}

type txboxer interface{
	GetGasPrice() uint64
	Push(tx *common.Transaction)
	Pop() *common.Transaction
	GetNonce() uint64
}

type boxes []txbox

func (b boxes) Len() int {
	return len(b)
}

func (b boxes) Less(i, j int) bool {
	return b[i].GasPrice < b[j].GasPrice
}

func (b boxes) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b txbox) GetFirstNonce() uint64 {
	return b.txs[0].Nonce
}

func (b txbox) GetLastNonce() uint64 {
	return b.txs[len(b.txs)-1].Nonce
}

func (b txbox) GetGasPrice() uint64 {
	return b.GasPrice
}

func (b *txbox) push(tx *common.Transaction) {
	b.txs = append(b.txs, tx)
}

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

func (b *txbox) replace(tx *common.Transaction) {
	for i, Tx := range b.txs {
		if tx.Nonce == Tx.Nonce {
			b.txs[i] = tx
			break
		}
	}
}

func (b txbox) GetAddress() common.Address {
	if len(b.txs) == 0 {
		return common.Address{}
	}
	return b.txs[0].From()
}



