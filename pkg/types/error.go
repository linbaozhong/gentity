package types

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorType 错误类型
type ErrorType int

const (
	// ErrorTypeUnknown 未知错误
	ErrorTypeUnknown ErrorType = -100 * (iota + 1)
	// ErrorTypeDB 数据库错误
	ErrorTypeDB
	// ErrorTypeValidation 验证错误
	ErrorTypeValidation
	// ErrorTypePermission 权限错误
	ErrorTypePermission
	// ErrorTypeNotFound 未找到错误
	ErrorTypeNotFound
	// ErrorTypeParam 参数错误
	ErrorTypeParam
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

// newError 创建新错误
func newError(typ ErrorType, op string, err error, messages ...string) *Error {
	e := &Error{
		Code: int(typ),
		Op:   op,
		Err:  err,
	}
	if len(messages) > 0 {
		e.Message = messages[0]
	}
	return e
}

// NewDB 创建数据库错误
func NewDB(op string, err error, message ...string) *Error {
	return newError(ErrorTypeDB, op, err, message...)
}

// NewValidation 创建验证错误
func NewValidation(op string, err error, message ...string) *Error {
	return newError(ErrorTypeValidation, op, err, message...)
}

// NewPermission 创建权限错误
func NewPermission(op string, err error, message ...string) *Error {
	return newError(ErrorTypePermission, op, err, message...)
}

// NewNotFound 创建未找到错误
func NewNotFound(op string, err error, message string) *Error {
	return newError(ErrorTypeNotFound, op, err, message)
}

// NewParamm 创建参数错误
func NewParam(op string, err error, message string) *Error {
	return newError(ErrorTypeParam, op, err, message)
}

// Wrap 包装错误，添加操作信息
func Wrap(err error, op string) *Error {
	if err == nil {
		return nil
	}
	return newError(ErrorTypeUnknown, op, err)
}

// WrapDB 包装数据库错误
func WrapDB(err error, op string) *Error {
	if err == nil {
		return nil
	}
	return newError(ErrorTypeDB, op, err)
}

// As 检查错误是否为 Error 类型
func As(err error, target **Error) bool {
	return errors.As(err, target)
}

// IsDBError 检查是否为数据库错误
func IsDBError(err error) bool {
	var appErr *Error
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Code == int(ErrorTypeDB)
	}
	return false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	var appErr *Error
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Code == int(ErrorTypeValidation)
	}
	return false
}

// IsPermissionError 检查是否为权限错误
func IsPermissionError(err error) bool {
	var appErr *Error
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Code == int(ErrorTypePermission)
	}
	return false
}

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	var appErr *Error
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Code == int(ErrorTypeNotFound)
	}
	return false
}

// IsParamError 检查是否为参数错误
func IsParamError(err error) bool {
	var appErr *Error
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Code == int(ErrorTypeParam)
	}
	return false
}
