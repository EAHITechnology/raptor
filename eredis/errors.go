package eredis

import (
	"errors"
)

var (
	ErrRespNil                = errors.New("nil returned")
	ErrRedisNotInit           = errors.New("Redis Not Init")
	ErrRedisLogNotInit        = errors.New("Redis Logger Not Init")
	ErrRedisNotConfigured     = errors.New("RedisName Is Not Configured")
	ErrRedisAddrNotConfigured = errors.New("Redis Addr Is Not Configured")
)
