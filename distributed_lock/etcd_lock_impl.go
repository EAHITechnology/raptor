package distributed_lock

import (
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

type EtcdDistributedLockManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	leaseID clientv3.LeaseID
	lease   clientv3.Lease
	cancel  context.CancelFunc
	conf    DistributedLockConfig
	key     string
}

func NewEtcdDistributedLockManager(ctx context.Context, conf DistributedLockConfig) (*EtcdDistributedLockManager, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.EtcdAddrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	lease := clientv3.NewLease(client)
	leaseResp, err := lease.Grant(ctx, conf.TTl)
	if err != nil {
		return nil, err
	}
	leaseID := leaseResp.ID

	ctx, cancel := context.WithCancel(ctx)
	if _, err := lease.KeepAlive(ctx, leaseID); err != nil {
		cancel()
		return nil, err
	}

	etcdDistributedLockManager := &EtcdDistributedLockManager{
		client:  client,
		key:     conf.Key,
		conf:    conf,
		lease:   lease,
		leaseID: leaseID,
		cancel:  cancel,
	}

	etcdDistributedLockManager.kv = clientv3.NewKV(client)

	return etcdDistributedLockManager, nil
}

func (e *EtcdDistributedLockManager) Lock(ctx context.Context) (string, error) {
	tx := e.kv.Txn(ctx)

	version := strconv.FormatInt(time.Now().Unix(), 10)

	tx.If(clientv3.Compare(clientv3.CreateRevision(e.key), "=", 0)).
		Then(clientv3.OpPut(e.key, version, clientv3.WithLease(e.leaseID)))
	txnResp, err := tx.Commit()
	if err != nil {
		return version, err
	}

	if !txnResp.Succeeded {
		return version, ErrLockFail
	}
	return version, nil
}

func (e *EtcdDistributedLockManager) Unlock(ctx context.Context, value string) error {
	tx := e.kv.Txn(ctx)

	tx.If(clientv3.Compare(clientv3.Value(e.key), "=", value)).
		Then(clientv3.OpDelete(e.key))

	txnResp, err := tx.Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return ErrUnLockFail
	}

	return nil
}

func (e *EtcdDistributedLockManager) Close(ctx context.Context) error {
	e.cancel()
	e.client.Close()
	return nil
}
