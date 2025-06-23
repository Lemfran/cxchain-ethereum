package txpool

import "cxchain-2023131015/common"

type TxPool struct {
	pending map[common.Address][]boxes
	queue map[common.Address]map[uint64]*common.Transaction
}

