package vmtest

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/kvstore/leveldb"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/trie/mpt"
	"cxchain-2023131015/txpool"
	"cxchain-2023131015/vm"
	"fmt"

	"testing"
)

func setupTestStateMachine() (*vm.StateMachine, func()) {
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		panic(err)
	}
	mptDB := mpt.NewMPT(db)
	statDB := statdb.NewStatDB(mptDB)
	txPool := txpool.NewTxPool(statDB)
	machine := vm.NewStateMachine(txPool)	
	return machine, func() {
		db.Close()
	}
}


func TestStateMachine_Execute(t *testing.T) {

	machine, closeMachine := setupTestStateMachine()
	defer closeMachine()
	
	accountinfo1,_:= common.GenerateAccount(100000000)
	fmt.Println(accountinfo1)
	machine.TxPool.StatDB.Store(accountinfo1.Address, accountinfo1.Account)
	accountinfo2,_:= common.GenerateAccount(100000000)
	fmt.Println(accountinfo2)
	machine.TxPool.StatDB.Store(accountinfo2.Address, accountinfo2.Account)

	actualBalance := machine.TxPool.StatDB.Load(accountinfo1.Address)
	actualBalance2 := machine.TxPool.StatDB.Load(accountinfo2.Address)

	fmt.Println(actualBalance)
	fmt.Println(actualBalance2)

	tx:=common.NewTransaction(accountinfo2.Address[:],100,100000,100000,100000,[]byte(""))
	tx.Sign(accountinfo1.PrivateKey)
	fmt.Println(tx)
	// 测试用例1: 正常交易
	t.Run("Normal transaction", func(t *testing.T) {
		machine.Execute(*tx)
		// 验证结果
		// ...
		// 验证账户余额变化
		expectedBalance := accountinfo1.Account.Balance-tx.Value-tx.GasPrice
		actualBalance := machine.TxPool.StatDB.Load(accountinfo1.Address).Balance
		if expectedBalance != actualBalance {
			t.Errorf("Balance mismatch, expected: %v, got: %v", expectedBalance, actualBalance)
		}
	})
	
	// 测试用例2: 余额不足
	// t.Run("Insufficient balance", func(t *testing.T) {
	// 	// ...
	// })
}