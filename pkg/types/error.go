package types

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Info    string `json:"info"`
}

func (e *Error) Error() string {
	if e.Info == "" {
		return e.Message
	}
	return e.Message + ":" + e.Info
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}
	if er, ok := err.(*Error); ok {
		return er.Code == e.Code
	}

	return false
}

// Join 将err加入到e中
func (e *Error) Join(err error) error {
	return errors.Join(e, err)
}

func (e *Error) SetInfo(i any) *Error {
	if i == nil {
		return e
	}
	if err, ok := i.(error); ok {
		e.Info = err.Error()
	} else if s, ok := i.(string); ok {
		e.Info = s
	} else {
		e.Info = fmt.Sprint(i)
	}
	return e
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
