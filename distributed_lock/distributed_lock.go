package distributed_lock

import (
	"golang.org/x/net/context"
)

type DistributedLockType string

const (
	EtcdDistributedLock      DistributedLockType = "etcd"
	ZkDistributedLock        DistributedLockType = "zk"
	RedLockDistributedLock   DistributedLockType = "red_lock"
	CustomizeDistributedLock DistributedLockType = "customize_lock"
)

type DistributedLockConfig struct {
	Type           DistributedLockType
	TTl            int64
	Key            string
	EtcdAddrs      []string
	ZkAddrs        []string
	RedisAddrs     []string
	CustomizeAddrs []string
}

type DistributedLockManager interface {
	Lock(context.Context) (string, error)
	Unlock(ctx context.Context, value string) error
}

func NewDistributedLockManager(ctx context.Context, conf DistributedLockConfig) (DistributedLockManager, error) {
	switch conf.Type {
	case EtcdDistributedLock:
		if len(conf.EtcdAddrs) == 0 {
			return nil, ErrEtcdAddrNil
		}
		return NewEtcdDistributedLockManager(ctx, conf)
	case ZkDistributedLock:
		if len(conf.ZkAddrs) == 0 {
			return nil, ErrZkAddrNil
		}
		return NewEtcdDistributedLockManager(ctx, conf)
	case RedLockDistributedLock:
		if len(conf.RedisAddrs) == 0 {
			return nil, ErrRedisAddrNil
		}
		return NewEtcdDistributedLockManager(ctx, conf)
	case CustomizeDistributedLock:
		if len(conf.CustomizeAddrs) == 0 {
			return nil, ErrCustomizeAddrNil
		}
		return NewEtcdDistributedLockManager(ctx, conf)
	default:
		return nil, ErrDistributedLockType
	}
}
