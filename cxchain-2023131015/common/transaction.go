package common

import (
	"math/big"
)

type Transaction struct {
	txdata
	signature
}

type txdata struct {
	To       []byte
	Value    uint64
	Nonce    uint64
	GasLimit uint64
	GasPrice uint64
	Input    []byte
}

type signature struct {
	R, S *big.Int
	V    uint8
}

type Transactioner interface {
	From() Address
	Sign(privkey []byte) error
}

func NewTransaction(to []byte, value uint64, nonce uint64, gasLimit uint64, gasPrice uint64, input []byte) *Transaction {
	return &Transaction{
		txdata: txdata{
			To:       to,
			Value:    value,
			Nonce:    nonce,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Input:    input,
		},
		signature: signature{},
	}
}

/*func Ecrecover(hash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(hash, sig)
}

func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.Unmarshal(S256(), s)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}*/

func (tx *Transaction) From() Address {
	/*txdata := tx.txdata
	toSign, _ := rlp.EncodeToBytes(txdata)
	msg := sha3.Keccak256(toSign)
	sig := make([]byte, 65)
	// 把RSV写进去
	copy(sig[32-len(tx.signature.R.Bytes()):], tx.signature.R.Bytes())
	copy(sig[64-len(tx.signature.S.Bytes()):], tx.signature.S.Bytes())
	sig[64] = tx.signature.V
	pubKey, err := secp256k1.RecoverPubkey(msg, sig)
	if err != nil {
		// TODO
		// return
	}
	return PubKeyToAddress(pubKey)*/
	//暂且默认输出
	return Address{0x01}
}


