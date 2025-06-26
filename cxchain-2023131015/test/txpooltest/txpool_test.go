package txpooltest

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/kvstore/leveldb"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/trie/mpt"
	"cxchain-2023131015/txpool"
	"fmt"

	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func setupTestDB() (*mpt.MPT, func()) {
	// 创建测试用的LevelDB和MPT
	// ...
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		panic(err)
	}
	mptDB := mpt.NewMPT(db)
	return mptDB, func() {
		db.Close()
	}
}

func setupTestTXPool() (*txpool.TxPool, func()) {
	mptDB, closeDB := setupTestDB()
	statDB := statdb.NewStatDB(mptDB)
	txPool := txpool.NewTxPool(statDB)
	return txPool, func() {
		closeDB()
	}
}

func TestNewTX(t *testing.T) {
	txPool, closeDB := setupTestTXPool()
	defer closeDB()
	
	privKey, _ := crypto.GenerateKey()
	fromAddr, _ := common.PrivateKeyToAddress(crypto.FromECDSA(privKey))
	

	fmt.Println(fromAddr)
	// 初始化账户并设置正确的nonce
	account := common.NewAccount()
	account.Nonce = 0  // 设置初始nonce
	txPool.StatDB.Store(fromAddr, *account)
	

	// 创建交易时使用匹配的nonce
	tx := common.NewTransaction(common.Address{}.Bytes(), 1, 1, 0, 0, nil)
	tx.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx)

	tx1 := common.NewTransaction(common.Address{}.Bytes(), 1, 2, 0, 123, nil)
	tx1.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx1)

	err := txPool.NewTX(tx)
	err1 := txPool.NewTX(tx1)
	if err != nil {
		t.Fatalf("NewTX failed: %v", err)
	}
	if err1 != nil {
		t.Fatalf("NewTX failed: %v", err1)
	}

	fmt.Println("ok2")

	txs := txPool.Pop()
	fmt.Println(txs)

	fmt.Println("okkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk2")

	txs1 := txPool.Pop()
	fmt.Println(txs1)

}

func TestNewTX2(t *testing.T) {
	txPool, closeDB := setupTestTXPool()
	defer closeDB()
	
	privKey, _ := crypto.GenerateKey()
	fromAddr, _ := common.PrivateKeyToAddress(crypto.FromECDSA(privKey))
	

	fmt.Println(fromAddr)
	// 初始化账户并设置正确的nonce
	account := common.NewAccount()
	account.Nonce = 0  // 设置初始nonce
	txPool.StatDB.Store(fromAddr, *account)
	

	// 创建交易时使用匹配的nonce
	tx := common.NewTransaction(common.Address{}.Bytes(), 1, 1, 0, 0, nil)
	tx.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx)

	tx1 := common.NewTransaction(common.Address{}.Bytes(), 1, 2, 0, 125, nil)
	tx1.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx1)

	tx2 := common.NewTransaction(common.Address{}.Bytes(), 1, 3, 0, 2, nil)
	tx2.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx2)

	tx3 := common.NewTransaction(common.Address{}.Bytes(), 1, 1, 0, 124, nil)
	tx3.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx3)

	err := txPool.NewTX(tx)
	err1 := txPool.NewTX(tx1)
	err2 := txPool.NewTX(tx2)

	err3 := txPool.NewTX(tx3)

	if err != nil {
		t.Fatalf("NewTX failed: %v", err)
	}
	if err1 != nil {
		t.Fatalf("NewTX failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("NewTX failed: %v", err2)
	}
	if err3 != nil {
		t.Fatalf("NewTX failed: %v", err3)
	}

	fmt.Println("ok2")

	txs := txPool.Pop()
	fmt.Println(txs)

	fmt.Println("okkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk2")

	txs1 := txPool.Pop()
	fmt.Println(txs1)

	txs3 := txPool.Pop()
	fmt.Println(txs3)

}

func TestNewTX3(t *testing.T) {
	txPool, closeDB := setupTestTXPool()
	defer closeDB()
	
	privKey, _ := crypto.GenerateKey()
	fromAddr, _ := common.PrivateKeyToAddress(crypto.FromECDSA(privKey))
	

	fmt.Println(fromAddr)
	// 初始化账户并设置正确的nonce
	account := common.NewAccount()
	account.Nonce = 0  // 设置初始nonce
	txPool.StatDB.Store(fromAddr, *account)
	

	// 创建交易时使用匹配的nonce
	tx := common.NewTransaction(common.Address{}.Bytes(), 1, 1, 0, 0, nil)
	tx.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx)

	tx1 := common.NewTransaction(common.Address{}.Bytes(), 1, 3, 0, 123, nil)
	tx1.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx1)

	err := txPool.NewTX(tx)
	err1 := txPool.NewTX(tx1)
	if err != nil {
		t.Fatalf("NewTX failed: %v", err)
	}
	if err1 != nil {
		t.Fatalf("NewTX failed: %v", err1)
	}

	tx2 := common.NewTransaction(common.Address{}.Bytes(), 1, 2, 0, 12, nil)
	tx2.Sign(crypto.FromECDSA(privKey))
	fmt.Println(tx2)

	err2 := txPool.NewTX(tx2)
	if err2 != nil {
		t.Fatalf("NewTX failed: %v", err2)
	}

	
	fmt.Println("ok2")

	txs := txPool.Pop()
	fmt.Println(txs)

	fmt.Println("okkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk2")

	txs1 := txPool.Pop()
	fmt.Println(txs1)

	txs2 := txPool.Pop()
	fmt.Println(txs2)

}