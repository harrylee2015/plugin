/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package executor

import (
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/gold5g/ptypes"
)

// roll back local db data
func (g *gold5g) execDelLocal(receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for _, log := range receiptData.Logs {
		switch log.Ty {
		case pt.TyLogGold5GStart, pt.TyLogGold5GBuy, pt.TyLogGold5GDraw:
			receipt := &pt.ReceiptGold5G{}
			if err := types.Decode(log.Log, receipt); err != nil {
				return nil, err
			}
			kv := g.rollbackLocalDB(receipt)
			dbSet.KV = append(dbSet.KV, kv...)
		}
	}
	return dbSet, nil
}
func (g *gold5g) ExecDelLocal_Start(payload *pt.Gold5GStart, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execDelLocal(receiptData)
}

func (g *gold5g) ExecDelLocal_Draw(payload *pt.Gold5GLuckyDraw, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execDelLocal(receiptData)
}

func (g *gold5g) ExecDelLocal_Buy(payload *pt.Gold5GBuyTicket, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return g.execDelLocal(receiptData)
}
