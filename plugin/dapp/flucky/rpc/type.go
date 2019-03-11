package rpc

import (
	rpctypes "github.com/33cn/chain33/rpc/types"
	ft "github.com/33cn/plugin/plugin/dapp/flucky/types"
)

// Jrpc json rpc struct
type Jrpc struct {
	cli *channelClient
}

// Grpc grpc struct
type Grpc struct {
	*channelClient
}

type channelClient struct {
	rpctypes.ChannelClient
}

// Init init grpc param
func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	grpc := &Grpc{channelClient: cli}
	cli.Init(name, s, &Jrpc{cli: cli}, grpc)
	ft.RegisterFluckyServer(s.GRPC(), grpc)
}
