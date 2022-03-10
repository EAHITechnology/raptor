package erpc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/EAHITechnology/raptor/erpc/balancer"
	"github.com/EAHITechnology/raptor/utils"
	"golang.org/x/net/context"
)

var HttpManager *HttpClientManager

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	HEAD   HttpMethod = "HEAD"
)

const (
	DefaultHttpIdleTimeout = 10
)

func (h *HttpMethod) getMethod() string {
	return string(*h)
}

type HttpClientConfig struct {
	// fuse flag todo
	// breakFlag bool
	BaseConfig RpcNetConfigInfo
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

	if conf.BaseConfig.IdleConnTimeout < DefaultHttpIdleTimeout {
		conf.BaseConfig.IdleConnTimeout = DefaultHttpIdleTimeout
	}

	h.conf = conf

	// http client
	h.client = &http.Client{
		Timeout: time.Millisecond * time.Duration(conf.BaseConfig.TimeOut),
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: time.Millisecond * time.Duration(conf.BaseConfig.DialTimeout),
			}).DialContext,
			MaxIdleConns:        conf.BaseConfig.MaxIdleConns,
			MaxConnsPerHost:     conf.BaseConfig.MaxConnsPerAddr,
			MaxIdleConnsPerHost: conf.BaseConfig.MaxIdleConnsPerAddr,
			IdleConnTimeout:     time.Second * time.Duration(conf.BaseConfig.IdleConnTimeout),
			TLSHandshakeTimeout: 5 * time.Second,
			ReadBufferSize:      conf.BaseConfig.ReadBufferSize,
			WriteBufferSize:     conf.BaseConfig.WriteBufferSize,
		},
	}

	// balancer
	balancetype, err := getBalancerTyp(conf.BaseConfig.Balancetype)
	if err != nil {
		return nil, err
	}

	balancerConfig := balancer.NewBalancerConfig()
	balancerConfig.SetBalancerTyp(balancetype)
	for idx, addr := range conf.BaseConfig.Addr {
		balancerConfig.SetItem(balancer.NewBalancerItem(addr, conf.BaseConfig.Wight[idx]))
	}

	balancer, err := balancer.NewBalancer((*balancerConfig))
	if err != nil {
		return nil, err
	}

	h.client.CloseIdleConnections()
	h.b = balancer
	return h, nil
}

func (h *HttpClient) Send(ctx context.Context, method HttpMethod, key []byte, query url.Values, header map[string]string, body io.Reader) (interface{}, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	hostInfo, err := h.b.Pick(key)
	if err != nil {
		return nil, err
	}

	url, err := utils.Write(hostInfo.GetAddr(), "?", query.Encode())
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method.getMethod(), url, body)
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

	if _, ok := hm.manager[conf.BaseConfig.ServiceName]; ok {
		return errors.New("service already exists")
	}

	client, err := NewHttpClient(conf)
	if err != nil {
		return err
	}

	hm.manager[conf.BaseConfig.ServiceName] = client

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

func NewSingleHttpClientManager(conf []*HttpClientConfig) error {
	manager := NewHttpClientManager()
	for _, val := range conf {
		if err := manager.NewHttpClient(val); err != nil {
			return err
		}
	}
	HttpManager = manager
	return nil
}
