package service_discovery

type ServiceDiscoveryConfig struct {
	EtcdAddr               []string
	ZkAddr                 []string
	CustomServiceDiscovery []string
	LocalServiceName       string
	Ip                     string
	Host                   string
}
