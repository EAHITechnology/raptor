package skip_list

import (
	"bytes"
	"errors"
	"math/rand"
	"sync"
)

const (
	defaultMaxLevel int     = 16
	defaultP        float32 = 0.25
)

var (
	ErrDataNotExist = errors.New("data does not exist")
)

type DefaultSkipListNode struct {
	Key     []byte
	Value   []byte
	forward []*DefaultSkipListNode
}

type DefaultSkipListImpl struct {
	level  int
	Lenth  int64
	lock   sync.RWMutex
	header *DefaultSkipListNode
}

type DefaultSkipListIterImpl struct {
	skl  *DefaultSkipListImpl
	iter *DefaultSkipListNode
}

func newForwards(level int) []*DefaultSkipListNode {
	forwards := []*DefaultSkipListNode{}

	for idx := level - 1; idx >= 0; idx-- {
		var forward *DefaultSkipListNode = nil
		forwards = append(forwards, forward)
	}

	return forwards
}

func NewDefaultSkipListNode(key, Value []byte, level int) *DefaultSkipListNode {
	return &DefaultSkipListNode{
		Key:     key,
		Value:   Value,
		forward: newForwards(level),
	}
}

func NewDefaultSkipListImpl() (*DefaultSkipListImpl, error) {
	header := NewDefaultSkipListNode([]byte(""), []byte(""), defaultMaxLevel)
	skl := &DefaultSkipListImpl{Lenth: 0, header: header, level: defaultMaxLevel}
	return skl, nil
}

func (d *DefaultSkipListImpl) Get(key []byte) ([]byte, bool, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	node := d.get(key, nil)
	if node == nil {
		return nil, false, nil
	}

	if bytes.Equal(node.Key, key) {
		return node.Value, true, nil
	}

	return nil, false, nil
}

func (d *DefaultSkipListImpl) get(key []byte, updateForwards *[]*DefaultSkipListNode) *DefaultSkipListNode {
	if len(key) == 0 {
		return nil
	}

	start := d.header
	for idx := d.level - 1; idx >= 0; idx-- {
		for start.forward[idx] != nil {
			compareValue := bytes.Compare(start.forward[idx].Key, key)
			if compareValue == -1 {
				start = start.forward[idx]
			} else {
				break
			}
		}

		if updateForwards != nil {
			(*updateForwards)[idx] = start
		}
	}
	return start.forward[0]
}

func (d *DefaultSkipListImpl) Put(key, value []byte) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	update := newForwards(defaultMaxLevel)

	node := d.get(key, &update)
	if node != nil && bytes.Equal(node.Key, key) {
		node.Value = value
		return nil
	}

	level := d.defaultRandomLevel()
	newNode := NewDefaultSkipListNode(key, value, level)
	for i := 0; i < level; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	d.Lenth++

	return nil
}

func (d *DefaultSkipListImpl) Del(key []byte) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	update := make([]*DefaultSkipListNode, defaultMaxLevel)

	node := d.get(key, &update)
	if node == nil {
		return ErrDataNotExist
	}

	if bytes.Equal(node.Key, key) {
		for idx := 0; idx < d.level; idx++ {
			if update[idx].forward[idx] != node {
				break
			}
			update[idx].forward[idx] = node.forward[idx]
		}
		d.Lenth--
	}

	return nil
}

func (d *DefaultSkipListImpl) Len() int64 {
	return d.Lenth
}

// TODO(EAHITechnology) iter checkpoint
func (d *DefaultSkipListImpl) GetIter() Iter {
	iter := &DefaultSkipListIterImpl{
		skl:  d,
		iter: d.header,
	}
	return iter
}

func (d *DefaultSkipListImpl) defaultRandomLevel() int {
	level := 1
	for rand.Float32() < defaultP && level < defaultMaxLevel {
		level++
	}
	return level
}

func (d *DefaultSkipListIterImpl) Seek(key []byte) error {
	d.skl.lock.RLock()
	defer d.skl.lock.RUnlock()

	node := d.skl.get(key, nil)
	if node == nil {
		return ErrDataNotExist
	}

	d.iter = node
	return nil
}

func (d *DefaultSkipListIterImpl) Get() ([]byte, error) {
	d.skl.lock.RLock()
	defer d.skl.lock.RUnlock()

	if d.iter == nil {
		return nil, ErrDataNotExist
	}

	return d.iter.Value, nil
}

func (d *DefaultSkipListIterImpl) Next() ([]byte, error) {
	if d.iter != nil && d.iter.forward[0] != nil {
		d.skl.lock.RLock()
		defer d.skl.lock.RUnlock()

		d.iter = d.iter.forward[0]
		return d.iter.Value, nil
	}

	return nil, ErrDataNotExist
}

func (d *DefaultSkipListIterImpl) HasNext() bool {
	return (d.iter != nil && d.iter.forward[0] != nil)
}

func (d *DefaultSkipListIterImpl) Close() {
	d.skl = nil
	d.iter = nil
}
