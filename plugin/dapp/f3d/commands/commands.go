/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package commands

import (
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	ptypes "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
	"github.com/spf13/cobra"
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
	)
	return cmd
}

func F3DInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info operation for f3d.",
	}
	cmd.AddCommand(
		userInfoCmd(),
		keyInfoCmd(),
		roundInfoCmd(),
		lastRoundInfoCmd(),
	)
	return cmd
}

func addF3DGameFlags(cmd *cobra.Command) {

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

	params := ptypes.GameDrawReq{Round: round}

	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DLuckyDrawTx", params, &res)
	ctx.RunWithoutMarshal()
}

func userInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Show the user info matched by address.",
		Run:   userInfoQuery,
	}
	addF3DGameFlags(cmd)
	return cmd
}

func keyInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Show the key info about the user matched by address.",
		Run:   keyInfoQuery,
	}
	addF3DGameFlags(cmd)
	return cmd
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

func lastRoundInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last_round",
		Short: "Show the last round info.",
		Run:   lastRoundInfoQuery,
	}
	addF3DGameFlags(cmd)
	return cmd
}

func userInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetInt64("round")
	addr, _ := cmd.Flags().GetString("addr")

	var params rpctypes.Query4Jrpc

	var rep interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryAddrInfo{
		Round: round,
		Addr:  addr,
	}
	params.FuncName = ptypes.Fu
	params.Payload = types.MustPBToJSON(&req)
	rep = &ptypes.ReplyAddrInfoList{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, rep)
	ctx.Run()
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
	params.FuncName = ptypes.FuncNameQueryRoundInfoByRound
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.ReplyF3D{}

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
	resp = &ptypes.ReplyF3D{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

func keyInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetString("round")
	addr, _ := cmd.Flags().GetString("addr")

	var params rpctypes.Query4Jrpc
	var rep interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryBuyRecordByRoundAndAddr{}
}
