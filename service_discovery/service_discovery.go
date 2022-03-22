package service_discovery

import "golang.org/x/net/context"

type Operate int8

const (
	Add    Operate = 0
	Delete Operate = 1
	Update Operate = 2
)

type ItermInfo struct {
	key  []byte
	info []byte
	oper Operate
}

type ServiceDiscoveryManager interface {
	ServiceRegister(ctx context.Context) error
	ServiceHeartbeat(ctx context.Context) error
	ServiceDiscovery(ctx context.Context, service string) (chan ItermInfo, error)
	Close(ctx context.Context) error
}

func NewServiceDiscoveryManager(ctx context.Context, serviceDiscoveryConfig ServiceDiscoveryConfig) (ServiceDiscoveryManager, error) {
	if len(serviceDiscoveryConfig.EtcdAddr) != 0 {
		return NewEtcdServiceDiscovery(ctx, serviceDiscoveryConfig)
	}
	return nil, nil
}

func (i *ItermInfo) GetKey() []byte {
	return i.key
}

func (i *ItermInfo) GetInfo() []byte {
	return i.info
}

func (i *ItermInfo) GetOperate() Operate {
	return i.oper
}
