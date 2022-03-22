package server

import (
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/EAHITechnology/raptor/balancer"
	"github.com/EAHITechnology/raptor/config"
	"github.com/EAHITechnology/raptor/elog"
	"github.com/EAHITechnology/raptor/emysql"
	"github.com/EAHITechnology/raptor/enet"
	"github.com/EAHITechnology/raptor/eredis"
	"github.com/EAHITechnology/raptor/erpc"
	"github.com/EAHITechnology/raptor/utils"
	"golang.org/x/net/context"
)

type DefaultServer struct {
	Typ ServerTyp

	serverConfigParser config.ServerConfigParser
	configParser       config.ConfigParser

	afterInitFunc  atferFuncObj
	beforeInitFunc beforeFuncObj
}

func NewDefaultServer(ctx context.Context, typ ServerTyp, configPath string) (*DefaultServer, error) {
	fileArray := strings.Split(configPath, ".")
	conf := config.ConfigInfo{
		Typ:  fileArray[len(fileArray)-1],
		Path: configPath,
	}

	serverConfigParser, err := config.NewServerConfigParser(conf)
	if err != nil {
		return nil, err
	}

	configParser, err := config.NewConfigParser(serverConfigParser.GetConfigCenterInfo())
	if err != nil {
		return nil, err
	}

	server := &DefaultServer{
		Typ:                typ,
		serverConfigParser: serverConfigParser,
		configParser:       configParser,
	}

	return server, nil
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
		elogFormat = elog.FORMAT_NORMAL
	}
	return elogFormat, nil
}

func (d *DefaultServer) mixLogConfig() (*elog.LogConfig, error) {
	logConfigInfo := d.serverConfigParser.GetLogConfigInfo()
	logLevel, err := logLevelCheckAndConvert(logConfigInfo.LogLevel)
	if err != nil {
		return nil, err
	}

	logFormat, err := logFormatCheckAndConvert(logConfigInfo.Format)
	if err != nil {
		return nil, err
	}

	if logConfigInfo.Prefix == "" {
		return nil, ErrLogPrefix
	}

	if logConfigInfo.Dir == "" {
		return nil, ErrLogDir
	}

	return &elog.LogConfig{
		LogTyp:         logConfigInfo.LogTyp,
		LogLevel:       logLevel,
		Prefix:         logConfigInfo.Prefix,
		Dir:            logConfigInfo.Dir,
		AutoClearHours: logConfigInfo.AutoClearHours,
		Depth:          logConfigInfo.Depth,
		Format:         logFormat,
	}, nil
}

func (d *DefaultServer) mixHttpserver() (*enet.EnetConfig, error) {
	serverConf := d.serverConfigParser.GetServerConfigInfo()
	if serverConf.ServiceName == "" {
		return nil, ErrServerNameNil
	}

	return &enet.EnetConfig{
		L:    elog.Elog,
		Host: serverConf.HttpPort,
	}, nil
}

func (d *DefaultServer) mixDB() ([]emysql.MConfigInfo, error) {
	mcs := []emysql.MConfigInfo{}
	databaseConfigs := d.configParser.GetDataBaseConfigs()

	for _, databaseConfig := range databaseConfigs {
		if databaseConfig.Master.Ip == "" {
			return nil, ErrMysqlIpNil
		}

		if databaseConfig.Name == "" {
			return nil, ErrMysqlNameNil
		}

		master := emysql.Account{
			Ip:       databaseConfig.Master.Ip,
			Username: databaseConfig.Master.Username,
			Password: databaseConfig.Master.Password,
		}

		slaves := []emysql.Account{}
		for _, slave := range databaseConfig.Slaves {
			slaves = append(slaves, emysql.Account{
				Ip:       slave.Ip,
				Username: slave.Username,
				Password: slave.Password,
			})
		}

		mc := emysql.MConfigInfo{
			Name:            databaseConfig.Name,
			Master:          master,
			Slaves:          slaves,
			Database:        databaseConfig.Database,
			Charset:         databaseConfig.Charset,
			ParseTime:       databaseConfig.ParseTime,
			Loc:             databaseConfig.Loc,
			ReadTimeout:     databaseConfig.ReadTimeout,
			MaxIdleConns:    databaseConfig.MaxIdleConns,
			MaxOpenConns:    databaseConfig.MaxOpenConns,
			ConnMaxLifetime: databaseConfig.ConnMaxLifetime,
			ConnMaxIdleTime: databaseConfig.ConnMaxIdleTime,
			LogMode:         databaseConfig.LogModel,
		}
		mcs = append(mcs, mc)
	}
	return mcs, nil
}

func (d *DefaultServer) mixRedis() ([]eredis.RedisInfo, error) {
	redisConfigs := d.configParser.GetRedisConfigs()
	redisInfos := []eredis.RedisInfo{}

	for _, r := range redisConfigs {
		if r.Name == "" {
			return nil, ErrRedisNameNil
		}
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
	return redisInfos, nil
}

func (d *DefaultServer) mixRpcCall() ([]*erpc.RpcNetConfigInfo, error) {
	rpcNetConfigInfo := d.configParser.GetRpcConfigs()
	rpcNetConfigs := []*erpc.RpcNetConfigInfo{}

	for _, r := range rpcNetConfigInfo {
		if r.ServiceName == "" {
			return nil, ErrServerNameNil
		}

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
	return rpcNetConfigs, nil
}

func (d *DefaultServer) initLogger() error {
	logConf, err := d.mixLogConfig()
	if err != nil {
		return err
	}

	if err := utils.CreateDir(logConf.Dir); err != nil {
		return err
	}

	return elog.NewSingleLogger(logConf)
}

func (d *DefaultServer) initHttpServer() error {
	if d.serverConfigParser.GetServerConfigInfo().HttpPort == "" {
		return nil
	}

	serverConf, err := d.mixHttpserver()
	if err != nil {
		return err
	}

	return enet.InitHttpServerSingle(serverConf)
}

func (d *DefaultServer) initDB() error {
	dbConfig, err := d.mixDB()
	if err != nil {
		return err
	}
	return emysql.NewMysqlSingle(dbConfig, elog.Elog)
}

func (d *DefaultServer) initRedis(ctx context.Context) error {
	redisConfig, err := d.mixRedis()
	if err != nil {
		return err
	}
	return eredis.InitRedis(ctx, redisConfig, nil, elog.Elog)
}

func (d *DefaultServer) initRpc(ctx context.Context) error {
	rpcConfigs, err := d.mixRpcCall()
	if err != nil {
		return err
	}

	httpManagerConfiges := []erpc.HttpManagerConfig{}
	for _, rpcConfig := range rpcConfigs {
		// service_discovery
		// sdc := service_discovery.ServiceDiscoveryConfig{}
		// if rpcConfig.EndpointsFrom != "file" && rpcConfig.EndpointsFrom != "" {
		// serviceDiscoveryConfig, ok := d.configParser.GetServiceDiscoveryConfig()
		// if !ok {
		// return errors.New("")
		// }
		// switch rpcConfig.EndpointsFrom {
		// case "zk":
		// sdc.ZkAddr = serviceDiscoveryConfig.ZkAddr
		// case "etcd":
		// sdc.EtcdAddr = serviceDiscoveryConfig.EtcdAddr
		// case "custom":
		// sdc.CustomServiceDiscovery = serviceDiscoveryConfig.CustomServiceDiscovery
		// default:
		// }

		// serviceDiscoveryManager, err := service_discovery.NewServiceDiscoveryManager(ctx, sdc)

		// rch, err := serviceDiscoveryManager.ServiceDiscovery(ctx, rpcConfig.ServiceName)
		// }

		// balancer
		balancetype, err := getBalancerTyp(rpcConfig.Balancetype)
		if err != nil {
			return err
		}
		balancerConfig := balancer.NewBalancerConfig()
		balancerConfig.SetBalancerTyp(balancetype)

		for idx, addr := range rpcConfig.Addr {
			balancerConfig.SetItem(balancer.NewBalancerItem(addr, rpcConfig.Wight[idx]))
		}

		balancer, err := balancer.NewBalancer((*balancerConfig))
		if err != nil {
			return err
		}

		if rpcConfig.Proto == "http" || rpcConfig.Proto == "https" {

			httpManagerConfige := erpc.HttpManagerConfig{
				Httpconf: &erpc.HttpClientConfig{
					BaseConfig: *rpcConfig,
				},
				Balancer: balancer,
			}
			httpManagerConfiges = append(httpManagerConfiges, httpManagerConfige)
		}
	}

	if err := erpc.NewSingleHttpClientManager(httpManagerConfiges); err != nil {
		return err
	}

	// etc...
	return nil
}

func (d *DefaultServer) initMq() error {
	return nil
}

func (d *DefaultServer) AfterInit(serverFunc ...ServerFunc) {
	if atomic.CompareAndSwapInt32(&d.afterInitFunc.flag, defaultAbleToWrite, defaultWriting) {
		d.afterInitFunc.serverFuncList = append(d.afterInitFunc.serverFuncList, serverFunc...)
	}
}

func (d *DefaultServer) BeforeInit(serverFunc ...ServerFunc) {
	if atomic.CompareAndSwapInt32(&d.beforeInitFunc.flag, defaultAbleToWrite, defaultWriting) {
		d.beforeInitFunc.serverFuncList = append(d.beforeInitFunc.serverFuncList, serverFunc...)
	}
}

func (d *DefaultServer) execAfter(ctx context.Context) {
	if atomic.CompareAndSwapInt32(&d.afterInitFunc.flag, defaultWriting, defaultEndWrite) {
		for _, serverFunc := range d.afterInitFunc.serverFuncList {
			go serverFunc(ctx)
		}
	}
}

func (d *DefaultServer) execBefore(ctx context.Context) {
	if atomic.CompareAndSwapInt32(&d.beforeInitFunc.flag, defaultWriting, defaultEndWrite) {
		for _, serverFunc := range d.beforeInitFunc.serverFuncList {
			go serverFunc(ctx)
		}
	}
}

func (d *DefaultServer) Run(ctx context.Context, cancel context.CancelFunc) error {

	// ---------------- init meta config ----------------
	if err := d.initLogger(); err != nil {
		return err
	}

	if err := d.initHttpServer(); err != nil {
		return err
	}

	d.execBefore(ctx)

	// ---------------- init config ----------------

	if err := d.initDB(); err != nil {
		return err
	}

	if err := d.initRedis(ctx); err != nil {
		return err
	}

	if err := d.initRpc(ctx); err != nil {
		return err
	}

	d.execAfter(ctx)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE, syscall.SIGUSR1)
	for {
		sig := <-sc
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			elog.Elog.Warnf("get os signal, close server signal: %s", sig.String())
			if enet.HttpWeb != nil {
				enet.HttpWeb.Close()
			}
			cancel()

			return nil
		case syscall.SIGUSR1:
			elog.Elog.Infof("ignore os signal: %s , start reload config", sig.String())
			if err := d.configParser.Reload(); err != nil {
				elog.Elog.Errorf("server reload config error:%v", err)
			}
		default:
			elog.Elog.Warnf("ignore os signal: %s", sig.String())
		}
	}
}
