package user

import (
	"context"
	"fmt"
	"reader/internal/model/dto"
)

func UserRegister(c context.Context, in *dto.UserRegisterReq, out *dto.UserRegisterResp) error {
	// todo: 在这里做用户注册，返回用户信息
	fmt.Println("UserRegister:", in)

	out.UserID = 12345678
	out.UserName = "哈利蔺特"

	return nil
}

func GetUser(c context.Context, in *dto.GetUserReq, out *dto.GetUserResp) error {
	// todo: 在这里做用户查询，返回用户信息
	fmt.Println("UserRegister:", in)

	out.UserID = 12345678
	out.UserName = "哈利蔺特"

	return nil
}
