// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package tblcompanyman

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id        = dialect.Field{Name: "id", Json: "id", Table: "company_man", Type: "types.BigInt"}
	AccountId = dialect.Field{Name: "account_id", Json: "account_id", Table: "company_man", Type: "types.BigInt"}
	CompanyId = dialect.Field{Name: "company_id", Json: "company_id", Table: "company_man", Type: "types.BigInt"}
	RealName  = dialect.Field{Name: "real_name", Json: "real_name", Table: "company_man", Type: "types.String"}
	Email     = dialect.Field{Name: "email", Json: "email", Table: "company_man", Type: "types.String"}
	Roles     = dialect.Field{Name: "roles", Json: "roles", Table: "company_man", Type: "types.String"}
	State     = dialect.Field{Name: "state", Json: "state", Table: "company_man", Type: "types.Int8"}
	Ctime     = dialect.Field{Name: "ctime", Json: "ctime", Table: "company_man", Type: "types.Time"}
	Utime     = dialect.Field{Name: "utime", Json: "utime", Table: "company_man", Type: "types.Time"}
	// 主键
	PrimaryKey = Id

	// 可写列
	WritableFields = []dialect.Field{
		AccountId,
		CompanyId,
		RealName,
		Email,
		Roles,
		State,
		Ctime,
		Utime,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		AccountId,
		CompanyId,
		RealName,
		Email,
		Roles,
		State,
		Ctime,
		Utime,
	}
)
