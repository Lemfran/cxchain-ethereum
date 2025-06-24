package txpool

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"errors"
	"sort"
)

type TxPool struct {
	StatDB  statdb.StatDB
	all map[common.Hash]bool
	pending map[common.Address]boxes
	queue map[common.Address]map[uint64]*common.Transaction
	Sortedboxes boxes
}

func (pool TxPool) NewTX(tx *common.Transaction) error {
	account := pool.StatDB.Load(tx.From())
	
	if account.Nonce >= tx.Nonce {
		return errors.New("nonce error")
	}
	nonce:=account.Nonce
	boxes:=pool.pending[tx.From()]
	if len(boxes) > 0 {
		last := boxes[len(boxes)-1]
		nonce = last.GetNonce()
	}
	if tx.Nonce > nonce+1 {
		pool.addQueueTx(tx)
		return nil
	} else if tx.Nonce == nonce+1 {
		// push
		pool.addPendingTx(tx)
		return nil
	} else {
		// 替换
		pool.replacePendingTx(tx)
		return nil
	}
}

func (pool TxPool) addQueueTx(tx *common.Transaction){
	list := pool.queue[tx.From()]
	if list == nil {
		list = make(map[uint64]*common.Transaction)
	}
	list[tx.Nonce] = tx
}

func (pool TxPool) addPendingTx(tx *common.Transaction){
	boxes := pool.pending[tx.From()]
	if len(boxes) == 0 {
		//加到pending中
		box:=txbox{
			txs: []*common.Transaction{tx},
			GasPrice: tx.GasPrice,
		}
		pool.pending[tx.From()]=append(pool.pending[tx.From()],box)
		pool.Sortedboxes = append(pool.Sortedboxes, box)
		sort.Sort(pool.Sortedboxes)
	} else {
		last := boxes[len(boxes)-1]
		if  tx.GasPrice >= last.GetGasPrice() {
			//加到pending中
			last.push(tx)
		} else {
			//新建一个txbox在pending中
			box:=txbox{
				txs: []*common.Transaction{tx},
				GasPrice: tx.GasPrice,
			}
			pool.pending[tx.From()]=append(pool.pending[tx.From()],box)
			pool.Sortedboxes = append(pool.Sortedboxes, box)
			sort.Sort(pool.Sortedboxes)
		}
	}
}

func (pool TxPool) replacePendingTx(tx *common.Transaction) {
	for _, txbox := range pool.pending[tx.From()] {
		if txbox.GetNonce() >= tx.Nonce {
			// replace
			if txbox.GetGasPrice() <= tx.GasPrice {
				txbox.replace(tx)
			}
			break
		}
	}
}

func (pool TxPool) Pop() *common.Transaction {
	boxes := pool.pending[pool.Sortedboxes[0].GetAddress()]
	if len(boxes) == 0 {
		return nil
	}

	tx := boxes[0].pop()
	if len(boxes[0].txs) == 0 {
		pool.pending[pool.Sortedboxes[0].GetAddress()] = boxes[1:]
	}
	return tx
}
