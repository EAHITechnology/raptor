package enet

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/EAHITechnology/raptor/context_trace"
	"github.com/EAHITechnology/raptor/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Context struct {
	*gin.Context
}

type HandlerFunc func(c *Context)

type NetLog interface {
	Debugf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Warnf(f string, args ...interface{})
	Errorf(f string, args ...interface{})
}

type RouterGroup struct {
	routerGroup *gin.RouterGroup
}

type HttpServer struct {
	engine      *gin.Engine
	listener    net.Listener
	closeCh     chan struct{}
	log         NetLog
	routerGroup map[string]*RouterGroup
}

func handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			c,
		}
		h(ctx)
	}
}

func InitHttpServer(ec *EnetConfig) (*HttpServer, error) {
	if utils.IsNil(ec.L) || ec.L == nil {
		return nil, ErrHttpLogNil
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	engine := gin.New()
	engine.Use(gin.Recovery())

	hs := &HttpServer{
		engine:      engine,
		closeCh:     make(chan struct{}),
		log:         ec.L,
		routerGroup: make(map[string]*RouterGroup),
	}

	engine.Use(hs.traceMiddle())
	engine.Use(hs.loggerMiddle())

	listener, err := net.Listen("tcp", ec.Host)
	if err != nil {
		return nil, err
	}
	hs.listener = listener

	engine.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))

	pprofRouteGroup := engine.Group("/debug/pprof")
	pprofRouteGroup.Any("/", gin.WrapF(pprof.Index))
	pprofRouteGroup.Any("/cmdline", gin.WrapF(pprof.Cmdline))
	pprofRouteGroup.Any("/profile", gin.WrapF(pprof.Profile))
	pprofRouteGroup.Any("/symbol", gin.WrapF(pprof.Symbol))
	pprofRouteGroup.Any("/trace", gin.WrapF(pprof.Trace))
	pprofRouteGroup.Any("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	pprofRouteGroup.Any("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	pprofRouteGroup.Any("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	pprofRouteGroup.Any("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
	pprofRouteGroup.Any("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	pprofRouteGroup.Any("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))

	return hs, nil
}

func InitHttpServerSingle(ec *EnetConfig) error {
	hs, err := InitHttpServer(ec)
	if err != nil {
		return err
	}
	HttpWeb = hs
	return nil
}

func (h *HttpServer) Close() {
	close(h.closeCh)
}

func (h *HttpServer) traceMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		c.Request, _, err = context_trace.GetHeaderTrace(c.Request)
		if err != nil {
			h.log.Errorf("traceMiddle GetHeaderTrace error:%s", err.Error())
		}
	}
}

func (h *HttpServer) loggerMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.Request.Host
		remoteIP, _ := c.RemoteIP()
		_, trace, err := context_trace.GetCtxTrace(c.Request.Context())
		if err != nil {
			h.log.Errorf("loggerMiddle GetCtxTrace err:%s", err.Error())
		}
		h.log.Infof("| %3d | %13v | %15s | %15s | %s | %s | %s |", statusCode, latencyTime, remoteIP, clientIP, reqMethod, reqUri, trace)
	}
}

func (h *HttpServer) NewGroup(path string) *RouterGroup {
	if group, ok := h.routerGroup[path]; ok {
		return group
	}

	rg := &RouterGroup{
		routerGroup: h.engine.Group(path),
	}

	h.routerGroup[path] = rg
	return rg

}

// ------------------- routerGroup -------------------

func (r *RouterGroup) Post(relativePath string, handlers ...HandlerFunc) {
	handlerFunces := []gin.HandlerFunc{}
	for _, handler := range handlers {
		handlerFunces = append(handlerFunces, handle(handler))
	}
	r.routerGroup.POST(relativePath, handlerFunces...)
}

func (r *RouterGroup) Get(relativePath string, handlers ...HandlerFunc) {
	handlerFunces := []gin.HandlerFunc{}
	for _, handler := range handlers {
		handlerFunces = append(handlerFunces, handle(handler))
	}
	r.routerGroup.GET(relativePath, handlerFunces...)
}

func (r *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) {
	handlerFunces := []gin.HandlerFunc{}
	for _, handler := range handlers {
		handlerFunces = append(handlerFunces, handle(handler))
	}
	r.routerGroup.Any(relativePath, handlerFunces...)
}

// ------------------- routerGroup end-------------------

func (h *HttpServer) Post(path string, handlers ...HandlerFunc) {
	handlerFunces := []gin.HandlerFunc{}
	for _, handler := range handlers {
		handlerFunces = append(handlerFunces, handle(handler))
	}
	h.engine.POST(path, handlerFunces...)
}

func (h *HttpServer) Get(path string, handlers ...HandlerFunc) {
	handlerFunces := []gin.HandlerFunc{}
	for _, handler := range handlers {
		handlerFunces = append(handlerFunces, handle(handler))
	}
	h.engine.GET(path, handlerFunces...)
}

func (h *HttpServer) Any(path string, handlers ...HandlerFunc) {
	handlerFunces := []gin.HandlerFunc{}
	for _, handler := range handlers {
		handlerFunces = append(handlerFunces, handle(handler))
	}
	h.engine.Any(path, handlerFunces...)
}

func CreateJsonResp(code int, msg string) CommonJsonResp {
	return CommonJsonResp{
		Code: code,
		Msg:  msg,
	}
}

func CreateSuccessJsonResp() CommonJsonResp {
	return CommonJsonResp{
		Code: http.StatusOK,
		Msg:  "success",
	}
}

func (h *HttpServer) Run() {
	defer func() {
		if h.listener != nil {
			if err := h.listener.Close(); err != nil {
				h.log.Errorf("listener close err:%v", err)
			}
		}
	}()

	errCh := make(chan error)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", h.engine)
		errCh <- http.Serve(h.listener, mux)
	}()

	select {
	case <-h.closeCh:
		h.log.Infof("closing http api server")
		return
	case err := <-errCh:
		h.log.Errorf("http api server exit on error:%v", err)
		return
	}

}
