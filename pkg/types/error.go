package types

import (
	"fmt"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Info    string `json:"info"`
}

func (e Error) Error() string {
	if e.Info == "" {
		return e.Message
	}
	return e.Message + ":" + e.Info
}

func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}
	if er, ok := err.(Error); ok {
		return er.Code == e.Code
	}
	return false
}

func (e Error) SetInfo(i any) Error {
	if i == nil {
		return e
	}
	ee := e
	if err, ok := i.(error); ok {
		ee.Info = err.Error()
	} else if s, ok := i.(string); ok {
		ee.Info = s
	} else {
		ee.Info = fmt.Sprint(i)
	}
	return ee
}

func NewError(code int, message string) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}
