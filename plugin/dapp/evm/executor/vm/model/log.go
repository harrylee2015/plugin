// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"math/big"
	"strings"

	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
)

// ContractLog 合约在日志，对应EVM中的Log指令，可以生成指定的日志信息
// 目前这些日志只是在合约执行完成时进行打印，没有其它用途
type ContractLog struct {
	// Address 合约地址
	Address common.Address

	// TxHash 对应交易哈希
	TxHash common.Hash

	// Index 日志序号
	Index int

	// Topics 此合约提供的主题信息
	Topics []common.Hash

	// Data 日志数据
	Data []byte
	// 区块高度
	BlockNumber uint64 `json:"blockNumber"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash"`
}

// PrintLog 合约日志打印格式
func (log *ContractLog) PrintLog(routerAbiStr string) {

	//routerAbiStr := mdb.GetAbi("1GMsmmzUPuQUkCJinuEyfkJoBoyJbiQKg")

	routerAbi, err := ethAbi.JSON(strings.NewReader(routerAbiStr))
	if err != nil {
		panic("Failed to read json")
	}

	eventDebug := routerAbi.Events["debug"].ID().Hex()
	eventLogRunData := routerAbi.Events["logRunData"].ID().Hex()

	if log.Topics[0].Hex() == eventDebug {
		type EventDebug struct {
			Des string
			Pos *big.Int
		}
		event := &EventDebug{}
		eventName := "debug"
		err = routerAbi.Unpack(event, eventName, log.Data)
		if err != nil {
			panic("Failed to unpack debug event")
		}
		log15.Debug("!Contract Log debug", "EventDebug", event)

	} else if log.Topics[0].Hex() == eventLogRunData {
		//logRunData(string func, address pair, address to, uint amountToken, uint amountETH, uint liquidity);
		type EventRunData struct {
			FuncName    string
			Pair        *common.Address
			To          *common.Address
			AmountToken *big.Int
			AmountEth   *big.Int
			Liquidity   *big.Int
		}
		event := &EventRunData{}
		eventName := "logRunData"
		err = routerAbi.Unpack(event, eventName, log.Data)
		if err != nil {
			panic("Failed to unpack debug event")
		}
		log15.Debug("!Contract Log debug", "logRunData", event)
	}

	log15.Debug("!Contract Log!", "Contract address", log.Address.String(), "TxHash", log.TxHash.Hex(), "Log Index", log.Index, "Log Topics", log.Topics, "Log Data", common.Bytes2Hex(log.Data))
}
