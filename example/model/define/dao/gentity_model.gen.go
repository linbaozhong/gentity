// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/cachego"
)

var (
	Where   = ace.Where
	Or      = ace.Or
	And     = ace.And
	getExec = ace.GetExec
	Create  = ace.Create
	Delete  = ace.Delete
	Update  = ace.Update
	Select  = ace.Select
	Sets    = ace.Sets
	Order   = ace.Order
	Asc     = ace.Asc
	Desc    = ace.Desc
)

// 生成综合条件cond缓存key
func GetCondHashKey(key any) string {
	return cachego.GetHashKey("c:", key)
}

// 生成id缓存key
func GetIdHashKey(key any) string {
	return cachego.GetHashKey("i:", key)
}

// 生成ids缓存key
func GetIdsHashKey(key any) string {
	return cachego.GetHashKey("s:", key)
}
