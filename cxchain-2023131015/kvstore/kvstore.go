package kvstore

import "io"
 
type KVStore interface {
 // Get retrieves the value associated with the given key.
 Get(key []byte) ([]byte, error)
 // Put stores the value with the given key.
 Put(key, value []byte) error
 // Delete removes the value associated with the given key.
 Delete(key []byte) error
 // Has checks if the key exists in the database.
 Has(key []byte) (bool, error)
 io.Closer
}

