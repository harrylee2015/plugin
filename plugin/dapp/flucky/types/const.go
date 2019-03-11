package types

import "github.com/33cn/chain33/types"
import "reflect"

const (
	FluckyActionBet = 3300
	TyLogFluckyBet  = 3333

	Decimal           = 100000000
	RandLuckyNum      = 10000
	RandLuckyBlockNum = 5

	FluckyBetAction = "Bet"

	FuncNameQueryBetInfo   = "QueryBetInfo"
	FuncNameQueryBetBatchInfo = "QueryBetInfoBatch"
	FuncNameQueryBetTimesInfo = "QueryBetTimes"
	FuncNameQueryBonusInfo = "QueryBonusInfo"
)

var (
	FluckyX      = "flucky"
	ExecerFlucky = []byte(FluckyX)

	actionName = map[string]int32{
		"Bet": FluckyActionBet,
	}

	logInfo = map[int64]*types.LogInfo{
		TyLogFluckyBet: {Ty: reflect.TypeOf(ReceiptFlucky{}), Name: "LogFluckyBet"},
	}
)
