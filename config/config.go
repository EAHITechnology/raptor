package config

import (
	"errors"

	"golang.org/x/net/context"
)

type ConfigParser interface {
	// 读取启动用元数据
	Read(path string) error

	// 加载业务配置
	InitConfig(ctx context.Context) error

	// Unmarshal 需要在调用 Read 获取完元数据之后再进行调用
	Unmarshal(obj interface{}) error
}

const (
	TomlConfigParserType = "toml"
	YamlConfigParserType = "yaml"
	YmlConfigParserType  = "yml"
)

var (
	ErrConfigParserType = errors.New("config parser error")
)

func NewConfigParser(typ string) (ConfigParser, error) {
	switch typ {
	case TomlConfigParserType:
		return NewTomlConfigManager(), nil
	case YamlConfigParserType, YmlConfigParserType:
		return nil, ErrConfigParserType
	default:
		return nil, ErrConfigParserType
	}
}
