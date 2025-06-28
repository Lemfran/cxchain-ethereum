package statdb

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/trie/mpt"
	"errors"
	"fmt"
)

// StatDB 结构体表示状态数据库，使用MPT(Merkle Patricia Trie)作为底层存储
// db字段指向一个MPT实例，用于存储账户状态
type StatDB struct {
	db *mpt.MPT
}

// StatDB_Interface 定义了状态数据库的接口
// 包含加载和存储账户状态的方法
type StatDB_Interface interface {
	Load(addr common.Address) common.Account
	Store(addr common.Address, account common.Account) error
}

// GetRoot 获取状态数据库的根哈希
// 返回MPT的根节点哈希值
func (s *StatDB) GetRoot() common.Hash {
	return common.Hash(s.db.Root.GetHash())
}

// NewStatDB 创建一个新的状态数据库实例
// 参数db是初始化好的MPT实例
// 返回一个指向StatDB的指针
func NewStatDB(db *mpt.MPT) *StatDB {
	return &StatDB{
		db: db,
	}
}

// Load 根据地址加载账户状态
// 从MPT中查询指定地址的账户数据并反序列化
// 如果查询失败或反序列化失败，返回空账户
func (s *StatDB) Load(addr common.Address) common.Account {
	//从mpt中根据地址获取用户状态
	addrBytes := addr[:]
	value, err := s.db.Has(addrBytes)
	if err != nil {
		fmt.Println("not found")
		return common.Account{}
	}
	var account common.Account
	err = account.Deserialize(value)
	if err != nil {
		return common.Account{}
	}
	fmt.Println("获取到用户状态")
	return account
}

// Store 存储账户状态到指定地址
// 将账户数据序列化后存入MPT
// 如果序列化或插入失败，返回错误
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
	return nil
}
