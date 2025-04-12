// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package db

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/linbaozhong/gentity/example/model/define/table/tblcompanyrole"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

const CompanyRoleTableName = "company_role"

var (
	companyrolePool = pool.New(app.Context, func() any {
		_obj := &CompanyRole{}
		_obj.UUID()
		return _obj
	})
)

func NewCompanyRole() *CompanyRole {
	_obj := companyrolePool.Get().(*CompanyRole)
	return _obj
}

// MarshalJSON
func (p *CompanyRole) MarshalJSON() ([]byte, error) {
	var _buf = bytes.NewBuffer(nil)
	_buf.WriteByte('{')
	if p.Id != 0 {
		_buf.WriteString(`"id":` + types.Marshal(p.Id) + `,`)
	}
	if p.CompanyId != 0 {
		_buf.WriteString(`"company_id":` + types.Marshal(p.CompanyId) + `,`)
	}
	if p.Name != "" {
		_buf.WriteString(`"name":` + types.Marshal(p.Name) + `,`)
	}
	if p.Descr != "" {
		_buf.WriteString(`"descr":` + types.Marshal(p.Descr) + `,`)
	}
	if p.Rules != "" {
		_buf.WriteString(`"rules":` + types.Marshal(p.Rules) + `,`)
	}
	if p.Type != 0 {
		_buf.WriteString(`"type":` + types.Marshal(p.Type) + `,`)
	}
	if p.State != 0 {
		_buf.WriteString(`"state":` + types.Marshal(p.State) + `,`)
	}
	if l := _buf.Len(); l > 1 {
		_buf.Truncate(l - 1)
	}
	_buf.WriteByte('}')
	return _buf.Bytes(), nil
}

// UnmarshalJSON
func (p *CompanyRole) UnmarshalJSON(data []byte) error {

	if !gjson.ValidBytes(data) {
		return errors.New("invalid json")
	}
	_result := gjson.ParseBytes(data)
	_result.ForEach(func(key, value gjson.Result) bool {
		var e error
		switch key.Str {
		case "id":
			p.Id = types.BigInt(value.Uint())
		case "company_id":
			p.CompanyId = types.BigInt(value.Uint())
		case "name":
			p.Name = types.String(value.Str)
		case "descr":
			p.Descr = types.String(value.Str)
		case "rules":
			p.Rules = types.String(value.Str)
		case "type":
			p.Type = types.Int8(value.Int())
		case "state":
			p.State = types.Int8(value.Int())
		}
		if e != nil {
			log.Error(e)
			return false
		}
		return true
	})
	return nil
}

// Free
func (p *CompanyRole) Free() {
	if p == nil {
		return
	}

	companyrolePool.Put(p)
}

// Reset
func (p *CompanyRole) Reset() {
	p.Id = 0
	p.CompanyId = 0
	p.Name = ""
	p.Descr = ""
	p.Rules = ""
	p.Type = 0
	p.State = 0

}

func (p *CompanyRole) TableName() string {
	return CompanyRoleTableName
}

func (p *CompanyRole) AssignPtr(args ...dialect.Field) []any {
	if len(args) == 0 {
		args = tblcompanyrole.ReadableFields
	}

	_vals := make([]any, 0, len(args))
	for _, col := range args {
		switch col {
		case tblcompanyrole.Id:
			_vals = append(_vals, &p.Id)
		case tblcompanyrole.CompanyId:
			_vals = append(_vals, &p.CompanyId)
		case tblcompanyrole.Name:
			_vals = append(_vals, &p.Name)
		case tblcompanyrole.Descr:
			_vals = append(_vals, &p.Descr)
		case tblcompanyrole.Rules:
			_vals = append(_vals, &p.Rules)
		case tblcompanyrole.Type:
			_vals = append(_vals, &p.Type)
		case tblcompanyrole.State:
			_vals = append(_vals, &p.State)
		}
	}

	return _vals
}

func (p *CompanyRole) Scan(rows *sql.Rows, args ...dialect.Field) ([]CompanyRole, bool, error) {
	defer rows.Close()
	company_roles := make([]CompanyRole, 0)

	if len(args) == 0 {
		args = tblcompanyrole.ReadableFields
	}

	for rows.Next() {
		_p := NewCompanyRole()
		_vals := _p.AssignPtr(args...)
		e := rows.Scan(_vals...)
		if e != nil {
			log.Error(e)
			return nil, false, e
		}
		company_roles = append(company_roles, *_p)
	}
	if e := rows.Err(); e != nil {
		log.Error(e)
		return nil, false, e
	}
	return company_roles, len(company_roles) > 0, nil
}

// RawAssignValues 向数据库写入数据前，为表列赋值。多用于批量插入和更新
// 如果 args 为空，则赋值所有可写字段
// 如果 args 不为空，则只赋值 args 中的字段
func (p *CompanyRole) RawAssignValues(args ...dialect.Field) ([]string, []any) {
	if len(args) == 0 {
		args = tblcompanyrole.WritableFields
	}
	return p.AssignValues(args...)
}

// AssignValues 向数据库写入数据前，为表列赋值。
// 如果 args 为空，则将非零值赋与可写字段
// 如果 args 不为空，则只赋值 args 中的字段
func (p *CompanyRole) AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		_lens = len(args)
		_cols []string
		_vals []any
	)

	if len(args) == 0 {
		args = tblcompanyrole.WritableFields
		_lens = len(args)
		_cols = make([]string, 0, _lens)
		_vals = make([]any, 0, _lens)
		for _, arg := range args {
			switch arg {
			case tblcompanyrole.Id:
				if p.Id == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrole.Id.Quote())
				_vals = append(_vals, p.Id)
			case tblcompanyrole.CompanyId:
				if p.CompanyId == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrole.CompanyId.Quote())
				_vals = append(_vals, p.CompanyId)
			case tblcompanyrole.Name:
				if p.Name == "" {
					continue
				}
				_cols = append(_cols, tblcompanyrole.Name.Quote())
				_vals = append(_vals, p.Name)
			case tblcompanyrole.Descr:
				if p.Descr == "" {
					continue
				}
				_cols = append(_cols, tblcompanyrole.Descr.Quote())
				_vals = append(_vals, p.Descr)
			case tblcompanyrole.Rules:
				if p.Rules == "" {
					continue
				}
				_cols = append(_cols, tblcompanyrole.Rules.Quote())
				_vals = append(_vals, p.Rules)
			case tblcompanyrole.Type:
				if p.Type == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrole.Type.Quote())
				_vals = append(_vals, p.Type)
			case tblcompanyrole.State:
				if p.State == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrole.State.Quote())
				_vals = append(_vals, p.State)
			}
		}
		return _cols, _vals
	}

	_cols = make([]string, 0, _lens)
	_vals = make([]any, 0, _lens)
	for _, arg := range args {
		switch arg {
		case tblcompanyrole.Id:
			_cols = append(_cols, tblcompanyrole.Id.Quote())
			_vals = append(_vals, p.Id)
		case tblcompanyrole.CompanyId:
			_cols = append(_cols, tblcompanyrole.CompanyId.Quote())
			_vals = append(_vals, p.CompanyId)
		case tblcompanyrole.Name:
			_cols = append(_cols, tblcompanyrole.Name.Quote())
			_vals = append(_vals, p.Name)
		case tblcompanyrole.Descr:
			_cols = append(_cols, tblcompanyrole.Descr.Quote())
			_vals = append(_vals, p.Descr)
		case tblcompanyrole.Rules:
			_cols = append(_cols, tblcompanyrole.Rules.Quote())
			_vals = append(_vals, p.Rules)
		case tblcompanyrole.Type:
			_cols = append(_cols, tblcompanyrole.Type.Quote())
			_vals = append(_vals, p.Type)
		case tblcompanyrole.State:
			_cols = append(_cols, tblcompanyrole.State.Quote())
			_vals = append(_vals, p.State)
		}
	}
	return _cols, _vals
}

//
func (p *CompanyRole) AssignKeys() (dialect.Field, any) {
	return tblcompanyrole.PrimaryKey, p.Id
}

//
func (p *CompanyRole) AssignPrimaryKeyValues(result sql.Result) error {
	_id, e := result.LastInsertId()
	if e != nil {
		return e
	}
	p.Id = types.BigInt(_id)
	return nil
}
