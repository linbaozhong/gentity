// Code generated by gentity. DO NOT EDIT.

package dto

import (
	"bytes"
	"errors"
	"github.com/linbaozhong/gentity/pkg/conv"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
	"github.com/linbaozhong/gentity/pkg/validator"
	"net/http"
)

var _ = conv.Any2String("")
var _ = bytes.NewBuffer(nil)

/*
	--- UserRegisterReq ---
*/
// Init
func (p *UserRegisterReq) Init() error {
	p.ID = types.NilUint64
	p.UserName = types.NilInt
	p.Password = types.NilString
	p.Email = types.NilString
	p.Content = types.NilFloat64
	p.Age = types.NilInt8
	p.AuthorIP = types.NilInt64
	p.Date = types.NilTime
	p.Get = nil
	p.Amount = types.NilInt64

	return nil
}

// Check
func (p *UserRegisterReq) Check() error {
	if p.ID == types.NilUint64 {
		return types.NewError(http.StatusBadRequest, "id is required")
	}
	if p.UserName == types.NilInt {
		return types.NewError(http.StatusBadRequest, "user_name is required")
	}
	if !validator.Range(conv.Any2String(p.UserName), "10\",\"25") {
		return types.NewError(http.StatusBadRequest, "Wrong user_name range")
	}
	if p.Password == types.NilString {
		return types.NewError(http.StatusBadRequest, "password is required")
	}
	if p.Email == types.NilString {
		return types.NewError(http.StatusBadRequest, "email is required")
	}
	if !validator.IsEmail(p.Email) {
		return types.NewError(http.StatusBadRequest, "Email格式错误")
	}
	if p.Content == types.NilFloat64 {
		return types.NewError(http.StatusBadRequest, " is required")
	}
	if !validator.RuneLength(p.Content.String(), "50\",\"100") {
		return types.NewError(http.StatusBadRequest, "Wrong  runelength")
	}
	if p.Age == types.NilInt8 {
		return types.NewError(http.StatusBadRequest, " is required")
	}
	if !validator.Range(conv.Any2String(p.Age), "18\",\"60") {
		return types.NewError(http.StatusBadRequest, "Wrong  range")
	}
	if !validator.IsIPv4(conv.Any2String(p.AuthorIP)) {
		return types.NewError(http.StatusBadRequest, "Wrong  format")
	}
	if p.Date == types.NilTime {
		return types.NewError(http.StatusBadRequest, " is required")
	}
	if p.Get == nil {
		return types.NewError(http.StatusBadRequest, "get is required")
	}
	if p.Amount == types.NilInt64 {
		return types.NewError(http.StatusBadRequest, " is required")
	}
	return nil
}

// UnmarshalJSON
func (p *UserRegisterReq) UnmarshalJSON(data []byte) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
			return
		}
	}()

	ok := gjson.ValidBytes(data)
	if !ok {
		return errors.New("invalid json")
	}
	_result := gjson.ParseBytes(data)
	var e error
	_result.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "id":
			p.ID = types.BigInt(value.Uint())
		case "user_name":
			p.UserName = int(value.Int())
		case "password":
			p.Password = types.String(value.Str)
		case "email":
			p.Email = value.Str
		case "get":
			e = types.Unmarshal(value, &p.Get, func(value gjson.Result) *UserRegisterResp {
				var obj *UserRegisterResp
				e := types.Unmarshal(value, &obj)
				if e != nil {
					panic(e)
				}
				return obj
			}(value))
		}
		if e != nil {
			log.Error(e)
			return false
		}
		return true
	})
	return nil
}

// UnmarshalValues
func (p *UserRegisterReq) UnmarshalValues(m map[string][]string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
			return
		}
	}()

	var e error
	for k, v := range m {
		value := gjson.Result{Type: gjson.String, Raw: v[0], Str: v[0]}
		switch k {
		case "id":
			p.ID = types.BigInt(value.Uint())
		case "user_name":
			p.UserName = int(value.Int())
		case "password":
			p.Password = types.String(value.Str)
		case "email":
			p.Email = value.Str
		case "get":
			e = types.Unmarshal(value, &p.Get, func(value gjson.Result) *UserRegisterResp {
				var obj *UserRegisterResp
				e := types.Unmarshal(value, &obj)
				if e != nil {
					panic(e)
				}
				return obj
			}(value))
		}
		if e != nil {
			log.Error(e)
			return e
		}
	}
	return nil
}

/*
	--- UserRegisterResp ---
*/
// MarshalJSON
func (p *UserRegisterResp) MarshalJSON() ([]byte, error) {
	var _buf = bytes.NewBuffer(nil)
	_buf.WriteByte('{')
	_buf.WriteString(`"user_id":` + types.Marshal(p.UserID) + `,`)
	_buf.WriteString(`"user_name":` + types.Marshal(p.UserName) + `,`)
	if p.Email != "" {
		_buf.WriteString(`"email":` + types.Marshal(p.Email) + `,`)
	}
	if l := _buf.Len(); l > 1 {
		_buf.Truncate(l - 1)
	}
	_buf.WriteByte('}')
	return _buf.Bytes(), nil
}

/*
	--- GetUserReq ---
*/
// Init
func (p *GetUserReq) Init() error {
	p.UserID = types.NilFloat64

	return nil
}

// Check
func (p *GetUserReq) Check() error {
	if p.UserID == types.NilFloat64 {
		return types.NewError(http.StatusBadRequest, "user_id is required")
	}
	return nil
}

// UnmarshalJSON
func (p *GetUserReq) UnmarshalJSON(data []byte) error {

	ok := gjson.ValidBytes(data)
	if !ok {
		return errors.New("invalid json")
	}
	_result := gjson.ParseBytes(data)
	var e error
	_result.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "user_id":
			p.UserID = value.Float()
		}
		if e != nil {
			log.Error(e)
			return false
		}
		return true
	})
	return nil
}

// UnmarshalValues
func (p *GetUserReq) UnmarshalValues(m map[string][]string) error {

	var e error
	for k, v := range m {
		value := gjson.Result{Type: gjson.String, Raw: v[0], Str: v[0]}
		switch k {
		case "user_id":
			p.UserID = value.Float()
		}
		if e != nil {
			log.Error(e)
			return e
		}
	}
	return nil
}

/*
	--- GetUserResp ---
*/
// MarshalJSON
func (p *GetUserResp) MarshalJSON() ([]byte, error) {
	var _buf = bytes.NewBuffer(nil)
	_buf.WriteByte('{')
	_buf.WriteString(`"user_id":` + types.Marshal(p.UserID) + `,`)
	_buf.WriteString(`"user_name":` + types.Marshal(p.UserName) + `,`)
	if p.Email != "" {
		_buf.WriteString(`"email":` + types.Marshal(p.Email) + `,`)
	}
	if l := _buf.Len(); l > 1 {
		_buf.Truncate(l - 1)
	}
	_buf.WriteByte('}')
	return _buf.Bytes(), nil
}

/*
	--- DispatchCompanyAddReq ---
*/
// Init
func (p *DispatchCompanyAddReq) Init() error {
	p.Name = types.NilString

	return nil
}

// Check
func (p *DispatchCompanyAddReq) Check() error {
	if p.Name == types.NilString {
		return types.NewError(http.StatusBadRequest, "name is required")
	}
	return nil
}

// UnmarshalJSON
func (p *DispatchCompanyAddReq) UnmarshalJSON(data []byte) error {

	ok := gjson.ValidBytes(data)
	if !ok {
		return errors.New("invalid json")
	}
	_result := gjson.ParseBytes(data)
	var e error
	_result.ForEach(func(key, value gjson.Result) bool {
		switch key.Str {
		case "name":
			p.Name = types.String(value.Str)
		}
		if e != nil {
			log.Error(e)
			return false
		}
		return true
	})
	return nil
}

// UnmarshalValues
func (p *DispatchCompanyAddReq) UnmarshalValues(m map[string][]string) error {

	var e error
	for k, v := range m {
		value := gjson.Result{Type: gjson.String, Raw: v[0], Str: v[0]}
		switch k {
		case "name":
			p.Name = types.String(value.Str)
		}
		if e != nil {
			log.Error(e)
			return e
		}
	}
	return nil
}
