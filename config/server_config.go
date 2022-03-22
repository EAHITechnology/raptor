package config

import "errors"

var (
	ErrConfigParserType = errors.New("config typ nil")
	ErrConfigParserPath = errors.New("config path nil")
)

const (
	TomlServerConfigParserType = "toml"
	YamlServerConfigParserType = "yaml"
	YmlServerConfigParserType  = "yml"
)

type ServerConfigParser interface {
	GetServerConfigInfo() ServerConfigInfo

	GetLogConfigInfo() LogConfigInfo

	GetConfigCenterInfo() ConfigCenterInfo

	Get(key string) interface{}
}

type ConfigInfo struct {
	Typ  string
	Path string
}

func NewServerConfigParser(conf ConfigInfo) (ServerConfigParser, error) {
	if conf.Typ == "" {
		return nil, ErrConfigParserType
	}
	if conf.Path == "" {
		return nil, ErrConfigParserPath
	}
	return NewServerConfigManager(conf)
}
