package flucky

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/chain33private/cezapara/plugin/dapp/flucky/commands"
	"github.com/chain33private/cezapara/plugin/dapp/flucky/executor"
	"github.com/chain33private/cezapara/plugin/dapp/flucky/rpc"
	"github.com/chain33private/cezapara/plugin/dapp/flucky/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.FluckyX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.FluckyCmd,
		RPC:      rpc.Init,
	})
}
