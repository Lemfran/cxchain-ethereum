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

type ChainConfig struct {
	Duration   time.Duration
	Coinbase   common.Address
	Difficulty uint64
}

type BlockMaker struct {
	Txpool  *txpool.TxPool
	Statedb *statdb.StatDB
	Exec    *vm.StateMachine

	Config ChainConfig
	Chain  *block.Blockchain

	NextHeader  *block.Header
	NextBody    *block.Body
	NextStateDB *statdb.StatDB

	interupt chan bool
}

func NewBlockMaker(txpool *txpool.TxPool, statedb *statdb.StatDB, exec *vm.StateMachine) *BlockMaker {
	chain :=block.NewBlockchain(txpool, statedb)
	return &BlockMaker{
		Txpool:  txpool,
		Statedb: statedb,
		Exec:    exec,
        interupt: make(chan bool, 1), // 初始化带缓冲的 channel
		Chain:   chain,
	}
}

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
        maker.NextHeader = block.NewHeader(&maker.Chain.CurrentHeader)
		maker.Chain.CurrentHeader = *maker.NextHeader
    }
    
    maker.NextHeader.Coinbase = maker.Config.Coinbase
}

func (maker *BlockMaker) Pack() {
	end := time.After(1 * time.Second)
	for {
		select {
		case <-maker.interupt:
			break
		case <-end:
			break
		default:
			if !maker.pack() {
				break
			}
		}
	}
}

func (maker *BlockMaker) pack() bool {
	tx := maker.Txpool.Pop()
	if tx == nil {
		time.Sleep(1 * time.Millisecond) // 添加短暂休眠避免CPU空转
		return false
	}
	
	receiption := maker.Exec.Execute1(*tx)
	if receiption == nil {
		return false
	}
	fmt.Println("-----------------------8-------------------------")
	maker.NextBody.Transactions = append(maker.NextBody.Transactions, tx)
	maker.NextBody.Receipts = append(maker.NextBody.Receipts, receiption)
	return true
}

func (maker *BlockMaker) Interupt() {
	maker.interupt <- true
}

func (maker *BlockMaker) Finalize() (*block.Header, *block.Body) {
	maker.NextHeader.Timestamp = uint64(maker.Config.Duration)
	maker.NextHeader.Nonce = 0
	// TODO
	// for n := 0; ; n++ {
	// 	maker.nextHeader.Nonce = 0
	// 	if maker.nextHeader.Hash() {
	// 		break
	// 	}
	// }

	return maker.NextHeader, maker.NextBody
}


//判断是否是创世区块
func (maker *BlockMaker) IsGenesisBlock() bool {
		var zeroHeader block.Header
		return maker.Chain.CurrentHeader == zeroHeader
}