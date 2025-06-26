package mpt

import (
	"bytes"
	"cxchain-2023131015/kvstore/leveldb"
	"fmt"
	"testing"
)

func TestMPTHasWithAddress(t *testing.T) {
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mpt := NewMPT(db)
	
	// 测试数据
	var address [20]byte
	copy(address[:], []byte{
	    0x01, 0x02, 0x03, 0x04, 0x05,
	    0x06, 0x07, 0x08, 0x09, 0x0a,
	    0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	    0x10, 0x11, 0x12, 0x13, 0x14,
	})
	value := []byte("test value")
	
	// 测试空树
	t.Run("EmptyTrie", func(t *testing.T) {
		_, err := mpt.Has(address[:])
		if err == nil || err.Error() != "empty trie" {
			t.Errorf("期望返回empty trie错误，但得到: %v", err)
		}
	})
	
	// 测试插入后查询
	t.Run("AfterInsert", func(t *testing.T) {
		if err := mpt.Insert(address[:], value); err != nil {
			t.Fatalf("插入失败: %v", err)
		}
		
		got, err := mpt.Has(address[:])
		if err != nil {
			t.Fatalf("Has方法出错: %v", err)
		}
		if !bytes.Equal(got, value) {
			t.Errorf("期望值: %v, 得到: %v", value, got)
		}
	})
	
	// 测试不存在的地址
	t.Run("NonExistAddress", func(t *testing.T) {
		if err := mpt.Insert(address[:], value); err != nil {
			t.Fatalf("插入失败: %v", err)
		}

		var nonExistAddr [20]byte
		copy(nonExistAddr[:], []byte{
	    0x32, 0x16, 0x17, 0x18, 0x19,
	    0x1a, 0x1b, 0xdc, 0x1d, 0x1e,
	    0x1f, 0x20, 0x21, 0x32, 0x23,
	    0x24, 0x25, 0x26, 0x27, 0x29,
		})


		v, err := mpt.Has(nonExistAddr[:])
		if err == nil {
			fmt.Println(v)
			t.Error("期望返回错误，但得到nil")
		}

		got, err := mpt.Has(address[:])
		if err != nil {
			t.Fatalf("Has方法出错: %v", err)
		}
		if !bytes.Equal(got, value) {
			t.Errorf("期望值: %v, 得到: %v", value, got)
		}
	})
}
func TestMPTWithAddress(t *testing.T) {
	// 初始化MPT
	db, err := leveldb.NewLevelDBStore("testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mpt := NewMPT(db)
	
	// 测试数据 - 使用[20]byte作为地址
	// 使用有效的十六进制地址
	var address1 [20]byte
	copy(address1[:], []byte{
	    0x01, 0x02, 0x03, 0x04, 0x05,
	    0x06, 0x07, 0x08, 0x09, 0x0a,
	    0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	    0x10, 0x11, 0x12, 0x13, 0x14,
	})
	var address2 [20]byte
	copy(address2[:], []byte{
	    0x15, 0x16, 0x17, 0x18, 0x19,
	    0x1a, 0x1b, 0x1c, 0x1d, 0x1e,
	    0x1f, 0x20, 0x21, 0x22, 0x23,
	    0x24, 0x25, 0x26, 0x27, 0x28,
	})
	var address3 [20]byte
	copy(address3[:], []byte{
	    0x23, 0x16, 0x17, 0x18, 0x19,
	    0x1a, 0x1b, 0x1c, 0x1d, 0x1e,
	    0x1f, 0x20, 0x21, 0x22, 0x23,
	    0x24, 0x25, 0x26, 0x27, 0x28,
	})
	
	value1 := []byte("value1")
	value2 := []byte("value2")
	
	// 测试插入
	t.Run("Insert", func(t *testing.T) {
		if err := mpt.Insert(address1[:], value1); err != nil {
			t.Fatalf("插入address1失败: %v", err)
		}
		
		if err := mpt.Insert(address2[:], value2); err != nil {
			t.Fatalf("插入address2失败: %v", err)
		}
		
		v, err := mpt.Has(address1[:])
		if err != nil {
			t.Fatalf("Has方法出错: %v", err)
		}
		if !bytes.Equal(v, value1) {
			t.Errorf("期望值: %v, 得到: %v", value1, v)
		}
	})	
}