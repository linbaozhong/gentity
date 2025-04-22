package dto

import (
	"github.com/linbaozhong/gentity/pkg/types"
	"time"
)

// checker
type UserRegisterReq struct {
	ID       types.BigInt      `json:"id" valid:"required"`
	UserName int               `json:"user_name" url:"user_name" form:"user_name" valid:"required,range(10|25)"`
	Password types.String      `json:"password" url:"password" form:"password" valid:"required"`
	Email    string            `json:"email" url:"email" form:"email" valid:"email~Email格式错误,required"`
	Content  types.Float64     `valid:"runelength(50|100),required"`
	Age      int8              `valid:"range(18|60),required"`
	AuthorIP int64             `valid:"ipv4"`
	Date     time.Time         `valid:"required"`
	Get      *UserRegisterResp `valid:"required" json:"get"`
	Amount   types.Money       `valid:"required"`
	// 约定访问者类型为Visitor
	// 注意：访问者类型必须实现Visiter接口，否则会报错
	Vis Visitor `json:"vis"`
}

// response
type UserRegisterResp struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email,omitempty"`
}
