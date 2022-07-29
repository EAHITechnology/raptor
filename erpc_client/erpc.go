package erpc

import "context"

type RpcClient interface {
	Send()
}

func InitErpc(ctx context.Context, rpcNetConfigs []RpcNetConfigInfo) error {
	return nil
}
