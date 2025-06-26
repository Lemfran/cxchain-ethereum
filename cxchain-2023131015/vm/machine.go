package vm

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/statdb"
)

type IMachine interface {
	Execute(state statdb.StatDB, tx common.Transaction)

	Execute1(state statdb.StatDB, tx common.Transaction) *common.Receipt
}

type StateMachine struct {
	StatDB statdb.StatDB
}

func (m *StateMachine) Execute(state statdb.StatDB, tx common.Transaction) {
	from := tx.From()
	to := tx.To
	value := tx.Value
	gasUsed := tx.GasPrice
	if tx.GasPrice < 21000 {
		return
	} else {
		gasUsed = 21000
	}
	gasUsed = gasUsed * tx.GasPrice
	cost := value + gasUsed

	fromAddr := common.Address(from)
	var accountfrom common.Account
	accountfrom = state.Load(fromAddr)
	
	if accountfrom.Balance < cost {
		return
	}

	accountfrom.Balance = accountfrom.Balance - cost

	state.Store(fromAddr, accountfrom)

	toAddr := common.Address(to)
	
	accountto := state.Load(toAddr)
	
	accountto.Balance= accountto.Balance + value

	state.Store(toAddr, accountto)
}
