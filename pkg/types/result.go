// @Title api响应数据对象
// @Description api接口返回的数据对象
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package types

import (
	"sync"
)

// JsonResult API响应数据结构
type JsonResult struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"msg"`
	Info    string `json:"info"` // dev消息
}

var resultPool = sync.Pool{
	New: func() any {
		return &JsonResult{
			Data: []int{},
		}
	},
}

// NewResult
func NewResult() *JsonResult {
	return resultPool.Get().(*JsonResult)
}

// Free
func (j *JsonResult) Free() {
	j.Code = 0
	j.Data = nil
	j.Message = ""
	j.Info = ""
	resultPool.Put(j)
}

//
// // SetCode 为响应数据设置code和msg,msg省略,取constant中定义的msg
// func (j *JsonResult) SetCode(e error, msg ...string) *JsonResult {
// 	if err, ok := e.(Error); ok {
// 		j.Code = err.Code
// 		j.Info = err.Info
// 	}
//
// 	if len(msg) > 0 {
// 		j.Message = msg[0]
// 		return j
// 	}
// 	if e != nil {
// 		j.Message = e.Error()
// 	}
// 	return j
// }
//
// func (j *JsonResult) SetData(d any) *JsonResult {
// 	j.Data = d
// 	return j
// }
//
// func (j *JsonResult) SetInfo(info string) *JsonResult {
// 	j.Info = info
// 	return j
// }
