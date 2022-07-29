package grpc_balancer

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

const listScheme = "list"

// ---------------- Resolver -----------------

type ListResolver struct {
	addrs []string
	cc    resolver.ClientConn
}

func (fr *ListResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (fr *ListResolver) Close() {}

func (fr *ListResolver) resolveByList() ([]resolver.Address, error) {
	var addresses = make([]resolver.Address, 0, len(fr.addrs))
	for _, v := range fr.addrs {
		address := resolver.Address{
			Addr:       v,
			ServerName: v,
			Type:       resolver.Backend,
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// ---------------- Builder -----------------

func listBuilderRegister() {
	resolver.Register(NewListBuilder())
}

type ListBuilder struct{}

func NewListBuilder() resolver.Builder {
	return &ListBuilder{}
}

// Scheme for list
func (lb *ListBuilder) Scheme() string {
	return listScheme
}

// Build
func (lb *ListBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := strings.Split(target.URL.Host, ",")

	fr := &ListResolver{
		cc:    cc,
		addrs: addrs,
	}

	result, err := fr.resolveByList()
	if err != nil {
		return nil, err
	}

	if err := fr.cc.UpdateState(resolver.State{Addresses: result}); err != nil {
		return nil, err
	}
	return fr, nil
}
