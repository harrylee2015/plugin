/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package unfreeze

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/gold5g/executor"
	"github.com/33cn/plugin/plugin/dapp/gold5g/commands"
	pt "github.com/33cn/plugin/plugin/dapp/gold5g/ptypes"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     pt.Gold5GX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      nil,
	})
}
