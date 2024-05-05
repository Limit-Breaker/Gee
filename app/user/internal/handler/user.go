package handler

import (
	"context"
	"log"
	"user/internal/service/pb"
)

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
