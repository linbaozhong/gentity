// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package tblapp

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id      = dialect.Field{Name: "id", Json: "id", Table: "app", Type: "types.BigInt"}
	Arch    = dialect.Field{Name: "arch", Json: "arch", Table: "app", Type: "types.String"}
	Version = dialect.Field{Name: "version", Json: "version", Table: "app", Type: "types.String"}
	Url     = dialect.Field{Name: "url", Json: "url", Table: "app", Type: "types.String"}
	State   = dialect.Field{Name: "state", Json: "state", Table: "app", Type: "types.Int8"}
	Force   = dialect.Field{Name: "force", Json: "force", Table: "app", Type: "types.Int8"}
	// 主键
	PrimaryKey = Id

	// 可写列
	WritableFields = []dialect.Field{
		Id,
		Arch,
		Version,
		Url,
		State,
		Force,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		Arch,
		Version,
		Url,
		State,
		Force,
	}
)