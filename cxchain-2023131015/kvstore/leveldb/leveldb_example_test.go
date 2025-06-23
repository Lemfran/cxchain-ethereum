package leveldb

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestLevelDBStore(t *testing.T) {
	dbPath := "testdb"
	defer os.RemoveAll(dbPath)
	store, err := NewLevelDBStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	// Test Put
	err = store.Put([]byte("key1"), []byte("value1"))
	assert.NoError(t, err)
	// Test Get
	value, err := store.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)
	// Test Has
	exists, err := store.Has([]byte("key1"))
	assert.NoError(t, err)
	assert.True(t, exists)
	// Test Delete
	err = store.Delete([]byte("key1"))
	assert.NoError(t, err)
	// Test Has after Delete
	exists, err = store.Has([]byte("key1"))
	assert.NoError(t, err)
	assert.False(t, exists)
	// Test Batch Write
	batch := new(leveldb.Batch)
	batch.Put([]byte("foo"), []byte("value"))
	batch.Put([]byte("bar"), []byte("another value"))
	batch.Delete([]byte("baz"))
	err = store.db.Write(batch, nil)
	assert.NoError(t, err)
	// Test Iterate
	iter := store.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("key: %s, value: %s\n", key, value)
	}
	iter.Release()
	err = iter.Error()
	assert.NoError(t, err)
}