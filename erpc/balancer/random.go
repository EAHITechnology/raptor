package balancer

import (
	"errors"
	"math/rand"
	"net"
	"sync"

	"github.com/EAHITechnology/raptor/utils/rand2"
)

type randomHostInfo struct {
	addr          string
	wight         int
	inflight      int
	clientSuccess int
	requestCount  int
	latency       float64
}

func (r *randomHostInfo) GetAddr() string {
	return r.addr
}

func (r *randomHostInfo) GetHost() (net.Conn, error) {
	return nil, nil
}

type randomBalancer struct {
	conf                     balancerConfig
	geometricProbabilityList []*randomHostInfo
	lock                     sync.RWMutex
	rand                     *rand2.Rand
}

func NewRandomBalancer(conf balancerConfig) (Balancer, error) {
	gplist := []*randomHostInfo{}

	if len(conf.balancerConfigs) == 0 {
		return nil, errors.New("addr nil")
	}

	for _, val := range conf.balancerConfigs {
		hostinfo := &randomHostInfo{
			addr:  val.addr,
			wight: val.wight,
		}

		for idx := val.wight; idx > 0; idx-- {
			gplist = append(gplist, hostinfo)
		}
	}

	return &randomBalancer{
		conf:                     conf,
		geometricProbabilityList: gplist,
		rand:                     rand2.New(rand.NewSource(1)),
	}, nil
}

func (r *randomBalancer) Pick() (HostInfo, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if len(r.geometricProbabilityList) == 0 {
		return nil, errors.New("list is null")
	}

	if len(r.geometricProbabilityList) == 1 {
		return r.geometricProbabilityList[0], nil
	}

	return r.geometricProbabilityList[rand.Int()%len(r.geometricProbabilityList)], nil
}
