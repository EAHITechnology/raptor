package erpc

import "context"

type RpcClient interface {
}

func InitDrpc(ctx context.Context, rpcNetConfigs []RpcNetConfigInfo) error {
	return nil
}
