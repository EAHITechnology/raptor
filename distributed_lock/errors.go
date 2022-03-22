package distributed_lock

import "errors"

var (
	ErrLockFail   = errors.New("lock fail")
	ErrUnLockFail = errors.New("Unlock fail")

	ErrDistributedLockType = errors.New("error distributed lock type")

	ErrEtcdAddrNil      = errors.New("etcd addr nil")
	ErrZkAddrNil        = errors.New("zk addr nil")
	ErrRedisAddrNil     = errors.New("redis addr nil")
	ErrCustomizeAddrNil = errors.New("customize addr nil")
)
