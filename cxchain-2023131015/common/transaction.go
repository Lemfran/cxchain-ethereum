package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type Txdata struct {
	To       []byte
	Value    uint64
	Nonce    uint64
	GasLimit uint64
	GasPrice uint64
	Input    []byte
}

type Signature struct {
	R, S *big.Int
	V    uint8
}

type Transaction struct {
	Txdata 
	Signature
}

type Transactioner interface {
	From() Address
	Sign(privkey []byte) error
}

func NewTransaction(to []byte, value uint64, nonce uint64, gasLimit uint64, gasPrice uint64, input []byte) *Transaction {
	return &Transaction{
		Txdata: Txdata{
			To:       to,
			Value:    value,
			Nonce:    nonce,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Input:    input,
		},
		Signature: Signature{},
	}
}

// Sign 实现Transactioner接口的Sign方法
func (tx *Transaction) Sign(privkey []byte) error {
	// 将私钥转换为ECDSA格式
	privateKey, err := crypto.ToECDSA(privkey)
	if err != nil {
		return err
	}
	
	// 计算交易哈希
	hash := tx.Hash()
	
	// 使用私钥签名
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		return err
	}
	
	// 解析签名结果
	tx.Signature.R = new(big.Int).SetBytes(signature[:32])
	tx.Signature.S = new(big.Int).SetBytes(signature[32:64])
	tx.Signature.V = signature[64] + 27 // 以太坊的V值调整
	
	return nil
}

// From 实现Transactioner接口的From方法，用于从签名恢复地址
func (tx *Transaction) From() Address {
	// 准备签名数据
	sig := make([]byte, 65)
	copy(sig[:32], tx.Signature.R.Bytes())
	copy(sig[32:64], tx.Signature.S.Bytes())
	sig[64] = tx.Signature.V - 27 // 还原V值
	
	// 计算交易哈希
	hash := tx.Hash()
	
	// 恢复公钥
	pubKey, err := crypto.Ecrecover(hash[:], sig)
	if err != nil {
		return Address{}
	}
	
	// 将公钥转换为地址
	return PubKeyToAddress(pubKey[1:]) // 去掉第一个字节(0x04)
}

// Hash 计算交易的哈希值
func (tx *Transaction) Hash() Hash {
	// 使用RLP编码交易数据
	encoded, err := rlp.EncodeToBytes(tx.Txdata)
	if err != nil {
		return Hash{}
	}
	// 计算Keccak256哈希
	// 假设Hash类型与github.com/ethereum/go-ethereum/common.Hash兼容，直接进行类型转换
	result := crypto.Keccak256Hash(encoded)
	return Hash(result)
}

func (tx *Transaction) Serialize() []byte {
	// 序列化交易数据
	data, err := rlp.EncodeToBytes(tx.Txdata)
	if err != nil {
		panic(err)
	}
	return data
}


