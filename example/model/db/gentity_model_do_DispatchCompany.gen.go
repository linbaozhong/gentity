// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package db

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/linbaozhong/gentity/example/model/define/table/tbldispatchcompany"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

const DispatchCompanyTableName = "dispatch_company"

var (
	dispatchcompanyPool = pool.New(app.Context, func() any {
		_obj := &DispatchCompany{}
		_obj.UUID()
		return _obj
	})
)

func NewDispatchCompany() *DispatchCompany {
	_obj := dispatchcompanyPool.Get().(*DispatchCompany)
	return _obj
}

// MarshalJSON
func (p *DispatchCompany) MarshalJSON() ([]byte, error) {
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
	if p.Address != "" {
		_buf.WriteString(`"address":` + types.Marshal(p.Address) + `,`)
	}
	if p.Creator != 0 {
		_buf.WriteString(`"creator":` + types.Marshal(p.Creator) + `,`)
	}
	if p.State != 0 {
		_buf.WriteString(`"state":` + types.Marshal(p.State) + `,`)
	}
	if !p.Ctime.IsZero() {
		_buf.WriteString(`"ctime":` + types.Marshal(p.Ctime) + `,`)
	}
	if !p.Utime.IsZero() {
		_buf.WriteString(`"utime":` + types.Marshal(p.Utime) + `,`)
	}
	if l := _buf.Len(); l > 1 {
		_buf.Truncate(l - 1)
	}
	_buf.WriteByte('}')
	return _buf.Bytes(), nil
}

// UnmarshalJSON
func (p *DispatchCompany) UnmarshalJSON(data []byte) error {

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
		case "address":
			p.Address = types.String(value.Str)
		case "creator":
			p.Creator = types.BigInt(value.Uint())
		case "state":
			p.State = types.Int8(value.Int())
		case "ctime":
			p.Ctime = types.Time{Time: value.Time()}
		case "utime":
			p.Utime = types.Time{Time: value.Time()}
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
func (p *DispatchCompany) Free() {
	if p == nil {
		return
	}

	dispatchcompanyPool.Put(p)
}

// Reset
func (p *DispatchCompany) Reset() {
	p.Id = 0
	p.CompanyId = 0
	p.Name = ""
	p.Address = ""
	p.Creator = 0
	p.State = 0
	p.Ctime = types.Time{}
	p.Utime = types.Time{}

}

func (p *DispatchCompany) TableName() string {
	return DispatchCompanyTableName
}

func (p *DispatchCompany) AssignPtr(args ...dialect.Field) []any {
	if len(args) == 0 {
		args = tbldispatchcompany.ReadableFields
	}

	_vals := make([]any, 0, len(args))
	for _, col := range args {
		switch col {
		case tbldispatchcompany.Id:
			_vals = append(_vals, &p.Id)
		case tbldispatchcompany.CompanyId:
			_vals = append(_vals, &p.CompanyId)
		case tbldispatchcompany.Name:
			_vals = append(_vals, &p.Name)
		case tbldispatchcompany.Address:
			_vals = append(_vals, &p.Address)
		case tbldispatchcompany.Creator:
			_vals = append(_vals, &p.Creator)
		case tbldispatchcompany.State:
			_vals = append(_vals, &p.State)
		case tbldispatchcompany.Ctime:
			_vals = append(_vals, &p.Ctime)
		case tbldispatchcompany.Utime:
			_vals = append(_vals, &p.Utime)
		}
	}

	return _vals
}

func (p *DispatchCompany) Scan(rows *sql.Rows, args ...dialect.Field) ([]DispatchCompany, bool, error) {
	defer rows.Close()
	dispatch_companys := make([]DispatchCompany, 0)

	if len(args) == 0 {
		args = tbldispatchcompany.ReadableFields
	}

	for rows.Next() {
		_p := NewDispatchCompany()
		_vals := _p.AssignPtr(args...)
		e := rows.Scan(_vals...)
		if e != nil {
			log.Error(e)
			return nil, false, e
		}
		dispatch_companys = append(dispatch_companys, *_p)
	}
	if e := rows.Err(); e != nil {
		log.Error(e)
		return nil, false, e
	}
	return dispatch_companys, len(dispatch_companys) > 0, nil
}

// RawAssignValues 向数据库写入数据前，为表列赋值。多用于批量插入和更新
// 如果 args 为空，则赋值所有可写字段
// 如果 args 不为空，则只赋值 args 中的字段
func (p *DispatchCompany) RawAssignValues(args ...dialect.Field) ([]string, []any) {
	if len(args) == 0 {
		args = tbldispatchcompany.WritableFields
	}
	return p.AssignValues(args...)
}

// AssignValues 向数据库写入数据前，为表列赋值。
// 如果 args 为空，则将非零值赋与可写字段
// 如果 args 不为空，则只赋值 args 中的字段
func (p *DispatchCompany) AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		_lens = len(args)
		_cols []string
		_vals []any
	)

	if len(args) == 0 {
		args = tbldispatchcompany.WritableFields
		_lens = len(args)
		_cols = make([]string, 0, _lens)
		_vals = make([]any, 0, _lens)
		for _, arg := range args {
			switch arg {
			case tbldispatchcompany.Id:
				if p.Id == 0 {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.Id.Quote())
				_vals = append(_vals, p.Id)
			case tbldispatchcompany.CompanyId:
				if p.CompanyId == 0 {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.CompanyId.Quote())
				_vals = append(_vals, p.CompanyId)
			case tbldispatchcompany.Name:
				if p.Name == "" {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.Name.Quote())
				_vals = append(_vals, p.Name)
			case tbldispatchcompany.Address:
				if p.Address == "" {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.Address.Quote())
				_vals = append(_vals, p.Address)
			case tbldispatchcompany.Creator:
				if p.Creator == 0 {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.Creator.Quote())
				_vals = append(_vals, p.Creator)
			case tbldispatchcompany.State:
				if p.State == 0 {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.State.Quote())
				_vals = append(_vals, p.State)
			case tbldispatchcompany.Ctime:
				if p.Ctime.IsZero() {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.Ctime.Quote())
				_vals = append(_vals, p.Ctime)
			case tbldispatchcompany.Utime:
				if p.Utime.IsZero() {
					continue
				}
				_cols = append(_cols, tbldispatchcompany.Utime.Quote())
				_vals = append(_vals, p.Utime)
			}
		}
		return _cols, _vals
	}

	_cols = make([]string, 0, _lens)
	_vals = make([]any, 0, _lens)
	for _, arg := range args {
		switch arg {
		case tbldispatchcompany.Id:
			_cols = append(_cols, tbldispatchcompany.Id.Quote())
			_vals = append(_vals, p.Id)
		case tbldispatchcompany.CompanyId:
			_cols = append(_cols, tbldispatchcompany.CompanyId.Quote())
			_vals = append(_vals, p.CompanyId)
		case tbldispatchcompany.Name:
			_cols = append(_cols, tbldispatchcompany.Name.Quote())
			_vals = append(_vals, p.Name)
		case tbldispatchcompany.Address:
			_cols = append(_cols, tbldispatchcompany.Address.Quote())
			_vals = append(_vals, p.Address)
		case tbldispatchcompany.Creator:
			_cols = append(_cols, tbldispatchcompany.Creator.Quote())
			_vals = append(_vals, p.Creator)
		case tbldispatchcompany.State:
			_cols = append(_cols, tbldispatchcompany.State.Quote())
			_vals = append(_vals, p.State)
		case tbldispatchcompany.Ctime:
			_cols = append(_cols, tbldispatchcompany.Ctime.Quote())
			_vals = append(_vals, p.Ctime)
		case tbldispatchcompany.Utime:
			_cols = append(_cols, tbldispatchcompany.Utime.Quote())
			_vals = append(_vals, p.Utime)
		}
	}
	return _cols, _vals
}

func (p *DispatchCompany) AssignKeys() (dialect.Field, any) {
	return tbldispatchcompany.PrimaryKey, p.Id
}

func (p *DispatchCompany) AssignPrimaryKeyValues(result sql.Result) error {
	_id, e := result.LastInsertId()
	if e != nil {
		return e
	}
	p.Id = types.BigInt(_id)
	return nil
}
