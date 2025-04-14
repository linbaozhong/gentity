// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package db

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/linbaozhong/gentity/example/model/define/table/tblcompanyrule"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/app"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

const CompanyRuleTableName = "company_rule"

var (
	companyrulePool = pool.New(app.Context, func() any {
		_obj := &CompanyRule{}
		_obj.UUID()
		return _obj
	})
)

func NewCompanyRule() *CompanyRule {
	_obj := companyrulePool.Get().(*CompanyRule)
	return _obj
}

// MarshalJSON
func (p *CompanyRule) MarshalJSON() ([]byte, error) {
	var _buf = bytes.NewBuffer(nil)
	_buf.WriteByte('{')
	if p.Id != 0 {
		_buf.WriteString(`"id":` + types.Marshal(p.Id) + `,`)
	}
	if p.Pid != 0 {
		_buf.WriteString(`"pid":` + types.Marshal(p.Pid) + `,`)
	}
	if p.Path != "" {
		_buf.WriteString(`"path":` + types.Marshal(p.Path) + `,`)
	}
	if p.Title != "" {
		_buf.WriteString(`"title":` + types.Marshal(p.Title) + `,`)
	}
	if p.Type != 0 {
		_buf.WriteString(`"type":` + types.Marshal(p.Type) + `,`)
	}
	if p.IsPrivate != 0 {
		_buf.WriteString(`"is_private":` + types.Marshal(p.IsPrivate) + `,`)
	}
	if p.State != 0 {
		_buf.WriteString(`"state":` + types.Marshal(p.State) + `,`)
	}
	if p.Descr != "" {
		_buf.WriteString(`"descr":` + types.Marshal(p.Descr) + `,`)
	}
	if p.Belong != 0 {
		_buf.WriteString(`"belong":` + types.Marshal(p.Belong) + `,`)
	}
	if l := _buf.Len(); l > 1 {
		_buf.Truncate(l - 1)
	}
	_buf.WriteByte('}')
	return _buf.Bytes(), nil
}

// UnmarshalJSON
func (p *CompanyRule) UnmarshalJSON(data []byte) error {

	if !gjson.ValidBytes(data) {
		return errors.New("invalid json")
	}
	_result := gjson.ParseBytes(data)
	_result.ForEach(func(key, value gjson.Result) bool {
		var e error
		switch key.Str {
		case "id":
			p.Id = types.BigInt(value.Uint())
		case "pid":
			p.Pid = types.BigInt(value.Uint())
		case "path":
			p.Path = types.String(value.Str)
		case "title":
			p.Title = types.String(value.Str)
		case "type":
			p.Type = types.Uint8(value.Uint())
		case "is_private":
			p.IsPrivate = types.Uint8(value.Uint())
		case "state":
			p.State = types.Int8(value.Int())
		case "descr":
			p.Descr = types.String(value.Str)
		case "belong":
			p.Belong = types.Int8(value.Int())
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
func (p *CompanyRule) Free() {
	if p == nil {
		return
	}

	companyrulePool.Put(p)
}

// Reset
func (p *CompanyRule) Reset() {
	p.Id = 0
	p.Pid = 0
	p.Path = ""
	p.Title = ""
	p.Type = 0
	p.IsPrivate = 0
	p.State = 0
	p.Descr = ""
	p.Belong = 0

}

func (p *CompanyRule) TableName() string {
	return CompanyRuleTableName
}

func (p *CompanyRule) AssignPtr(args ...dialect.Field) []any {
	if len(args) == 0 {
		args = tblcompanyrule.ReadableFields
	}

	_vals := make([]any, 0, len(args))
	for _, col := range args {
		switch col {
		case tblcompanyrule.Id:
			_vals = append(_vals, &p.Id)
		case tblcompanyrule.Pid:
			_vals = append(_vals, &p.Pid)
		case tblcompanyrule.Path:
			_vals = append(_vals, &p.Path)
		case tblcompanyrule.Title:
			_vals = append(_vals, &p.Title)
		case tblcompanyrule.Type:
			_vals = append(_vals, &p.Type)
		case tblcompanyrule.IsPrivate:
			_vals = append(_vals, &p.IsPrivate)
		case tblcompanyrule.State:
			_vals = append(_vals, &p.State)
		case tblcompanyrule.Descr:
			_vals = append(_vals, &p.Descr)
		case tblcompanyrule.Belong:
			_vals = append(_vals, &p.Belong)
		}
	}

	return _vals
}

func (p *CompanyRule) Scan(rows *sql.Rows, args ...dialect.Field) ([]CompanyRule, bool, error) {
	defer rows.Close()
	company_rules := make([]CompanyRule, 0)

	if len(args) == 0 {
		args = tblcompanyrule.ReadableFields
	}

	for rows.Next() {
		_p := NewCompanyRule()
		_vals := _p.AssignPtr(args...)
		e := rows.Scan(_vals...)
		if e != nil {
			log.Error(e)
			return nil, false, e
		}
		company_rules = append(company_rules, *_p)
	}
	if e := rows.Err(); e != nil {
		log.Error(e)
		return nil, false, e
	}
	return company_rules, len(company_rules) > 0, nil
}

// RawAssignValues 向数据库写入数据前，为表列赋值。多用于批量插入和更新
// 如果 args 为空，则赋值所有可写字段
// 如果 args 不为空，则只赋值 args 中的字段
func (p *CompanyRule) RawAssignValues(args ...dialect.Field) ([]string, []any) {
	if len(args) == 0 {
		args = tblcompanyrule.WritableFields
	}
	return p.AssignValues(args...)
}

// AssignValues 向数据库写入数据前，为表列赋值。
// 如果 args 为空，则将非零值赋与可写字段
// 如果 args 不为空，则只赋值 args 中的字段
func (p *CompanyRule) AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		_lens = len(args)
		_cols []string
		_vals []any
	)

	if len(args) == 0 {
		args = tblcompanyrule.WritableFields
		_lens = len(args)
		_cols = make([]string, 0, _lens)
		_vals = make([]any, 0, _lens)
		for _, arg := range args {
			switch arg {
			case tblcompanyrule.Id:
				if p.Id == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Id.Quote())
				_vals = append(_vals, p.Id)
			case tblcompanyrule.Pid:
				if p.Pid == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Pid.Quote())
				_vals = append(_vals, p.Pid)
			case tblcompanyrule.Path:
				if p.Path == "" {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Path.Quote())
				_vals = append(_vals, p.Path)
			case tblcompanyrule.Title:
				if p.Title == "" {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Title.Quote())
				_vals = append(_vals, p.Title)
			case tblcompanyrule.Type:
				if p.Type == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Type.Quote())
				_vals = append(_vals, p.Type)
			case tblcompanyrule.IsPrivate:
				if p.IsPrivate == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrule.IsPrivate.Quote())
				_vals = append(_vals, p.IsPrivate)
			case tblcompanyrule.State:
				if p.State == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrule.State.Quote())
				_vals = append(_vals, p.State)
			case tblcompanyrule.Descr:
				if p.Descr == "" {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Descr.Quote())
				_vals = append(_vals, p.Descr)
			case tblcompanyrule.Belong:
				if p.Belong == 0 {
					continue
				}
				_cols = append(_cols, tblcompanyrule.Belong.Quote())
				_vals = append(_vals, p.Belong)
			}
		}
		return _cols, _vals
	}

	_cols = make([]string, 0, _lens)
	_vals = make([]any, 0, _lens)
	for _, arg := range args {
		switch arg {
		case tblcompanyrule.Id:
			_cols = append(_cols, tblcompanyrule.Id.Quote())
			_vals = append(_vals, p.Id)
		case tblcompanyrule.Pid:
			_cols = append(_cols, tblcompanyrule.Pid.Quote())
			_vals = append(_vals, p.Pid)
		case tblcompanyrule.Path:
			_cols = append(_cols, tblcompanyrule.Path.Quote())
			_vals = append(_vals, p.Path)
		case tblcompanyrule.Title:
			_cols = append(_cols, tblcompanyrule.Title.Quote())
			_vals = append(_vals, p.Title)
		case tblcompanyrule.Type:
			_cols = append(_cols, tblcompanyrule.Type.Quote())
			_vals = append(_vals, p.Type)
		case tblcompanyrule.IsPrivate:
			_cols = append(_cols, tblcompanyrule.IsPrivate.Quote())
			_vals = append(_vals, p.IsPrivate)
		case tblcompanyrule.State:
			_cols = append(_cols, tblcompanyrule.State.Quote())
			_vals = append(_vals, p.State)
		case tblcompanyrule.Descr:
			_cols = append(_cols, tblcompanyrule.Descr.Quote())
			_vals = append(_vals, p.Descr)
		case tblcompanyrule.Belong:
			_cols = append(_cols, tblcompanyrule.Belong.Quote())
			_vals = append(_vals, p.Belong)
		}
	}
	return _cols, _vals
}

//
func (p *CompanyRule) AssignKeys() (dialect.Field, any) {
	return tblcompanyrule.PrimaryKey, p.Id
}

//
func (p *CompanyRule) AssignPrimaryKeyValues(result sql.Result) error {
	_id, e := result.LastInsertId()
	if e != nil {
		return e
	}
	p.Id = types.BigInt(_id)
	return nil
}
