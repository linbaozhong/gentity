// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package tblcompany

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id               = dialect.Field{Name: "id", Json: "id", Table: "company", Type: "types.BigInt"}
	LongName         = dialect.Field{Name: "long_name", Json: "long_name", Table: "company", Type: "types.String"}
	ShortName        = dialect.Field{Name: "short_name", Json: "short_name", Table: "company", Type: "types.String"}
	Address          = dialect.Field{Name: "address", Json: "address", Table: "company", Type: "types.String"}
	Email            = dialect.Field{Name: "email", Json: "email", Table: "company", Type: "types.String"}
	ContactName      = dialect.Field{Name: "contact_name", Json: "contact_name", Table: "company", Type: "types.String"}
	ContactTelephone = dialect.Field{Name: "contact_telephone", Json: "contact_telephone", Table: "company", Type: "types.String"}
	ContactMobile    = dialect.Field{Name: "contact_mobile", Json: "contact_mobile", Table: "company", Type: "types.String"}
	ContactEmail     = dialect.Field{Name: "contact_email", Json: "contact_email", Table: "company", Type: "types.String"}
	LegalName        = dialect.Field{Name: "legal_name", Json: "legal_name", Table: "company", Type: "types.String"}
	Creator          = dialect.Field{Name: "creator", Json: "creator", Table: "company", Type: "types.BigInt"}
	State            = dialect.Field{Name: "state", Json: "state", Table: "company", Type: "types.Int8"}
	Status           = dialect.Field{Name: "status", Json: "status", Table: "company", Type: "types.Int8"}
	Ctime            = dialect.Field{Name: "ctime", Json: "ctime", Table: "company", Type: "types.Time"}
	Utime            = dialect.Field{Name: "utime", Json: "utime", Table: "company", Type: "types.Time"}
	// 主键
	PrimaryKey = Id

	// 可写列
	WritableFields = []dialect.Field{
		LongName,
		ShortName,
		Address,
		Email,
		ContactName,
		ContactTelephone,
		ContactMobile,
		ContactEmail,
		LegalName,
		Creator,
		State,
		Status,
		Ctime,
		Utime,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		LongName,
		ShortName,
		Address,
		Email,
		ContactName,
		ContactTelephone,
		ContactMobile,
		ContactEmail,
		LegalName,
		Creator,
		State,
		Status,
		Ctime,
		Utime,
	}
)
