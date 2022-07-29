package grpc_balancer

import (
	"time"
)

type BalancerType string

type BalancerConfigFunc func(*balancerConfig) *balancerConfig

const (
	Nonetype   BalancerType = ""
	RandomType BalancerType = "random"
	RangeType  BalancerType = "round_robin"
)

type balancerConfig struct {
	maxSendMsgSize int
	maxRecvMsgSize int
	dialBlock      bool
	dialTimeOut    time.Duration
	connPoolSize   int
	balancerType   BalancerType

	fuseFailureRate int64
	fuseMinQps      int64
	fuseOpenMsTime  int64
}

func newDefaultOption() *balancerConfig {
	return &balancerConfig{
		balancerType:    RangeType,
		maxSendMsgSize:  defaultMaxSendRecvSize,
		maxRecvMsgSize:  defaultMaxSendRecvSize,
		dialBlock:       false,
		dialTimeOut:     defaultDialTimeOut,
		connPoolSize:    1,
		fuseFailureRate: defaultFailRate,
		fuseMinQps:      1,
		fuseOpenMsTime:  5000,
	}
}

func WithMaxRecvMsgSize(maxRecvMsgSize int) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.maxRecvMsgSize = maxRecvMsgSize
		return b
	}
}

func WithMaxSendMsgSize(maxSendMsgSize int) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.maxSendMsgSize = maxSendMsgSize
		return b
	}
}

// if dialBlock is true.The program will wait for the rpc
// connection complete before starting the service,
// otherwise the program will establish the connection asynchronously.
func WithDialBlock(flag bool) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.dialBlock = flag
		return b
	}
}

func WithDialTimeOut(dialTimeOut time.Duration) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.dialTimeOut = dialTimeOut
		return b
	}
}

// conn pool config.
// Only valid for fseaddr schema.
func WithConnPoolSize(connPoolSize int) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.connPoolSize = connPoolSize
		return b
	}
}

// the balancerType is Nonetype(RangeType) by default,
// you can set balancerType with RandomType or RangeType.
func WithConnBalancerType(balancerType BalancerType) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.balancerType = balancerType
		return b
	}
}

func WithFuseFailureRate(failureRate int64) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.fuseFailureRate = failureRate
		return b
	}
}

func WithFuseMinQps(minQps int64) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.fuseMinQps = minQps
		return b
	}
}

func WithFuseOpenMsTime(openMsTime int64) BalancerConfigFunc {
	return func(b *balancerConfig) *balancerConfig {
		b.fuseOpenMsTime = openMsTime
		return b
	}
}
