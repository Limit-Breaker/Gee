## grpc服务端

在app目录下新建一个目录user，这将是rpc的服务端
目前项目的目录结构如下
```
goProject
├─app
│  ├─apiProxy
│  │  ├─config
│  │  ├─internal
│  │  │  ├─handler
│  │  │  │  ├─greet
│  │  │  │  └─user
│  │  │  └─service
│  │  │      └─pb
│  │  ├─middleware
│  │  └─router
│  └─user
```

#### 服务端接口定义
客户端和服务端使用相同的接口定义文件 <br>
在user/internal/service目录下新建pb目录，将apiProxy/internal/service/pb目录下的userService.pb.go和userService_grpc.pb.go文件拷贝到pb目录下

#### 服务端处理函数
服务端需要实现pb.UnimplementedUserServiceServer接口中的GetCaptcha方法<br>
这里简单实现下逻辑
```
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) GetCaptcha(ctx context.Context, rq *pb.GetCaptchaRequest) (*pb.GetCaptchaResponse, error) {
	// 校验参数 TODO
	mobile := rq.Mobile
	// 生成验证码
	code := "123456"
	//调用短信平台 TODO
	log.Printf("往手机: %s 发送验证码[%s]", mobile, code)
	return &pb.GetCaptchaResponse{Code: code}, nil
}
```

#### 服务端的启动
在user目录下新建main.go函数, 完成服务端的启动
```
func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", config.DefaultConfig.Grpc.Addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer() // 创建gRPC服务器
	defer s.Stop()
	// 在gRPC服务端注册服务
	pb.RegisterUserServiceServer(s, &handler.UserServiceServer{})
	log.Printf("start user rpc server")
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		log.Printf("failed to serve: %v", err)
		return
	}
}
```

## 验证
在user目录下go run main.go
在apiProxy目录下 go run main.go
发送请求
```
// POST http://localhost:8080/user/login
// username: makabaka
// password: 123456
// mobile: 13211111111

{
    "code": "123456"
}
```
