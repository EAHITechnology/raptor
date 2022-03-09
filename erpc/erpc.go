package erpc

import "context"

var HttpManager *HttpClientManager

type RpcClient interface {
}

func InitDrpc(ctx context.Context, rpcNetConfigs []RpcNetConfigInfo) error {
	return nil
}
