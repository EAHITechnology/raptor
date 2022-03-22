package config

type DatabaseAccount struct {
	Ip       string `mapstructure:"ip"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type DatabaseConfigInfo struct {
	Name            string            `mapstructure:"name"`
	Master          DatabaseAccount   `mapstructure:"master"`
	Slaves          []DatabaseAccount `mapstructure:"slaves"`
	Database        string            `mapstructure:"database"`
	Charset         string            `mapstructure:"charset"`
	ParseTime       string            `mapstructure:"parseTime"`
	Loc             string            `mapstructure:"loc"`
	ReadTimeout     string            `mapstructure:"readTimeout"`
	MaxIdleConns    int               `mapstructure:"maxIdleConns"`
	MaxOpenConns    int               `mapstructure:"maxOpenConns"`
	ConnMaxIdleTime int               `mapstructure:"connMaxIdletime"`
	ConnMaxLifetime int               `mapstructure:"connMaxLifetime"`
	DiscoverFlag    bool              `mapstructure:"discover_flag"`
	LogModel        bool              `mapstructure:"log_model"`
}

type RedisConfigInfo struct {
	Name           string `mapstructure:"name"`
	Addr           string `mapstructure:"addr"`
	Password       string `mapstructure:"password"`
	MaxIdle        int    `mapstructure:"max_idle"`
	IdleTimeout    int64  `mapstructure:"max_idletimeout"`
	MaxActive      int    `mapstructure:"max_active"`
	ReadTimeout    int64  `mapstructure:"read_timeout"`
	WriteTimeout   int64  `mapstructure:"write_timeout"`
	SlowTime       int64  `mapstructure:"slow_time"`
	ConnectTimeout int64  `mapstructure:"connect_time"`
	Wait           bool   `mapstructure:"wait"`
	Database       int    `mapstructure:"databases"`
}

// remote call
type RpcNetConfigInfo struct {

	// Service Discovery Name.
	ServiceName string `mapstructure:"service_name"`

	// Protocol type, HTTP, HTTPS and grpc are currently supported.
	Proto string `mapstructure:"proto"`

	// Service discovery type，eg: etcd, zk, apollo, list.
	// If we choose list，We will get the remote call service address from the "Addr" configuration.
	EndpointsFrom string `mapstructure:"endpoints_from"`

	// Address list. See "endpointsfrom" for details.
	// The "Addr" can also be competent for the task of service discovery.
	Addr []string `mapstructure:"addr"`

	// Weights
	Wight []int `mapstructure:"wight"`

	// Load balancing type.
	// eg: consistency_hash, p2c, random, range.
	Balancetype string `mapstructure:"balancetype"`

	// rpc dial time out.(Millisecond default 0.)
	DialTimeout int `mapstructure:"dial_timeout"`

	// rpc total time out.(Millisecond default 0.)
	TimeOut int `mapstructure:"timeout"`

	// back off retry times.
	RetryTimes int `mapstructure:"retry_times"`

	// every addr max conns num.
	MaxConnsPerAddr int `mapstructure:"max_conns_per_addr"`

	// every addr max idle conns num.
	MaxIdleConnsPerAddr int `mapstructure:"max_idleconns_per_addr"`

	// all addr max idle conns num.
	MaxIdleConns int `mapstructure:"max_idleconns"`

	// idle timeout. (second default 10.)
	IdleConnTimeout int `mapstructure:"idleconn_timeout"`

	// ReadBufferSize (bytes).
	ReadBufferSize int `mapstructure:"readbuffer_size"`

	// WriteBufferSize (bytes).
	WriteBufferSize int `mapstructure:"writebuffer_size"`
}

type ServiceDiscovery struct {
	EtcdAddr               []string `mapstructure:"etcd_addr"`
	ZkAddr                 []string `mapstructure:"zk_addr"`
	CustomServiceDiscovery []string `mapstructure:"custom_service_discovery"`
}

type Config struct {
	ServiceDiscovery ServiceDiscovery     `mapstructure:"service_discovery"`
	DatabaseConfigs  []DatabaseConfigInfo `mapstructure:"database"`
	RedisConfigs     []RedisConfigInfo    `mapstructure:"redis"`
	RpcNetConfigs    []RpcNetConfigInfo   `mapstructure:"rpc_server_client"`
}
