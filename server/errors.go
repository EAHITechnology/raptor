package server

import "errors"

var (
	ErrLogLevel      = errors.New("log level is illegal")
	ErrLogPrefix     = errors.New("log prefix is nil")
	ErrLogDir        = errors.New("log dir is nil")
	ErrServerNameNil = errors.New("server name is nil")
	ErrConfigNil     = errors.New("config nil")
	ErrRedisNameNil  = errors.New("redis name is nil")
	ErrMysqlNameNil  = errors.New("mysql name is nil")
	ErrMysqlIpNil    = errors.New("mysql ip is nil")
)
