package statdbtest

import (
	"bytes"
	"cxchain-2023131015/common"
	"cxchain-2023131015/kvstore/leveldb"
	"cxchain-2023131015/statdb"
	"cxchain-2023131015/trie/mpt"
	"fmt"
	"os"
	"testing"
)

// setupTestDB 创建测试用的LevelDB (重命名函数，不再以Test开头)
func setupTestDB() (*leveldb.LevelDBStore, func()) {
	db, err := leveldb.NewLevelDBStore("test/testdb")
	if err != nil {
		panic(err)
	}
	return db, func() {
		db.Close()
		os.RemoveAll("testdb")
	}
}

func TestStatDB_BasicOperations(t *testing.T) {
	db, cleanup := setupTestDB() // 使用重命名后的函数
	defer cleanup()

	mpt := mpt.NewMPT(db)
	statDB := statdb.NewStatDB(mpt)

	// 测试账户
	addr := common.Address{1}
	account := common.Account{
		Nonce:    1,
		Balance:  100,
		Codehash: []byte("codehash"),
		Root:     []byte("root"),
	}

	// 测试Store方法
	err := statDB.Store(addr, account)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// 测试Load方法
	loadedAccount := statDB.Load(addr)
	fmt.Println(loadedAccount)
	if !equalAccounts(loadedAccount, account) {
		t.Errorf("Loaded account mismatch, expected: %+v, got: %+v", account, loadedAccount)
	}

	// 测试更新账户
	account.Nonce = 2
	account.Balance = 200
	err = statDB.Store(addr, account)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// 验证更新
	updatedAccount := statDB.Load(addr)
	if !equalAccounts(updatedAccount, account) {
		t.Errorf("Updated account mismatch, expected: %+v, got: %+v", account, updatedAccount)
	}

	// 测试不存在的地址
	nonExistAddr := common.Address{2}
	nonExistAccount := statDB.Load(nonExistAddr)
	if !equalAccounts(nonExistAccount, common.Account{}) {
		t.Errorf("Non-existent address should return empty account, got: %+v", nonExistAccount)
	}
}

func TestStatDB_EdgeCases(t *testing.T) {
	db, cleanup := setupTestDB() // 使用重命名后的函数
	defer cleanup()

	mpt := mpt.NewMPT(db)
	statDB := statdb.NewStatDB(mpt)

	// 测试空地址
	emptyAddr := common.Address{}
	emptyAccount := common.Account{}
	err := statDB.Store(emptyAddr, emptyAccount)
	if err != nil {
		t.Fatalf("Store empty address failed: %v", err)
	}

	// 测试空账户
	addr := common.Address{3}
	err = statDB.Store(addr, common.Account{})
	if err != nil {
		t.Fatalf("Store empty account failed: %v", err)
	}

	// 测试损坏的数据
	corruptedAddr := common.Address{4}
	err = db.Put(corruptedAddr[:], []byte("corrupted data"))
	if err != nil {
		t.Fatalf("Put corrupted data failed: %v", err)
	}
	corruptedAccount := statDB.Load(corruptedAddr)
	if !equalAccounts(corruptedAccount, common.Account{}) {
		t.Errorf("Corrupted data should return empty account, got: %+v", corruptedAccount)
	}
}

func equalAccounts(a, b common.Account) bool {
	if a.Nonce != b.Nonce {
		return false
	}
	if a.Balance != b.Balance {
		return false
	}
	if !bytes.Equal(a.Codehash, b.Codehash) {
		return false
	}
	if !bytes.Equal(a.Root, b.Root) {
		return false
	}
	return true
}