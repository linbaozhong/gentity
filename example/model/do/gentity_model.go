// Code generated by gentity. DO NOT EDIT.

package do

import (
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
	"github.com/linbaozhong/gentity/pkg/types"
)

// tablename app
type App struct {
	atype.AceModel
	Id      types.BigInt    `json:"id,omitempty" db:"'id' pk"`        //
	Arch    types.AceString `json:"arch,omitempty" db:"'arch'"`       // 操作系统架构
	Version types.AceString `json:"version,omitempty" db:"'version'"` // 版本号
	Url     types.AceString `json:"url,omitempty" db:"'url'"`         // 应用下载地址
	State   types.AceInt8   `json:"state,omitempty" db:"'state'"`     //
	Force   types.AceInt8   `json:"force,omitempty" db:"'force'"`     //
}

// tablename user
type User struct {
	atype.AceModel
	Id    types.BigInt    `json:"id,omitempty" db:"'id' pk auto"` //
	Uuid  types.AceString `json:"uuid,omitempty" db:"'uuid'"`     // 用户识别码
	Ctime types.AceTime   `json:"ctime,omitempty" db:"'ctime'"`   //
}

// tablename user_log
type UserLog struct {
	atype.AceModel
	Id         types.BigInt    `json:"id,omitempty" db:"'id' pk auto"`           //
	UserId     types.BigInt    `json:"user_id,omitempty" db:"'user_id'"`         //
	LoginTime  types.AceTime   `json:"login_time,omitempty" db:"'login_time'"`   // 登录时间
	Device     types.AceString `json:"device,omitempty" db:"'device'"`           // 登录终端参数
	Os         types.AceString `json:"os,omitempty" db:"'os'"`                   //
	OsVersion  types.AceString `json:"os_version,omitempty" db:"'os_version'"`   //
	AppName    types.AceString `json:"app_name,omitempty" db:"'app_name'"`       //
	AppVersion types.AceString `json:"app_version,omitempty" db:"'app_version'"` //
	Ip         types.AceString `json:"ip,omitempty" db:"'ip'"`                   // ip地址
}
