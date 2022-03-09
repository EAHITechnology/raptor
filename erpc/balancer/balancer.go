package balancer

import (
	"errors"
	"net"
)

type BalancerTyp string

const (
	RandomType          BalancerTyp = "random"
	P2cType             BalancerTyp = "p2c"
	ConsistencyHashType BalancerTyp = "consistency_hash"
	RangeType           BalancerTyp = "range"
)

type balancerItem struct {
	addr  string
	wight int
}

type balancerConfig struct {
	balancerTyp     BalancerTyp
	balancerConfigs []balancerItem
}

type HostInfo interface {
	GetAddr() string
	GetHost() (net.Conn, error)
}

type Balancer interface {
	// Pick well get a HostInfo interface.
	// Then HostInfo contains link information.
	Pick() (HostInfo, error)
}

/*
the wight should be greater than or equal to 1.
if wight == 0 , load balancing behavior does not take effect.
*/
func NewBalancerItem(addr string, wight int) balancerItem {
	return balancerItem{
		addr:  addr,
		wight: wight,
	}
}

func NewBalancerConfig() *balancerConfig {
	return &balancerConfig{}
}

func (b *balancerConfig) SetItem(conf ...balancerItem) {
	b.balancerConfigs = append(b.balancerConfigs, conf...)
}

func (b *balancerConfig) SetBalancerTyp(typ BalancerTyp) {
	b.balancerTyp = typ
}

func NewBalancer(conf balancerConfig) (Balancer, error) {
	switch conf.balancerTyp {
	case RandomType:
		return NewRandomBalancer(conf)
	case P2cType:
		return nil, nil
	case ConsistencyHashType:
		return nil, nil
	case RangeType:
		return nil, nil
	default:
		return nil, errors.New("illegal type")
	}
}
