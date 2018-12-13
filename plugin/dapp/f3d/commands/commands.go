/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package commands

import (
	"fmt"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	ptypes "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
	"github.com/spf13/cobra"
	"time"
)

func F3DCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "f3d",
		Short: "F3D contracts operation",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		F3DGameCmd(),
		F3DInfoCmd(),
	)

	return cmd
}

func F3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "game",
		Short: "game operation for f3d.",
	}
	cmd.AddCommand(
		createF3DGameCmd(),
		luckyDrawF3DGameCmd(),
		buyKeyF3DGameCmd(),
	)
	return cmd
}

func F3DInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info operation for f3d.",
	}
	cmd.AddCommand(
		recordInfoCmd(),
		roundInfoCmd(),
		roundsInfoCmd(),
		lastRoundInfoCmd(),
	)
	return cmd
}

func addF3DGameFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("round", "r", 0, "Game Round")
	cmd.MarkFlagRequired("round")
}

func createF3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a new f3d game.",
		Run:   create,
	}
	addF3DGameFlags(cmd)
	return cmd
}

func create(cmd *cobra.Command, args []string) {
	round, _ := cmd.Flags().GetInt64("round")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	params := ptypes.GameStartReq{Round: round}

	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DStartTx", params, &res)
	ctx.RunWithoutMarshal()
}

func luckyDrawF3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draw",
		Short: "Send the bonus to the last user that buys the key",
		Run:   luckyDraw,
	}
	addF3DGameFlags(cmd)
	return cmd
}

func luckyDraw(cmd *cobra.Command, args []string) {
	round, _ := cmd.Flags().GetInt64("round")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	// 检测是否到了开奖时间
	//go func() {
	var interval time.Duration
	for {
		roundInfo := getLastRoundInfo(rpcLaddr)
		if !roundCheck(round, roundInfo) {
			return
		}
		if remainTimeCheck(roundInfo, &interval) {
			fmt.Println("Begin to luckydraw, time:", time.Now().Unix())
			params := ptypes.GameDrawReq{Round: round}

			var res string
			ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DLuckyDrawTx", params, &res)
			ctx.RunWithoutMarshal()
			break
		}
		fmt.Println("It 's not time to luckydraw , remainTime is ", roundInfo.RemainTime, ", interval:", interval*time.Second, ", now:", time.Now().Unix())
		time.Sleep(interval * time.Second)
		continue
	}
	//}()
	// 开启新一轮
	params := ptypes.GameStartReq{Round: round + 1}

	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DStartTx", params, &res)
	ctx.RunWithoutMarshal()
}

func getLastRoundInfo(rpcAddr string) *ptypes.RoundInfo {
	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	params.FuncName = ptypes.FuncNameQueryLastRoundInfo
	params.Payload = types.MustPBToJSON(&ptypes.QueryF3DLastRound{})
	resp = &ptypes.RoundInfo{}

	ctx := jsonclient.NewRPCCtx(rpcAddr, "Chain33.Query", params, resp)
	ctx.Run()

	return resp.(*ptypes.RoundInfo)
}

func roundCheck(round int64, info *ptypes.RoundInfo) bool {
	if round > 0 && round == info.Round {
		return true
	}
	return false
}

func remainTimeCheck(info *ptypes.RoundInfo, interval *time.Duration) bool {
	currentTime := time.Now().Unix()
	remainTime := info.RemainTime + info.UpdateTime - currentTime

	if remainTime > 600 {
		*interval = 600
	} else if remainTime > 300 {
		*interval = 60
	} else if remainTime > 60 {
		*interval = 30
	} else if remainTime > 0 {
		*interval = 10
	} else {
		return true
	}

	return false
}

func buyKeyF3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy",
		Short: "Buy some keys",
		Run:   buyKeys,
	}
	addNumberFlags(cmd)
	return cmd
}

func buyKeys(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	num, _ := cmd.Flags().GetInt64("num")

	params := ptypes.GameBuyKeysReq{
		Num: num,
	}

	var resp string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DBuyKeysTx", params, &resp)
	ctx.RunWithoutMarshal()
}

func addNumberFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("num", "n", 0, "the number of keys that you want to buy.")
	cmd.MarkFlagRequired("num")
}

func recordInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Show the buy record about the user matched by address and round.",
		Run:   buyRecordInfoQuery,
	}
	addRecordInfoQueryFlags(cmd)
	return cmd
}

func addRecordInfoQueryFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("round", "r", 0, "Game Round")
	cmd.Flags().Int64P("index", "i", 0, "index")
	cmd.Flags().StringP("addr", "a", "", "user addr")
	cmd.MarkFlagRequired("round")
}

func roundInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "round",
		Short: "Show the round info matched by round.",
		Run:   roundInfoQuery,
	}
	addF3DGameFlags(cmd)
	return cmd
}

func roundsInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rounds",
		Short: "Show the round info matched by round.",
		Run:   roundsInfoQuery,
	}
	addRoundsInfoQueryFlag(cmd)
	return cmd
}

func addRoundsInfoQueryFlag(cmd *cobra.Command) {
	cmd.Flags().Int64P("startRound", "s", 0, "start round")
	cmd.Flags().Int32P("direction", "d", 0, "query direction, 0: desc  1:asc")
	cmd.Flags().Int32P("count", "c", 0, "query amount")

	cmd.MarkFlagRequired("start round")
}

func lastRoundInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last_round",
		Short: "Show the last round info.",
		Run:   lastRoundInfoQuery,
	}

	return cmd
}

func roundInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetInt64("round")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryF3DByRound{
		Round: round,
	}
	fmt.Println(req)

	params.FuncName = ptypes.FuncNameQueryRoundInfoByRound
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.RoundInfo{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

func roundsInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	startRound, _ := cmd.Flags().GetInt64("startRound")
	direction, _ := cmd.Flags().GetInt32("direction")
	count, _ := cmd.Flags().GetInt32("count")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryF3DListByRound{
		StartRound: startRound,
		Direction:  direction,
		Count:      count,
	}
	params.FuncName = ptypes.FuncNameQueryRoundsInfoByRounds
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.ReplyF3DList{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

func lastRoundInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryF3DLastRound{}
	params.FuncName = ptypes.FuncNameQueryLastRoundInfo
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.RoundInfo{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

func buyRecordInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetInt64("round")
	addr, _ := cmd.Flags().GetString("addr")
	index, _ := cmd.Flags().GetInt64("index")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryBuyRecordByRoundAndAddr{
		Round: round,
		Addr:  addr,
		Index: index,
	}
	params.FuncName = ptypes.FuncNameQueryBuyRecordByRoundAndAddr
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.ReplyBuyRecord{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}
