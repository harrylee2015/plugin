package types

import (
	"fmt"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
)

var opt = &table.Option{
	Prefix:  "LODB",
	Name:    "flucky",
	Primary: "addr:index",
	Index:   []string{"addr", "index", "time", "amount", "randnum", "maxnum", "bonus"},
}

//NewTable 新建表
func NewTable(kvdb db.KV) *table.Table {
	rowmeta := NewFluckyRow()
	table, err := table.NewTable(rowmeta, kvdb, opt)
	if err != nil {
		panic(err)
	}
	return table
}

//OracleRow table meta 结构
type FluckyRow struct {
	*ReceiptFlucky
}

//NewOracleRow 新建一个meta 结构
func NewFluckyRow() *FluckyRow {
	return &FluckyRow{ReceiptFlucky: &ReceiptFlucky{}}
}

//CreateRow 新建数据行(注意index 数据一定也要保存到数据中)
func (tx *FluckyRow) CreateRow() *table.Row {
	return &table.Row{Data: &ReceiptFlucky{}}
}

//SetPayload 设置数据
func (tx *FluckyRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*QueryBetInfo); ok {
		tx.ReceiptFlucky.Index = txdata.Idx
		tx.ReceiptFlucky.Addr = txdata.Addr
		return nil
	} else if txdata, ok := data.(*ReceiptFlucky); ok {
		tx.ReceiptFlucky = txdata
		return nil
	} else if txdata, ok := data.(*QueryBetInfoBatch); ok {
		tx.ReceiptFlucky.Addr = txdata.GetAddr()
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *FluckyRow) Get(key string) ([]byte, error) {
	if key == "addr" {
		return []byte(tx.Addr), nil
	} else if key == "index" {
		return []byte(fmt.Sprintf("%08d", tx.Index)), nil
	} else if key == "time" {
		return []byte(fmt.Sprintf("%d", tx.Time)), nil
	} else if key == "amount" {
		return []byte(fmt.Sprintf("%2d", tx.Amount)), nil
	} else if key == "randnum" {
		return []byte(fmt.Sprintf("%2d", tx.RandNum)), nil
	} else if key == "maxnum" {
		return []byte(fmt.Sprintf("%2d", tx.MaxNum)), nil
	} else if key == "bonus" {
		return []byte(fmt.Sprintf("%2d", tx.Bonus)), nil
	} else if key == "addr:index" {
		return [] byte(fmt.Sprintf("%s:%08d", tx.Addr, tx.Index)), nil
	} else if key == "addrprefix" {
		return [] byte(fmt.Sprintf("%s:", tx.Addr)), nil
	}
	return nil, types.ErrNotFound
}
