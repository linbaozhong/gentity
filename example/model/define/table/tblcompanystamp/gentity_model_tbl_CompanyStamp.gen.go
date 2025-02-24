// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package tblcompanystamp

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id          = dialect.Field{Name: "id", Json: "id", Table: "company_stamp", Type: "types.Money"}
	CompanyId   = dialect.Field{Name: "company_id", Json: "company_id", Table: "company_stamp", Type: "types.Money"}
	Url         = dialect.Field{Name: "url", Json: "url", Table: "company_stamp", Type: "types.String"}
	Genre       = dialect.Field{Name: "genre", Json: "genre", Table: "company_stamp", Type: "types.Int8"}
	IsDefault   = dialect.Field{Name: "is_default", Json: "is_default", Table: "company_stamp", Type: "types.Int8"}
	Creator     = dialect.Field{Name: "creator", Json: "creator", Table: "company_stamp", Type: "types.BigInt"}
	CreatorName = dialect.Field{Name: "creator_name", Json: "creator_name", Table: "company_stamp", Type: "types.String"}
	Department  = dialect.Field{Name: "department", Json: "department", Table: "company_stamp", Type: "types.String"}
	Position    = dialect.Field{Name: "position", Json: "position", Table: "company_stamp", Type: "types.String"}
	State       = dialect.Field{Name: "state", Json: "state", Table: "company_stamp", Type: "types.Int8"}
	Status      = dialect.Field{Name: "status", Json: "status", Table: "company_stamp", Type: "types.Int8"}
	Ctime       = dialect.Field{Name: "ctime", Json: "ctime", Table: "company_stamp", Type: "types.Time"}
	Utime       = dialect.Field{Name: "utime", Json: "utime", Table: "company_stamp", Type: "types.Time"}
	// 主键
	PrimaryKey = Id

	// 可写列
	WritableFields = []dialect.Field{
		Id,
		CompanyId,
		Url,
		Genre,
		IsDefault,
		Creator,
		CreatorName,
		Department,
		Position,
		State,
		Status,
		Ctime,
		Utime,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		CompanyId,
		Url,
		Genre,
		IsDefault,
		Creator,
		CreatorName,
		Department,
		Position,
		State,
		Status,
		Ctime,
		Utime,
	}
)
