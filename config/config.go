package config

import (
	"errors"
)

var (
	ErrConfigParserTyp = errors.New("config parser error")
)

const (
	TomlConfigParserType = "toml"
	YamlConfigParserType = "yaml"
	YmlConfigParserType  = "yml"
)

type ConfigParser interface {
	Get(key string) interface{}

	GetServiceDiscoveryConfig() (ServiceDiscovery, bool)
	GetDataBaseConfig(key string) (DatabaseConfigInfo, bool)
	GetRedisConfig(key string) (RedisConfigInfo, bool)
	GetRpcConfig(key string) (RpcNetConfigInfo, bool)

	GetDataBaseConfigs() []DatabaseConfigInfo
	GetRedisConfigs() []RedisConfigInfo
	GetRpcConfigs() []RpcNetConfigInfo

	Unmarshal(obj interface{}) error

	Reload() error
}

func NewConfigParser(configCenter ConfigCenterInfo) (ConfigParser, error) {
	if configCenter.FilePath != "" {
		return NewFileConfigParser(configCenter)
	}

	if len(configCenter.ApolloAddr) != 0 {
		return NewApolloConfigParser(configCenter)
	}

	if len(configCenter.EtcdAddrs) != 0 {
		return NewEtcdConfigParser(configCenter)
	}

	return nil, ErrConfigParserTyp
}
