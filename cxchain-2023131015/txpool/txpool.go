package txpool

import (
	"cxchain-2023131015/common"
)

type pool interface {
	AddTransaction(tx *common.Transaction) error
	PopTransaction() (*common.Transaction, error)
	SetStatRoot(root []byte)
	NotifyTxEvent(txs []*common.Transaction)
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

func (b txbox) GetNonce() uint64 {
	return b.txs[0].Nonce
}

func (b txbox) GetGasPrice() uint64 {
	return b.GasPrice
}

func (b txbox) push(tx *common.Transaction) {
	b.txs = append(b.txs, tx)
}

func (b txbox) pop() *common.Transaction {
	tx := b.txs[0]
	b.txs = b.txs[1:]
	return tx
}

func (b txbox) replace(tx *common.Transaction) {
	for i, Tx := range b.txs {
		if tx.Nonce == Tx.Nonce {
			b.txs[i] = tx
			break
		}
	}
}

func (b txbox) GetAddress() common.Address {
	return b.txs[0].From()
}

