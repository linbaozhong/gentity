// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package tbluser_log

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id         = dialect.Field{Name: "id", Json: "id", Table: "user_log", Type: "types.BigInt"}
	UserId     = dialect.Field{Name: "user_id", Json: "user_id", Table: "user_log", Type: "types.BigInt"}
	LoginTime  = dialect.Field{Name: "login_time", Json: "login_time", Table: "user_log", Type: "types.Time"}
	Device     = dialect.Field{Name: "device", Json: "device", Table: "user_log", Type: "types.String"}
	Os         = dialect.Field{Name: "os", Json: "os", Table: "user_log", Type: "types.String"}
	OsVersion  = dialect.Field{Name: "os_version", Json: "os_version", Table: "user_log", Type: "types.String"}
	AppName    = dialect.Field{Name: "app_name", Json: "app_name", Table: "user_log", Type: "types.String"}
	AppVersion = dialect.Field{Name: "app_version", Json: "app_version", Table: "user_log", Type: "types.String"}
	Ip         = dialect.Field{Name: "ip", Json: "ip", Table: "user_log", Type: "types.String"}
	// 主键
	PrimaryKey = Id

	// 可写列
	WritableFields = []dialect.Field{
		UserId,
		LoginTime,
		Device,
		Os,
		OsVersion,
		AppName,
		AppVersion,
		Ip,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		UserId,
		LoginTime,
		Device,
		Os,
		OsVersion,
		AppName,
		AppVersion,
		Ip,
	}
)
