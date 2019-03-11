package rpc

import (
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
	context "golang.org/x/net/context"
)

func (c *channelClient) Bet(ctx context.Context, head *ft.FluckyBet) (*types.UnsignTx, error) {
	val := &ft.FluckyAction{
		Ty:    ft.FluckyActionBet,
		Value: &ft.FluckyAction_Bet{Bet: head},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	data, err := types.FormatTxEncode(string(ft.ExecerFlucky), tx)
	if err != nil {
		return nil, err
	}
	return &types.UnsignTx{Data: data}, nil
}
