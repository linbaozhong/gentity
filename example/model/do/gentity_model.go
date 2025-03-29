package do

import (
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/types"
)

// tablename users
type Users struct {
	pool.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`                //
	From       types.String `json:"from,omitempty" db:"'from' size:10"`               // 微信openid
	OpenId     types.String `json:"open_id,omitempty" db:"'open_id' size:64"`         // 支付宝id
	UnionId    types.String `json:"union_id,omitempty" db:"'union_id' size:64"`       // 微信unionid
	SessionKey types.String `json:"session_key,omitempty" db:"'session_key' size:64"` //
	Mobile     types.String `json:"mobile,omitempty" db:"'mobile' size:20"`           //
	Email      types.String `json:"email,omitempty" db:"'email' size:100"`            //
	Pwd        types.String `json:"pwd,omitempty" db:"'pwd' size:32"`                 //
	Avatar     types.String `json:"avatar,omitempty" db:"'avatar' size:100"`          //
	Nick       types.String `json:"nick,omitempty" db:"'nick' size:50"`               //
	Gender     types.String `json:"gender,omitempty" db:"'gender' size:1"`            //
	State      types.Int8   `json:"state,omitempty" db:"'state' size:3"`              // 状态:1=可用 0=禁用 -1=删除
	Ctime      types.Time   `json:"ctime,omitempty" db:"'ctime'"`                     //
}

// tablename users_info
type UsersInfo struct {
	pool.Model
	Id       types.BigInt `json:"id,omitempty" db:"'id' pk size:20"`            //
	Name     types.String `json:"name,omitempty" db:"'name' size:45"`           // 真实姓名
	IdNumber types.String `json:"id_number,omitempty" db:"'id_number' size:20"` // 身份证号
	Country  types.String `json:"country,omitempty" db:"'country' size:5"`      //
	Province types.String `json:"province,omitempty" db:"'province' size:10"`   //
	City     types.String `json:"city,omitempty" db:"'city' size:10"`           //
	Utime    types.Time   `json:"utime,omitempty" db:"'utime'"`                 //
}
