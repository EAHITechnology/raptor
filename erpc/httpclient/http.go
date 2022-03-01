package httpclient

// todo

import (
	"net"
	"net/http"
	"time"
)

type HttpClientConfig struct {
	TimeOut             time.Duration
	DialTimeout         time.Duration
	IdleConnTimeout     time.Duration
	KeepAlive           time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	ReadBufferSize      int
	WriteBufferSize     int
}

type HttpClient struct {
	conf   *HttpClientConfig
	client *http.Client
}

func NewHttpClient(conf *HttpClientConfig) *HttpClient {
	c := &HttpClient{}
	c.conf = conf
	c.client = &http.Client{
		Timeout: conf.TimeOut,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   conf.DialTimeout,
				KeepAlive: conf.KeepAlive,
			}).DialContext,
			MaxIdleConns:        conf.MaxIdleConns,
			MaxConnsPerHost:     conf.MaxConnsPerHost,
			MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
			IdleConnTimeout:     conf.IdleConnTimeout,
			ReadBufferSize:      conf.ReadBufferSize,
			WriteBufferSize:     conf.WriteBufferSize,
		},
	}
	return c
}
