package balancer

import (
	"errors"
	"net"
	"sync"
)

type consistencyHashHostInfo struct {
	addr          string
	wight         int
	inflight      int
	clientSuccess int
	requestCount  int
	latency       float64
}

func (r *consistencyHashHostInfo) GetAddr() string {
	return r.addr
}

func (r *consistencyHashHostInfo) GetHost() (net.Conn, error) {
	return nil, nil
}

type consistencyHashBalancer struct {
	conf                     balancerConfig
	geometricProbabilityList []*consistencyHashHostInfo
	lock                     sync.RWMutex
}

func NewConsistencyHashBalancer(conf balancerConfig) (Balancer, error) {
	gplist := []*consistencyHashHostInfo{}

	if len(conf.balancerConfigs) == 0 {
		return nil, errors.New("addr nil")
	}

	for _, val := range conf.balancerConfigs {
		hostinfo := &consistencyHashHostInfo{
			addr:  val.addr,
			wight: val.wight,
		}

		for idx := val.wight; idx > 0; idx-- {
			gplist = append(gplist, hostinfo)
		}
	}

	return &consistencyHashBalancer{
		conf:                     conf,
		geometricProbabilityList: gplist,
	}, nil
}

func (r *consistencyHashBalancer) Pick() (HostInfo, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if len(r.geometricProbabilityList) == 0 {
		return nil, errors.New("list is null")
	}

	if len(r.geometricProbabilityList) == 1 {
		return r.geometricProbabilityList[0], nil
	}

	return nil, nil
}
