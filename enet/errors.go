package enet

import (
	"errors"
)

var (
	ErrHttpGroupNil = errors.New("http group nil")
	ErrHttpLogNil   = errors.New("log is nil")
)
