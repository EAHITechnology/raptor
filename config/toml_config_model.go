package config

type ServerConfigInfo struct {
	ServiceName string `toml:"service_name"`
	HttpPort    string `toml:"http_port"`
	RpcPort     string `toml:"rpc_port"`
}

type LogConfigInfo struct {
	LogTyp         string `toml:"log_type"`
	Dir            string `toml:"dir"`
	LogLevel       string `toml:"log_level"`
	Prefix         string `toml:"prefix"`
	AutoClearHours int    `toml:"auto_clear_hours"`
	Depth          int    `toml:"depth"`
	Format         string `toml:"format"`
}

type TomlConfigCenterInfo struct {
	FilePath   string   `toml:"file_path,omitempty"`
	EtcdAddrs  []string `toml:"etcd_addrs,omitempty"`
	ApolloAddr string   `toml:"apollo_addr,omitempty"`
}

type DatabaseAccount struct {
	Ip       string `toml:"ip"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type DatabaseConfigInfo struct {
	Name            string            `toml:"name"`
	Master          DatabaseAccount   `toml:"master"`
	Slaves          []DatabaseAccount `toml:"slaves"`
	Database        string            `toml:"database"`
	Charset         string            `toml:"charset"`
	ParseTime       string            `toml:"parseTime"`
	Loc             string            `toml:"loc"`
	ReadTimeout     string            `toml:"readTimeout"`
	MaxIdleConns    int               `toml:"maxIdleConns"`
	MaxOpenConns    int               `toml:"maxOpenConns"`
	ConnMaxIdleTime int               `toml:"connMaxIdletime"`
	ConnMaxLifetime int               `toml:"connMaxLifetime"`
	DiscoverFlag    bool              `toml:"discover_flag"`
	LogModel        bool              `toml:"log_model"`
}

type RedisConfigInfo struct {
	Name           string `toml:"name"`
	Addr           string `toml:"addr"`
	Password       string `toml:"password"`
	MaxIdle        int    `toml:"max_idle"`
	IdleTimeout    int64  `toml:"max_idletimeout"`
	MaxActive      int    `toml:"max_active"`
	ReadTimeout    int64  `toml:"read_timeout"`
	WriteTimeout   int64  `toml:"write_timeout"`
	SlowTime       int64  `toml:"slow_time"`
	ConnectTimeout int64  `toml:"connect_time"`
	Wait           bool   `toml:"wait"`
	Database       int    `toml:"databases"`
}

// remote call
type RpcNetConfigInfo struct {

	// Service Discovery Name.
	ServiceName string `toml:"service_name"`

	// Protocol type, HTTP, HTTPS and grpc are currently supported.
	Proto string `toml:"proto"`

	// Service discovery type，eg: etcd, zk, apollo, list.
	// If we choose list，We will get the remote call service address from the "Addr" configuration.
	EndpointsFrom string `toml:"endpoints_from"`

	// Address list. See "endpointsfrom" for details.
	// The "Addr" can also be competent for the task of service discovery.
	Addr []string `toml:"addr"`

	// Weights
	Wight []int `toml:"wight"`

	// Load balancing type.
	// eg: consistency_hash, p2c, random, range.
	Balancetype string `toml:"balancetype"`

	// rpc dial time out
	DialTimeout int `toml:"dial_timeout"`

	// rpc total time out
	TimeOut int `toml:"timeout"`

	// back off retry times
	RetryTimes int `toml:"retry_times"`

	// every addr max conns num
	MaxConnsPerAddr int `toml:"max_conns_per_addr"`

	// every addr max idle conns num
	MaxIdleConnsPerAddr int `toml:"max_idleconns_per_addr"`

	// all addr max idle conns num
	MaxIdleConns int `toml:"max_idleconns"`

	// idle timeout
	IdleConnTimeout int `toml:"idleconn_timeout"`

	// ReadBufferSize
	ReadBufferSize int `toml:"readbuffer_size"`

	// WriteBufferSize
	WriteBufferSize int `toml:"writebuffer_size"`
}
