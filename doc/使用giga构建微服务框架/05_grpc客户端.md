## grpc简明教程
https://zhuanlan.zhihu.com/p/660634947


## grpc客户端
以一个例子说明怎样在本框架中使用grpc
用户发送登陆请求，apiProxy在收到用户请求后，调用grpc接口，将用户手机号发送给grpc服务端，服务端生成验证码，调用短信平台将验证码发送到用户手机

#### 基于protobuf定义服务接口
在internal目录下新建service目录，service目录下建立pb目录，在Pb目录下新建userService.proto文件
在userService.proto文件中定义rpc调用函数和输入输出结构体
```
// 定义服务
service UserService {
  // GetCaptcha 方法
  rpc GetCaptcha (GetCaptchaRequest) returns (GetCaptchaResponse) {}
}

// 请求消息
message GetCaptchaRequest {
  string mobile = 1;
}

// 响应消息
message GetCaptchaResponse {
  string code = 1;
}
```
GetCaptchaRequest是输入结构体， GetCaptchaResponse是grpc服务端返回的结构体

在pb目录下使用命令
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative 在userService.proto
```
生成userService.pb.go和userService_grpc.pb.go文件 <br>
userService_grpc.pb.go 文件包含的内容主要是需要实现和使用的go接口代码，server端需要用的函数、待实现的接口以及client端需要用的函数都在里面。
userService.pb.go文件生成的是要用到的数据结构，比如proto文件中定义的GetCaptchaRequest就被编译成了go中的struct结构

#### grpc客户端
在apiProxy/router/router.go目录下，增加grpc客户端的初始化函数
```
// UserServiceClient 维护UserService客户端实例
var UserServiceClient pb.UserServiceClient

func InitRpcClient() {
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial("127.0.0.1:8972",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	UserServiceClient = pb.NewUserServiceClient(conn)
}
```

#### 使用grpc客户端
如何让请求的处理函数使用grpc调用函数？
可以考虑使用中间件，将grpc客户端实例放入context中。请求的处理函数可以在context中将需要的实例拿出
在giga.context中增加
```
Keys map[string]interface{}
```
在middleware目录下新增一个rpc的middleware
```
// MiddlewareRpc 将rpc client实例存在Keys中
func MiddlewareRpc(services map[string]interface{}) giga.HandlerFunc {
	return func(c *giga.Context) {
		// 将rpc client实例存在Keys中
		c.Keys = make(map[string]interface{})
		for k, v := range services {
			c.Keys[k] = v
		}
		c.Next()
	}
}
```
之后，在路由分组中使用这个中间件
```
func (l *RouterUser) Route(r *giga.Engine) {
	h := user.HandlerUser{}
	// 设置路由分组
	user := r.Group("/user")
	// 将rpc实例设置到context中
	m := make(map[string]interface{})
	m["user"] = UserServiceClient

	user.Use(middleware.MiddlewareRpc(m))
	{

		user.POST("/register", h.UserRegister)
		user.POST("/login", h.UserLogin)
	}
}
```
请求处理函数需要调用rpc接口的时，取出对应的客户端实例
```
func (h *HandlerUser) UserLogin(c *giga.Context) {
	// 取出参数
	mobile := c.PostForm("mobile")
	// Context.Key中取出服务实例
	userService, ok := c.Keys["user"].(pb.UserServiceClient)
	if !ok {
		log.Fatalf("could not get rpc client")
		return
	}
	// 设置超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 执行RPC调用并打印收到的响应数据
	res, err := userService.GetCaptcha(ctx, &pb.GetCaptchaRequest{Mobile: mobile})
	if err != nil {
		log.Fatalf("could not get captcha: %v", err)
		return
	}
	c.JSON(http.StatusOK, giga.H{"code": res.Code})
}
```
#### 修改配置文件
在启动grpc客户端时，grpc的配置同样可以通过读取配置文件的方式拿到
在config.yaml中增加grpc的配置
```
grpc:
  addr: "127.0.0.1:8972"
```
在config.go中增加读取grpc配置的函数
```
type Config struct {
	viper *viper.Viper
	App
	Grpc
}

type Grpc struct {
	Addr string
}

func (c *Config) LoadGrpcConfig() {
	grpc := Grpc{}
	grpc.Addr = c.viper.GetString("grpc.addr")
	c.Grpc = grpc
}
```



