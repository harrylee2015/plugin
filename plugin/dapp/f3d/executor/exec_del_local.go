package executor

import (
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
)

// roll back local db data
func (f *f3d) execDelLocal(receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for _, log := range receiptData.Logs {
		switch log.Ty {
		case pt.TyLogf3dStart, pt.TyLogf3dBuy, pt.TyLogf3dDraw:
			receipt := &pt.ReceiptF3D{}
			if err := types.Decode(log.Log, receipt); err != nil {
				return nil, err
			}
			kv := f.rollbackLocalDB(receipt)
			dbSet.KV = append(dbSet.KV, kv...)
		}
	}
	return dbSet, nil
}
func (f *f3d) ExecDelLocal_Start(payload *pt.F3DStart, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{}, nil
}

func (f *f3d) ExecDelLocal_Draw(payload *pt.F3DLuckyDraw, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{}, nil
}

func (f *f3d) ExecDelLocal_Buy(payload *pt.F3DBuyKey, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{}, nil
}
