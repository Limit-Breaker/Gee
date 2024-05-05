package giga

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type HandlerFunc func(*Context)

// Engine
type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // support middleware
		engine      *Engine
	}

	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func NewEngine() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: group.engine,
	}
	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

// Use 增加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s, group.prefix:%s ", method, pattern, group.prefix)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET 新增Get请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 新增Post请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 通过路由找到对应的中间件函数，并将其加入到context里
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

// Run defines the method to start a http server
//func (engine *Engine) Run(addr string) {
//	if err := http.ListenAndServe(addr, engine); err != nil && err != http.ErrServerClosed {
//		log.Fatalf("server listen err:%s \n", err)
//	}
//}

func (engine *Engine) Run(srvName string, addr string) {
	server := http.Server{
		Addr:    addr,
		Handler: engine,
	}
	// 优雅启停
	go func() {
		log.Printf("server %s, running in %s\n", srvName, server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server %s listen err:%s \n", srvName, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// 接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	// SIGINT 用户发送中断(Ctrl+C)
	// SIGTERM 终止进程 软件终止信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在此阻塞
	<-quit
	log.Printf("server %s shutting down...\n", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server %s shutdown err, cause by: %s\n", srvName, err)
	}

	select {
	case <-ctx.Done():
		log.Printf("server %s wait shutdown timeout \n", srvName)
	}
	log.Printf("server %s exiting...\n\n", srvName)
}
