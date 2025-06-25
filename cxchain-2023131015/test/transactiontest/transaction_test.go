package test

import (
	"math/big"
	"testing"
	
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	
	"cxchain-2023131015/common"
)

func TestTransactionSignAndFrom(t *testing.T) {
	// 生成测试私钥
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	
	// 准备测试交易
	txData := common.Txdata{
        To:       []byte{0x00, 0x11},
        Value:    100,
        Nonce:    1,
        GasLimit: 2000,
        GasPrice: 200,
        Input:    []byte{0x01, 0x02},
    }

    sig := common.Signature{
        R: new(big.Int).SetInt64(1),
        S: new(big.Int).SetInt64(2),
        V: 27,
    }

    // 创建Transaction实例
    tx := common.Transaction{
        Txdata:   txData,
        Signature: sig,
    }
	
	// 测试签名
	err = tx.Sign(crypto.FromECDSA(privateKey))
	assert.NoError(t, err)
	assert.NotNil(t, tx.Signature.R)
	assert.NotNil(t, tx.Signature.S)
	assert.NotZero(t, tx.Signature.V)
	
	// 测试From方法恢复地址
	recoveredAddr := tx.From()
	expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	commonAddr := common.Address(expectedAddr)
	assert.Equal(t, commonAddr, common.Address(recoveredAddr))
}

func TestTransactionHash(t *testing.T) {
	// 准备测试交易
	tx := common.Transaction{
		Txdata: common.Txdata{
		// 由于 common.HexToAddress 未定义，手动将十六进制字符串转换为字节切片
		To:       []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23},
			Value:    100,
			Nonce:    1,
			GasLimit: 21000,
			GasPrice: 10,
		},
	}
	
	// 测试哈希计算
	hash1 := tx.Hash()
	hash2 := tx.Hash()
	assert.Equal(t, hash1, hash2) // 相同交易哈希应相同
	
	// 修改交易内容后哈希应不同
	tx.Txdata.Value = 200
	hash3 := tx.Hash()
	assert.NotEqual(t, hash1, hash3)
}

func TestInvalidSignature(t *testing.T) {
	// 准备测试交易
	tx := common.Transaction{
		Txdata: common.Txdata{
// 由于 common.HexToAddress 未定义，手动将十六进制字符串转换为字节切片
To:       []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x23},
			Value:    100,
			Nonce:    1,
		},
	}
	
	// 测试无效私钥签名
	err := tx.Sign([]byte("invalid_private_key"))
	assert.Error(t, err)
	
	// 测试恢复无效签名
	tx.Signature = common.Signature{
		R: big.NewInt(0),
		S: big.NewInt(0),
		V: 0,
	}
	addr := tx.From()
	assert.Equal(t, common.Address{}, addr)
}