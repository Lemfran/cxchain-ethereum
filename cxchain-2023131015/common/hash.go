package common

import (
	"crypto/sha3"
)

type Hash [32]byte

func NewHash(value []byte) Hash {
	hasher := sha3.New256()
	hasher.Write(value)
	return Hash(hasher.Sum(nil))
}
