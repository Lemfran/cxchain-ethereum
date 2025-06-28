package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDBStore 结构体封装了LevelDB数据库实例
// db字段指向一个LevelDB数据库实例
type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore 创建一个新的LevelDB存储实例
// 参数dbPath是数据库文件路径
// 使用Bloom过滤器(误判率10%)优化查询性能
// 返回LevelDBStore指针和可能的错误
func NewLevelDBStore(dbPath string) (*LevelDBStore, error) {
	options := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(dbPath, options)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{db: db}, nil
}

// Get 根据key获取对应的value
// 参数key是要查询的键
// 返回对应的值和可能的错误
func (store *LevelDBStore) Get(key []byte) ([]byte, error) {
	return store.db.Get(key, nil)
}

// Put 存储键值对到数据库
// 参数key是要存储的键
// 参数value是要存储的值
// 返回可能的错误
func (store *LevelDBStore) Put(key, value []byte) error {
	return store.db.Put(key, value, nil)
}

// Delete 从数据库中删除指定的key
// 参数key是要删除的键
// 返回可能的错误
func (store *LevelDBStore) Delete(key []byte) error {
	return store.db.Delete(key, nil)
}

// Has 检查数据库中是否存在指定的key
// 参数key是要检查的键
// 返回bool表示是否存在和可能的错误
func (store *LevelDBStore) Has(key []byte) (bool, error) {
	_, err := store.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Close 关闭数据库连接
// 返回可能的错误
func (store *LevelDBStore) Close() error {
	return store.db.Close()
}