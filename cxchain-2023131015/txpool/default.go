package txpool

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"errors"
	"sort"
)

type TxPool struct {
	StatDB  *statdb.StatDB
	pending map[common.Address]boxes
	queue map[common.Address]map[uint64]*common.Transaction
	Sortedboxes boxes
}

func (pool *TxPool) NewTX(tx *common.Transaction) error {
	// 由于 pool.StatDB 是指针类型，需要先解引用再调用方法
	// 假设 statdb.StatDB 接口有一个 Load 方法，这里直接调用
	account := (*pool.StatDB).Load(tx.From())

	if account.Nonce >= tx.Nonce {
		return errors.New("nonce error")
	}
	nonce:=account.Nonce
	boxes:=pool.pending[tx.From()]
	if len(boxes) > 0 {
		last := boxes[len(boxes)-1]
		nonce = last.GetLastNonce()
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

func (pool *TxPool) addQueueTx(tx *common.Transaction){
	list := pool.queue[tx.From()]
	if list == nil {
		list = make(map[uint64]*common.Transaction)
	}
	list[tx.Nonce] = tx
	pool.queue[tx.From()] = list
}

func (pool *TxPool) addPendingTx(tx *common.Transaction){
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

	//看Queue(map类型)中是否要更新出来交易到pending中，其中Nonce值是tx.Nonce+1，并且有个循环检查是否把符合的输出完全
	if pool.queue[tx.From()] != nil {
		nonce := tx.Nonce + 1
		for {
			if txq, ok := pool.queue[tx.From()][nonce]; ok {
				delete(pool.queue[tx.From()], nonce)
				pool.addPendingTx(txq)
				nonce++
			} else {
				break
			}
		}
	}
	
}

func (pool *TxPool) replacePendingTx(tx *common.Transaction) {
	var flag int
	for i, txbox := range pool.pending[tx.From()] {
		if txbox.GetFirstNonce() <= tx.Nonce && txbox.GetLastNonce() >= tx.Nonce {
			// replace
			if txbox.GetGasPrice() <= tx.GasPrice {
				pool.pending[tx.From()][i].replace(tx)
			}
			flag = i
			break
		}
	}

	//pending[tx.From()]中tx的gas值更新有两种情况，一种是在盒子中间，一种是在盒子首部
	//1.盒子中间更新就是修改的gas大就替换，修改gas值小就不管了
	//2.盒子首部更新分为两种种情况，一种是修改的gas大于等于前面盒子，一种是gas值小于前面盒子
	//3.第一种情况就要把盒子中的交易的gas挨个与前盒子的gas比较，大合并到前面的盒子，小就不动，并且是按Nonce值顺序检查，Nonce小的优先，直到gas小于前面盒子的gas，就不动了
	//4.第二种情况就直接替换，无需合并盒子
	//总结：1，4不用管因为已经替换了，只要对3进行检查合并就行了

	if tx.Nonce == pool.pending[tx.From()][flag].GetFirstNonce() && flag == 0{
		length:=len(pool.pending[tx.From()][flag].txs)
		txchange:=pool.pending[tx.From()][flag].pop()
				box:=txbox{
					txs: []*common.Transaction{txchange},
					GasPrice: txchange.GasPrice,
				}
				//flag+1之后为原切片，增加新切片于flag处
				//先扩容
				pool.pending[tx.From()]=append(pool.pending[tx.From()],box)
				copy(pool.pending[tx.From()][flag+1:],pool.pending[tx.From()][flag:])
				pool.pending[tx.From()][flag]=box
				pool.pending[tx.From()][flag].GasPrice=tx.GasPrice

				flag++
				tx = pool.pending[tx.From()][flag].txs[0]
		for i:=1;i<length;i++ {
			if tx.GasPrice >= pool.pending[tx.From()][flag-1].GetGasPrice() {
				txchange:=pool.pending[tx.From()][flag].pop()
				pool.pending[tx.From()][flag-1].push(txchange)
				tx = pool.pending[tx.From()][flag].txs[0]
			} else if len(pool.pending[tx.From()][flag].txs) > 1 && tx.GasPrice < pool.pending[tx.From()][flag-1].GetGasPrice(){
				txchange:=pool.pending[tx.From()][flag].pop()
				box:=txbox{
					txs: []*common.Transaction{txchange},
					GasPrice: txchange.GasPrice,
				}
				//flag+1之后为原切片，增加新切片于flag处
				//先扩容
				pool.pending[tx.From()]=append(pool.pending[tx.From()],box)
				copy(pool.pending[tx.From()][flag+1:],pool.pending[tx.From()][flag:])
				pool.pending[tx.From()][flag]=box

				flag++
				tx = pool.pending[tx.From()][flag].txs[0]
			} 
		}
	}

	if tx.Nonce == pool.pending[tx.From()][flag].GetFirstNonce() && flag > 0 {
		length:=len(pool.pending[tx.From()][flag].txs)
		for i:=0;i<length;i++ {
			if tx.GasPrice >= pool.pending[tx.From()][flag-1].GetGasPrice() {
				txchange:=pool.pending[tx.From()][flag].pop()
				pool.pending[tx.From()][flag-1].push(txchange)
				tx = pool.pending[tx.From()][flag].txs[0]
			} else if len(pool.pending[tx.From()][flag].txs) > 1 && tx.GasPrice < pool.pending[tx.From()][flag-1].GetGasPrice(){
				txchange:=pool.pending[tx.From()][flag].pop()
				box:=txbox{
					txs: []*common.Transaction{txchange},
					GasPrice: txchange.GasPrice,
				}
				//flag+1之后为原切片，增加新切片于flag处
				//先扩容
				pool.pending[tx.From()]=append(pool.pending[tx.From()],box)
				copy(pool.pending[tx.From()][flag+1:],pool.pending[tx.From()][flag:])
				pool.pending[tx.From()][flag]=box
				pool.pending[tx.From()][flag].GasPrice=tx.GasPrice

				flag++
				tx = pool.pending[tx.From()][flag].txs[0]
			} 
		}
	}

	/*如果新交易的 Nonce 等于当前交易盒子的首部 Nonce，并且当前交易盒子不是第一个盒子（flag > 0），则进入循环。
 	在循环中，如果新交易的 GasPrice 大于等于前一个交易盒子的 GasPrice，则从当前交易盒子中弹出一个交易（pop），并将其推入前一个交易盒子（push）。
	更新 tx 为当前交易盒子的第一个交易。
	如果新交易的 GasPrice 小于前一个交易盒子的 GasPrice，则退出循环。*/

}

func (pool *TxPool) Pop() *common.Transaction {
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
