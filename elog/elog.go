package elog

import (
	"context"
)

type LogLevel int32
type LogFormat string

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	FORMAT_JSON LogFormat = "json"
)

const (
	LOG_DIR = "./log"

	NULL_AUTO_CLEAR_HOURS = 0
	AUTO_CLEAR_HOURS      = 1
)

type Logger interface {
	DebugfCtx(ctx context.Context, f string, args ...interface{})
	InfofCtx(ctx context.Context, f string, args ...interface{})
	WarnfCtx(ctx context.Context, f string, args ...interface{})
	ErrorfCtx(ctx context.Context, f string, args ...interface{})

	Printf(f string, args ...interface{})
	Debugf(f string, args ...interface{})
	Errorf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Warnf(f string, args ...interface{})

	Debug(msg string)
	Error(msg string)
	Info(msg string)
	Warn(msg string)
}

func checkLogConfig(lc *LogConfig) error {
	if lc.Dir == "" {
		lc.Dir = LOG_DIR
	}

	if lc.Prefix == "" {
		return ErrLogPrefixNil
	}

	if lc.LogTyp == "" {
		lc.LogTyp = "default"
	}

	if lc.AutoClearHours == NULL_AUTO_CLEAR_HOURS {
		lc.AutoClearHours = AUTO_CLEAR_HOURS
	}

	return nil
}

func NewLogger(lc *LogConfig) (Logger, error) {
	if err := checkLogConfig(lc); err != nil {
		return nil, err
	}

	switch lc.LogTyp {
	case "default":
		return NewDefaultLogger(lc)
	default:
		return nil, ErrLogTypeIllegal
	}
}

func NewSingleLogger(lc *LogConfig) error {
	if err := checkLogConfig(lc); err != nil {
		return err
	}
	switch lc.LogTyp {
	case "default":
		logger, err := NewDefaultLogger(lc)
		if err != nil {
			return err
		}
		Elog = logger
	default:
		return ErrLogTypeIllegal
	}
	return nil
}
