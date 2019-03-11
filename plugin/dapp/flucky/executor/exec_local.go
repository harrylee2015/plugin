package executor

import (
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
)

// save receiptData to local db
func (f *Flucky) execLocal(receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	table := ft.NewTable(f.GetLocalDB())
	for _, log := range receiptData.Logs {
		switch log.Ty {
		case ft.TyLogFluckyBet:
			var receipt ft.ReceiptFlucky
			if err := types.Decode(log.Log, &receipt); err != nil {
				return nil, err
			}
			if err := table.Replace(&receipt); err != nil {
				return nil, err
			}
			kvs, err := table.Save()
			if err != nil {
				return nil, err
			}
			dbSet.KV = append(dbSet.KV, kvs...)
		}
	}
	return dbSet, nil
}
func (f *Flucky) ExecLocal_Bet(payload *ft.FluckyBet, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return f.execLocal(receiptData)
}
