package rpc

import (
	"context"
	"encoding/hex"

	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
)

// BlackwhiteCreateTx 创建游戏RPC接口
func (c *Jrpc) FluckyBetTx(parm *ft.FluckyBet, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}
	head := &ft.FluckyBet{
		Amount: parm.Amount,
	}
	reply, err := c.cli.Bet(context.Background(), head)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}
