package dto

// checker
type UserRegisterReq struct {
	UserName string `json:"user_name" url:"user_name" form:"user_name" valid:"required"`
	Password string `json:"password" url:"password" form:"password" valid:"required"`
	Email    string `json:"email" url:"email" form:"email" valid:"required,email"`
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
