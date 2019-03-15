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

func (g *gold5g) Query_QueryLastRoundInfo(in *pt.QueryGold5GLastRound) (types.Message, error) {
	return queryList(g.GetLocalDB(), g.GetStateDB(), in)
}

func (g *gold5g) Query_QueryRoundInfoByRound(in *pt.QueryGold5GByRound) (types.Message, error) {
	return queryList(g.GetLocalDB(), g.GetStateDB(), in)
}

func (g *gold5g) Query_QueryRoundsInfoByRounds(in *pt.QueryGold5GListByRound) (types.Message, error) {
	return queryList(g.GetLocalDB(), g.GetStateDB(), in)
}

func (g *gold5g) Query_QueryTicketInfoByRoundAndAddr(in *pt.QueryTicketInfoByRoundAndAddr) (types.Message, error) {
	return queryList(g.GetLocalDB(), g.GetStateDB(), in)
}

func (g *gold5g) Query_QueryBuyRecordByRoundAndAddr(in *pt.QueryBuyRecordByRoundAndAddr) (types.Message, error) {
	return queryList(g.GetLocalDB(), g.GetStateDB(), in)
}
