package erpc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/EAHITechnology/raptor/balancer"
	"github.com/EAHITechnology/raptor/utils"
	"golang.org/x/net/context"
)

var HttpManager *HttpClientManager

var (
	ErrBalancerNil          = errors.New("balancer nil")
	ErrServiceNotExists     = errors.New("service not exists")
	ErrServiceAlreadyExists = errors.New("service already exists")
)

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

type HttpManagerConfig struct {
	Httpconf *HttpClientConfig
	Balancer balancer.Balancer
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

func NewHttpClient(conf *HttpClientConfig, b balancer.Balancer) (*HttpClient, error) {
	h := &HttpClient{}

	if conf.BaseConfig.IdleConnTimeout < DefaultHttpIdleTimeout {
		conf.BaseConfig.IdleConnTimeout = DefaultHttpIdleTimeout
	}

	h.conf = conf

	// http client
	h.client = &http.Client{
		Timeout: time.Millisecond * time.Duration(conf.BaseConfig.TimeOut),
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConns:        conf.BaseConfig.MaxIdleConns,
			MaxConnsPerHost:     conf.BaseConfig.MaxConnsPerAddr,
			MaxIdleConnsPerHost: conf.BaseConfig.MaxIdleConnsPerAddr,
			IdleConnTimeout:     time.Second * time.Duration(conf.BaseConfig.IdleConnTimeout),
			TLSHandshakeTimeout: 5 * time.Second,
			ReadBufferSize:      conf.BaseConfig.ReadBufferSize,
			WriteBufferSize:     conf.BaseConfig.WriteBufferSize,
		},
	}

	h.b = b
	return h, nil
}

func (h *HttpClient) Send(ctx context.Context, method HttpMethod, key []byte, uri string, query url.Values, header map[string]string, body io.Reader) ([]byte, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	hostInfo, err := h.b.Pick(key)
	if err != nil {
		return nil, err
	}

	url, err := utils.Write(hostInfo.GetAddr(), uri, "?", query.Encode())
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

func (h *HttpClient) AddAddr(addr string, wight int) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	item := balancer.NewBalancerItem(addr, wight)
	if err := h.b.Add(item); err != nil {
		return err
	}

	return nil
}

func (h *HttpClient) RemoveAddr(addr string) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if err := h.b.Remove(addr); err != nil {
		return err
	}
	return nil
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

func (hm *HttpClientManager) NewHttpClient(conf *HttpClientConfig, b balancer.Balancer) error {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	if _, ok := hm.manager[conf.BaseConfig.ServiceName]; ok {
		return ErrServiceAlreadyExists
	}

	client, err := NewHttpClient(conf, b)
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
		return nil, ErrServiceNotExists
	}

	return client, nil
}

func NewSingleHttpClientManager(conf []HttpManagerConfig) error {
	manager := NewHttpClientManager()
	for _, val := range conf {
		if val.Balancer == nil {
			return ErrBalancerNil
		}
		if err := manager.NewHttpClient(val.Httpconf, val.Balancer); err != nil {
			return err
		}
	}
	HttpManager = manager
	return nil
}
