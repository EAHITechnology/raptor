package balancer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsistencyHashBalancer(t *testing.T) {

	list := []balancerItem{}
	item0 := balancerItem{
		addr:  "xx.xxx.xxx.01",
		wight: 1,
	}

	item1 := balancerItem{
		addr:  "xx.xxx.xxx.02",
		wight: 1,
	}

	item2 := balancerItem{
		addr:  "xx.xxx.xxx.03",
		wight: 1,
	}

	item3 := balancerItem{
		addr:  "xx.xxx.xxx.04",
		wight: 1,
	}

	list = append(list, item0, item1, item2, item3)

	conf := balancerConfig{
		balancerTyp:     ConsistencyHashType,
		balancerConfigs: list,
	}
	balancer, err := NewConsistencyHashBalancer(conf)
	assert.Nil(t, err)
	host, err := balancer.Pick([]byte("23472"))
	assert.Nil(t, err)

	fmt.Println("addr:", host.GetAddr())
}
