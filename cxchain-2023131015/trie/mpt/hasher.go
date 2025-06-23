package mpt

import "crypto/sha3"

func Hash(value []byte) []byte {
	hasher :=sha3.New256()
	hasher.Write(value)
	return hasher.Sum(nil)
}