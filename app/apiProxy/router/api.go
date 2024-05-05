package router

import (
	"apiProxy/middleware"
	"giga"

	"apiProxy/internal/handler/greet"
	"apiProxy/internal/handler/user"
)

type RouterGreet struct {
}

func (g *RouterGreet) Route(r *giga.Engine) {
	greet := greet.HandlerGreet{}
	r.GET("/hello/:name", greet.Hello)
}

type RouterUser struct {
}

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
