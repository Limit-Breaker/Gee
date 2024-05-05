## 1 使用net/http库，启动web服务
Go语言内置的net/http库，封装了HTTP网络编程的基础的接口<br>
通过net/http库，可以很快搭建起一个http服务

```
func main() {
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
```

使用apifox或者postman访问该接口
```
Get http://localhost:8080/hello
```
可以得到返回
```
Header["Accept"] = ["*/*"]
Header["Accept-Encoding"] = ["gzip, deflate, br"]
Header["Connection"] = ["keep-alive"]
Header["User-Agent"] = ["Apifox/1.0.0 (https://apifox.com)"]
```
说明设置的路由处理函数已经被正确调用

http.ListenAndServe(":8080", nil)用来启动http服务，它可以传入两个参数，第一个参数":8080"是web服务监听的地址， 第二个参数是处理进入服务的http请求的实例， 设置为nil表示使用net/http自带的默认实例处理
```
// net/http库源码
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
    ...
	handler.ServeHTTP(rw, req)
}
```

## 2 实现http.Handler接口

ListenAndServe的第二个参数，是实现自定义web框架的入口
```
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```
handler是一个接口，只需要实现ServeHTTP方法
```
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
根据go语言的鸭子模型，只要实现ServeHTTP方法，即实现了该接口

定义一个结构体
```
// Engine 处理http请求的实例
type Engine struct{}
```
实现ServeHTTP 方法
```
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // 根据不同的路由匹配对应的处理函数
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
```
启动web服务的时候，传入该实例
```
func main() {
	engine := &Engine{}
	log.Fatal(http.ListenAndServe(":8080", engine))
}
```
访问
```
Get http://localhost:8080/hello
```
同样得到正确的返回结果
```
Header["Connection"] = ["keep-alive"]
Header["User-Agent"] = ["Apifox/1.0.0 (https://apifox.com)"]
Header["Accept"] = ["*/*"]
Header["Accept-Encoding"] = ["gzip, deflate, br"]
```

## 3 搭建giga框架的雏形
首先，giga需要作为一个单独包被其他程序调用，因此将原代码目录结构进行拆分
```
giga/
  |--giga.go
  |--go.mod
main.go
go.mod
```
之后对giga框架的开发都在giga目录下实现，主程序main.go只负责服务启停<br>
由于giga目录存放在本地，所以在goProject/go.mod中加入
```
require giga v0.0.0

replace giga => ./giga
```
对路由的注册和匹配这部分需要与单独解耦出，方便后期的扩展<br>
整体思路是在Engine 内部维护一个map，将路由和对应的处理函数以K,v的形式存放，在匹配路由时候根据请求组成key,在map中寻找对应的处理函数
```
type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

// New new Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRoute 新增路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}
```

同时，为了方便用户使用，对常见的请求进行简单的封装
```
// GET 新增Get请求
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 新增Post请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
```

对http服务的启动进行封装
```
// Run 启动http服务
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
```

由于路由的注册方式发送变动，因此ServeHTTP需要改写，在engine的router Map中匹配路由
```
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
```

使用新的giga框架启动服务

```
func main() {
	r := giga.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	r.Run(":8080")
}
```