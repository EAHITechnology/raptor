package config

type ApolloConfigParser struct{}

func NewApolloConfigParser(configCenter ConfigCenterInfo) (ConfigParser, error) {
	return nil, nil
}

func (a *ApolloConfigParser) Get(key string) interface{} {
	return nil
}

func (a *ApolloConfigParser) GetDataBaseConfig(key string) (DatabaseConfigInfo, bool) {
	return DatabaseConfigInfo{}, false
}

func (a *ApolloConfigParser) GetRedisConfig(key string) (RedisConfigInfo, bool) {
	return RedisConfigInfo{}, false
}

func (a *ApolloConfigParser) GetRpcConfig(key string) (RpcNetConfigInfo, bool) {
	return RpcNetConfigInfo{}, false
}

func (a *ApolloConfigParser) GetDataBaseConfigs() []DatabaseConfigInfo {
	return nil
}

func (a *ApolloConfigParser) GetRedisConfigs() []RedisConfigInfo {
	return nil
}

func (a *ApolloConfigParser) GetRpcConfigs() []RpcNetConfigInfo {
	return nil
}

func (a *ApolloConfigParser) Unmarshal(obj interface{}) error {
	return nil
}

func (e *ApolloConfigParser) Reload() error {
	return nil
}
