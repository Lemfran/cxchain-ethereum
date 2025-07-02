package makertest

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/kvstore/leveldb"
	"cxchain-2023131015/maker"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/trie/mpt"
	"cxchain-2023131015/txpool"
	"cxchain-2023131015/vm"
	"fmt"
	"testing"
	"time"
)


func TestBlockMaker_Finalize(t *testing.T) {
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		panic(err)
	}
	mptDB := mpt.NewMPT(db)
	statDB := statdb.NewStatDB(mptDB)
	txPool := txpool.NewTxPool(statDB)
	machine := vm.NewStateMachine(txPool)
	maker := maker.NewBlockMaker(txPool, statDB, machine)

	// 初始化测试账户
	account, _ := common.GenerateAccount(100000000)
	maker.Statedb = statdb.NewStatDB(mpt.NewMPT(db))
	maker.Statedb.Store(account.Address, account.Account)
	defer db.Close()

	// 创建新区块
	maker.NewBlock()

	// 测试Finalize方法
	t.Run("Finalize block", func(t *testing.T) {
		header, body := maker.Finalize()

		if header == nil {
			t.Error("Header should not be nil")
		}
		if body == nil {
			t.Error("Body should not be nil")
		}
	})
}

func TestBlockMaker_Pack(t *testing.T) {
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		panic(err)
	}
	mptDB := mpt.NewMPT(db)
	statDB := statdb.NewStatDB(mptDB)
	txPool := txpool.NewTxPool(statDB)
	machine := vm.NewStateMachine(txPool)
	maker := maker.NewBlockMaker(txPool, statDB, machine)
	defer db.Close()


	

	// 添加测试交易
	account1, _ := common.GenerateAccount(100000000)
	account2, _ := common.GenerateAccount(100000000)
	tx := common.NewTransaction(account2.Address[:], 100, 1, 100000, 100000, []byte(""))
	tx.Sign(account1.PrivateKey)
	maker.Statedb.Store(account1.Address, account1.Account)
	maker.Statedb.Store(account2.Address, account2.Account)
	fmt.Println("---------------------------放入交易------------------------------")
	txPool.NewTX(tx)
	fmt.Println("----------------------------开始打包----------------------------------")

	// 测试Pack方法
	maker.NewBlock()
	go maker.Pack()
	
	
	// 等待打包完成
	time.Sleep(1 * time.Second)
	maker.Interupt()
	fmt.Println("--------------------------打包完成---------------------------")
	
	if len(maker.NextBody.Transactions) == 0 {
		t.Error("Should pack at least one transaction")
	}
}

func TestBlockMaker_PackBlock(t *testing.T) {
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		panic(err)
	}
	mptDB := mpt.NewMPT(db)
	statDB := statdb.NewStatDB(mptDB)
	txPool := txpool.NewTxPool(statDB)
	machine := vm.NewStateMachine(txPool)
	maker := maker.NewBlockMaker(txPool, statDB, machine)
	defer db.Close()


	// 添加测试交易
	//生成账户并存入状态树
	account1, _ := common.GenerateAccount(100000000)
	account2, _ := common.GenerateAccount(100000000)
	tx := common.NewTransaction(account2.Address[:], 100, 1, 100000, 100000, []byte(""))
	tx2 := common.NewTransaction(account2.Address[:], 100, 2, 100000, 120000, []byte(""))
	tx3 := common.NewTransaction(account2.Address[:], 100, 3, 100000, 130000, []byte(""))
	tx.Sign(account1.PrivateKey)
	tx2.Sign(account1.PrivateKey)
	tx3.Sign(account1.PrivateKey)


	maker.Statedb.Store(account1.Address, account1.Account)
	maker.Statedb.Store(account2.Address, account2.Account)
	fmt.Println("---------------------------放入交易------------------------------")
	//添加交易
	txPool.NewTX(tx)
	txPool.NewTX(tx2)
	fmt.Println("----------------------------开始打包----------------------------------")

	// 测试PackBlock方法
	//第一次打包
	maker.PackBlock()

	//添加交易3
	txPool.NewTX(tx3)

	//第二次打包
	maker.PackBlock()
}
