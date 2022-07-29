package enet

type EnetConfig struct {
	L    NetLog
	Host string
}

type CommonJsonResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
