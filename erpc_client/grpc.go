package erpc

import (
	"context"
	balancer "github.com/EAHITechnology/raptor/erpc_balancer/grpc_balancer"
	"google.golang.org/grpc"
	"sync/atomic"
)

type GrpcClient struct {
	conf  *RpcNetConfigInfo
	bs    []*balancer.Balancer
	index uint32
}

func NewGrpcClient(conf *RpcNetConfigInfo) (*GrpcClient, error) {
	bs := []*balancer.Balancer{}

	for idx := 0; idx < conf.MaxConnsPerAddr; idx++ {
		b, err := balancer.NewBalancer(conf.Addr[0])
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}

	return &GrpcClient{
		conf: conf,
		bs:   bs,
	}, nil
}

func (g *GrpcClient) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return g.bs[int(atomic.AddUint32(&g.index, 1))%len(g.bs)].GetConn().Invoke(ctx, method, args, reply, opts...)
}

func (g *GrpcClient) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return g.bs[int(atomic.AddUint32(&g.index, 1))%len(g.bs)].GetConn().NewStream(ctx, desc, method, opts...)
}
