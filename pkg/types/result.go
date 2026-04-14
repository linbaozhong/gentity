// @Title api响应数据对象
// @Description api接口返回的数据对象
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package types

import (
	"bytes"
	"sync"
)

// JsonResult API响应数据结构
type JsonResult struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"msg"`
}

var resultPool = sync.Pool{
	New: func() any {
		return &JsonResult{
			// Data: []int{},
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
	resultPool.Put(j)
}

func (j JsonResult) MarshalJSON() ([]byte, error) {
	var _buf = bytes.NewBuffer(nil)
	_buf.WriteByte('{')
	_buf.WriteString(`"code":` + Marshal(j.Code) + `,`)
	_buf.WriteString(`"data":` + Marshal(j.Data) + `,`)
	_buf.WriteString(`"msg":` + Marshal(j.Message))
	_buf.WriteByte('}')
	return _buf.Bytes(), nil
}
