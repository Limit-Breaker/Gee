package router

import (
	"apiProxy/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"giga"

	"apiProxy/internal/service/pb"
)

var UserServiceClient pb.UserServiceClient

type Router interface {
	Route(r *giga.Engine)
}

type RegisterRouter struct {
	engine *giga.Engine
}

func NewRegister(r *giga.Engine) *RegisterRouter {
	return &RegisterRouter{
		engine: r,
	}
}

func (register *RegisterRouter) AddRoute(rt Router) {
	rt.Route(register.engine)
}

func InitRouter(r *giga.Engine) {
	register := NewRegister(r)
	register.AddRoute(&RouterGreet{})
	register.AddRoute(&RouterUser{})
}

func InitRpcClient() {
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial(config.DefaultConfig.Grpc.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	UserServiceClient = pb.NewUserServiceClient(conn)
}
