package erpc

type RpcNetConfigInfo struct {

	// Service Discovery Name, compatible dirpc.
	ServiceName string

	// Protocol type, HTTP, HTTPS and grpc are currently supported.
	Proto string

	// Service discovery type，eg: etcd, zk, apollo, list.
	// If we choose list，We will get the remote call service address from the "Addr" configuration.
	EndpointsFrom string

	// Address list. See "endpointsfrom" for details.
	// The "Addr" can also be competent for the task of service discovery.
	Addr []string

	// Weights
	Wight []int

	// Load balancing type.
	// eg: consistency_hash, p2c, random, range.
	Balancetype string

	// rpc dial time out
	DialTimeout int

	// rpc total time out
	TimeOut int

	// back off retry times
	RetryTimes int

	// every addr max conns num
	MaxConnsPerAddr int

	// every addr max idle conns num
	MaxIdleConnsPerAddr int

	// all addr max idle conns num
	MaxIdleConns int

	// idle timeout
	IdleConnTimeout int

	// ReadBufferSize
	ReadBufferSize int

	// WriteBufferSize
	WriteBufferSize int
}
