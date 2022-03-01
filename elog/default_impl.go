package elog

import (
	"context"
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type defaultLogger struct {
	logger *zap.Logger
	lc     *LogConfig
}

// NewDefaultLogger function returns a `Logger` object. We need put in logger config.
func NewDefaultLogger(lc *LogConfig) (*defaultLogger, error) {
	if lc.Format == FORMAT_JSON {
		lc.Prefix = lc.Prefix + "_json"
	}

	now := time.Now()
	hook := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s-%04d%02d%02d%02d.log", lc.Dir, lc.Prefix, now.Year(), now.Month(), now.Day(), now.Hour()),
		MaxSize:    100,
		MaxBackups: 500,
		MaxAge:     7, //days
		Compress:   false,
	}

	hookErr := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s-%04d%02d%02d%02d_error.log", lc.Dir, lc.Prefix, now.Year(), now.Month(), now.Day(), now.Hour()),
		MaxSize:    100,
		MaxBackups: 500,
		MaxAge:     7, //days
		Compress:   false,
	}

	var syncer = zapcore.AddSync(hook)
	var syncerErr = zapcore.AddSync(hookErr)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,         // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder = zapcore.NewConsoleEncoder(encoderConfig)
	if lc.Format == FORMAT_JSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, syncer, infoLevel),
		zapcore.NewCore(encoder, syncerErr, warnLevel),
	)

	logger := zap.New(core).WithOptions(zap.AddCallerSkip(lc.Depth), zap.AddCaller())

	return &defaultLogger{
		logger: logger,
		lc:     lc,
	}, nil
}

func (l *defaultLogger) DebugfCtx(ctx context.Context, f string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) InfofCtx(ctx context.Context, f string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) WarnfCtx(ctx context.Context, f string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) ErrorfCtx(ctx context.Context, f string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) Printf(f string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) Debugf(f string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) Errorf(f string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) Infof(f string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) Warnf(f string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(f, args...))
}

func (l *defaultLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *defaultLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *defaultLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *defaultLogger) Warn(msg string) {
	l.logger.Warn(msg)
}
