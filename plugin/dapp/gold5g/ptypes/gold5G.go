/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package types

import (
	"reflect"

	"github.com/33cn/chain33/types"
)

// action for executor
const (
	Gold5GActionStart = iota + 1
	Gold5GActionDraw
	Gold5GActionBuy
)

const (
	TyLogGold5GUnknown = iota + 110
	TyLogGold5GStart
	TyLogGold5GDraw
	TyLogGold5GBuy
)

// query func name
const (
	FuncNameQueryLastRoundInfo            = "QueryLastRoundInfo"
	FuncNameQueryRoundInfoByRound         = "QueryRoundInfoByRound"
	FuncNameQueryRoundsInfoByRounds       = "QueryRoundsInfoByRounds"
	FuncNameQueryTicketInfoByRoundAndAddr = "QueryTicketInfoByRoundAndAddr"
	FuncNameQueryBuyRecordByRoundAndAddr  = "QueryBuyRecordByRoundAndAddr"
)

var (
	logMap = map[string]int32{
		"Start": Gold5GActionStart,
		"Draw":  Gold5GActionDraw,
		"Buy":   Gold5GActionBuy,
	}

	typeMap      = map[int64]*types.LogInfo{}
	Gold5GX      = "gold5g"
	ExecerGold5G = []byte(Gold5GX)
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(Gold5GX))
	types.RegistorExecutor(Gold5GX, NewType())
	types.RegisterDappFork(Gold5GX, "Enable", 0)
}

type Gold5GType struct {
	types.ExecTypeBase
}

func NewType() *Gold5GType {
	c := &Gold5GType{}
	c.SetChild(c)
	return c
}

func (t *Gold5GType) GetPayload() types.Message {
	return &Gold5GAction{}
}

func (t *Gold5GType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Start": Gold5GActionStart,
		"Draw":  Gold5GActionDraw,
		"Buy":   Gold5GActionBuy,
	}
}

func (t *Gold5GType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogGold5GStart: {Ty: reflect.TypeOf(ReceiptGold5G{}), Name: "LogStartGold5G"},
		TyLogGold5GDraw:  {Ty: reflect.TypeOf(ReceiptGold5G{}), Name: "LogDrawGold5G"},
		TyLogGold5GBuy:   {Ty: reflect.TypeOf(ReceiptGold5G{}), Name: "LogBuyGold5G"},
	}
}
