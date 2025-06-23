package leveldb

import (
	"fmt"
	"testing"
)

func BenchmarkPut(b *testing.B) {
	dbPath := "testdb"
	store, err := NewLevelDBStore(dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	for i := 0; i < b.N; i++ {
		err = store.Put([]byte(fmt.Sprintf("key%d", i)), []byte("value"))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	dbPath := "testdb"
	store, err := NewLevelDBStore(dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	// Populate the database
	for i := 0; i < 1000; i++ {
		key := []byte(fmt.Sprintf("key%d", i))
		err = store.Put(key, []byte("value"))
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key%d", i % 1000)) // Use modulo to avoid too many keys
		_, err = store.Get(key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	dbPath := "testdb"
	store, err := NewLevelDBStore(dbPath)
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	// Populate the database
	for i := 0; i < 1000; i++ {
		err = store.Put([]byte(fmt.Sprintf("key%d", i)), []byte("value"))
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = store.Delete([]byte(fmt.Sprintf("key%d", i)))
		if err != nil {
			b.Fatal(err)
		}
	}
}