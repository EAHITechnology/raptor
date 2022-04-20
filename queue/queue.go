package queue

import (
	"errors"
)

type Queue interface {
	Put(v interface{}) error
	Pop() (interface{}, error)
	Len() int64
	Cap() int64
}

type QueueConfig struct {
	Typ   string
	Cap   int64
	Lenth int64
}

const (
	NilqueueType     = ""
	DefaultqueueType = "default"
)

var (
	ErrQueueType = errors.New("queue type error")
)

func NewQueue(conf QueueConfig) (Queue, error) {
	switch conf.Typ {
	case NilqueueType, DefaultqueueType:
		return NewDefaultqueue(), nil
	default:
		return nil, ErrQueueType
	}
}
