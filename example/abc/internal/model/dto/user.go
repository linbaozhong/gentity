package dto

import "abc/internal/constant"

// UserRegisterReq 用户注册请求数据
// 字段标签说明(可选项)：
// json：用于json序列化和反序列化,解析Content-Type为application/json时使用
// url：用于url参数解析,解析URL Query时使用
// form：用于表单参数解析,解析Content-Type为application/x-www-form-urlencoded和multipart/form-data时使用
// param：用于url动态路径参数解析
// valid：用于数据校验
type UserRegisterReq struct {
	UserName string `json:"user_name" url:"user_name" form:"user_name" param:"user_name" valid:"required"`
	Password string `json:"password" url:"password" form:"password" param:"password" valid:"required"`
	Email    string `json:"email" url:"email" form:"email" param:"email" valid:"email~邮箱格式错误"`
}

func (u *UserRegisterReq) Check() error {
	if u.UserName == "" {
		return constant.ErrUserName
	}
	if u.Password == "" {
		return constant.ErrPassword
	}
	return nil
}

type UserRegisterResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}

//
type GetUserReq struct {
	UserID uint64 `json:"user_id" url:"user_id" form:"user_id"`
}

func (u *GetUserReq) Check() error {
	if u.UserID == 0 {
		return constant.ErrUserID
	}
	return nil
}

type GetUserResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}
