package double_buffer

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/EAHITechnology/raptor/utils"
)

type DoubleBufLog interface {
	Debugf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Warnf(f string, args ...interface{})
	Errorf(f string, args ...interface{})
}

type Reloader interface {
	ReloadBuf(ctx context.Context) (interface{}, error)
}

type DoubleBufferOpts struct {
	Reloader   Reloader
	RelaodTime time.Duration
	Log        DoubleBufLog
}

type DoubleBuffer struct {
	flag      int32
	buf       [2]interface{}
	closeChan chan struct{}
	opts      DoubleBufferOpts
}

var (
	ErrDoubleBufferLogNil      = errors.New("double_buffer log nil")
	ErrDoubleBufferReloaderNil = errors.New("double_buffer Reloader nil")

	DefaultRelaodTime = time.Minute * 1
)

func NewDoubleBuffer(ctx context.Context, opts DoubleBufferOpts) (*DoubleBuffer, error) {
	if opts.Log == nil || utils.IsNil(opts.Log) {
		return nil, ErrDoubleBufferLogNil
	}

	if opts.Reloader == nil || utils.IsNil(opts.Reloader) {
		return nil, ErrDoubleBufferReloaderNil
	}

	if opts.RelaodTime < DefaultRelaodTime {
		opts.RelaodTime = DefaultRelaodTime
	}

	buf, err := opts.Reloader.ReloadBuf(ctx)
	if err != nil {
		return nil, err
	}

	doubleBuffer := &DoubleBuffer{
		flag:      0,
		closeChan: make(chan struct{}),
		opts:      opts,
	}
	doubleBuffer.buf[0] = buf
	doubleBuffer.buf[1] = buf

	go doubleBuffer.runReload(ctx)

	return doubleBuffer, nil
}

func (d *DoubleBuffer) runReload(ctx context.Context) {
	ticker := time.NewTicker(d.opts.RelaodTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			newBuf, err := d.opts.Reloader.ReloadBuf(ctx)
			if err != nil {
				d.opts.Log.Errorf("runReload ReloadBuf err:%v", err)
			}
			d.buf[(atomic.LoadInt32(&d.flag)+1)%2] = newBuf
			atomic.AddInt32(&d.flag, 1)
		}
	}
}

func (d *DoubleBuffer) GetBuf() interface{} {
	return d.buf[atomic.LoadInt32(&d.flag)%2]
}

func (d *DoubleBuffer) Close() {
	close(d.closeChan)
}
