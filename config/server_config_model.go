package config

type ServerConfigInfo struct {
	ServiceName string `mapstructure:"service_name"`
	HttpPort    string `mapstructure:"http_port"`
	RpcPort     string `mapstructure:"rpc_port"`
}

type LogConfigInfo struct {
	LogTyp         string `mapstructure:"log_type"`
	Dir            string `mapstructure:"dir"`
	LogLevel       string `mapstructure:"log_level"`
	Prefix         string `mapstructure:"prefix"`
	AutoClearHours int    `mapstructure:"auto_clear_hours"`
	Depth          int    `mapstructure:"depth"`
	Format         string `mapstructure:"format"`
}

type ConfigCenterInfo struct {
	FileType   string   `mapstructure:"file_type"`
	FilePath   string   `mapstructure:"file_path"`
	EtcdAddrs  []string `mapstructure:"etcd_addrs"`
	ApolloAddr string   `mapstructure:"apollo_addr"`
}

type ServerConfig struct {
	ServerConfig ServerConfigInfo `mapstructure:"server"`
	LogConfig    LogConfigInfo    `mapstructure:"log"`
	ConfigCenter ConfigCenterInfo `mapstructure:"config_center"`
}
