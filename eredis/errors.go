package eredis

import (
	"errors"
)

var (
	REDIS_NOT_INIT_ERR            = errors.New("Redis Not Init")
	REDIS_LOG_NOT_INIT_ERR        = errors.New("Redis Logger Not Init")
	REDIS_NOT_CONFIGURED_ERR      = errors.New("RedisName Is Not Configured")
	REDIS_ADDR_NOT_CONFIGURED_ERR = errors.New("Redis Addr Is Not Configured")

	ErrNil = errors.New("nil returned")
)
