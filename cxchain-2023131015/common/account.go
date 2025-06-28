package common

import (
	"bytes"
	"encoding/json"
)

type Account struct {
	Nonce    uint64 `json:"nonce"`
	Balance  uint64 `json:"balance"`
	Codehash []byte `json:"codehash"`
	Root     []byte `json:"root"`
}

func NewAccount() *Account {
	return &Account{
		Nonce:    0,
		Balance:  0,
		Codehash: nil,
		Root:     nil,
	}
}

func (a *Account) Serialize() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Account) Deserialize(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, a)
}

//判断账户是否相等
func EqualAccounts(a, b Account) bool {
	if a.Nonce != b.Nonce {
		return false
	}
	if a.Balance != b.Balance {
		return false
	}
	if !bytes.Equal(a.Codehash, b.Codehash) {
		return false
	}
	if !bytes.Equal(a.Root, b.Root) {
		return false
	}
	return true
}