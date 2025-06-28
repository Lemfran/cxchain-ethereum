package common

import (
	"github.com/ethereum/go-ethereum/crypto"
)

// Address 类型表示一个20字节的区块链地址
type Address [20]byte

// Bytes 将Address类型转换为字节切片
func (a Address) Bytes() []byte {
	return a[:]
}

// Addresser 接口定义了地址相关的操作
// PrivateKeyToPubkey: 从私钥生成公钥
// PubKeyToAddress: 从公钥生成地址
// PrivateKeyToAddress: 从私钥直接生成地址
type Addresser interface {
	PrivateKeyToPubkey(privkey []byte) ([]byte, error)
	PubKeyToAddress(pubkey []byte) Address
	PrivateKeyToAddress(privkey []byte) (Address, error)
}

// PrivateKeyToPubkey 从私钥生成公钥
// privkey: 私钥字节切片
// 返回: 公钥字节切片或错误
func PrivateKeyToPubkey(privkey []byte) ([]byte, error) {
	// 将字节切片转换为ECDSA私钥
	privateKey, err := crypto.ToECDSA(privkey)
	if err != nil {
		return nil, err
	}

	// 拼接X和Y坐标作为公钥
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return pubKey, nil
}

// PubKeyToAddress 从公钥生成地址
// pubkey: 公钥字节切片(64或65字节)
// 返回: 生成的地址
func PubKeyToAddress(pubkey []byte) Address {
	// 检查公钥长度是否合法(64或65字节)
	if len(pubkey) != 64 && len(pubkey) != 65 {
		return Address{}  // 返回空地址
	}

	// 如果是65字节的公钥(含0x04前缀)，去掉第一个字节
	if len(pubkey) == 65 {
		pubkey = pubkey[1:]
	}

	// 使用Keccak256算法哈希公钥
	hash := crypto.Keccak256(pubkey)
	var addr Address
	// 取哈希的最后20字节作为地址
	copy(addr[:], hash[12:])
	return addr
}

// PrivateKeyToAddress 从私钥直接生成地址
// privkey: 私钥字节切片
// 返回: 生成的地址或错误
func PrivateKeyToAddress(privkey []byte) (Address, error) {
	// 先通过私钥生成公钥
	pubkey, err := PrivateKeyToPubkey(privkey)
	if err != nil {
		return Address{}, err
	}
	// 再从公钥生成地址
	return PubKeyToAddress(pubkey), nil
}
