package service_discovery

import (
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

type EtcdServiceDiscovery struct {
	config clientv3.Config
	client *clientv3.Client

	leaseId clientv3.LeaseID
	lease   clientv3.Lease

	cancel context.CancelFunc

	serviceMap map[string]chan ItermInfo

	serviceDiscoveryConfig ServiceDiscoveryConfig
	lock                   sync.RWMutex
	log                    ServiceDiscoveryLog
}

func NewEtcdServiceDiscovery(ctx context.Context, serviceDiscoveryConfig ServiceDiscoveryConfig, log ServiceDiscoveryLog) (*EtcdServiceDiscovery, error) {
	etcdConfig := clientv3.Config{
		Endpoints:   serviceDiscoveryConfig.EtcdAddr,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, err
	}

	// -------------  lease  -------------
	lease := clientv3.NewLease(client)
	leaseGrantResp, err := lease.Grant(ctx, 10)
	if err != nil {
		return nil, err
	}
	leaseId := leaseGrantResp.ID
	ctx, cancel := context.WithCancel(ctx)
	if _, err := lease.KeepAlive(ctx, leaseId); err != nil {
		cancel()
		return nil, err
	}

	etcdServiceDiscovery := &EtcdServiceDiscovery{
		serviceMap:             make(map[string]chan ItermInfo),
		client:                 client,
		config:                 etcdConfig,
		lease:                  lease,
		leaseId:                leaseId,
		cancel:                 cancel,
		serviceDiscoveryConfig: serviceDiscoveryConfig,
		log:                    log,
	}

	return etcdServiceDiscovery, nil
}

func (e *EtcdServiceDiscovery) ServiceRegister(ctx context.Context) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	kv := clientv3.NewKV(e.client)

	// key : service_name/ip:host
	if _, err := kv.Put(ctx,
		e.serviceDiscoveryConfig.LocalServiceName+"/"+e.serviceDiscoveryConfig.Ip+":"+e.serviceDiscoveryConfig.Host,
		"",
		clientv3.WithLease(e.leaseId),
	); err != nil {
		return err
	}
	return nil
}

func (e *EtcdServiceDiscovery) ServiceHeartbeat(ctx context.Context) error {
	return nil
}

func getEtcdOper(event *clientv3.Event) Operate {
	var oper Operate = Add
	switch event.Type {
	case mvccpb.PUT:
		if event.IsCreate() {
			oper = Add
		}
		if event.IsModify() {
			oper = Update
		}
	case mvccpb.DELETE:
		oper = Delete
	}
	return oper
}

func (e *EtcdServiceDiscovery) watchService(ctx context.Context, service string) chan ItermInfo {
	itermInfoChan := make(chan ItermInfo, 1024)

	go func(ctx context.Context, itermInfoChan chan ItermInfo) {
		resp, err := e.client.Get(ctx, service+"/", clientv3.WithPrefix())
		if err != nil {
			return
		}

		for _, kv := range resp.Kvs {
			select {
			case <-ctx.Done():
				close(itermInfoChan)
				return
			default:
				itermInfoChan <- ItermInfo{
					key:  kv.Key,
					info: kv.Value,
					oper: Add,
				}
			}
		}

		for {
			rch := e.client.Watch(ctx, service+"/", clientv3.WithPrefix())
			for wresp := range rch {
				if err := wresp.Err(); err != nil {
					e.log.Errorf("watchService Watch key:%s Error:%s", service+"/", err.Error())
					// TODO backoff
					time.Sleep(time.Second * 1)
					break
				}

				for _, ev := range wresp.Events {
					select {
					case <-ctx.Done():
						close(itermInfoChan)
						return
					default:
						itermInfoChan <- ItermInfo{
							key:  ev.Kv.Key,
							info: ev.Kv.Value,
							oper: getEtcdOper(ev),
						}
					}
				}
			}
		}
	}(ctx, itermInfoChan)

	return itermInfoChan
}

func (e *EtcdServiceDiscovery) ServiceDiscovery(ctx context.Context, service string) (chan ItermInfo, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, ok := e.serviceMap[service]; ok {
		return nil, ErrServiceAlreadyExists
	}

	rch := e.watchService(ctx, service)
	e.serviceMap[service] = rch

	return rch, nil
}

func (e *EtcdServiceDiscovery) Close(ctx context.Context) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.cancel()

	if _, err := e.lease.Revoke(ctx, e.leaseId); err != nil {
		return err
	}
	if err := e.client.Close(); err != nil {
		return err
	}
	e.serviceMap = nil
	return nil
}
