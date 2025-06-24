package common

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/crypto"
)


type Address [20]byte

type Addresser interface {
	PrivateKeyToPubkey(privkey []byte) ([]byte, error)
	PubKeyToAddress(pubkey []byte) Address
	PrivateKeyToAddress(privkey []byte) (Address, error)
}

func PrivateKeyToPubkey(privkey []byte) ([]byte, error) {
	privateKey, err := crypto.ToECDSA(privkey)
	if err != nil {
		return nil, err
	}
	
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return pubKey, nil
}

func PubKeyToAddress(pubkey []byte) Address {
	if len(pubkey) != 64 {
		return Address{}
	}
	
	hash := sha256.Sum256(pubkey)
	var addr Address
	copy(addr[:], hash[12:])
	return addr
}

func PrivateKeyToAddress(privkey []byte) (Address, error) {
	pubkey, err := PrivateKeyToPubkey(privkey)
	if err != nil {
		return Address{}, err
	}
	return PubKeyToAddress(pubkey), nil
}

