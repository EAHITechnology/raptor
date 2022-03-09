package erpc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/EAHITechnology/raptor/erpc/balancer"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	HEAD   HttpMethod = "HEAD"
)

func (h *HttpMethod) getMethod() string {
	return string(*h)
}

type HttpClientConfig struct {
	TimeOut             time.Duration
	IdleConnTimeout     time.Duration
	KeepAlive           time.Duration
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	baseConfig          RpcNetConfigInfo
}

type HttpClient struct {
	conf   *HttpClientConfig
	client *http.Client
	lock   sync.RWMutex
	b      balancer.Balancer
}

type HttpClientManager struct {
	manager map[string]*HttpClient
	lock    sync.RWMutex
}

func getBalancerTyp(typ string) (balancer.BalancerTyp, error) {
	var balancetype balancer.BalancerTyp
	switch typ {
	case "", "random":
		balancetype = balancer.RandomType
	case "p2c":
		balancetype = balancer.P2cType
	case "consistency_hash":
		balancetype = balancer.ConsistencyHashType
	case "range":
		balancetype = balancer.RangeType
	default:
		return "", errors.New("illegal type")
	}
	return balancetype, nil
}

func NewHttpClient(conf *HttpClientConfig) (*HttpClient, error) {
	h := &HttpClient{}
	h.conf = conf

	// http client
	h.client = &http.Client{
		Timeout: conf.TimeOut,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: time.Duration(conf.baseConfig.DialTimeout),
			}).DialContext,
			MaxIdleConns:        conf.baseConfig.MaxIdleConns,
			MaxConnsPerHost:     conf.MaxConnsPerHost,
			MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
			IdleConnTimeout:     conf.IdleConnTimeout,
			TLSHandshakeTimeout: 10 * time.Second,
			ReadBufferSize:      conf.baseConfig.ReadBufferSize,
			WriteBufferSize:     conf.baseConfig.WriteBufferSize,
		},
	}

	// balancer
	balancetype, err := getBalancerTyp(conf.baseConfig.Balancetype)
	if err != nil {
		return nil, err
	}

	balancerConfig := balancer.NewBalancerConfig()
	balancerConfig.SetBalancerTyp(balancetype)
	for _, addr := range conf.baseConfig.Addr {
		balancerConfig.SetItem(balancer.NewBalancerItem(addr, 1))
	}

	balancer, err := balancer.NewBalancer((*balancerConfig))
	if err != nil {
		return nil, err
	}

	h.client.CloseIdleConnections()
	h.b = balancer
	return h, nil
}

func (h *HttpClient) Send(method HttpMethod, query string, header map[string]string, body io.Reader) (interface{}, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	hostInfo, err := h.b.Pick()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method.getMethod(), hostInfo.GetAddr(), body)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		request.Header.Add(k, v)
	}

	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll err:%v", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http resp status : %d,  msg: %s", response.StatusCode, string(respB))
	}

	return respB, nil
}

func (h *HttpClient) Close() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.client.CloseIdleConnections()
}

func NewHttpClientManager() *HttpClientManager {
	return &HttpClientManager{
		manager: make(map[string]*HttpClient),
	}
}

func (hm *HttpClientManager) NewHttpClient(conf *HttpClientConfig) error {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	if _, ok := hm.manager[conf.baseConfig.ServiceName]; ok {
		return errors.New("service already exists")
	}

	client, err := NewHttpClient(conf)
	if err != nil {
		return err
	}

	hm.manager[conf.baseConfig.ServiceName] = client

	return nil
}

func (hm *HttpClientManager) GetClient(serviceName string) (*HttpClient, error) {
	hm.lock.RLock()
	defer hm.lock.RUnlock()

	client, ok := hm.manager[serviceName]
	if !ok {
		return nil, errors.New("service not exists")
	}

	return client, nil
}
