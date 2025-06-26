package statdb

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/trie/mpt"
	"errors"
	"fmt"
)

type StatDB struct {
	db *mpt.MPT
}

type StatDB_Interface interface {
	Load(addr common.Address) common.Account
	Store(addr common.Address, account common.Account) error
}


func NewStatDB(db *mpt.MPT) *StatDB {
	return &StatDB{
		db: db,
	}
}

func (s *StatDB) Load(addr common.Address) common.Account {
	//从mpt中根据地址获取用户状态
	addrBytes := addr[:]

	value, err := s.db.Has(addrBytes)
	if err != nil {
		fmt.Println("not found")
		return common.Account{}
	}
	fmt.Println("3")
	var account common.Account
	err = account.Deserialize(value)
	if err != nil {
		return common.Account{}
	}
	fmt.Println("1")
	return account
}

func (s *StatDB) Store(addr common.Address, account common.Account) error {
	node, err := account.Serialize()
	addrBytes := addr[:]
	if err != nil {
		return errors.New("Serialize Error")
	}
	//将用户状态序列化后存储到mpt中
	err = s.db.Insert(addrBytes, node)
	if err != nil {
		return err
	}
	fmt.Println("2")
	return nil
}
