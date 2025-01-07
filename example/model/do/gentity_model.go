package do

import (
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/types"
)

// tablename user_log
type UserLog struct {
	ace.Model
	Id         types.BigInt `json:"id,omitempty" db:"'id' pk auto size(20)"`           //
	UserId     types.BigInt `json:"user_id,omitempty" db:"'user_id' size(20)"`         //
	LoginTime  types.Time   `json:"login_time,omitempty" db:"'login_time'"`            // 登录时间
	Device     types.String `json:"device,omitempty" db:"'device' size(255)"`          // 登录终端参数
	Os         types.String `json:"os,omitempty" db:"'os' size(10)"`                   //
	OsVersion  types.String `json:"os_version,omitempty" db:"'os_version' size(10)"`   //
	AppName    types.String `json:"app_name,omitempty" db:"'app_name' size(10)"`       //
	AppVersion types.String `json:"app_version,omitempty" db:"'app_version' size(10)"` //
	Ip         types.String `json:"ip,omitempty" db:"'ip' size(50)"`                   // ip地址

	// User User `json:"user,omitempty" db:"ref:user_id fk:id"`
}

// tablename app
type App struct {
	ace.Model
	Id      types.BigInt `json:"id,omitempty" db:"'id' pk size(20)"`        //
	Arch    types.String `json:"arch,omitempty" db:"'arch' size(10)"`       // 操作系统架构
	Version types.String `json:"version,omitempty" db:"'version' size(10)"` // 版本号
	Url     types.String `json:"url,omitempty" db:"'url' size(145)"`        // 应用下载地址
	State   types.Int8   `json:"state,omitempty" db:"'state' size(3)"`      //
	Force   types.Int8   `json:"force,omitempty" db:"'force' size(3)"`      //
}

// tablename user
type User struct {
	ace.Model
	Id    types.BigInt `json:"id,omitempty" db:"'id' pk auto size(20)"` //
	Uuid  types.String `json:"uuid,omitempty" db:"'uuid' size(45)"`     // 用户识别码
	Ctime types.Time   `json:"ctime,omitempty" db:"'ctime'"`            //

	// UserLogs []UserLog `json:"user_logs,omitempty" db:"ref:id fk:user_id"` //
}
