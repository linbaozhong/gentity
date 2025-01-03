package dto

import "github.com/linbaozhong/gentity/pkg/types"

// checker
type UserRegisterReq struct {
	UserName int          `json:"user_name" url:"user_name" form:"user_name" valid:"required,range(10|25)"`
	Password types.String `json:"password" url:"password" form:"password" valid:"required"`
	Email    string       `json:"email" url:"email" form:"email" valid:"email~Email格式错误,required"`
	Content  types.String `valid:"runelength(50|100),required"`
	Age      int          `valid:"range(18|60),required"`
	AuthorIP int64        `valid:"ipv4"`
	Date     string       `valid:"-"`
}

type UserRegisterResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}

// checker
type GetUserReq struct {
	UserID uint64 `json:"user_id" url:"user_id" form:"user_id" valid:"required"`
}

type GetUserResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}

// checker
type DispatchCompanyAddReq struct {
	Name string `json:"name" url:"name" form:"name" valid:"required"`
}
