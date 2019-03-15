/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/gold5g/ptypes"
)

var (
	flog = log.New("module", "execs.gold5g")
)

var driverName = pt.Gold5GX

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&gold5g{}))
}

func Init(name string, sub []byte) {
	var config pt.Config
	if sub != nil {
		types.MustDecode(sub, &config)
	}
	pt.SetConfig(&config)

	drivers.Register(GetName(), newgold5g, types.GetDappFork(driverName, "Enable"))
}

type gold5g struct {
	drivers.DriverBase
}

func newgold5g() drivers.Driver {
	t := &gold5g{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

func GetName() string {
	return newgold5g().GetName()
}

func (g *gold5g) GetDriverName() string {
	return driverName
}
func (g *gold5g) ExecutorOrder() int64 {
	return drivers.ExecLocalSameTime
}
func (g *gold5g) updateLocalDB(r *pt.ReceiptGold5G) (kvs []*types.KeyValue) {
	switch r.Action {
	case pt.Gold5GActionStart:
		start := &pt.Gold5GStartRound{
			Round: r.Round,
		}
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GStartRound(r.Round), Value: types.Encode(start)})
	case pt.Gold5GActionBuy:
		addrInfo := &pt.AddrInfo{
			Addr:     r.Addr,
			Round:    r.Round,
			BuyCount: r.BuyCount,
		}
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GAddrRound(r.Round, r.Addr), Value: types.Encode(addrInfo)})
		buyRecord := &pt.Gold5GBuyRecord{
			Round: r.Round,
			Addr:  r.Addr,
			Index: r.Index,
		}
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GBuyRound(r.Round, r.Addr, r.Index), Value: types.Encode(buyRecord)})
	case pt.Gold5GActionDraw:
		draw := &pt.Gold5GDrawRound{
			Round: r.Round,
		}
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GDrawRound(r.Round), Value: types.Encode(draw)})
	}
	return kvs
}

func (g *gold5g) rollbackLocalDB(r *pt.ReceiptGold5G) (kvs []*types.KeyValue) {
	switch r.Action {
	case pt.Gold5GActionStart:
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GStartRound(r.Round), Value: nil})
	case pt.Gold5GActionBuy:
		if r.BuyCount <= 1 {
			kvs = append(kvs, &types.KeyValue{Key: calcGold5GAddrRound(r.Round, r.Addr), Value: nil})
		} else {
			addrInfo := &pt.AddrInfo{
				Addr:     r.Addr,
				Round:    r.Round,
				BuyCount: r.BuyCount - 1,
			}
			kvs = append(kvs, &types.KeyValue{Key: calcGold5GAddrRound(r.Round, r.Addr), Value: types.Encode(addrInfo)})
		}
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GBuyRound(r.Round, r.Addr, r.Index), Value: nil})
	case pt.Gold5GActionDraw:
		kvs = append(kvs, &types.KeyValue{Key: calcGold5GDrawRound(r.Round), Value: nil})
	}
	return kvs
}

// GetPayloadValue get payload value
func (g *gold5g) GetPayloadValue() types.Message {
	return &pt.Gold5GAction{}
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (g *gold5g) CheckReceiptExecOk() bool {
	return true
}
