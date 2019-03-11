package executor

import (
	"fmt"
	"github.com/33cn/chain33/common/db/table"
	"strconv"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/client"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
	"time"
)

func calcFluckyUserTimesKey(addr string) string {
	key := fmt.Sprintf("mavl-flucky-user-times:%s", addr)
	return key
}

func calcFluckyBonusInfoKey() string {
	key := fmt.Sprintf("mavl-flucky-bonul-info")
	return key
}

type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	api          client.QueueProtocolAPI
	localDB      dbm.Lister
	index        int
}

// NewAction 创建action
func NewAction(pb *Flucky, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()

	return &Action{pb.GetCoinsAccount(), pb.GetStateDB(), hash, fromaddr,
		pb.GetBlockTime(), pb.GetHeight(), dapp.ExecAddress(string(tx.Execer)), pb.GetAPI(), pb.GetLocalDB(), index}
}

func (action *Action) FluckyBet(bet *ft.FluckyBet) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	// balance check
	if !action.CheckAccountBalance(action.fromaddr, int64(bet.GetAmount()), 0) {
		flog.Error("FluckyBet", "checkExecAccountBalance", action.fromaddr, "execaddr", action.execaddr, "err", types.ErrNoBalance.Error())
		return nil, types.ErrNoBalance
	}

	receipt, err := action.coinsAccount.ExecTransfer(action.fromaddr, action.execaddr, action.execaddr, int64(bet.GetAmount()*ft.Decimal))
	if err != nil {
		flog.Error("FluckyBet", "ExecTransfer", action.fromaddr, "execaddr", action.execaddr, "err", err)
		return nil, err
	}

	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	// 获取当前奖池信息
	var bonus ft.BonusInfo
	bonusInfo, err := getBonulPoolInfo(action.db, []byte(calcFluckyBonusInfoKey()))
	if err == types.ErrNotFound {
		bonus.UserCount = 0
		bonus.BonusPool += float32(bet.GetAmount())
	} else if err != nil {
		flog.Error("FluckyBet", fmt.Errorf("can't find the bonus info!"))
		return nil, fmt.Errorf("can't find the bonus info!")
	} else {
		bonus.UserCount = bonusInfo.GetUserCount()
		bonus.BonusPool = bonusInfo.BonusPool + float32(bet.GetAmount())
	}

	var index ft.BetReq
	idx, err := getIdxInfo(action.db, []byte(calcFluckyUserTimesKey(action.fromaddr)))
	if err == types.ErrNotFound {
		// 第一次购买，用户数量+1
		bonus.UserCount += 1
		index.Index = 1
	} else if err != nil {
		return nil, fmt.Errorf("get index info error.")
	} else {
		index.Index = idx.GetIndex() + 1
	}

	key := []byte(calcFluckyUserTimesKey(action.fromaddr))
	value := types.Encode(&index)
	//更新stateDB缓存
	action.db.Set(key, value)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})

	// 根据购买的keys生成随机数
	var reward ft.BetInfo
	var randInfo []int64

	//betRule := ft.GetBetTimesRule(bet.GetAmount())
	//randInfo = action.getRandNum(betRule.GetTimes())
	randInfo = action.getRandNumOnce()

	reward.RandNum = randInfo
	UpdateBonusInfo(&reward, float32(bet.GetAmount()), bonus.BonusPool)
	bonus.BonusPool -= reward.GetBonus()

	receipt, err = action.coinsAccount.ExecTransfer(action.execaddr, action.fromaddr, action.execaddr, int64(reward.Bonus*ft.Decimal))
	if err != nil {
		flog.Error("FluckyBet", "ExecTransfer", action.fromaddr, "execaddr", action.execaddr, "err", err)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receiptLog := action.GetBetReceiptLog(index.GetIndex(), uint64(time.Now().Unix()), bet.GetAmount(), &reward)
	logs = append(logs, receiptLog)

	// 每轮结束对bonus进行检查
	if bonus.BonusPool >= ft.GetMaxBonus() {
		flog.Debug("Begin to transfer bunds from bonus to platform...")
		receipt, err = action.coinsAccount.ExecTransfer(action.execaddr, ft.GetPlatFormAddr(), action.execaddr, int64(ft.GetBonusToPlatform()*ft.Decimal))
		if err != nil {
			flog.Error("FluckyBet ExecTransfer failed", "from", action.execaddr, "to", ft.GetPlatFormAddr(), "execaddr", action.execaddr, "err", err)
		} else {
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)
			bonus.BonusPool -= ft.GetBonusToPlatform()
		}
	}

	if bonus.BonusPool <= ft.GetMinBonus() {
		flog.Debug("Begin to transfer bunds from platform to bonus...")
		receipt, err = action.coinsAccount.ExecTransfer(ft.GetPlatFormAddr(), action.execaddr, action.execaddr, int64(ft.GetPlatformToBonus()*ft.Decimal))
		if err != nil {
			flog.Error("FluckyBet ExecTransfer failed", "from", ft.GetPlatFormAddr(), "to", action.execaddr, "execaddr", action.execaddr, "err", err)
		} else {
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)
			bonus.BonusPool += ft.GetPlatformToBonus()
		}
	}

	kvset, _ := action.GetKVSet(&bonus)
	kv = append(kv, kvset...)
	kvset, _ = action.GetKVSet(&index)
	kv = append(kv, kvset...)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *Action) CheckAccountBalance(from string, amount, frozen int64) bool {
	acc := action.coinsAccount.LoadExecAccount(from, action.execaddr)
	if acc.GetBalance() >= amount && acc.GetFrozen() >= frozen {
		return true
	}
	return false
}

func (action *Action) RandNum(preRand int64) (num int64, err error) {
	newmodify := fmt.Sprintf("%s:%d:%d", string(action.txhash), action.index+int(preRand), action.blocktime)

	modify := common.Sha256([]byte(newmodify))
	baseNum, err := strconv.ParseUint(common.ToHex(modify[0:4]), 0, 64)
	if err != nil {
		flog.Error("RandNum parse uint failed", "error", err)
		return -1, err
	}

	num = int64(baseNum) % ft.GetModeNum()
	return
}

func (action *Action) GetKVSet(param interface{}) (kvset []*types.KeyValue, result interface{}) {
	if bonusInfo, ok := param.(*ft.BonusInfo); ok {
		value := types.Encode(bonusInfo)
		//更新stateDB缓存
		action.db.Set([]byte(calcFluckyBonusInfoKey()), value)
		kvset = append(kvset, &types.KeyValue{Key: []byte(calcFluckyBonusInfoKey()), Value: value})
	}
	if betInfo, ok := param.(*ft.BetInfo); ok {
		value := types.Encode(betInfo)
		action.db.Set([]byte(calcFluckyUserTimesKey(action.fromaddr)), value)
		kvset = append(kvset, &types.KeyValue{Key: []byte(calcFluckyUserTimesKey(action.fromaddr)), Value: value})
	}
	return kvset, nil
}

func (action *Action) GetBetReceiptLog(idx int64, time uint64, amount int64, betInfo *ft.BetInfo) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = ft.TyLogFluckyBet
	r := &ft.ReceiptFlucky{
		Index:   idx,
		Addr:    action.fromaddr,
		Time:    time,
		Amount:  amount,
		RandNum: betInfo.GetRandNum(),
		MaxNum:  betInfo.GetMaxNum(),
		Bonus:   betInfo.GetBonus(),
		Action:  ft.FluckyActionBet,
	}
	log.Log = types.Encode(r)
	return log
}

func getIdxInfo(db dbm.KV, key []byte) (*ft.BetReq, error) {
	value, err := db.Get(key)
	if err != nil {
		flog.Error("getIdxInfo", "can't get value from db,key:", string(key), "err", err.Error())
		return nil, err
	}

	var index ft.BetReq
	err = types.Decode(value, &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

func getMaxRandNum(arr []int64) int64 {
	max := arr[0]
	for _, rand := range arr {
		if rand > max {
			max = rand
		}
	}
	return max
}

func (action *Action) getRandNumOnce() []int64 {
	var randInfo []int64
	var randNum int64
	randNum, err := action.RandNum(0)
	if err != nil {
		flog.Error("get rand num failed", "error", err)
		return nil
	}
	randInfo = append(randInfo, randNum)

	return randInfo
}

func (action *Action) getRandNum(times int64) []int64 {
	var randInfo []int64
	var randNum int64
	randNum, err := action.RandNum(0)
	if err != nil {
		flog.Error("get rand num failed", "error", err)
		return nil
	}
	randInfo = append(randInfo, randNum)
	for i := 1; i < int(times); i++ {
		randNum, err = action.RandNum(randNum)
		if err != nil {
			flog.Error("get rand num failed", "error", err)
			return nil
		}
		randInfo = append(randInfo, randNum)
	}
	return randInfo
}

func rewardBetAmount(maxNum int64, bet *ft.BetInfo, betAmount float32) {
	if maxNum >= 0 && maxNum <= 5000 {
		bet.Bonus = betAmount * 0.1
	} else if maxNum > 5000 && maxNum <= 8000 {
		bet.Bonus = betAmount * 0.5
	} else if maxNum > 8000 && maxNum <= 9000 {
		bet.Bonus = betAmount
	}
}

func rewardBonusPool(maxNum int64, bet *ft.BetInfo, bonusAmount float32) {
	if maxNum > 9000 && maxNum <= 9500 {
		bet.Bonus = bonusAmount * 0.005
	} else if maxNum > 9500 && maxNum <= 9900 {
		bet.Bonus = bonusAmount * 0.01
	} else if maxNum > 9900 && maxNum <= 9950 {
		bet.Bonus = bonusAmount * 0.05
	} else if maxNum > 9950 && maxNum <= 9990 {
		bet.Bonus = bonusAmount * 0.1
	} else if maxNum > 9990 && maxNum < 9999 {
		bet.Bonus = bonusAmount * 0.2
	} else if maxNum == 9999 {
		bet.Bonus = bonusAmount * 0.5
	}
}

func getBonulPoolInfo(db dbm.KV, key []byte) (*ft.BonusInfo, error) {
	value, err := db.Get(key)
	if err != nil {
		flog.Error("db getBonulPoolInfo", "can't get value from db,key:", string(key), "err", err.Error())
		return nil, err
	}

	var bonusInfo ft.BonusInfo
	err = types.Decode(value, &bonusInfo)
	if err != nil {
		return nil, err
	}
	return &bonusInfo, nil
}

func UpdateBonusInfo(bet *ft.BetInfo, betAmount float32, bonusAmount float32) {
	maxNum := getMaxRandNum(bet.GetRandNum())
	bet.MaxNum = maxNum
	rule := ft.GetRewardRule(maxNum)
	switch rule.Type {
	case 0:
		bet.Bonus = betAmount * rule.Percent
		break
	case 1:
		bet.Bonus = bonusAmount * rule.Percent
		break
	}
}

func findBetInfo(db dbm.KV, addr string) (*ft.BetInfo, error) {
	data, err := db.Get([]byte(calcFluckyUserTimesKey(addr)))
	if err != nil {
		flog.Debug("findLottery", "get", err)
		return nil, err
	}
	var bet ft.BetInfo
	//decode
	err = types.Decode(data, &bet)
	if err != nil {
		flog.Debug("findLottery", "decode", err)
		return nil, err
	}
	return &bet, nil
}

func QueryBetTimes(db dbm.KV, query *ft.QueryBetTimes) (types.Message, error) {
	var timeInfo ft.ReplyBetTimes
	info, err := db.Get([]byte(calcFluckyUserTimesKey(query.GetAddr())))
	if err == types.ErrNotFound {
		timeInfo.Times = 0
		return &timeInfo, nil
	} else if err != nil {
		flog.Debug("QueryBetTimes get key failed", "error", err)
		return nil, err
	}

	err = types.Decode(info, &timeInfo)
	if err != nil {
		flog.Debug("QueryBetTimes decode failed", "Decode", err)
		return nil, err
	}
	return &timeInfo, nil
}

func QueryBetInfo(db dbm.KVDB, param *ft.QueryBetInfo) (types.Message, error) {
	query := ft.NewTable(db).GetQuery(db)
	row, err := query.ListOne("", param, nil)
	if err != nil {
		flog.Debug("QueryBetInfo get failed", "error", err)
		return nil, err
	}
	if info, ok := row.Data.(*ft.ReceiptFlucky); ok {
		return info, nil
	}
	return nil, types.ErrNotFound
}

func QueryBonusInfo(db dbm.KV) (types.Message, error) {
	info, err := db.Get([]byte(calcFluckyBonusInfoKey()))
	if err != nil {
		flog.Debug("QueryBonusInfo get key failed", "error", err)
		return nil, err
	}

	var bonusInfo ft.ReplyBonusInfo
	err = types.Decode(info, &bonusInfo)
	if err != nil {
		flog.Debug("QueryBonusInfo decode failed", "Decode", err)
		return nil, err
	}
	return &bonusInfo, nil
}

func List(db dbm.KVDB, stateDB dbm.KV, param *ft.QueryBetInfoBatch) (types.Message, error) {
	return QueryBetListByPage(db, stateDB, param)
}

func QueryBetListByPage(db dbm.KVDB, stateDB dbm.KV, param *ft.QueryBetInfoBatch) (types.Message, error) {
	query := ft.NewTable(db).GetQuery(db)
	var rows []*table.Row
	var err error

	prefix := getAddrPrefix(param.GetAddr())
	count := param.GetCount()
	direction := param.GetDirection()
	if param.Index == 0 {
		rows, err = query.ListIndex("addr", prefix, nil, count, direction)
	} else {
		rows, err = query.ListIndex("addr", prefix, getAddrInfo(param.GetAddr(), param.GetIndex()), count, direction)
	}

	if err != nil {
		flog.Error("QueryBetListByPage list", "error", err)
		return nil, err
	}

	var infos ft.ReplyBetInfoBatch
	for _, val := range rows {
		infos.Bets = append(infos.Bets, val.Data.(*ft.ReceiptFlucky))
	}
	return &infos, nil
}

func getAddrPrefix(addr string) []byte {
	return []byte( fmt.Sprintf("%s", addr))
}

func getAddrInfo(addr string, index int64) []byte {
	return []byte (fmt.Sprintf("%s:%08d", addr, index))
}
