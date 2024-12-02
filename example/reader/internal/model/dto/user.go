package dto

import (
	"reader/internal/constant/err"
)

type UserRegisterReq struct {
	UserName string `json:"user_name" url:"user_name" form:"user_name"`
	Password string `json:"password" url:"password" form:"password"`
	Email    string `json:"email" url:"email" form:"email"`
}

func (u *UserRegisterReq) Check() error {
	if u.UserName == "" {
		return err.ErrUserName
	}
	if u.Password == "" {
		return err.ErrPassword
	}
	return nil
}

type UserRegisterResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}

type GetUserReq struct {
	UserID uint64 `json:"user_id" url:"user_id" form:"user_id"`
}

func (u *GetUserReq) Check() error {
	if u.UserID == 0 {
		return err.ErrUserID
	}
	return nil
}

type GetUserResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}
