package vm

import (
	"cxchain-2023131015/common"
	"cxchain-2023131015/txpool"
	"fmt"
)

// 包含两个执行方法：Execute和Execute1
type IMachine interface {
	Execute(tx common.Transaction) error
	Execute1(tx common.Transaction) *common.Receipt
}

// StateMachine 结构体实现了IMachine接口
// TxPool字段指向交易池实例
type StateMachine struct {
	TxPool *txpool.TxPool
}

// NewStateMachine 创建一个新的状态机实例
// 参数txPool是交易池实例
// 返回StateMachine指针
func NewStateMachine(txPool *txpool.TxPool) *StateMachine {
	return &StateMachine{TxPool: txPool}
}

// Execute 执行交易的核心方法
// 处理交易转账逻辑，包括：
// 1. 检查gas价格是否足够
// 2. 计算交易成本
// 3. 更新发送方和接收方账户余额
// 4. 更新nonce值
func (m *StateMachine) Execute(tx common.Transaction) error {
	// 获取交易基本信息
	from := tx.From()
	to := tx.To
	value := tx.Value
	gasUsed := tx.GasPrice
	fmt.Println(from)
	fmt.Println(gasUsed)

	// 检查gas价格是否足够
	if tx.GasPrice < 21000 {
		return fmt.Errorf("gas price is too low")
	} else {
		gasUsed = 1
	}

	// 计算交易总成本
	gasUsed = gasUsed * tx.GasPrice
	cost := value + gasUsed

	// 处理发送方账户
	fromAddr := common.Address(from)
	var accountfrom common.Account
	accountfrom = m.TxPool.StatDB.Load(fromAddr)
	
	fmt.Println(accountfrom)
	if accountfrom.Balance < cost {
		return fmt.Errorf("balance is not enough")
	}

	// 更新发送方余额和nonce
	accountfrom.Balance = accountfrom.Balance - cost
	accountfrom.Nonce = accountfrom.Nonce + 1
	m.TxPool.StatDB.Store(fromAddr, accountfrom)

	// 处理接收方账户
	toAddr := common.Address(to)
	accountto := m.TxPool.StatDB.Load(toAddr)
	fmt.Println(accountto)
	if common.EqualAccounts(accountto, common.Account{}) {
		// 如果接收方账户不存在，则创建新账户
		accountto = common.Account{
			Balance: value,
			Nonce:   0,
		}
	}

	// 更新接收方余额
	accountto.Balance = accountto.Balance + value
	m.TxPool.StatDB.Store(toAddr, accountto)
	return nil
}

// Execute1 执行交易并返回收据
// 在Execute基础上增加了收据生成
// 返回交易执行结果收据
func (m *StateMachine) Execute1(tx common.Transaction) *common.Receipt {
	var receipt common.Receipt
	err := m.Execute(tx)
	if err != nil {
		return nil
	}
	receipt.Status = 1
	receipt.TransactionHash = tx.Hash()
	return &receipt
}

