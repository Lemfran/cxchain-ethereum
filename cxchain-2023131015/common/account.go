package common


type Account struct {
	Nonce    uint64
	Balance  uint64
	Codehash []byte
	Root     []byte
}


func NewAccount() *Account {
	return &Account{
		Nonce:    0,
		Balance:  0,
		Codehash: nil,
		Root:     nil,
	}
}
