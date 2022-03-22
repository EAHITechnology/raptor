package config

import (
	"bytes"
	"io/ioutil"

	"github.com/spf13/viper"
)

type ServerConfigManager struct {
	metaPath     string
	ServerConfig ServerConfig
	parser       *viper.Viper
}

func NewServerConfigManager(conf ConfigInfo) (ServerConfigParser, error) {
	scm := &ServerConfigManager{
		metaPath: conf.Path,
	}

	sc := ServerConfig{}
	sc.LogConfig = LogConfigInfo{}
	sc.ServerConfig = ServerConfigInfo{}
	sc.ConfigCenter = ConfigCenterInfo{}

	fileB, err := ioutil.ReadFile(conf.Path)
	if err != nil {
		return nil, err
	}

	parser := viper.New()
	parser.SetConfigType(conf.Typ)
	if err := parser.ReadConfig(bytes.NewBuffer(fileB)); err != nil {
		return nil, err
	}

	if err := parser.Unmarshal(&sc); err != nil {
		return nil, err
	}

	scm.ServerConfig = sc
	scm.parser = parser

	return scm, nil
}

func (scm *ServerConfigManager) Get(key string) interface{} {
	return scm.parser.Get(key)
}

func (scm *ServerConfigManager) GetServerConfigInfo() ServerConfigInfo {
	return scm.ServerConfig.ServerConfig
}

func (scm *ServerConfigManager) GetLogConfigInfo() LogConfigInfo {
	return scm.ServerConfig.LogConfig
}

func (scm *ServerConfigManager) GetConfigCenterInfo() ConfigCenterInfo {
	return scm.ServerConfig.ConfigCenter
}
