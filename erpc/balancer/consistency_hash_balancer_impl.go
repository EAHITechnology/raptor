package balancer

import (
	"errors"
	"net"
	"sort"
	"sync"

	"github.com/EAHITechnology/raptor/utils"
)

const (
	MAX_HSAHNODE_NUM = 1024
)

type consistencyHashHostInfo struct {
	addr      string
	wight     int
	hashValue uint64

	//todo
	inflight      int
	clientSuccess int
	requestCount  int
	latency       float64
}

func (c *consistencyHashHostInfo) GetAddr() string {
	return c.addr
}

func (c *consistencyHashHostInfo) GetHost() (net.Conn, error) {
	return nil, nil
}

type consistencyHashHostInfoList []consistencyHashHostInfo

func (c consistencyHashHostInfoList) Len() int {
	return len(c)
}

func (c consistencyHashHostInfoList) Less(i, j int) bool {
	return c[i].hashValue < c[j].hashValue
}

func (c consistencyHashHostInfoList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type consistencyHashBalancer struct {
	conf     balancerConfig
	hashRing []consistencyHashHostInfo
	lock     sync.RWMutex
}

func NewConsistencyHashBalancer(conf balancerConfig) (Balancer, error) {
	ring := []consistencyHashHostInfo{}

	if len(conf.balancerConfigs) == 0 {
		return nil, errors.New("addr nil")
	}

	if len(conf.balancerConfigs) > MAX_HSAHNODE_NUM {
		return nil, errors.New("addr nil")
	}

	for _, val := range conf.balancerConfigs {
		hostinfo := consistencyHashHostInfo{
			addr:      val.addr,
			wight:     val.wight,
			hashValue: utils.MurmurHash64A([]byte(val.addr)),
		}

		ring = append(ring, hostinfo)
	}

	sort.Sort(consistencyHashHostInfoList(ring))

	return &consistencyHashBalancer{
		conf:     conf,
		hashRing: ring,
	}, nil
}

func (c *consistencyHashBalancer) Pick(key []byte) (HostInfo, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if len(c.hashRing) == 0 {
		return nil, errors.New("list is null")
	}

	if len(c.hashRing) == 1 {
		return &c.hashRing[0], nil
	}

	hashval := utils.MurmurHash64A(key)

	idx := sort.Search(len(c.hashRing), func(i int) bool {
		return c.hashRing[i].hashValue >= hashval
	})

	return &c.hashRing[idx], nil
}

func (c *consistencyHashBalancer) Add(conf ...balancerItem) error {
	return nil
}

func (c *consistencyHashBalancer) Remove(addr string) error {
	return nil
}
