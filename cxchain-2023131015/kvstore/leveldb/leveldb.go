package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelDBStore struct {
	db *leveldb.DB
}

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

func (store *LevelDBStore) Get(key []byte) ([]byte, error) {
	return store.db.Get(key, nil)
}

func (store *LevelDBStore) Put(key, value []byte) error {
	return store.db.Put(key, value, nil)
}

func (store *LevelDBStore) Delete(key []byte) error {
	return store.db.Delete(key, nil)
}

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

func (store *LevelDBStore) Close() error {
	return store.db.Close()
}