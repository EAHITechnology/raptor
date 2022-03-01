package config

import "errors"

var (
	ErrLogPrefix = errors.New("log prefix is nil")
	ErrLogDir    = errors.New("log dir is nil")
	ErrLogLevel  = errors.New("log level is illegal")

	ErrServerNameNil = errors.New("server name is nil")

	ErrConfigNil = errors.New("config nil")
)
