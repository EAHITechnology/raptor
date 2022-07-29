package grpc_balancer

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc/resolver"
)

const fileScheme = "file"

// ---------------- Resolver -----------------

type SpaceResolver struct {
	path        string
	serviceName string
	reloadTime  time.Duration
	addrMap     map[string]struct{}
	cc          resolver.ClientConn
	ctx         context.Context
	cancel      context.CancelFunc
	lock        sync.Mutex
}

func (fr *SpaceResolver) ResolveNow(resolver.ResolveNowOptions) {
	fr.lock.Lock()
	defer fr.lock.Unlock()
	result, err := fr.resolveByFile()
	if err != nil {
		fr.cc.ReportError(err)
		return
	}

	if len(result) == 0 {
		return
	}

	addrMap := make(map[string]struct{})
	for _, r := range result {
		addrMap[r.Addr] = struct{}{}
	}

	if len(fr.addrMap) == len(addrMap) {
		// expect sameFlag is true.
		// if sameFlag == true , we don't need to compare.
		// UpdateState function has redundant locks,
		// so we don't want the comparison logic to sink into grpc.
		var sameFlag bool = true
		for key := range addrMap {
			if _, ok := fr.addrMap[key]; !ok {
				sameFlag = false
				break
			}
		}

		if sameFlag {
			return
		}
	}

	fr.addrMap = addrMap
	fr.cc.UpdateState(resolver.State{Addresses: result})
}

func (fr *SpaceResolver) Close() {
	fr.cancel()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

// resolveByFile function parse a naming file.
func (fr *SpaceResolver) resolveByFile() ([]resolver.Address, error) {
	vp := viper.New()
	vp.SetConfigType("json")
	vp.SetConfigFile(fr.path)
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	items := vp.GetStringSlice(fr.serviceName)
	var addresses = make([]resolver.Address, 0, len(items))
	for _, v := range items {
		address := resolver.Address{
			Addr:       v,
			ServerName: v,
			Type:       resolver.Backend,
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// ---------------- watcher -----------------

// watcher monitors file changes.
// TODO(:malikhou) event driven.
func (fr *SpaceResolver) watcher(ctx context.Context) {
	ticker := time.NewTicker(fr.reloadTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fr.ResolveNow(resolver.ResolveNowOptions{})
		}
	}
}

// ---------------- Builder -----------------

// SpaceBuilderRegister register a builder in grpc resolver.
func spaceBuilderRegister() {
	resolver.Register(NewSpaceBuilder())
}

type SpaceBuilder struct{}

func NewSpaceBuilder() resolver.Builder {
	return &SpaceBuilder{}
}

// Scheme for space
func (mb *SpaceBuilder) Scheme() string {
	return fileScheme
}

// Build
func (mb *SpaceBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	fileName := target.URL.Path
	serviceName := target.URL.Query().Get("service_name")
	reloadTimeStr := target.URL.Query().Get("reload_time")

	if !pathExists(fileName) {
		return nil, ErrEmptyShopeeFilePath
	}

	if serviceName == "" {
		return nil, ErrEmptyService
	}

	reloadTime := defaultReloadTime
	if reloadTimeStr != "" {
		reloadTimeInt, err := strconv.ParseInt(reloadTimeStr, 10, 64)
		if err != nil {
			return nil, err
		}
		reloadTime = time.Second * time.Duration(reloadTimeInt)
	}

	ctx, cancel := context.WithCancel(context.Background())
	fr := &SpaceResolver{
		cc:          cc,
		path:        fileName,
		serviceName: serviceName,
		reloadTime:  reloadTime,
		ctx:         ctx,
		cancel:      cancel,
		addrMap:     make(map[string]struct{}),
	}

	fr.ResolveNow(resolver.ResolveNowOptions{})
	go fr.watcher(fr.ctx)
	return fr, nil
}
