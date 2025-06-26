package common

import (
	"github.com/ethereum/go-ethereum/crypto"
)


type AccountInfo struct {
	PrivateKey []byte
	Address    Address
	Account    Account
}

func GenerateAccount(balance uint64) (*AccountInfo, error) {

	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privKeyBytes := crypto.FromECDSA(privKey)

	addr ,_:= PrivateKeyToAddress(privKeyBytes)
	account := Account{
		Balance: balance,
		Nonce:   0,
	}

	return &AccountInfo{
		PrivateKey: privKeyBytes,
		Address:    Address(addr),
		Account:    account,
	}, nil
}



