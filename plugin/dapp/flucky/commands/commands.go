package commands

import (
	"fmt"

	jsonrpc "github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
	"github.com/spf13/cobra"
)

func FluckyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flucky",
		Short: "flucky game management",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		FluckyGameCmd(),
		FluckyInfoCmd(),
	)

	return cmd
}

func FluckyGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "game",
		Short: "flucky game operation",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		FluckyBetCmd(),
	)

	return cmd
}

func FluckyInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "flucky info query",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		FluckyBetInfoQueryCmd(),
		FluckyBetInfoBatchQueryCmd(),
		FluckyBetTimesQueryCmd(),
		FluckyBonusInfoQueryCmd(),
	)

	return cmd
}

func FluckyBetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bet",
		Short: "flucky bet operation",
		Run:   fluckybet,
	}

	addFluckyBetCmdFlag(cmd)

	return cmd
}

func addFluckyBetCmdFlag(cmd *cobra.Command) {
	cmd.Flags().Int64P("amount", "a", 1, "the amount of the bet operation")
}

func fluckybet(cmd *cobra.Command, args []string) {
	amount, _ := cmd.Flags().GetInt64("amount")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	if amount <= 0 {
		amount = 1
	}

	payload := fmt.Sprintf("{\"amount\":\"%d\"}", amount)
	params := &rpctypes.CreateTxIn{
		Execer:     types.ExecName(ft.FluckyX),
		ActionName: ft.FluckyBetAction,
		Payload:    []byte(payload),
	}

	var res string
	ctx := jsonrpc.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", params, &res)
	ctx.RunWithoutMarshal()
}

func FluckyBetInfoQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bet",
		Short: "bet info query",
		Run:   fluckyBetInfoQuery,
	}

	addFluckyBetInfoFlag(cmd)

	return cmd
}

func addFluckyBetInfoFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("addr", "a", "", "addr want to query")
	cmd.MarkFlagRequired("addr")

	cmd.Flags().Int64P("index", "i", 1, "index want to query")
	cmd.MarkFlagRequired("index")
}

func fluckyBetInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	addr, _ := cmd.Flags().GetString("addr")
	index, _ := cmd.Flags().GetInt64("index")

	var params rpctypes.Query4Jrpc
	//var resp interface{}

	params.Execer = ft.FluckyX
	req := ft.QueryBetInfo{
		Addr: addr,
		Idx:  index,
	}
	params.FuncName = ft.FuncNameQueryBetInfo
	params.Payload = types.MustPBToJSON(&req)
	var resp ft.BetInfo

	ctx := jsonrpc.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func FluckyBetInfoBatchQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bets",
		Short: "bet info batch query",
		Run:   fluckyBetInfoBatchQuery,
	}
	addFluckyBetBatchInfoFlag(cmd)
	return cmd
}

func addFluckyBetBatchInfoFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("addr", "a", "", "addr for batch info query")
	cmd.Flags().Int64P("index", "i", 0, "index for batch info query")
	cmd.Flags().Int32P("count", "c", 1, "count for batch info query")
	cmd.Flags().Int32P("direction", "d", 1, "direction for batch info query")
	cmd.MarkFlagRequired("addr")
}

func fluckyBetInfoBatchQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	addr, _ := cmd.Flags().GetString("addr")
	index, _ := cmd.Flags().GetInt64("index")
	count, _ := cmd.Flags().GetInt32("count")
	direction, _ := cmd.Flags().GetInt32("direction")

	var params rpctypes.Query4Jrpc
	//var resp interface{}

	params.Execer = ft.FluckyX
	req := ft.QueryBetInfoBatch{
		Addr:      addr,
		Index:     index,
		Count:     count,
		Direction: direction,
	}
	params.FuncName = ft.FuncNameQueryBetBatchInfo
	params.Payload = types.MustPBToJSON(&req)
	var resp ft.ReplyBetInfoBatch

	ctx := jsonrpc.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func FluckyBetTimesQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "times",
		Short: "bet times",
		Run:   fluckyBetTimesQuery,
	}

	addFluckyBetTimesFlag(cmd)
	return cmd
}

func addFluckyBetTimesFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("addr", "a", "", "addr for the query request")
	cmd.MarkFlagRequired("addr")
}

func fluckyBetTimesQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	addr, _ := cmd.Flags().GetString("addr")

	var params rpctypes.Query4Jrpc

	params.Execer = ft.FluckyX
	req := ft.QueryBetTimes{
		Addr: addr,
	}
	params.FuncName = ft.FuncNameQueryBetTimesInfo
	params.Payload = types.MustPBToJSON(&req)
	var resp ft.ReplyBetTimes
	ctx := jsonrpc.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}

func FluckyBonusInfoQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bonus",
		Short: "flucky bonus info query",
		Run:   fluckyBonusInfoQuery,
	}
	return cmd
}

func fluckyBonusInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var params rpctypes.Query4Jrpc
	//var resp interface{}

	params.Execer = ft.FluckyX
	params.FuncName = ft.FuncNameQueryBonusInfo

	var resp ft.BonusInfo

	ctx := jsonrpc.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &resp)
	ctx.Run()
}
