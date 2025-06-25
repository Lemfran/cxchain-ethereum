package common

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type Address [20]byte

func (a Address) Bytes() any {
	panic("unimplemented")
}

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
	// 接受65字节(含0x04前缀)或64字节(不含前缀)的公钥
	if len(pubkey) != 64 && len(pubkey) != 65 {
		return Address{}
	}

	// 如果是65字节，去掉第一个字节(0x04)
	if len(pubkey) == 65 {
		pubkey = pubkey[1:]
	}

	// 使用Keccak256哈希算法
	hash := crypto.Keccak256(pubkey)
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
