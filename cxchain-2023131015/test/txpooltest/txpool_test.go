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

	txs := txPool.Pop()
	fmt.Println(txs)

	txs1 := txPool.Pop()
	fmt.Println(txs1)

}


func TestReplacePendingTx(t *testing.T) {
	// 测试替换pending交易
	// ...
}