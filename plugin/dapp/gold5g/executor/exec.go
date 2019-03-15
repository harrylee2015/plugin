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

func (g *gold5g) Exec_Start(payload *pt.Gold5GStart, tx *types.Transaction, index int) (*types.Receipt, error) {
	action, err := NewAction(g, tx, index)
	if err != nil {
		return nil, err
	}
	return action.Gold5GStart(payload)
}

func (g *gold5g) Exec_Draw(payload *pt.Gold5GLuckyDraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action, err := NewAction(g, tx, index)
	if err != nil {
		return nil, err
	}
	return action.Gold5GLuckyDraw(payload)
}

func (g *gold5g) Exec_Buy(payload *pt.Gold5GBuyTicket, tx *types.Transaction, index int) (*types.Receipt, error) {
	action, err := NewAction(g, tx, index)
	if err != nil {
		return nil, err
	}
	return action.Gold5GBuyTicket(payload)
}
