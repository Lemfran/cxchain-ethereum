package statdb

import (
	"cxchain-2023131015/trie"
	"cxchain-2023131015/common"
)

type StatDB interface {
	Load(addr common.Address) common.Account
	Store(addr common.Address, account common.Account) error
}

type statDB struct {
	db *mpt.MPT
}

func (s *statDB) Load(addr common.Address) common.Account {
	//从mpt中根据地址获取用户状态
	node, err := s.mpt.Has(addr)
	if err != nil {
		return common.Account{}
	}

	var account common.Account
	err = common.Deserialize(node, &account)
	if err != nil {
		return common.Account{}
	}
	return account
}

func (s *statDB) Store(addr common.Address, account common.Account) error {
	node, err := account.Serialize()
	if err != nil {
		return err
	}
	//将用户状态序列化后存储到mpt中
	err = s.mpt.Insert(addr, node)
	if err != nil {
		return err
	}
	return nil
}
