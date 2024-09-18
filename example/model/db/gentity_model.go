// Code generated by gentity. DO NOT EDIT.

package db

import (
	"time"
)

// tablename company
type Company struct {
	Id               uint64    `json:"id" db:"'id' pk"`                              // 授权方企业本地id
	Platform         string    `json:"platform" db:"'platform'"`                     // 授权方企业的平台方
	CorpId           string    `json:"corp_id" db:"'corp_id'"`                       // 平台授权企业id
	CorpType         int8      `json:"corp_type" db:"'corp_type'"`                   // 企业类型
	FullCorpName     string    `json:"full_corp_name" db:"'full_corp_name'"`         // 企业全称
	CorpType2        int8      `json:"corp_type2" db:"'corp_type2'"`                 // 0 是普通组织\n1 是项目\n2是圈子\n3没有业务表现形式\n4是自建班级群\n10是敏捷组织\n11是培训群敏捷组织
	CorpName         string    `json:"corp_name" db:"'corp_name'"`                   // 企业简称
	Industry         string    `json:"industry" db:"'industry'"`                     // 行业类型
	IsAuthenticated  bool      `json:"is_authenticated" db:"'is_authenticated'"`     // 是否认证
	LicenseCode      string    `json:"license_code" db:"'license_code'"`             // 序列号
	CorpLogoUrl      string    `json:"corp_logo_url" db:"'corp_logo_url'"`           // 企业logo
	InviteUrl        string    `json:"invite_url" db:"'invite_url'"`                 // 企业邀请链接
	InviteCode       string    `json:"invite_code" db:"'invite_code'"`               // 邀请码，只有自己邀请的企业才会返回邀请码，可用该邀请码统计不同渠道的拉新，否则值为空字符串
	IsEcologicalCorp bool      `json:"is_ecological_corp" db:"'is_ecological_corp'"` //
	AuthLevel        int8      `json:"auth_level" db:"'auth_level'"`                 // 企业认证等级：\n\n0：未认证\n1：高级认证\n2：中级认证\n3：初级认证
	AuthChannel      string    `json:"auth_channel" db:"'auth_channel'"`             // 渠道码
	AuthChannelType  string    `json:"auth_channel_type" db:"'auth_channel_type'"`   // 渠道类型。为了避免渠道码重复，可与渠道码共同确认渠道。可能为空，非空时当前只有满天星类型，值为STAR_ACTIVITY
	State            int8      `json:"state" db:"'state'"`                           // 系统状态：-1：已删除；0：禁用；1：可用
	StateTime        time.Time `json:"state_time" db:"'state_time'"`                 // 系统状态时间
	CreatedTime      time.Time `json:"created_time" db:"'created_time'"`             // 创建时间
}
