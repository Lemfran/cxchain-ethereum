package txpool

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"errors"
	"fmt"
	"sort"
)

// TxPool 结构体表示交易池，用于管理待处理的交易
// StatDB: 状态数据库指针，用于查询账户状态
// pending: 按地址分组的待处理交易盒子
// queue: 按地址和nonce分组的排队交易
// Sortedboxes: 按GasPrice排序的交易盒子
type TxPool struct {
	StatDB  *statdb.StatDB
	pending map[common.Address]boxes
	queue map[common.Address]map[uint64]*common.Transaction
	Sortedboxes boxes
}

// NewTxPool 创建并初始化一个新的交易池
// statDB: 状态数据库指针
// 返回: 初始化后的TxPool指针
func NewTxPool(statDB *statdb.StatDB) *TxPool {
	return &TxPool{
		StatDB: statDB,
		pending: make(map[common.Address]boxes),
		queue: make(map[common.Address]map[uint64]*common.Transaction),
		Sortedboxes: make(boxes,0),
	}
}

// NewTX 向交易池中添加新交易
// tx: 要添加的交易
// 返回: 错误信息(如果有)
func (pool *TxPool) NewTX(tx *common.Transaction) error {
	// 从状态数据库加载发送方账户信息
	account := (*pool.StatDB).Load(tx.From())
	
	fmt.Println("-----------------------3-------------------------")
	if account.Nonce >= tx.Nonce {
		return errors.New("nonce error.")
	}
	
	// 获取当前nonce值
	nonce:=account.Nonce
	boxes:=pool.pending[tx.From()]
	if len(boxes) > 0 {
		last := boxes[len(boxes)-1]
		nonce = last.GetLastNonce()
	}
	
	// 根据nonce值决定如何处理交易
	if tx.Nonce > nonce+1 {
		fmt.Println("-----------------------4-------------------------")
		pool.addQueueTx(tx)  // nonce值过大，加入队列
		fmt.Println("-----------------------5-------------------------")
		return nil
	} else if tx.Nonce == nonce+1 {
		pool.addPendingTx(tx)  // 连续nonce，直接加入待处理
		return nil
	} else {
		pool.replacePendingTx(tx)  // 替换现有交易
		return nil
	}
}

// addQueueTx 将交易添加到队列中
// tx: 要添加的交易
func (pool *TxPool) addQueueTx(tx *common.Transaction){
	fmt.Println("-----------------------6-------------------------")
	// 初始化或获取该地址的交易队列
	list := pool.queue[tx.From()]
	if list == nil {
		list = make(map[uint64]*common.Transaction)
	}
	// 按nonce存储交易
	list[tx.Nonce] = tx
	pool.queue[tx.From()] = list
}

// addPendingTx 将交易添加到待处理池中
// tx: 要添加的交易
func (pool *TxPool) addPendingTx(tx *common.Transaction){
	// 获取该地址的待处理交易盒子
	boxes := pool.pending[tx.From()]

	if len(boxes) == 0 {
		// 创建新的交易盒子并添加到待处理池
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
			// GasPrice足够高，可以合并到现有盒子
			last.push(tx)
			pool.pending[tx.From()][len(boxes)-1] = last
			
			// 重新构建并排序所有交易盒子
			var newboxs []txbox
			for _, box := range pool.pending {
				newboxs = append(newboxs, box...)
			}
			pool.Sortedboxes = newboxs
			sort.Sort(pool.Sortedboxes)
			
		} else {
			// GasPrice太低，创建新的交易盒子
			box:=txbox{
				txs: []*common.Transaction{tx},
				GasPrice: tx.GasPrice,
			}
			pool.pending[tx.From()]=append(pool.pending[tx.From()],box)
			pool.Sortedboxes = append(pool.Sortedboxes, box)
			sort.Sort(pool.Sortedboxes)
		}
	}

	// 检查队列中是否有连续nonce的交易可以移出
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
	//3.第一种情况就要把盒子中的交易的gas挨个与前盒子的gas比较，大合并到前面的盒子，小就不动，并且是按Nonce值顺序检查，Nonce小的优先，一个一个检查，小于前一个就分裂
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
		//重新给所有盒子排序
		pool.Allsort()
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
			} else if len(pool.pending[tx.From()][flag].txs) == 1 {
				pool.pending[tx.From()][flag].GasPrice=tx.GasPrice
			}
		}
		pool.Allsort()
	}

	/*如果新交易的 Nonce 等于当前交易盒子的首部 Nonce，并且当前交易盒子不是第一个盒子（flag > 0），则进入循环。
 	在循环中，如果新交易的 GasPrice 大于等于前一个交易盒子的 GasPrice，则从当前交易盒子中弹出一个交易（pop），并将其推入前一个交易盒子（push）。
	更新 tx 为当前交易盒子的第一个交易。
	如果新交易的 GasPrice 小于前一个交易盒子的 GasPrice，则退出循环。*/

}

// Pop 从交易池中取出并返回一组交易
// 返回: 取出的交易列表，如果没有交易则返回nil
func (pool *TxPool) Pop() []*common.Transaction {
	var txs []*common.Transaction
	fmt.Println("-----------------------pop前状态-------------------------")
	fmt.Println("全部交易盒子：",pool.Sortedboxes)

	// 检查交易池是否为空
	if len(pool.Sortedboxes) == 0 || len(pool.pending) == 0 || len(pool.Sortedboxes[0].txs) == 0 {
		return nil
	}
	
	// 获取GasPrice最高的交易地址
	addr := pool.Sortedboxes[0].GetAddress()
	fmt.Println("-----------------------pop前状态-------------------------")
	
	// 获取该地址的所有交易盒子
	boxes, exists := pool.pending[addr]
	if !exists || len(boxes) == 0 {
		return nil
	}
	
	// 获取第一个交易盒子中的交易数量
	txsLen := len(boxes[0].txs)
	fmt.Println(pool.pending[addr])
	// 循环取出所有交易
	for i:=0;i<txsLen;i++ {
		fmt.Println(txsLen)
		fmt.Println("取出一笔交易")
		
		// 从交易盒子中取出交易
		fmt.Println(len(pool.Sortedboxes[0].txs))
		if len(pool.Sortedboxes[0].txs) == 0 {
			fmt.Println("交易盒子已空")
			pool.Sortedboxes = pool.Sortedboxes[1:]
			pool.pending[addr] = pool.pending[addr][1:]
			return txs
		}
		tx := pool.pending[addr][0].pop()
		
		// 更新Sortedboxes中的交易列表
		pool.Sortedboxes[0].txs=pool.Sortedboxes[0].txs[1:]
		
		// 将交易添加到返回列表
		txs = append(txs, tx)
	}
	fmt.Println(txs)

	// 如果Sortedboxes中的交易已全部取出，则更新pending
	if len(pool.Sortedboxes[0].txs) == 0 {
		fmt.Println("交易盒子已空")
		pool.Sortedboxes = pool.Sortedboxes[1:]
		pool.pending[addr] = pool.pending[addr][1:]
	}
	
	return txs
}

//给所有交易盒子重新排序
func (pool *TxPool)Allsort() {
	var newboxs []txbox
			for _, box := range pool.pending {
				newboxs = append(newboxs, box...)
			}
		pool.Sortedboxes = newboxs
		sort.Sort(pool.Sortedboxes)
}
