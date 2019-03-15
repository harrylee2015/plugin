/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package executor

import (
	"fmt"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/gold5g/ptypes"
	tk "github.com/33cn/plugin/plugin/dapp/token/types"
)

// round status
const (
	decimal = float64(100000000) //1e8
	// ListDESC  desc query
	ListDESC = int32(0)
	// ListASC  asc query
	ListASC         = int32(1)
	Gold5GRoundLast = "round-last"
	// DefaultCount 默认一次取多少条记录
	DefaultCount = int32(20)

	// MaxCount 最多取100条
	MaxCount = int32(100)
)

// GetReceiptLog get receipt log
func (action *Action) GetStartReceiptLog(roundInfo *pt.RoundInfo) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pt.TyLogGold5GStart
	r := &pt.ReceiptGold5G{
		Addr:   action.fromaddr,
		Round:  roundInfo.Round,
		Index:  action.GetIndex(),
		Action: pt.TyLogGold5GStart,
	}
	log.Log = types.Encode(r)

	return log
}
func (action *Action) GetBuyReceiptLog(addrInfo *pt.AddrInfo) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pt.TyLogGold5GBuy
	r := &pt.ReceiptGold5G{
		Addr:     action.fromaddr,
		Round:    addrInfo.Round,
		Index:    action.GetIndex(),
		BuyCount: addrInfo.BuyCount,
		Action:   pt.Gold5GActionBuy,
	}
	log.Log = types.Encode(r)

	return log
}
func (action *Action) GetDrawReceiptLog(roundInfo *pt.RoundInfo) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pt.TyLogGold5GDraw
	r := &pt.ReceiptGold5G{
		Addr:   action.fromaddr,
		Round:  roundInfo.Round,
		Index:  action.GetIndex(),
		Action: pt.TyLogGold5GDraw,
	}
	log.Log = types.Encode(r)

	return log
}

//GetIndex get index
func (action *Action) GetIndex() int64 {
	return action.height*types.MaxTxsPerBlock + int64(action.index)
}

func calcGold5GByRound(round int64) string {
	key := fmt.Sprintf("roundInfo-%010d", round)
	return key
}

func calcGold5GUserAddrs(round int64, addr string) string {
	key := fmt.Sprintf("user-addrs-%010d-%s", round, addr)
	return key
}
func calcGold5GUserTickets(round int64, addr string, index int64) string {
	key := fmt.Sprintf("user-tickets-%010d-%s-%018d", round, addr, index)
	return key
}

//GetKVSet get kv set
func (action *Action) GetKVSet(param interface{}) (kvset []*types.KeyValue, result interface{}) {
	if roundInfo, ok := param.(*pt.RoundInfo); ok {
		value := types.Encode(roundInfo)
		//更新stateDB缓存
		action.db.Set(Key(calcGold5GByRound(roundInfo.Round)), value)
		action.db.Set(Key(Gold5GRoundLast), value)
		kvset = append(kvset, &types.KeyValue{Key: Key(calcGold5GByRound(roundInfo.Round)), Value: value})
		kvset = append(kvset, &types.KeyValue{Key: Key(Gold5GRoundLast), Value: value})
	}
	if ticketInfo, ok := param.(*pt.TicketInfo); ok {
		value := types.Encode(ticketInfo)
		action.db.Set(Key(calcGold5GUserTickets(ticketInfo.Round, ticketInfo.Addr, action.GetIndex())), value)
		kvset = append(kvset, &types.KeyValue{Key: Key(calcGold5GUserTickets(ticketInfo.Round, ticketInfo.Addr, action.GetIndex())), Value: value})
		addrInfo, err := getGold5GAddrInfo(action.db, Key(calcGold5GUserAddrs(ticketInfo.Round, ticketInfo.Addr)))
		if err != nil {
			flog.Warn("Gold5G db getGold5GAddrInfo", "can't get value from db,key:", calcGold5GUserAddrs(ticketInfo.Round, ticketInfo.Addr))
			var addr pt.AddrInfo
			addr.Addr = action.fromaddr
			addr.TicketNum = ticketInfo.TicketNum
			addr.BuyCount = 1
			addr.Round = ticketInfo.Round
			addr.TotalCost = float64(ticketInfo.TicketNum) * ticketInfo.TicketPrice
			value := types.Encode(&addr)
			action.db.Set(Key(calcGold5GUserAddrs(ticketInfo.Round, ticketInfo.Addr)), value)
			kvset = append(kvset, &types.KeyValue{Key: Key(calcGold5GUserAddrs(ticketInfo.Round, ticketInfo.Addr)), Value: value})
			return kvset, &addr
		} else {
			addrInfo.Addr = action.fromaddr
			addrInfo.BuyCount = addrInfo.BuyCount + 1
			addrInfo.Round = ticketInfo.Round
			addrInfo.TicketNum = addrInfo.TicketNum + ticketInfo.TicketNum
			addrInfo.TotalCost = addrInfo.TotalCost + float64(ticketInfo.TicketNum)*ticketInfo.TicketPrice
			value := types.Encode(addrInfo)
			action.db.Set(Key(calcGold5GUserAddrs(ticketInfo.Round, ticketInfo.Addr)), value)
			kvset = append(kvset, &types.KeyValue{Key: Key(calcGold5GUserAddrs(ticketInfo.Round, ticketInfo.Addr)), Value: value})
			return kvset, addrInfo
		}

	}
	if addrInfo, ok := param.(*pt.AddrInfo); ok {
		value := types.Encode(addrInfo)
		action.db.Set(Key(calcGold5GUserAddrs(addrInfo.Round, addrInfo.Addr)), value)
		kvset = append(kvset, &types.KeyValue{Key: Key(calcGold5GUserAddrs(addrInfo.Round, addrInfo.Addr)), Value: value})
		return kvset, addrInfo
	}
	return kvset, nil
}

func (action *Action) updateCount(status int32, addr string) (kvset []*types.KeyValue) {

	return kvset
}

// Key gameId to save key
func Key(id string) (key []byte) {
	key = append(key, []byte("mavl-"+pt.Gold5GX+"-")...)
	key = append(key, []byte(id)...)
	return key
}

// Action action struct
type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	localDB      dbm.Lister
	index        int
}

// NewAction new action
func NewAction(g *gold5g, tx *types.Transaction, index int) (*Action, error) {
	hash := tx.Hash()
	fromaddr := tx.From()
	if pt.GetTokenSymbol() == "" {
		return &Action{g.GetCoinsAccount(), g.GetStateDB(), hash, fromaddr,
			g.GetBlockTime(), g.GetHeight(), dapp.ExecAddress(string(tx.Execer)), g.GetLocalDB(), index}, nil
	}
	tokenAccount, err := account.NewAccountDB(tk.TokenX, pt.GetTokenSymbol(), g.GetStateDB())
	return &Action{tokenAccount, g.GetStateDB(), hash, fromaddr,
		g.GetBlockTime(), g.GetHeight(), dapp.ExecAddress(string(tx.Execer)), g.GetLocalDB(), index}, err

}

func (action *Action) checkExecAccountBalance(fromAddr string, active, frozen int64) bool {
	acc := action.coinsAccount.LoadExecAccount(fromAddr, action.execaddr)
	if acc.GetBalance() >= active && acc.GetFrozen() >= frozen {
		return true
	}
	return false
}

//F3d start game
func (action *Action) Gold5GStart(g *pt.Gold5GStart) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var startRound int64
	//addr check
	if action.fromaddr != pt.Getgold5gManagerAddr() {
		flog.Error("Gold5GStart", "manager addr not match.", "err", pt.ErrGold5GManageAddr.Error())
		return nil, pt.ErrGold5GManageAddr
	}
	lastRound, err := getGold5GRoundInfo(action.db, Key(Gold5GRoundLast))
	if err == nil && lastRound != nil {
		if lastRound.EndTime == 0 {
			flog.Error("Gold5GStart", "start round", startRound)
			return nil, pt.ErrGold5GStartRound
		}
		startRound = lastRound.Round
	}

	account := action.coinsAccount.LoadExecAccount(action.fromaddr, action.execaddr)
	roundInfo := &pt.RoundInfo{
		Round:           startRound + 1,
		BeginTime:       action.blocktime,
		LastTicketPrice: pt.Getgold5gTicketPrice(),
		LastTicketTime:  action.blocktime,
		RemainTime:      pt.Getgold5gTimeLife(),
		//TODO is the floating-point precision here accurate?
		BonusPool:  float64(account.Frozen) / decimal,
		UpdateTime: action.blocktime,
	}

	receiptLog := action.GetStartReceiptLog(roundInfo)
	logs = append(logs, receiptLog)
	kvset, _ := action.GetKVSet(roundInfo)
	kv = append(kv, kvset...)
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//Gold5G start game
func (action *Action) Gold5GBuyTicket(buy *pt.Gold5GBuyTicket) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//addr check
	if action.fromaddr == pt.Getgold5gManagerAddr() {
		flog.Error("Gold5GBuyTicket", "manager can't buy ticket.", "err", pt.ErrGold5GManageBuyTicket.Error())
		return nil, pt.ErrGold5GManageBuyTicket
	}

	lastRound, err := getGold5GRoundInfo(action.db, Key(Gold5GRoundLast))
	if err != nil || lastRound == nil {
		flog.Error("Gold5GBuyTicket", "last round", lastRound.Round, "err", fmt.Errorf("not found the last round info!"))
		return nil, fmt.Errorf("not found the last round info!")
	}
	//round game status check
	if lastRound.EndTime != 0 {
		flog.Error("Gold5GBuyTicket", "last round", lastRound.Round, "err", pt.ErrGold5GBuyTicket.Error())
		return nil, pt.ErrGold5GBuyTicket
	}
	// remainTime check
	if lastRound.UpdateTime+lastRound.RemainTime < action.blocktime {
		flog.Error("Gold5GBuyTicket", "time out", "err", pt.ErrGold5GBuyTicketTimeOut.Error())
		return nil, pt.ErrGold5GBuyTicketTimeOut
	}
	// balance check
	if !action.checkExecAccountBalance(action.fromaddr, buy.GetTicketNum()*int64(lastRound.GetLastTicketPrice()*decimal), 0) {
		flog.Error("Gold5GBuyTicket", "checkExecAccountBalance", action.fromaddr, "execaddr", action.execaddr, "err", types.ErrNoBalance.Error())
		return nil, types.ErrNoBalance
	}
	// supportFund
	receipt, err := action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gSupportFundAddr(), action.execaddr, buy.GetTicketNum()*int64(pt.Getgold5gBonusSupportFund()*decimal))
	if err != nil {
		flog.Error("Gold5GBuyTicket.ExecTransfer supportFund", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", buy.GetTicketNum()*int64(pt.Getgold5gBonusSupportFund()*decimal))
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	//partner
	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gPartnerAddr(), action.execaddr, buy.GetTicketNum()*int64(pt.Getgold5gBonusPartner()*decimal))
	if err != nil {
		flog.Error("Gold5GBuyTicket.ExecTransfer partner", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", buy.GetTicketNum()*int64(pt.Getgold5gBonusPartner()*decimal))
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	//developer
	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gDeveloperAddr(), action.execaddr, buy.GetTicketNum()*int64(pt.Getgold5gBonusDeveloper()*decimal))
	if err != nil {
		flog.Error("Gold5GBuyTicket.ExecTransfer developer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", buy.GetTicketNum()*int64(pt.Getgold5gBonusDeveloper()*decimal))
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	//Promotion
	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gPromotionAddr(), action.execaddr, buy.GetTicketNum()*int64(pt.Getgold5gBonusPromotion()*decimal))
	if err != nil {
		flog.Error("Gold5GBuyTicket.ExecTransfer promotion", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", buy.GetTicketNum()*int64(pt.Getgold5gBonusDeveloper()*decimal))
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	//bonusPool
	remain := buy.GetTicketNum() * int64((pt.Getgold5gTicketPrice()-pt.Getgold5gBonusSupportFund()-pt.Getgold5gBonusPartner()-pt.Getgold5gBonusDeveloper()-pt.Getgold5gBonusPromotion())*decimal)
	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gManagerAddr(), action.execaddr, remain)
	if err != nil {
		flog.Error("Gold5GBuyTicket.ExecTransfer promotion", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", remain)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.coinsAccount.ExecFrozen(pt.Getgold5gManagerAddr(), action.execaddr, remain)
	if err != nil {
		flog.Error("Gold5GBuyTicket.Frozen", "addr", pt.Getgold5gManagerAddr(), "execaddr", action.execaddr, "amount", remain)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// first update ticketInfo
	ticketInfo := &pt.TicketInfo{}
	ticketInfo.TicketNum = buy.TicketNum
	ticketInfo.Round = lastRound.Round
	ticketInfo.Addr = action.fromaddr
	ticketInfo.TicketPrice = lastRound.LastTicketPrice
	ticketInfo.BuyTicketTime = action.blocktime
	ticketInfo.BuyTicketTxHash = common.ToHex(action.txhash)
	ticketInfo.Index = action.GetIndex()
	kvset, v := action.GetKVSet(ticketInfo)
	kv = append(kv, kvset...)
	if addrInfo, ok := v.(*pt.AddrInfo); ok {
		// first buy
		if addrInfo.BuyCount == 1 {
			lastRound.UserCount = lastRound.UserCount + 1
		}
		receiptLog := action.GetBuyReceiptLog(addrInfo)
		logs = append(logs, receiptLog)
	}

	lastRound, outList := action.calculateOverflowBonus(lastRound, buy)
	//outTicket
	for _, outTicket := range outList {
		bonus := 105 * outTicket.TicketNum * int64(decimal)
		if !action.checkExecAccountBalance(pt.Getgold5gManagerAddr(), 0, bonus) {
			flog.Error("Gold5GLuckyDraw", "checkExecAccountBalance", action.fromaddr, "execaddr", action.execaddr, "err", types.ErrNoBalance.Error())
			return nil, types.ErrNoBalance
		}
		receipt, err := action.coinsAccount.ExecActive(pt.Getgold5gManagerAddr(), action.execaddr, bonus)
		if err != nil {
			flog.Error("Gold5GLuckyDraw.ExecActive", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", bonus)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
		//pay bonus for firstBonus
		receipt, err = action.coinsAccount.ExecTransfer(pt.Getgold5gManagerAddr(), outTicket.Addr, action.execaddr, bonus)
		if err != nil {
			flog.Error("Gold5GLuckyDraw.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", bonus)
			return nil, err
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)

		//update everyone sortPoolBonus
		info, err := getGold5GAddrInfo(action.db, Key(calcGold5GUserAddrs(lastRound.Round, outTicket.Addr)))
		if err != nil {
			continue
		}
		info.SortPoolBonus += float64(bonus) / decimal
		kvset, _ := action.GetKVSet(info)
		kv = append(kv, kvset...)
	}

	lastRound.RemainTime = lastRound.RemainTime + lastRound.UpdateTime - action.blocktime
	lastRound.TicketCount += buy.TicketNum
	lastRound.LastTicketTime = action.blocktime
	lastRound.UpdateTime = action.blocktime

	addTime := pt.Getgold5gTimeTicket() * buy.TicketNum
	if lastRound.RemainTime+addTime >= pt.Getgold5gTimeLife() {
		lastRound.RemainTime = pt.Getgold5gTimeLife()
	} else {
		lastRound.RemainTime = lastRound.RemainTime + addTime
	}
	//Todo  add addr and nums
	kvset, _ = action.GetKVSet(lastRound)
	kv = append(kv, kvset...)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (action *Action) calculateOverflowBonus(lastRound *pt.RoundInfo, buy *pt.Gold5GBuyTicket) (*pt.RoundInfo, []*pt.OutTicketInfo) {
	var outTicketList []*pt.OutTicketInfo
	//var sortPoolInfo []*pt.InTicketInfo

	if len(lastRound.SortPoolInfo) == 0 {
		in := &pt.InTicketInfo{
			TicketNum: buy.TicketNum,
			Addr:      action.fromaddr,
		}
		lastRound.SortPoolInfo = append(lastRound.SortPoolInfo, in)
		lastRound.BonusSortPool += float64(80 * buy.TicketNum)
		return lastRound, outTicketList
	}
	outNum := int64(lastRound.BonusSortPool+float64(80*buy.TicketNum)-float64(lastRound.SortPoolInfo[0].TicketNum*80)) / 25
	addBonus := float64(80 * buy.TicketNum)

	if outNum <= 0 {
		in := &pt.InTicketInfo{
			TicketNum: buy.TicketNum,
			Addr:      action.fromaddr,
		}
		lastRound.SortPoolInfo = append(lastRound.SortPoolInfo, in)
		lastRound.BonusSortPool += addBonus
		return lastRound, outTicketList
	}
HERE:
	if 1 <= outNum && outNum < lastRound.SortPoolInfo[0].TicketNum {
		out := &pt.OutTicketInfo{
			TicketNum: outNum,
			Addr:      lastRound.SortPoolInfo[0].Addr,
		}
		outTicketList = append(outTicketList, out)
		lastRound.SortPoolInfo[0].TicketNum -= outNum
		lastRound.BonusSortPool = lastRound.BonusSortPool + addBonus - float64(105*outNum)
	}
	if outNum >= lastRound.SortPoolInfo[0].TicketNum {
		out := &pt.OutTicketInfo{
			TicketNum: lastRound.SortPoolInfo[0].TicketNum,
			Addr:      lastRound.SortPoolInfo[0].Addr,
		}
		outTicketList = append(outTicketList, out)
		lastRound.SortPoolInfo = lastRound.SortPoolInfo[1:]
		lastRound.BonusSortPool = lastRound.BonusSortPool + addBonus - float64(105*out.TicketNum)
		outNum = int64(lastRound.BonusSortPool-float64(lastRound.SortPoolInfo[0].TicketNum*80)) / 25
		addBonus = 0
		goto HERE
	}
	//add this buy ticket
	in := &pt.InTicketInfo{
		TicketNum: buy.TicketNum,
		Addr:      action.fromaddr,
	}
	lastRound.SortPoolInfo = append(lastRound.SortPoolInfo, in)
	return lastRound, outTicketList
}

//Gold5G lucky draws
func (action *Action) Gold5GLuckyDraw(buy *pt.Gold5GLuckyDraw) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//addr check
	if action.fromaddr != pt.Getgold5gManagerAddr() {
		flog.Error("Gold5GLuckyDraw", "manager addr not match.", "err", pt.ErrGold5GManageAddr.Error())
		return nil, pt.ErrGold5GManageAddr
	}

	lastRound, err := getGold5GRoundInfo(action.db, Key(Gold5GRoundLast))

	if err != nil || lastRound == nil {
		flog.Error("Gold5GLuckyDraw", "last round", lastRound.Round, "err", fmt.Errorf("not found the last round info!"))
		return nil, fmt.Errorf("not found the last round info!")
	}
	// remainTime check
	if lastRound.UpdateTime+lastRound.RemainTime > action.blocktime {
		flog.Error("Gold5GLuckyDraw", "remain time not be zerio", "err", pt.ErrGold5GDrawRemainTime.Error())
		return nil, pt.ErrGold5GDrawRemainTime
	}
	//round game status check
	if lastRound.EndTime != 0 {
		flog.Error("Gold5GLuckyDraw", "last round", lastRound.Round, "err", pt.ErrGold5GDrawRepeat.Error())
		return nil, pt.ErrGold5GDrawRepeat
	}

	//round info check,when no one buy keys,just finish the game
	if lastRound.UserCount == 0 {
		lastRound.RemainTime = 0
		lastRound.EndTime = action.blocktime
		lastRound.UpdateTime = action.blocktime
		receiptLog := action.GetDrawReceiptLog(lastRound)
		logs = append(logs, receiptLog)
		kvset, _ := action.GetKVSet(lastRound)
		kv = append(kv, kvset...)
		return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
	}

	firstBonus := lastRound.BonusPool * pt.Getgold5gBonusFirstPrize() * decimal
	secondaryBonus := lastRound.BonusPool * pt.Getgold5gBonusSecondaryPrize() * decimal
	sortPoolBonus := lastRound.BonusSortPool * decimal

	// balance check
	if !action.checkExecAccountBalance(action.fromaddr, 0, int64(firstBonus+secondaryBonus+sortPoolBonus)) {
		flog.Error("Gold5GLuckyDraw", "checkExecAccountBalance", action.fromaddr, "execaddr", action.execaddr, "err", types.ErrNoBalance.Error())
		return nil, types.ErrNoBalance
	}
	receipt, err := action.coinsAccount.ExecActive(action.fromaddr, action.execaddr, int64(firstBonus+secondaryBonus))
	if err != nil {
		flog.Error("Gold5GLuckyDraw.ExecActive", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", (firstBonus+secondaryBonus)/decimal)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)
	//pay bonus for firstBonus
	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gFirstPrizeAddr(), action.execaddr, int64(firstBonus))
	if err != nil {
		flog.Error("Gold5GLuckyDraw.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", firstBonus/decimal)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	//pay bonus for secondaryBonus
	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gSecondaryPrizeAddr(), action.execaddr, int64(secondaryBonus))
	if err != nil {
		flog.Error("Gold5GLuckyDraw.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", secondaryBonus/decimal)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.coinsAccount.ExecTransfer(action.fromaddr, pt.Getgold5gUnassignedAddr(), action.execaddr, int64(sortPoolBonus))
	if err != nil {
		flog.Error("Gold5GLuckyDraw.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", sortPoolBonus/decimal)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	lastRound.RemainTime = 0
	lastRound.EndTime = action.blocktime
	lastRound.UpdateTime = action.blocktime
	receiptLog := action.GetDrawReceiptLog(lastRound)
	logs = append(logs, receiptLog)
	kvset, _ := action.GetKVSet(lastRound)
	kv = append(kv, kvset...)
	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}
func getGold5GRoundInfo(db dbm.KV, key []byte) (*pt.RoundInfo, error) {
	value, err := db.Get(key)
	if err != nil {
		flog.Error("Gold5G db getGold5GRoundInfo", "can't get value from db,key:", key, "err", err.Error())
		return nil, err
	}

	var roundInfo pt.RoundInfo
	err = types.Decode(value, &roundInfo)
	if err != nil {
		return nil, err
	}
	return &roundInfo, nil

}

func getGold5GAddrInfo(db dbm.KV, key []byte) (*pt.AddrInfo, error) {
	value, err := db.Get(key)
	if err != nil {
		flog.Error("Gold5G db getGold5GAddrInfo", "can't get value from db,key:", key, "err", err.Error())
		return nil, err
	}

	var info pt.AddrInfo
	err = types.Decode(value, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func getGold5GBuyRecord(db dbm.KV, key []byte) (*pt.TicketInfo, error) {
	value, err := db.Get(key)
	if err != nil {
		flog.Error("Gold5G db getGold5GBuyRecord", "can't get value from db,key:", key, "err", err.Error())
		return nil, err
	}

	var info pt.TicketInfo
	err = types.Decode(value, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func queryList(db dbm.Lister, stateDB dbm.KV, param interface{}) (types.Message, error) {
	direction := ListDESC
	count := DefaultCount
	if query, ok := param.(*pt.QueryAddrInfo); ok {
		direction = query.GetDirection()
		if 0 < query.GetCount() && query.GetCount() <= MaxCount {
			count = query.GetCount()
		}
		if query.Round == 0 {
			return nil, fmt.Errorf("round can't be zero!")
		}

		var values [][]byte
		var err error
		if query.GetAddr() == "" { //第一次查询
			values, err = db.List(calcGold5GAddrPrefix(query.Round), nil, count, direction)
		} else {
			values, err = db.List(calcGold5GAddrPrefix(query.Round), calcGold5GAddrRound(query.Round, query.Addr), count, direction)
		}
		if err != nil {
			return nil, err
		}
		var addrList []*pt.AddrInfo
		for _, value := range values {
			var addrInfo pt.AddrInfo
			err := types.Decode(value, &addrInfo)
			if err != nil {
				continue
			}
			addrList = append(addrList, &addrInfo)
		}
		return &pt.ReplyAddrInfoList{AddrInfoList: addrList}, nil
	}
	//query last round info
	if _, ok := param.(*pt.QueryGold5GLastRound); ok {
		lastRound, err := getGold5GRoundInfo(stateDB, Key(Gold5GRoundLast))
		if err != nil {
			flog.Error("Gold5G db queryList", "can't get lastRound:err", err.Error())
			return nil, err
		}
		return lastRound, nil
	}
	//query round info by round
	if query, ok := param.(*pt.QueryGold5GByRound); ok {
		if query.Round == 0 {
			return nil, fmt.Errorf("round can't be zero!")
		}
		round, err := getGold5GRoundInfo(stateDB, Key(calcGold5GByRound(query.Round)))
		if err != nil {
			flog.Error("Gold5G db queryList", "can't get lastRound:err", err.Error())
			return nil, err
		}
		return round, nil
	}
	//query addr info
	if query, ok := param.(*pt.QueryTicketInfoByRoundAndAddr); ok {
		if query.Round == 0 || query.Addr == "" {
			return nil, fmt.Errorf("round can't be zero,addr can't be empty!")
		}
		addrInfo, err := getGold5GAddrInfo(stateDB, Key(calcGold5GUserAddrs(query.Round, query.Addr)))
		if err != nil {
			flog.Error("F3D db queryList", "can't get addr Info,err", err.Error())
			return nil, err
		}
		return addrInfo, nil
	}
	//query buy record
	if query, ok := param.(*pt.QueryBuyRecordByRoundAndAddr); ok {
		if query.Round == 0 || query.Addr == "" {
			return nil, fmt.Errorf("round can't be zero,addr can't be empty!")
		}

		var values [][]byte
		var err error
		if query.Index == 0 { //第一次查询
			values, err = db.List(calcGold5GBuyPrefix(query.Round, query.Addr), nil, count, direction)
		} else {
			values, err = db.List(calcGold5GBuyPrefix(query.Round, query.Addr), calcGold5GBuyRound(query.Round, query.Addr, query.Index), count, direction)
		}
		if err != nil {
			return nil, err
		}
		var recordList []*pt.TicketInfo
		for _, value := range values {
			var r pt.Gold5GBuyRecord
			err := types.Decode(value, &r)
			if err != nil {
				continue
			}
			record, err := getGold5GBuyRecord(stateDB, Key(calcGold5GUserTickets(r.Round, r.Addr, r.Index)))
			if err != nil {
				flog.Error("Gold5G db queryList", "can't get buy record,err", err.Error())
				continue
			}
			recordList = append(recordList, record)
		}
		return &pt.ReplyBuyRecord{RecordList: recordList}, nil
	}
	return nil, fmt.Errorf("this query can't be supported!")
}
