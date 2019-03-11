package executor

import (
	"fmt"

	"github.com/33cn/chain33/common/address"
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
)

var flog = log.New("module", "execs.flucky")

var fluckyAddr = address.ExecAddress(ft.FluckyX)

var driverName = ft.FluckyX

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Flucky{}))
}

// Init 重命名执行器名称
func Init(name string, sub []byte) {
	var config ft.FluckyCfg
	if sub != nil {
		types.MustDecode(sub, &config)
	}

	ft.SetConfig(&config)

	driverName = name
	ft.FluckyX = driverName
	ft.ExecerFlucky = []byte(driverName)
	drivers.Register(name, newFlucky, types.GetDappFork(driverName, "Enable"))
}

// Blackwhite 几类执行器结构体
type Flucky struct {
	drivers.DriverBase
}

func newFlucky() drivers.Driver {
	c := &Flucky{}
	c.SetChild(c)
	c.SetExecutorType(types.LoadExecutorType(driverName))
	return c
}

// GetName 获取执行器别名
func GetName() string {
	return newFlucky().GetName()
}

// GetDriverName 获取执行器名字
func (f *Flucky) GetDriverName() string {
	return driverName
}

func calcFluckyUserHistoryKey(addr string, index int64) string {
	key := fmt.Sprintf("LODB-flucky-user-history:%s:%d", addr, index)
	return key
}

func calcFluckyUserHistoryPrefix(addr string) string {
	key := fmt.Sprintf("LODB-flucky-user-history:%s", addr)
	return key
}

func (f *Flucky) updateLocalDB(r *ft.ReceiptFlucky) (kvs []*types.KeyValue) {
	switch r.Action {
	case ft.FluckyActionBet:
		start := &ft.ReceiptFlucky{
			Index:   r.GetIndex(),
			Addr:    r.GetAddr(),
			Amount:  r.GetAmount(),
			Time:    r.GetTime(),
			RandNum: r.GetRandNum(),
			MaxNum:  r.GetMaxNum(),
			Bonus:   r.GetBonus(),
		}
		kvs = append(kvs, &types.KeyValue{Key: []byte(calcFluckyUserHistoryKey(r.GetAddr(), r.GetIndex())), Value: types.Encode(start)})
	}
	return kvs
}
