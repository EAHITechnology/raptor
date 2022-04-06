package service_discovery

import (
	"errors"

	"golang.org/x/net/context"
)

type Operate int8

const (
	Add    Operate = 0
	Delete Operate = 1
	Update Operate = 2
)

var (
	ErrLogNil = errors.New("")
)

type ItermInfo struct {
	key  []byte
	info []byte
	oper Operate
}

type ServiceDiscoveryLog interface {
	Debugf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Warnf(f string, args ...interface{})
	Errorf(f string, args ...interface{})
}

type ServiceDiscoveryManager interface {
	// Service registe
	ServiceRegister(ctx context.Context) error

	// The ServiceHeartbeat method maintains the
	// heartbeat behavior of the service and registry
	ServiceHeartbeat(ctx context.Context) error

	// If you choose the service discovery method,
	// you must consume the chan include `ItermInfo` it returns.
	ServiceDiscovery(ctx context.Context, service string) (chan ItermInfo, error)

	// Close ServiceDiscoveryManager
	Close(ctx context.Context) error
}

func NewServiceDiscoveryManager(ctx context.Context, serviceDiscoveryConfig ServiceDiscoveryConfig, log ServiceDiscoveryLog) (ServiceDiscoveryManager, error) {
	if log == nil {
		return nil, ErrLogNil
	}

	if len(serviceDiscoveryConfig.EtcdAddr) != 0 {
		return NewEtcdServiceDiscovery(ctx, serviceDiscoveryConfig, log)
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
