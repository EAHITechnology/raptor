package queue

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

const (
	DefaultReTryTimes = 3
)

var (
	ErrPutFail = errors.New("put fail")
	ErrPopNil  = errors.New("pop nil")
	ErrPopFail = errors.New("pop fail")
)

type Defaultqueue struct {
	head  unsafe.Pointer
	tail  unsafe.Pointer
	lenth int64
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}

func NewDefaultqueue() *Defaultqueue {
	n := unsafe.Pointer(&node{})
	return &Defaultqueue{head: n, tail: n}
}

// Rpush puts the given value v at the tail of the queue.
func (d *Defaultqueue) Put(v interface{}) error {
	n := &node{value: v}
	retryTimes := DefaultReTryTimes
	for retryTimes > 0 {
		retryTimes--
		tail := load(&d.tail)
		next := load(&tail.next)
		if tail == load(&d.tail) {
			if next == nil {
				if cas(&tail.next, next, n) {
					cas(&d.tail, tail, n)
					d.lenth++
					return nil
				}
			} else {
				cas(&d.tail, tail, next)
			}
		}
	}
	return ErrPutFail
}

// Lpop removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (d *Defaultqueue) Pop() (interface{}, error) {
	retryTimes := DefaultReTryTimes
	for retryTimes > 0 {
		retryTimes--
		head := load(&d.head)
		tail := load(&d.tail)
		next := load(&head.next)
		if head == load(&d.head) {
			if head == tail {
				if next == nil {
					return nil, ErrPopNil
				}
				cas(&d.tail, tail, next)
			} else {
				v := next.value
				if cas(&d.head, head, next) {
					cas(&head.next, next, nil)
					d.lenth--
					return v, nil
				}
			}
		}
	}
	return nil, ErrPopFail
}

func (d *Defaultqueue) Len() int64 {
	return d.lenth
}

func (d *Defaultqueue) Cap() int64 {
	return d.lenth
}
