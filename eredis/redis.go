package eredis

import (
	"context"
	"sync"
	"time"

	"github.com/EAHITechnology/raptor/utils"
	"github.com/gomodule/redigo/redis"
)

const (
	MAX_IDLE        = 50
	MAX_ACTIVE      = 100
	IDLE_TIMEOUT    = 300
	READ_TIMEOUT    = 500
	WRITE_TIMEOUT   = 500
	CONNECT_TIMEOUT = 1000
)

type RedisLog interface {
	Debugf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Warnf(f string, args ...interface{})
	Errorf(f string, args ...interface{})
}

type RedisDynamicConfiguration interface {
	GetRedisConfig() ([]RedisInfo, error)
	GetRedisConfigTime() time.Duration
}

type Redis struct {
	pool *redis.Pool
	rc   RedisInfo
}

type RedisClients struct {
	Rs     map[string]*Redis
	lock   *sync.RWMutex
	logger RedisLog
}

func checkConfig(r *RedisInfo) error {
	if r.RedisName == "" {
		return REDIS_NOT_CONFIGURED_ERR
	}
	if r.Addr == "" {
		return REDIS_ADDR_NOT_CONFIGURED_ERR
	}
	if r.MaxIdle == 0 {
		r.MaxIdle = MAX_IDLE
	}
	if r.MaxActive == 0 {
		r.MaxActive = MAX_ACTIVE
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = IDLE_TIMEOUT
	}
	if r.ReadTimeout == 0 {
		r.ReadTimeout = READ_TIMEOUT
	}
	if r.WriteTimeout == 0 {
		r.WriteTimeout = WRITE_TIMEOUT
	}
	if r.ConnectTimeout == 0 {
		r.ConnectTimeout = CONNECT_TIMEOUT
	}
	return nil
}

func NewRedis(ctx context.Context, i RedisInfo) (*Redis, error) {
	if err := checkConfig(&i); err != nil {
		return nil, err
	}
	return &Redis{
		pool: &redis.Pool{
			MaxIdle:     i.MaxIdle,
			MaxActive:   i.MaxActive,
			IdleTimeout: time.Duration(i.IdleTimeout) * time.Second,
			Wait:        i.Wait,
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					i.Addr,
					redis.DialPassword(i.Password),
					redis.DialReadTimeout(time.Duration(i.ReadTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(i.WriteTimeout)*time.Millisecond),
					redis.DialConnectTimeout(time.Duration(i.ConnectTimeout)*time.Millisecond),
					redis.DialDatabase(i.Database),
				)
			},
		},
		rc: i,
	}, nil
}

func loadRedis(ctx context.Context, i RedisInfo) error {
	red, err := NewRedis(ctx, i)
	if err != nil {
		return err
	}

	redisclients.lock.Lock()
	defer redisclients.lock.Unlock()
	redisclients.Rs[i.RedisName] = red

	return nil
}

func InitRedis(ctx context.Context, redisInfos []RedisInfo, rdc RedisDynamicConfiguration, logger RedisLog) error {
	if utils.IsNil(logger) || logger == nil {
		return REDIS_LOG_NOT_INIT_ERR
	}

	redisclients = &RedisClients{
		Rs:     make(map[string]*Redis),
		lock:   new(sync.RWMutex),
		logger: logger,
	}

	for _, i := range redisInfos {
		if err := loadRedis(ctx, i); err != nil {
			return err
		}
	}

	if rdc != nil && !utils.IsNil(rdc) {
		go reLoadConfig(ctx, rdc)
	}

	return nil
}

func (r *Redis) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, len(args)+1)
	parmas[0] = key

	if len(args) > 0 {
		for idx, v := range args {
			parmas[idx+1] = v
		}
	}
	return con.Do(cmd, parmas...)
}

func (r *Redis) Del(args ...interface{}) (count int64, err error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()
	var reply interface{}
	reply, err = con.Do("Del", args...)
	if err != nil {
		return
	}

	if reply == nil {
		return 0, ErrNil
	}

	return reply.(int64), nil
}

func (r *Redis) Script(keyCount int, data string, args []string) (interface{}, error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()

	script := redis.NewScript(keyCount, data)
	script.Load(con)

	rargs := redis.Args{}.AddFlat(args)
	reply, err := script.Do(con, rargs...)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *Redis) Set(key string, value interface{}) error {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return err
	}
	defer con.Close()

	reply, err := con.Do("SET", key, value)
	if err != nil {
		return err
	}

	if reply == nil {
		return ErrNil
	}

	if _, err := redis.String(reply, err); err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetString(key string) (string, error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return "", err
	}
	defer con.Close()

	reply, err := con.Do("GET", key)
	if err != nil {
		return "", err
	}

	if reply == nil {
		return "", ErrNil
	}

	return redis.String(reply, err)
}

func (r *Redis) GetInt64(key string) (int64, error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return 0, err
	}
	defer con.Close()

	reply, err := con.Do("GET", key)
	if err != nil {
		return 0, err
	}

	if reply == nil {
		return 0, ErrNil
	}

	return redis.Int64(reply, err)
}

func (r *Redis) GetStringMap(key string) (map[string]string, error) {
	con := r.pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()

	reply, err := con.Do("GET", key)
	if err != nil {
		return nil, err
	}

	if reply == nil {
		return nil, ErrNil
	}

	return redis.StringMap(reply, err)
}

/*
关闭一个链接池
*/
func (r *Redis) Close() error {
	return r.pool.Close()
}

func Close() {
	redisclients.lock.Lock()
	defer redisclients.lock.Unlock()

	for _, rs := range redisclients.Rs {
		if err := rs.Close(); err != nil {
			redisclients.logger.Errorf("redis close err:%v", err)
		}
		redisclients.logger.Infof("redis close name:%s", rs.rc.RedisName)
	}
}

func getClients() *RedisClients {
	return redisclients
}

func GetClient(redisName string) (*Redis, error) {
	redisclients := getClients()
	redisclients.lock.RLock()
	defer redisclients.lock.RUnlock()

	r, ok := redisclients.Rs[redisName]
	if !ok {
		return nil, REDIS_NOT_INIT_ERR
	}
	return r, nil
}

func reLoadConfig(ctx context.Context, rdc RedisDynamicConfiguration) {
	fun := "reLoadConfig -->"
	d := rdc.GetRedisConfigTime()
	if d < reloadTime {
		d = reloadTime
	}
	ticker := time.NewTicker(d)
	defer ticker.Stop()

selectLoop:
	for {
		select {
		case <-ticker.C:
			redisInfos, err := rdc.GetRedisConfig()
			if err != nil {
				redisclients.logger.Errorf("%s GetRedisConfig err:%v", fun, err)
				continue
			}

			redisclients.lock.Lock()
			for _, redisInfo := range redisInfos {
				if rs, ok := redisclients.Rs[redisInfo.RedisName]; ok {
					// No pointer type can be directly compared
					if rs.rc != redisInfo {
						// First need to close the resource, but need add count
						rs.Close()
					}
				}
				if err := loadRedis(ctx, redisInfo); err != nil {
					redisclients.logger.Errorf("%s loadRedis err:%v", fun, err)
				}
			}
			redisclients.lock.Unlock()
		case <-ctx.Done():
			break selectLoop
		}
	}
}
