package grpc_balancer

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/EAHITechnology/raptor/breaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultMaxSendRecvSize int           = 1024 * 1024 * 32
	defaultWindowSize      int32         = 1024 * 1024 * 1024
	defaultDialTimeOut     time.Duration = time.Second * 3
	defaultReloadTime      time.Duration = time.Second * 10
	defaultConnPoolSize    int           = 1
	defaultFailRate        int64         = 40
)

var (
	ErrEmptyReloader       = errors.New("reloader is nil")
	ErrEmptyShopeeFilePath = errors.New("shopee file is nil")
	ErrEmptyService        = errors.New("service is nil")
)

type Balancer struct {
	conn *fseConn
}

func init() {
	// resolver
	spaceBuilderRegister()
	listBuilderRegister()

	// balancer
	randomRegister()

	// metric
	grpc_prometheus.EnableClientHandlingTimeHistogram()
}

type fseConn struct {
	// rpc conn
	proxyConn *grpc.ClientConn

	// break
	breaker *breaker.CircuitBreaker
}

func (fc *fseConn) Breaker(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return fc.breaker.Do(ctx,
		func(ctx context.Context) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		},
		nil,
	)
}

func getLoadBalancingPolicy(opts *balancerConfig) string {
	var loadBalancingPolicy = `{"LoadBalancingPolicy": "%s"}`
	switch opts.balancerType {
	case RangeType:
		loadBalancingPolicy = fmt.Sprintf(loadBalancingPolicy, RangeType)
	case RandomType:
		loadBalancingPolicy = fmt.Sprintf(loadBalancingPolicy, RandomType)
	default:
		loadBalancingPolicy = fmt.Sprintf(loadBalancingPolicy, RangeType)
	}
	return loadBalancingPolicy
}

func checkNamingPolicy(proxyAddr string) error {
	if !strings.Contains(proxyAddr, "://") {
		return nil
	}

	if _, err := url.Parse(proxyAddr); err != nil {
		return err
	}

	return nil
}

func newGrpcConn(proxyAddr string, opts *balancerConfig) (*fseConn, error) {
	fc := &fseConn{}
	connect := func(proxyAddr string, opts *balancerConfig, fc *fseConn) error {
		if err := checkNamingPolicy(proxyAddr); err != nil {
			return err
		}

		fc.breaker = breaker.NewCircuitBreaker(
			breaker.NewCircuitBreakerConfig().
				SetFailureRateThreshold(opts.fuseFailureRate).
				SetMinQPS(opts.fuseMinQps).
				SetOpenStatusDurationMs(opts.fuseOpenMsTime).
				SetCellIntervalMs(1000).
				SetSize(10))

		var grpcDialOpts = []grpc.DialOption{}
		grpcDialOpts = append(grpcDialOpts,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithInitialConnWindowSize(defaultWindowSize),
			grpc.WithInitialWindowSize(defaultWindowSize),
			grpc.WithChainUnaryInterceptor(fc.Breaker),
			grpc.WithDefaultServiceConfig(getLoadBalancingPolicy(opts)),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(opts.maxRecvMsgSize),
				grpc.MaxCallSendMsgSize(opts.maxSendMsgSize),
			),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  100 * time.Millisecond,
					Multiplier: 1.6,
					Jitter:     0.2,
					MaxDelay:   1 * time.Second,
				},
			}),
			grpc.WithChainUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		)

		if opts.dialBlock {
			grpcDialOpts = append(grpcDialOpts, grpc.WithBlock())
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, opts.dialTimeOut)
		defer cancel()

		proxyConn, err := grpc.DialContext(ctx,
			proxyAddr,
			grpcDialOpts...,
		)
		if err != nil {
			return err
		}

		fc.proxyConn = proxyConn

		return nil
	}

	if err := connect(proxyAddr, opts, fc); err != nil {
		return nil, err
	}

	return fc, nil
}

func NewBalancer(proxyAddr string, opts ...BalancerConfigFunc) (*Balancer, error) {
	conf := newDefaultOption()
	for _, opt := range opts {
		opt(conf)
	}

	fseconn, err := newGrpcConn(proxyAddr, conf)
	if err != nil {
		return nil, err
	}

	return &Balancer{
		conn: fseconn,
	}, nil
}

func (b *Balancer) GetConn() *grpc.ClientConn {
	return b.conn.proxyConn
}

func (b *Balancer) Close() {
	b.conn.proxyConn.Close()
}
