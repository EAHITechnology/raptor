package elog

import (
	"errors"
)

var (
	ErrLogPrefixNil   = errors.New("log Prefix is nil")
	ErrLogTypeIllegal = errors.New("log Type Illegal")
)
