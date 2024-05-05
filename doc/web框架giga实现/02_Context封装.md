## 1 封装context

处理http请求，需要根据请求*http.Request的参数，处理之后，构造响应http.ResponseWriter。<br>
但是这两个对象提供的接口粒度太细，用户在使用的时候需要设置的细节太多，比如我们要构造一个完整的响应，需要考虑消息头(Header)和消息体(Body)，而 Header 包含了状态码(StatusCode)，消息类型(ContentType)等几乎每次请求都需要设置的信息。<br>
其次，如果不进行有效的封装，那么框架的用户将需要对每个请求写大量重复，繁杂的代码，而且容易出错。
第三，出于对框架功能扩展的需求，比如实现中间件，动态路由等功能


因此，参gin框架的实现，加入context的使用，新建context.go
```
giga/
  |--context.go
  |--giga.go
  |--go.mod
main.go
go.mod
```

定义context结构体
```
type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
```
context里面暂时比较简单，包含了http.ResponseWriter和*http.Request，另外提供了对 Method 和 Path 这两个常用属性的直接访问。<br>
同时提供了新建context的方法

在此基础上，封装对请求中几种常用属性的访问和设置方法
```
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
```
进一步地，还可以封装提供快速构造String/Data/JSON/HTML响应的方法
```
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
```

## 改造route

将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go，方便对router的功能进行增强，例如提供动态路由的支持。<br>
由于想利用使用context对请求进行了封装，因此想通过请求的method和path组成Key需要使用context。<br>
在router 增加handle方法，传入Context。
```
giga/
  |--context.go
  |--giga.go
  |--go.mod
  |--router.go
main.go
go.mod
```
router.go代码
```
type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
```

