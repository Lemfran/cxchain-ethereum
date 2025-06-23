package txpool

import (
	"cxchain-2023131015/common"
)

type pool interface {
	AddTransaction(tx *common.Transaction) error
	PopTransaction() (*common.Transaction, error)
	
}

type txbox struct {
	txs []*common.Transaction
	GasPrice uint64
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


