package eredis

import (
	"time"
)

var (
	redisclients *RedisClients

	reloadTime = 1 * time.Minute
)
