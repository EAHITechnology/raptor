package config

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/EAHITechnology/raptor/elog"
	"github.com/EAHITechnology/raptor/emysql"
	"github.com/EAHITechnology/raptor/enet"
	"github.com/EAHITechnology/raptor/eredis"
	"github.com/EAHITechnology/raptor/erpc"
	"golang.org/x/net/context"
)

type TomlServerConfig struct {
	ServerConfig     ServerConfigInfo     `toml:"server"`
	LogConfig        LogConfigInfo        `toml:"log"`
	TomlConfigCenter TomlConfigCenterInfo `toml:"config_center"`
}

type TomlConfig struct {
	DatabaseConfigs []DatabaseConfigInfo `toml:"database,omitempty"`
	RedisConfigs    []RedisConfigInfo    `toml:"redis,omitempty"`
	RpcNetConfigs   []RpcNetConfigInfo   `toml:"rpc_server_client,omitempty"`
}

type TomlConfigManager struct {
	conf             []byte
	tomlServerConfig *TomlServerConfig
	tomlConfig       *TomlConfig
}

func NewTomlConfigManager() *TomlConfigManager {
	return &TomlConfigManager{}
}

/*
Read 函数需要传入一个 server config 文件的路径
*/
func (t *TomlConfigManager) Read(path string) error {
	t.tomlServerConfig = &TomlServerConfig{}
	t.tomlConfig = &TomlConfig{}

	if err := t.tomlServerConfig.read(path); err != nil {
		return err
	}

	if _, err := t.getConfInfo(); err != nil {
		return err
	}

	if err := t.tomlConfig.read(t.conf); err != nil {
		return err
	}

	return nil
}

/*
Unmarshal 方法需要发生在 Read 方法后, Unmarshal 方法需要传入一个用于反序列化对象的指针, 将 server 基础配置解析到 obj 对象上
*/
func (t *TomlConfigManager) Unmarshal(obj interface{}) error {
	return toml.Unmarshal(t.conf, obj)
}

func (t *TomlConfigManager) getConfInfo() ([]byte, error) {
	conf := []byte{}
	if t.tomlServerConfig.TomlConfigCenter.FilePath != "" {
		var err error
		conf, err = ioutil.ReadFile(t.tomlServerConfig.TomlConfigCenter.FilePath)
		if err != nil {
			return conf, err
		}
	}

	if t.tomlServerConfig.TomlConfigCenter.ApolloAddr != "" {
		//todo
	}

	if len(t.tomlServerConfig.TomlConfigCenter.EtcdAddrs) != 0 {
		//todo
	}

	t.conf = conf
	return conf, nil
}

/*
InitConfig 方法需要发生在 Read 方法后, InitConfig 方法会加载业务配置
*/
func (t *TomlConfigManager) InitConfig(ctx context.Context) error {
	// ---------------------- init server config ----------------------

	// init log
	if err := t.tomlServerConfig.initLogger(); err != nil {
		return err
	}

	// init web http
	if err := t.tomlServerConfig.initHttpServer(); err != nil {
		return err
	}

	// ---------------------- init base config ----------------------

	//init mysql/tidb client
	if t.tomlConfig.DatabaseConfigs != nil || len(t.tomlConfig.DatabaseConfigs) != 0 {
		if err := emysql.NewMysqlSingle(t.tomlConfig.mixMysql(), elog.Elog); err != nil {
			panic("InitConfig NewMysql error: " + err.Error())
		}
		elog.Elog.Infof("InitConfig NewMysql success")
	}

	//init redis/codis client
	if t.tomlConfig.RedisConfigs != nil || len(t.tomlConfig.RedisConfigs) != 0 {
		if err := eredis.InitRedis(ctx, t.tomlConfig.mixRedis(), nil, elog.Elog); err != nil {
			panic("InitConfig InitRedis error: " + err.Error())
		}
		elog.Elog.Infof("InitConfig InitRedis success")
	}

	// call
	// I haven't figured out how to abstract the remote call module in this area.
	if t.tomlConfig.RpcNetConfigs != nil || len(t.tomlConfig.RpcNetConfigs) != 0 {
		rpcConfigs := t.tomlConfig.mixRpcCall()

		httpConfig := []*erpc.HttpClientConfig{}
		for _, rpcConfig := range rpcConfigs {
			if rpcConfig.Proto == "http" || rpcConfig.Proto == "https" {
				httpConfig = append(httpConfig, &erpc.HttpClientConfig{
					BaseConfig: *rpcConfig,
				})
			}
		}

		if err := erpc.NewSingleHttpClientManager(httpConfig); err != nil {
			return err
		}
	}

	// mq

	// etc
	return nil
}

func (t *TomlServerConfig) read(path string) error {
	if _, err := toml.DecodeFile(path, t); err != nil {
		return err
	}
	return nil
}

func logLevelCheckAndConvert(logLevel string) (elog.LogLevel, error) {
	var elogLevel elog.LogLevel
	switch logLevel {
	case "DEBUG":
		elogLevel = elog.DEBUG
	case "INFO":
		elogLevel = elog.INFO
	case "ERROR":
		elogLevel = elog.ERROR
	case "WARNING":
		elogLevel = elog.WARNING
	default:
		return -1, ErrLogLevel
	}
	return elogLevel, nil
}

func logFormatCheckAndConvert(logFormat string) (elog.LogFormat, error) {
	var elogFormat elog.LogFormat
	switch logFormat {
	case "json":
		elogFormat = elog.FORMAT_JSON
	default:
	}
	return elogFormat, nil
}

func (t *TomlServerConfig) mixLogConfig() (*elog.LogConfig, error) {
	logLevel, err := logLevelCheckAndConvert(t.LogConfig.LogLevel)
	if err != nil {
		return nil, err
	}

	logFormat, err := logFormatCheckAndConvert(t.LogConfig.Format)
	if err != nil {
		return nil, err
	}

	if t.LogConfig.Prefix == "" {
		return nil, ErrLogPrefix
	}
	if t.LogConfig.Dir == "" {
		return nil, ErrLogDir
	}

	return &elog.LogConfig{
		LogTyp:         t.LogConfig.LogTyp,
		LogLevel:       logLevel,
		Prefix:         t.LogConfig.Prefix,
		Dir:            t.LogConfig.Dir,
		AutoClearHours: t.LogConfig.AutoClearHours,
		Depth:          t.LogConfig.Depth,
		Format:         logFormat,
	}, nil
}

func (t *TomlServerConfig) initLogger() error {
	lc, err := t.mixLogConfig()
	if err != nil {
		return err
	}

	if err := createDir(t.LogConfig.Dir); err != nil {
		return err
	}

	return elog.NewSingleLogger(lc)
}

func (t *TomlServerConfig) mixHttpserver() (*enet.DnetConfig, error) {
	if t.ServerConfig.ServiceName == "" {
		return nil, ErrServerNameNil
	}

	// todo check host
	return &enet.DnetConfig{
		L:    elog.Elog,
		Host: t.ServerConfig.HttpPort,
	}, nil
}

func (t *TomlServerConfig) initHttpServer() error {
	hsc, err := t.mixHttpserver()
	if err != nil {
		return err
	}

	if hsc.Host == "" {
		return nil
	}

	return enet.InitHttpServerSingle(hsc)
}

func (t *TomlConfig) read(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if _, err := toml.DecodeReader(bytes.NewReader(data), t); err != nil {
		return err
	}
	return nil
}

func (t *TomlConfig) mixMysql() []emysql.MConfigInfo {
	mcs := []emysql.MConfigInfo{}
	for _, m := range t.DatabaseConfigs {
		if m.Master.Ip == "" {
			panic("mixMysql splitConnect Master Ip Nil")
		}

		master := emysql.Account{
			Ip:       m.Master.Ip,
			Username: m.Master.Username,
			Password: m.Master.Password,
		}

		slaves := []emysql.Account{}
		for _, slave := range m.Slaves {
			slaves = append(slaves, emysql.Account{
				Ip:       slave.Ip,
				Username: slave.Username,
				Password: slave.Password,
			})
		}

		mc := emysql.MConfigInfo{
			Name:            m.Name,
			Master:          master,
			Slaves:          slaves,
			Database:        m.Database,
			Charset:         m.Charset,
			ParseTime:       m.ParseTime,
			Loc:             m.Loc,
			ReadTimeout:     m.ReadTimeout,
			MaxIdleConns:    m.MaxIdleConns,
			MaxOpenConns:    m.MaxOpenConns,
			ConnMaxLifetime: m.ConnMaxLifetime,
			ConnMaxIdleTime: m.ConnMaxIdleTime,
			LogMode:         m.LogModel,
		}
		mcs = append(mcs, mc)
	}
	return mcs
}

func (t *TomlConfig) mixRedis() []eredis.RedisInfo {
	redisInfos := []eredis.RedisInfo{}
	for _, r := range t.RedisConfigs {
		redisInfo := eredis.RedisInfo{
			RedisName:      r.Name,
			Addr:           r.Addr,
			MaxIdle:        r.MaxIdle,
			MaxActive:      r.MaxActive,
			IdleTimeout:    r.IdleTimeout,
			ReadTimeout:    r.ReadTimeout,
			WriteTimeout:   r.WriteTimeout,
			ConnectTimeout: r.ConnectTimeout,
			Password:       r.Password,
			Wait:           r.Wait,
		}
		redisInfos = append(redisInfos, redisInfo)
	}
	return redisInfos
}

func (t *TomlConfig) mixRpcCall() []*erpc.RpcNetConfigInfo {
	rpcNetConfigs := []*erpc.RpcNetConfigInfo{}
	for _, r := range t.RpcNetConfigs {
		rpcNetConfig := &erpc.RpcNetConfigInfo{
			ServiceName:         r.ServiceName,
			Proto:               r.Proto,
			EndpointsFrom:       r.EndpointsFrom,
			Addr:                r.Addr,
			Wight:               r.Wight,
			Balancetype:         r.Balancetype,
			DialTimeout:         r.DialTimeout,
			TimeOut:             r.TimeOut,
			RetryTimes:          r.RetryTimes,
			MaxConnsPerAddr:     r.MaxConnsPerAddr,
			MaxIdleConnsPerAddr: r.MaxIdleConnsPerAddr,
			MaxIdleConns:        r.MaxIdleConns,
			IdleConnTimeout:     r.IdleConnTimeout,
			ReadBufferSize:      r.ReadBufferSize,
			WriteBufferSize:     r.WriteBufferSize,
		}
		rpcNetConfigs = append(rpcNetConfigs, rpcNetConfig)
	}
	return rpcNetConfigs
}
