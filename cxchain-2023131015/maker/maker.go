package maker

import (
	"cxchain-2023131015/block"
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/txpool"
	"cxchain-2023131015/vm"
	"fmt"
	"time"
)

// ChainConfig 定义区块链的配置参数
// Duration: 区块生成间隔时间
// Coinbase: 矿工地址
// Difficulty: 挖矿难度
type ChainConfig struct {
	Duration   time.Duration
	Coinbase   common.Address
	Difficulty uint64
}

// BlockMaker 负责创建和打包新区块
// 包含交易池、状态数据库、虚拟机执行器等组件
type BlockMaker struct {
	Txpool  *txpool.TxPool    // 交易池
	Statedb *statdb.StatDB    // 状态数据库
	Exec    *vm.StateMachine  // 虚拟机执行器

	Config ChainConfig        // 区块链配置
	Chain  *block.Blockchain  // 区块链

	NextHeader  *block.Header  // 下一个区块头
	NextBody    *block.Body    // 下一个区块体
	NextStateDB *statdb.StatDB // 下一个状态数据库

	interupt chan bool         // 中断通道
}

// NewBlockMaker 创建新的区块生成器
func NewBlockMaker(txpool *txpool.TxPool, statedb *statdb.StatDB, exec *vm.StateMachine) *BlockMaker {
	chain := block.NewBlockchain(txpool, statedb)
	return &BlockMaker{
		Txpool:  txpool,
		Statedb: statedb,
		Exec:    exec,
        interupt: make(chan bool, 1), // 初始化带缓冲的 channel
		Chain:   chain,
	}
}

// NewBlock 创建新的区块
func (maker *BlockMaker) NewBlock() {
	maker.NextBody = block.NewBody()
    // 判断是否是创世区块
    if maker.IsGenesisBlock() {
        maker.NextHeader = &block.Header{
			Root:       maker.Statedb.GetRoot(),
			ParentHash: common.Hash{},
			Height:     0,
		}
		maker.Chain.CurrentHeader = *maker.NextHeader
    } else {
		// 非创世区块，根据当前区块创建新的区块
        maker.NextHeader = block.NewHeader(&maker.Chain.CurrentHeader)
		maker.Chain.CurrentHeader = *maker.NextHeader
    }
    maker.NextHeader.Coinbase = maker.Config.Coinbase
}

// Pack 打包交易到区块中
func (maker *BlockMaker) Pack() {
	end := time.After(1 * time.Second)
	//这里为了方便测试减少循环次数为5次
	for i:=0;i<5;i++{
		select {
		case <-maker.interupt:  // 检查中断信号
			break
		case <-end:  // 超时检查
			break
		default:
			if !maker.pack() {
				break
			}
		}
	}
}

// pack 内部打包方法，处理单个交易
func (maker *BlockMaker) pack() bool {
	txs := maker.Txpool.Pop()
	if txs == nil {
		time.Sleep(1 * time.Millisecond) // 添加短暂休眠避免CPU空转
		return false
	}
	
	for i:=0;i<len(txs);i++ {
		fmt.Println("-------------------开始执行---------------------")
		receiption := maker.Exec.Execute1(*txs[i])
		if receiption == nil {
			return false
		}
		fmt.Println("-----------------------8-------------------------",txs[i])
		maker.NextBody.Transactions = append(maker.NextBody.Transactions, txs[i])
		maker.NextBody.Receipts = append(maker.NextBody.Receipts, receiption)
	}
	return true
}

// Interupt 发送中断信号
func (maker *BlockMaker) Interupt() {
	maker.interupt <- true
}

// Finalize 完成区块的最终处理
func (maker *BlockMaker) Finalize() (*block.Header, *block.Body) {
	maker.NextHeader.Timestamp = uint64(maker.Config.Duration)
	maker.NextHeader.Nonce = 0
	for {
		if maker.validNonce() {
			break
		}
	}
	return maker.NextHeader, maker.NextBody
}

// validNonce 验证Nonce值是否有效
func (maker *BlockMaker) validNonce() bool {
	return true
}


// IsGenesisBlock 判断当前区块是否是创世区块
func (maker *BlockMaker) IsGenesisBlock() bool {
		var zeroHeader block.Header
		return maker.Chain.CurrentHeader == zeroHeader
}

// PackBlock 打包区块的完整流程
func (maker *BlockMaker) PackBlock() {
	maker.Config.Duration = 3 * time.Second
	maker.NewBlock()
	maker.Pack()
	header, _:= maker.Finalize()
	maker.Chain.CurrentHeader = *header
	fmt.Printf("打包完成")
}
