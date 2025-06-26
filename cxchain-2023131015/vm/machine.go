package vm

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/txpool"
	"fmt"
)

type IMachine interface {
	Execute(tx common.Transaction) error


	Execute1(tx common.Transaction) *common.Receipt
}

type StateMachine struct {
	TxPool *txpool.TxPool
}

func NewStateMachine(txPool *txpool.TxPool) *StateMachine {
	return &StateMachine{TxPool: txPool}
}

func (m *StateMachine) Execute( tx common.Transaction) {

	from := tx.From()
	to := tx.To
	value := tx.Value
	gasUsed := tx.GasPrice
	fmt.Println(from)
	fmt.Println(gasUsed)
	if tx.GasPrice < 21000 {
		return
	} else {
		gasUsed = 1
	}
	gasUsed = gasUsed * tx.GasPrice
	cost := value + gasUsed

	fromAddr := common.Address(from)
	var accountfrom common.Account
	accountfrom = m.TxPool.StatDB.Load(fromAddr)
	
	fmt.Println(accountfrom)
	if accountfrom.Balance < cost {
		return
	}

	accountfrom.Balance = accountfrom.Balance - cost

	accountfrom.Nonce = accountfrom.Nonce + 1

	m.TxPool.StatDB.Store(fromAddr, accountfrom)


	toAddr := common.Address(to)
	
	accountto := m.TxPool.StatDB.Load(toAddr)
	fmt.Println(accountto)
	if common.EqualAccounts(accountto, common.Account{}) {
		accountto = common.Account{
			Balance: value,
			Nonce:   0,
		}
	}

	accountto.Balance= accountto.Balance + value

	m.TxPool.StatDB.Store(toAddr, accountto)
}
