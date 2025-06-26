package vm

import (
	"cxchain-2023131015/block"
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/txpool"
	"time"
)

type ChainConfig struct {
	Duration   time.Duration
	Coinbase   common.Address
	Difficulty uint64
}

type BlockMaker struct {
	txpool  *txpool.TxPool
	stateDB *statdb.StatDB
	exec    IMachine

	config ChainConfig
	chain  block.Blockchain

	nextHeader  *block.Header
	nextBody    *block.Body
	nextStateDB *statdb.StatDB

	interupt chan bool
}

// func NewBlockMaker(txpool txpool.TxPool, stateDB *statdb.StatDB, exec IMachine) *BlockMaker {
// 	return &BlockMaker{
// 		txpool:  &txpool,
// 		stateDB: stateDB,
// 		exec:    exec,
// 	}
// }

// func (maker *BlockMaker) Pack() {
// 	end := time.After(maker.config.Duration)
// 	for {
// 		select {
// 		case <-maker.interupt:
// 			break
// 		case <-end:
// 			break
// 		default:
// 			maker.pack()
// 		}
// 	}
// }

// func (maker *BlockMaker) pack() {
// 	tx := maker.txpool.Pop()
// 	receiption := maker.exec.Execute1(, *tx)
// 	maker.nextBody.Transactions = append(maker.nextBody.Transactions, tx)
// 	maker.nextBody.Receipts = append(maker.nextBody.Receipts, receiption)
// }

// func (maker *BlockMaker) Interupt() {
// 	maker.interupt <- true
// }

// func (maker *BlockMaker) Finalize() (*block.Header, *block.Body) {
// 	maker.nextHeader.Timestamp = uint64(time.Now().Unix())
// 	maker.nextHeader.Nonce = 0
// 	// TODO
// 	// for n := 0; ; n++ {
// 	// 	maker.nextHeader.Nonce = 0
// 	// 	if maker.nextHeader.Hash() {
// 	// 		break
// 	// 	}
// 	// }

// 	return maker.nextHeader, maker.nextBody
// }
