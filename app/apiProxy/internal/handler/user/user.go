package user

import (
	"context"
	"giga"
	"log"
	"net/http"
	"time"

	"apiProxy/internal/service/pb"
)

type HandlerUser struct {
}

func (h *HandlerUser) UserRegister(c *giga.Context) {
	c.JSON(http.StatusOK, giga.H{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
		"age":      c.PostForm("age"),
		"mobile":   c.PostForm("mobile"),
	})
}

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
