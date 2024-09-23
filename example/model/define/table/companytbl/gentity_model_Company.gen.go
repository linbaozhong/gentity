// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package companytbl

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
	Id               = dialect.Field{Name: "id", Table: "company", Type: "uint64"}
	Platform         = dialect.Field{Name: "platform", Table: "company", Type: "string"}
	CorpId           = dialect.Field{Name: "corp_id", Table: "company", Type: "string"}
	CorpType         = dialect.Field{Name: "corp_type", Table: "company", Type: "int8"}
	FullCorpName     = dialect.Field{Name: "full_corp_name", Table: "company", Type: "string"}
	CorpType2        = dialect.Field{Name: "corp_type2", Table: "company", Type: "int8"}
	CorpName         = dialect.Field{Name: "corp_name", Table: "company", Type: "string"}
	Industry         = dialect.Field{Name: "industry", Table: "company", Type: "string"}
	IsAuthenticated  = dialect.Field{Name: "is_authenticated", Table: "company", Type: "bool"}
	LicenseCode      = dialect.Field{Name: "license_code", Table: "company", Type: "string"}
	CorpLogoUrl      = dialect.Field{Name: "corp_logo_url", Table: "company", Type: "string"}
	InviteUrl        = dialect.Field{Name: "invite_url", Table: "company", Type: "string"}
	InviteCode       = dialect.Field{Name: "invite_code", Table: "company", Type: "string"}
	IsEcologicalCorp = dialect.Field{Name: "is_ecological_corp", Table: "company", Type: "bool"}
	AuthLevel        = dialect.Field{Name: "auth_level", Table: "company", Type: "int8"}
	AuthChannel      = dialect.Field{Name: "auth_channel", Table: "company", Type: "string"}
	AuthChannelType  = dialect.Field{Name: "auth_channel_type", Table: "company", Type: "string"}
	State            = dialect.Field{Name: "state", Table: "company", Type: "int8"}
	StateTime        = dialect.Field{Name: "state_time", Table: "company", Type: "time.Time"}
	CreatedTime      = dialect.Field{Name: "created_time", Table: "company", Type: "time.Time"}
	// 主键
	PrimaryKeys = []dialect.Field{
		Id,
	}

	// 可写列
	WritableFields = []dialect.Field{
		Id,
		Platform,
		CorpId,
		CorpType,
		FullCorpName,
		CorpType2,
		CorpName,
		Industry,
		IsAuthenticated,
		LicenseCode,
		CorpLogoUrl,
		InviteUrl,
		InviteCode,
		IsEcologicalCorp,
		AuthLevel,
		AuthChannel,
		AuthChannelType,
		State,
		StateTime,
		CreatedTime,
	}
	// 可读列
	ReadableFields = []dialect.Field{
		Id,
		Platform,
		CorpId,
		CorpType,
		FullCorpName,
		CorpType2,
		CorpName,
		Industry,
		IsAuthenticated,
		LicenseCode,
		CorpLogoUrl,
		InviteUrl,
		InviteCode,
		IsEcologicalCorp,
		AuthLevel,
		AuthChannel,
		AuthChannelType,
		State,
		StateTime,
		CreatedTime,
	}
)
