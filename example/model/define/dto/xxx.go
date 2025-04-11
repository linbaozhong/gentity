package dto

import "github.com/linbaozhong/gentity/pkg/types"

// checker
type GetUserReq struct {
	UserID float64 `json:"user_id" url:"user_id" form:"user_id" valid:"required"`
}

// response
type GetUserResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}

// checker
type DispatchCompanyAddReq struct {
	Name types.String `json:"name" url:"name" form:"name" valid:"required"`
}
