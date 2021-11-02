package test

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"homeShoppingMall/userServer/dao/redis"
	"homeShoppingMall/userServer/handler"
	"homeShoppingMall/userServer/proto"
	"log"
	"testing"
)

var conn *grpc.ClientConn
var err error
var userClient proto.UserSeverClient

func GrpcInit() {
	conn, err = grpc.Dial(
		"consul://120.26.67.141:8500/userServer?wait=14s",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 生成grpc的client并调用接口
	userClient = proto.NewUserSeverClient(conn)
}

func TestGetAllUserList(t *testing.T) {
	GrpcInit()
	response, err := userClient.GetAllUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 1,
	})
	fmt.Println(response)
	if err != nil {
		fmt.Printf("TestAllUserList userClient.GetAllUserList failed, err: %v", err)
	}
	for _, user := range response.UserList {
		fmt.Println(1)
		checkResponse, err := userClient.CheckPassword(context.Background(), &proto.PasswordInfo{
			Password:          "123456",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			zap.L().Error("test userClient.CheckPassword failed", zap.Error(err))
			return
		}
		fmt.Println(checkResponse.Ok)
	}
}

func TestGetUserInfoByMobile(t *testing.T) {
	GrpcInit()
	response, err := userClient.GetUserInfoByMobile(context.Background(), &proto.MobileInfo{
		Mobile: "18096777821",
	})
	if err != nil {
		zap.L().Error("test userClient.GetUserInfoByMobilefailed", zap.Error(err))
	}
	fmt.Println(response)
}

func TestGetUserInfoByUserID(t *testing.T) {
	GrpcInit()
	response, err := userClient.GetUserInfoByUserID(context.Background(), &proto.UserID{
		UserID: 100087493478584320,
	})
	if err != nil {
		zap.L().Error("test userClient.TestGetUserInfoByUserID failed", zap.Error(err))
	}
	fmt.Println(response)
}

func TestCreateUser(t *testing.T) {
	GrpcInit()
	response, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		UserID:   0,
		Mobile:   "18096777833",
		Password: "123456",
	})
	if err != nil {
		zap.L().Error("test userClient.TestGetUserInfoByUserID failed", zap.Error(err))
	}
	fmt.Println("success")
	fmt.Println(response)
}

func TestCheckAuthCode(t *testing.T) {
	GrpcInit()
	str := handler.RandChar()
	redis.Rdb.Set("18096777815", str, 6000)
	ok, err := userClient.CheckAuthCode(context.Background(), &proto.AuthCodeInfo{
		Mobile:       "18096777815",
		UserAuthCode: str,
	})
	if err != nil {
		print(111)
		return
	}
	fmt.Println(ok.Ok)
}
