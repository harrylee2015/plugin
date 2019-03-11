package executor

import (
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
)

func (f *Flucky) Query_QueryBetTimes(in *ft.QueryBetTimes) (types.Message, error) {
	return QueryBetTimes(f.GetStateDB(), in)
}

func (f *Flucky) Query_QueryBetInfoBatch(in *ft.QueryBetInfoBatch) (types.Message, error) {
	return QueryBetListByPage(f.GetLocalDB(), f.GetStateDB(), in)
}

func (f *Flucky) Query_QueryBetInfo(in *ft.QueryBetInfo) (types.Message, error) {
	return QueryBetInfo(f.GetLocalDB(), in)
}

func (f *Flucky) Query_QueryBonusInfo(in *ft.QueryBonusInfo) (types.Message, error) {
	return QueryBonusInfo(f.GetStateDB())
}
