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

	// Service Discovery Name, compatible dirpc.
	ServiceName string `toml:"service_name"`

	// Protocol type, HTTP, HTTPS and grpc are currently supported.
	Proto string `toml:"proto"`

	// Service discovery type，eg: disf,list.
	// If we choose list，We will get the remote call service address from the "Addr" configuration.
	EndpointsFrom string `toml:"endpoints_from"`

	// Address list. See "endpointsfrom" for details.
	// The "Addr" can also be competent for the task of service discovery.
	Addr []string `toml:"addr"`

	// Load balancing type
	// eg: consistency_hash, hash, range
	Balancetype string `toml:"balancetype"`

	// rpc dial time out
	DialTimeout int `toml:"dial_timeout"`

	// rpc read time out
	ReadTimeout int `toml:"read_timeout"`

	// back off retry times
	RetryTimes int `toml:"retry_times"`

	// rpc
	MaxSize int `toml:"max_size"`

	// max free conn
	MaxIdleConn int `toml:"max_idle_conn"`

	// fuse flag
	BreakFlag bool `toml:"break_flag"`
}
