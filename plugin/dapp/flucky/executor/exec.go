package executor

import (
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
)

// Exec_Create 创建游戏
func (c *Flucky) Exec_Bet(payload *ft.FluckyBet, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(c, tx, index)
	return action.FluckyBet(payload)
}
