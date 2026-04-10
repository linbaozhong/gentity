package err

import (
	"errors"
	"strings"
)

// ErrorType 错误类型
type ErrorType int

const (
	// ErrorTypeUnknown 未知错误
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeDB 数据库错误
	ErrorTypeDB
	// ErrorTypeValidation 验证错误
	ErrorTypeValidation
	// ErrorTypePermission 权限错误
	ErrorTypePermission
	// ErrorTypeNotFound 未找到错误
	ErrorTypeNotFound
)

// AppError 应用错误
type AppError struct {
	Type    ErrorType
	Op      string // 操作名称
	Err     error  // 原始错误
	Message string // 用户可见的消息
}

func (e *AppError) Error() string {
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
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// New 创建新错误
func New(typ ErrorType, op string, err error, message string) *AppError {
	return &AppError{
		Type:    typ,
		Op:      op,
		Err:     err,
		Message: message,
	}
}

// NewDB 创建数据库错误
func NewDB(op string, err error) *AppError {
	return New(ErrorTypeDB, op, err, "数据库操作失败")
}

// NewValidation 创建验证错误
func NewValidation(op string, err error, message string) *AppError {
	return New(ErrorTypeValidation, op, err, message)
}

// NewPermission 创建权限错误
func NewPermission(op string, err error, message string) *AppError {
	return New(ErrorTypePermission, op, err, message)
}

// NewNotFound 创建未找到错误
func NewNotFound(op string, err error, message string) *AppError {
	return New(ErrorTypeNotFound, op, err, message)
}

// Wrap 包装错误，添加操作信息
func Wrap(err error, op string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Type: ErrorTypeUnknown,
		Op:   op,
		Err:  err,
	}
}

// WrapDB 包装数据库错误
func WrapDB(err error, op string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Type: ErrorTypeDB,
		Op:   op,
		Err:  err,
	}
}

// GetOp 获取操作名称
func GetOp(err error) string {
	var appErr *AppError
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Op
	}
	return ""
}

// As 检查错误是否为 AppError 类型
func As(err error, target **AppError) bool {
	return errors.As(err, target)
}

// IsDBError 检查是否为数据库错误
func IsDBError(err error) bool {
	var appErr *AppError
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Type == ErrorTypeDB
	}
	return false
}

// IsValidationError 检查是否为验证错误
func IsValidationError(err error) bool {
	var appErr *AppError
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Type == ErrorTypeValidation
	}
	return false
}

// IsPermissionError 检查是否为权限错误
func IsPermissionError(err error) bool {
	var appErr *AppError
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Type == ErrorTypePermission
	}
	return false
}

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	var appErr *AppError
	if ok := As(err, &appErr); ok && appErr != nil {
		return appErr.Type == ErrorTypeNotFound
	}
	return false
}
