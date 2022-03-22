package config

import (
	"bytes"
	"io/ioutil"
	"sync"

	"github.com/spf13/viper"
)

type FileConfigParser struct {
	configCenter ConfigCenterInfo

	config      Config
	parser      *viper.Viper
	databaseMap map[string]DatabaseConfigInfo
	redisMap    map[string]RedisConfigInfo
	rpcMap      map[string]RpcNetConfigInfo

	lock sync.RWMutex
}

func NewFileConfigParser(configCenter ConfigCenterInfo) (ConfigParser, error) {
	fc := &FileConfigParser{}
	fc.configCenter = configCenter
	fc.databaseMap = make(map[string]DatabaseConfigInfo)
	fc.redisMap = make(map[string]RedisConfigInfo)
	fc.rpcMap = make(map[string]RpcNetConfigInfo)

	if err := fc.loadConfig(); err != nil {
		return nil, err
	}

	return fc, nil
}

func (f *FileConfigParser) loadConfig() error {
	fileB, err := ioutil.ReadFile(f.configCenter.FilePath)
	if err != nil {
		return err
	}

	parser := viper.New()
	parser.SetConfigType(f.configCenter.FileType)
	if err := parser.ReadConfig(bytes.NewBuffer(fileB)); err != nil {
		return err
	}

	c := Config{}
	if err := parser.Unmarshal(&c); err != nil {
		return err
	}

	f.config = c
	f.parser = parser

	databaseMap := make(map[string]DatabaseConfigInfo)
	for _, v := range f.config.DatabaseConfigs {
		databaseMap[v.Name] = v
	}
	f.databaseMap = databaseMap

	redisMap := make(map[string]RedisConfigInfo)
	for _, v := range f.config.RedisConfigs {
		redisMap[v.Name] = v
	}
	f.redisMap = redisMap

	rpcMap := make(map[string]RpcNetConfigInfo)
	for _, v := range f.config.RpcNetConfigs {
		rpcMap[v.ServiceName] = v
	}
	f.rpcMap = rpcMap

	return nil
}

func (f *FileConfigParser) Get(key string) interface{} {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.parser.Get(key)
}

func (f *FileConfigParser) GetServiceDiscoveryConfig() (ServiceDiscovery, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if len(f.config.ServiceDiscovery.EtcdAddr) != 0 {
		return f.config.ServiceDiscovery, true
	}

	if len(f.config.ServiceDiscovery.ZkAddr) != 0 {
		return f.config.ServiceDiscovery, true
	}

	if len(f.config.ServiceDiscovery.CustomServiceDiscovery) != 0 {
		return f.config.ServiceDiscovery, true
	}

	return f.config.ServiceDiscovery, false
}

func (f *FileConfigParser) GetDataBaseConfig(key string) (DatabaseConfigInfo, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	val, ok := f.databaseMap[key]
	return val, ok
}

func (f *FileConfigParser) GetRedisConfig(key string) (RedisConfigInfo, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	val, ok := f.redisMap[key]
	return val, ok
}

func (f *FileConfigParser) GetRpcConfig(key string) (RpcNetConfigInfo, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	val, ok := f.rpcMap[key]
	return val, ok
}

func (f *FileConfigParser) GetDataBaseConfigs() []DatabaseConfigInfo {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.config.DatabaseConfigs
}

func (f *FileConfigParser) GetRedisConfigs() []RedisConfigInfo {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.config.RedisConfigs
}

func (f *FileConfigParser) GetRpcConfigs() []RpcNetConfigInfo {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.config.RpcNetConfigs
}

func (f *FileConfigParser) Unmarshal(obj interface{}) error {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.parser.Unmarshal(obj)
}

func (f *FileConfigParser) Reload() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.loadConfig()
}
