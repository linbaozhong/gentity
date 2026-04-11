package types

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	Op      string `json:"-"`    // 操作名称
	Err     error  `json:"-"`    // 原始错误
	Code    int    `json:"code"` // 用户可见错误码
	Message string `json:"msg"`  // 用户可见的消息
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	var parts []string
	if e.Op != "" {
		parts = append(parts, e.Op)
	}
	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}
	return strings.Join(parts, ": ")
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Join 将err加入到e中
func (e *Error) Join(err error) error {
	return errors.Join(e, err)
}

// SetMessage 设置消息
func (e *Error) SetMessage(i any) *Error {
	if i == nil {
		return e
	}
	if err, ok := i.(error); ok {
		e.Message = err.Error()
	} else if s, ok := i.(string); ok {
		e.Message = s
	} else {
		e.Message = fmt.Sprint(i)
	}
	return e
}

// SetOp 设置操作名称
func (e *Error) SetOp(op string) *Error {
	e.Op = op
	return e
}

// NewError 创建一个新的错误
func NewError(code int, message string, ops ...string) *Error {
	e := &Error{
		Code:    code,
		Message: message,
	}
	if len(ops) > 0 {
		e.Op = ops[0]
	}
	return e
}
