package test

import (
	"testing"
	
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	
	"cxchain-2023131015/common"
)

func TestPrivateKeyToPubkey(t *testing.T) {
	// 生成测试私钥
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	// 测试私钥转公钥
	privBytes := crypto.FromECDSA(privateKey)
	pubkey, err := common.PrivateKeyToPubkey(privBytes)
	assert.NoError(t, err)
	assert.Equal(t, 64, len(pubkey))
	
	// 验证公钥是否正确
	expectedPubkey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	assert.Equal(t, expectedPubkey, pubkey)
}

func TestPubKeyToAddress(t *testing.T) {
	// 生成测试私钥
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	// 获取公钥
	pubkey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	
	// 测试公钥转地址
	address := common.PubKeyToAddress(pubkey)
	
	// 验证地址是否正确
	expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	// 转换为common.Address类型
	commonAddr := common.Address(expectedAddr)
	assert.Equal(t, commonAddr, common.Address(address))
}

func TestPrivateKeyToAddress(t *testing.T) {
	// 生成测试私钥
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	// 测试私钥转地址
	privBytes := crypto.FromECDSA(privateKey)
	address, err := common.PrivateKeyToAddress(privBytes)
	assert.NoError(t, err)
	
	// 验证地址是否正确
	expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	// 转换为common.Address类型
	commonAddr := common.Address(expectedAddr)
	// 转换为common.Address类型
	commonAddr2 := common.Address(address)
	assert.Equal(t, commonAddr2, commonAddr)
}

func TestInvalidPrivateKey(t *testing.T) {
	// 测试无效私钥
	_, err := common.PrivateKeyToPubkey([]byte("invalid_private_key"))
	assert.Error(t, err)
	
	_, err = common.PrivateKeyToAddress([]byte("invalid_private_key"))
	assert.Error(t, err)
}

func TestInvalidPubkey(t *testing.T) {
	// 测试无效公钥
	addr := common.PubKeyToAddress([]byte("invalid_pubkey"))
	assert.Equal(t, common.Address{}, addr)
}