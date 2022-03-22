package server

import (
	"errors"

	"github.com/EAHITechnology/raptor/balancer"
	"golang.org/x/net/context"
)

type ServerTyp string

const (
	NULL_TYPE    ServerTyp = ""
	DEFAULT_TYPE ServerTyp = "default"
)

const (
	defaultAbleToWrite int32 = 0
	defaultWriting     int32 = 1
	defaultEndWrite    int32 = 2
)

type ServerFunc func(ctx context.Context)

type atferFuncObj struct {
	serverFuncList []ServerFunc
	flag           int32
}

type beforeFuncObj struct {
	serverFuncList []ServerFunc
	flag           int32
}

type Server interface {
	// Execute after Run function.
	// AfterInit is only suitable for initialization functions,
	// not suitable executing complex logic.
	AfterInit(serverFunc ...ServerFunc)

	// Execute before Run function.
	// The framework may still load some necessary functions
	// before Run executes(like logger, service register).
	// BeforeInit is only suitable for initialization functions,
	// not suitable for executing complex logic.
	BeforeInit(serverFunc ...ServerFunc)

	// run frame
	// The function will block until an error is
	// reported in the process.
	Run(ctx context.Context, cancel context.CancelFunc) error
}

func NewServer(ctx context.Context, typ ServerTyp, configPath string) (Server, error) {
	switch typ {
	case NULL_TYPE, DEFAULT_TYPE:
		return NewDefaultServer(ctx, typ, configPath)
	default:
		return NewDefaultServer(ctx, typ, configPath)
	}
}

func getBalancerTyp(typ string) (balancer.BalancerTyp, error) {
	var balancetype balancer.BalancerTyp
	switch typ {
	case "", "random":
		balancetype = balancer.RandomType
	case "p2c":
		balancetype = balancer.P2cType
	case "consistency_hash":
		balancetype = balancer.ConsistencyHashType
	case "range":
		balancetype = balancer.RangeType
	default:
		return "", errors.New("illegal type")
	}
	return balancetype, nil
}
