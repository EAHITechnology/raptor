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

	CONSISTENCY_HASH_WAIT_SORT = -1
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
	addrMap  map[string]int
}

func getHashRing(items []balancerItem) ([]consistencyHashHostInfo, map[string]int) {
	ring, addrMap := []consistencyHashHostInfo{}, make(map[string]int)

	for idx, item := range items {
		hostinfo := consistencyHashHostInfo{
			addr:      item.addr,
			wight:     item.wight,
			hashValue: utils.MurmurHash64A([]byte(item.addr)),
		}
		ring = append(ring, hostinfo)
		addrMap[item.addr] = idx
	}
	sort.Sort(consistencyHashHostInfoList(ring))
	return ring, addrMap
}

func NewConsistencyHashBalancer(conf balancerConfig) (Balancer, error) {
	if len(conf.balancerConfigs) == 0 {
		return nil, errors.New("addr nil")
	}

	if len(conf.balancerConfigs) > MAX_HSAHNODE_NUM {
		return nil, errors.New("addr nil")
	}

	ring, addrMap := getHashRing(conf.balancerConfigs)

	return &consistencyHashBalancer{
		conf:     conf,
		hashRing: ring,
		addrMap:  addrMap,
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
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, val := range conf {
		if _, ok := c.addrMap[val.addr]; ok {
			continue
		}

		c.conf.balancerConfigs = append(c.conf.balancerConfigs, val)
		c.addrMap[val.addr] = CONSISTENCY_HASH_WAIT_SORT
	}

	c.hashRing, c.addrMap = getHashRing(c.conf.balancerConfigs)
	return nil
}

func (c *consistencyHashBalancer) Remove(addr string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.addrMap[addr]
	if !ok {
		return errors.New("addr none")
	}

	tmpFront := c.conf.balancerConfigs[:c.addrMap[addr]]
	tmpBackend := []balancerItem{}
	if c.addrMap[addr] < len(c.conf.balancerConfigs)-1 {
		tmpBackend = c.conf.balancerConfigs[c.addrMap[addr]+1:]
	}
	c.conf.balancerConfigs = []balancerItem{}
	c.conf.balancerConfigs = append(c.conf.balancerConfigs, tmpFront...)
	c.conf.balancerConfigs = append(c.conf.balancerConfigs, tmpBackend...)

	c.hashRing, c.addrMap = getHashRing(c.conf.balancerConfigs)
	return nil
}
