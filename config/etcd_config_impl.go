package config

type EtcdConfigParser struct{}

func NewEtcdConfigParser(configCenter ConfigCenterInfo) (ConfigParser, error) {
	return nil, nil
}

func (e *EtcdConfigParser) Get(key string) interface{} {
	return nil
}

func (e *EtcdConfigParser) GetDataBaseConfig(key string) (DatabaseConfigInfo, bool) {
	return DatabaseConfigInfo{}, false
}

func (e *EtcdConfigParser) GetRedisConfig(key string) (RedisConfigInfo, bool) {
	return RedisConfigInfo{}, false
}

func (e *EtcdConfigParser) GetRpcConfig(key string) (RpcNetConfigInfo, bool) {
	return RpcNetConfigInfo{}, false
}

func (e *EtcdConfigParser) GetDataBaseConfigs() []DatabaseConfigInfo {
	return nil
}

func (e *EtcdConfigParser) GetRedisConfigs() []RedisConfigInfo {
	return nil
}

func (e *EtcdConfigParser) GetRpcConfigs() []RpcNetConfigInfo {
	return nil
}

func (e *EtcdConfigParser) Unmarshal(obj interface{}) error {
	return nil
}

func (e *EtcdConfigParser) Reload() error {
	return nil
}
