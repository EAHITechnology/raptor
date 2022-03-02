package eredis

import (
	"time"
)

var (
	redisclients *RedisClients

	reloadTime = 1 * time.Minute
)

const (
	UnlockLua = `
local key = KEYS[1]
local old_version = ARGV[1]
local version = redis.call("GET",key)
if (version == old_version)
then
	return redis.call("DEL",key)
else
	return 0
end
`
)
