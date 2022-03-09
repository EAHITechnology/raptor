package balancer

import (
	"errors"
	"math/rand"
	"net"
	"sync"

	"github.com/EAHITechnology/raptor/utils/rand2"
)

type randomHostInfo struct {
	addr  string
	wight int

	//todo
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
	addrMap                  map[string]int
}

func getGplist(items []balancerItem) ([]*randomHostInfo, map[string]int) {
	gplist, addrMap := []*randomHostInfo{}, make(map[string]int)

	for idx, val := range items {
		hostinfo := &randomHostInfo{
			addr:  val.addr,
			wight: val.wight,
		}

		addrMap[val.addr] = idx

		for idx := val.wight; idx > 0; idx-- {
			gplist = append(gplist, hostinfo)
		}
	}
	return gplist, addrMap
}

func NewRandomBalancer(conf balancerConfig) (Balancer, error) {

	if len(conf.balancerConfigs) == 0 {
		return nil, errors.New("addr nil")
	}

	gplist, addrMap := getGplist(conf.balancerConfigs)

	return &randomBalancer{
		conf:                     conf,
		geometricProbabilityList: gplist,
		rand:                     rand2.New(rand.NewSource(1)),
		addrMap:                  addrMap,
	}, nil
}

func (r *randomBalancer) Pick(key []byte) (HostInfo, error) {
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

func (r *randomBalancer) Add(conf ...balancerItem) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, val := range conf {
		if i, ok := r.addrMap[val.addr]; ok {
			r.conf.balancerConfigs[i].wight = val.wight
			continue
		}

		r.conf.balancerConfigs = append(r.conf.balancerConfigs, balancerItem{
			addr:  val.addr,
			wight: val.wight,
		})
	}

	r.geometricProbabilityList, r.addrMap = getGplist(r.conf.balancerConfigs)
	return nil
}

func (r *randomBalancer) Remove(addr string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.addrMap[addr]; !ok {
		return errors.New("addr none")
	}

	tmpFront := r.conf.balancerConfigs[:r.addrMap[addr]]
	tmpBackend := []balancerItem{}
	if r.addrMap[addr] < len(r.conf.balancerConfigs)-1 {
		tmpBackend = r.conf.balancerConfigs[r.addrMap[addr]+1:]
	}
	r.conf.balancerConfigs = []balancerItem{}
	r.conf.balancerConfigs = append(r.conf.balancerConfigs, tmpFront...)
	r.conf.balancerConfigs = append(r.conf.balancerConfigs, tmpBackend...)

	r.geometricProbabilityList, r.addrMap = getGplist(r.conf.balancerConfigs)
	return nil
}
