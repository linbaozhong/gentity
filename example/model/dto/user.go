package dto

import (
	"github.com/linbaozhong/gentity/pkg/types"
	"time"
)

// @checker
type UserRegisterReq struct {
	ID       *types.BigInt     `json:"id" valid:"required"`
	UserName *int              `json:"user_name" url:"user_name" form:"user_name" valid:"required,range(10|25),min(11)"`
	Password *types.String     `json:"password" url:"password" form:"password" valid:"required"`
	Email    *string           `json:"email" url:"email" form:"email" valid:"email~Email格式错误,required"`
	Content  *types.String     `valid:"runelength(50|100),in(\"1.1.1.1\"|\"2.2.2.2\")"`
	Age      *int8             `valid:"required,in(18|60)"`
	AuthorIP *string           `valid:"ipv4,minstringlength(10)"`
	Date     *time.Time        `valid:"required"`
	Get      *UserRegisterResp `valid:"required" json:"get"`
	Amount   *types.Money      `valid:"max(100)"`
}

// @response
type UserRegisterResp struct {
	UserID   *uint64 `json:"user_id"`
	UserName *string `json:"user_name"`
	Email    *string `json:"email,omitempty"`
}

type Visitor struct {
	Name *string `json:"name"`
}
