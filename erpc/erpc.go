package erpc

import "context"

type RpcClient interface {
}

func InitErpc(ctx context.Context, rpcNetConfigs []RpcNetConfigInfo) error {
	return nil
}
