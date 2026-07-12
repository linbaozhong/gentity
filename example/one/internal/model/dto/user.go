package dto

// UserRegisterReq 用户注册请求数据
// @request
// 字段标签说明(可选项)：
// json：用于json序列化和反序列化,解析Content-Type为application/json、application/x-www-form-urlencoded和multipart/form-data时使用
// valid：用于数据校验
type UserRegisterReq struct {
	UserName *string `json:"user_name" valid:"required"`
	Password *string `json:"password" valid:"required"`
	Email    *string `json:"email" valid:"email~邮箱格式错误"`
}

// UserRegisterResp 用户注册响应数据
// @response
type UserRegisterResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}

// @request
type GetUserReq struct {
	UserID *uint64 `json:"user_id"`
}

// @response
type GetUserResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}
