package erpc

type RpcNetConfigInfo struct {

	// Service Discovery Name, compatible dirpc.
	ServiceName string

	// Protocol type, HTTP, HTTPS and grpc are currently supported.
	Proto string

	// Service discovery type，eg: disf,list.
	// If we choose list，We will get the remote call service address from the "Addr" configuration.
	EndpointsFrom string

	// Address list. See "endpointsfrom" for details.
	// The "Addr" can also be competent for the task of service discovery.
	Addr []string

	// Load balancing type
	// eg: consistency_hash, hash, range
	Balancetype string

	// rpc dial time out
	DialTimeout int

	// rpc read time out
	ReadTimeout int

	// back off retry times
	RetryTimes int

	// rpc
	MaxSize int

	// max free conn
	MaxIdleConn int

	// fuse flag
	BreakFlag bool
}
