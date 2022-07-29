package grpc_balancer

import (
	"math/rand"
	"time"

	"github.com/EAHITechnology/raptor/utils/rand2"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const RANDOM_LB_NAME = "random"

func randomRegister() {
	balancer.Register(
		base.NewBalancerBuilder(RANDOM_LB_NAME, &RandomBuilder{
			rand: rand2.New(rand.NewSource(time.Now().UnixNano())),
		}, base.Config{HealthCheck: true}))
}

type RandomBuilder struct {
	rand *rand2.Rand
}

func (rb *RandomBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for subConn := range info.ReadySCs {
		scs = append(scs, subConn)
	}
	return &RandomPicker{
		subConn: scs,
		rand:    rb.rand,
	}
}

type RandomPicker struct {
	subConn []balancer.SubConn
	rand    *rand2.Rand
}

func (rp *RandomPicker) Pick(info balancer.PickInfo) (r balancer.PickResult, err error) {
	if len(rp.subConn) == 0 {
		return r, balancer.ErrNoSubConnAvailable
	}

	if len(rp.subConn) == 1 {
		r.SubConn = rp.subConn[0]
		return r, nil
	}

	r.SubConn = rp.subConn[rp.rand.Int()%len(rp.subConn)]
	return r, nil
}
