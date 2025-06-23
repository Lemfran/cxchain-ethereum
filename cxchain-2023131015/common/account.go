package common

import "encoding/json"

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