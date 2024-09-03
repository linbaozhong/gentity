// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

var model_str = `
{{- $tablename := .TableName}}
{{- $structName := .StructName}}
{{- $keys := .Keys}}
package {{.PackageName}}

import (
	"database/sql"
	"sync"
	{{- if .HasTime }}
	"time"
	{{- end}}
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

const {{$structName}}TableName = "{{.TableName}}"

var (
	{{lower .StructName}}Pool = sync.Pool{
		New: func() interface{} {
			return &{{.StructName}}{}
		},
	}
)

func New{{.StructName}}() *{{.StructName}} {
	return {{lower .StructName}}Pool.Get().(*{{.StructName}})
}

//Free
func (p *{{.StructName}}) Free() {
	if p == nil {
		return
	}
	{{range $key, $value := .Columns}}p.{{$key}} = {{getTypeValue $value}}
	{{end}}
	{{lower .StructName}}Pool.Put(p)
}


func (p *{{.StructName}}) TableName() string {
	return {{$structName}}TableName
}

func (p *{{.StructName}}) Scan(rows *sql.Rows, args ...atype.Field) ([]*{{.StructName}}, error) {
	defer rows.Close()
	{{.TableName}}s := make([]*{{.StructName}}, 0)

	if len(args) == 0 {
		args = {{$structName}}ReadableFields
	}

	for rows.Next() {
		p := New{{.StructName}}()
		vals := make([]any, 0, len(args))
		for _, col := range args {
			switch col {
			{{- range $key, $value := .Columns}}
			case {{$structName}}Tbl.{{ $key }}:
				vals = append(vals, &p.{{ $key }})
			{{- end}}
			}
		}
		err := rows.Scan(vals...)
		if err != nil {
			return nil, err
		}
		{{.TableName}}s = append({{.TableName}}s, p)
	}
	return {{.TableName}}s, nil
}

func (p *{{.StructName}})AssignValues(args ...atype.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = {{$structName}}WritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			{{- range $key, $value := .Columns}}
			case {{$structName}}Tbl.{{ $key }}:
				if p.{{ $key }}{{getZeroValue $value}} {
					continue
				}
				cols = append(cols, {{$structName}}Tbl.{{ $key }}.Quote())
				vals = append(vals, p.{{ $key }})
			{{- end}}
			}
		}
		return cols, vals
	}
	
	cols = make([]string, 0, lens)
	vals = make([]any, 0, lens)
	for _, arg := range args {
		switch arg {
		{{- range $key, $value := .Columns}}
		case {{$structName}}Tbl.{{ $key }}:
			cols = append(cols, {{$structName}}Tbl.{{ $key }}.Quote())
			vals = append(vals, p.{{ $key }})
		{{- end}}
		}
	}
	return cols, vals
}

//
func (p *{{.StructName}}) AssignKeys() ([]atype.Field, []any) {
	return {{$structName}}PrimaryKeys, []any{
	{{- range $key,$value := .Keys}}
		p.{{$value}},
	{{- end}}
	}
}


var (
	{{$structName}}Tbl = struct {
	{{- range $key, $value := .Columns}}
		{{ $key }} atype.Field
	{{- end}}
	}{
	{{- range $key, $value := .Columns}}
		{{ $key }}: atype.Field{Name: "{{index $value 0}}",Table: {{$structName}}TableName},
	{{- end}}
	}

{{- if .HasPrimaryKey}}
	// 主键
	{{$structName}}PrimaryKeys = []atype.Field{
	{{- range $key,$value := .Keys}}
		{{$structName}}Tbl.{{$value}},
	{{- end}}
	}
{{- end}}

	// 可写列
	{{$structName}}WritableFields = []atype.Field{
	{{- range $key, $value := .Columns}}
		{{- if or (eq (index $value 4) "->") (eq (index $value 4) "")}}
		{{$structName}}Tbl.{{$key}},
		{{- end}}
	{{- end}}
	}
	// 可读列
	{{$structName}}ReadableFields = []atype.Field{
	{{- range $key, $value := .Columns}}
		{{- if or (eq (index $value 4) "<-") (eq (index $value 4) "")}}
		{{$structName}}Tbl.{{$key}},
		{{- end}}
	{{- end}}
	}

)
`
var (
	//
	tableTpl = `
{{- $tablename := .TableName}}
package {{ .TableName }}

import (
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

var (
{{- range $key, $value := .Columns}}
	{{ $key }} = atype.Field{Name: "{{index $value 0}}",Table: "{{ $tablename }}"}
{{- end}}

{{- if .HasPrimaryKey}}
	// 主键
	PrimaryKeys = []atype.Field{
	{{- range $key,$value := .Keys}}
		{{$value}},
	{{- end}}
	}
{{- end}}

	// 可写列
	WritableFields = []atype.Field{
{{- range $key, $value := .Columns}}
	{{- if or (eq (index $value 4) "->") (eq (index $value 4) "")}}
	{{$key}},
	{{- end}}
{{- end}}
	}
	// 可读列
	ReadableFields = []atype.Field{
{{- range $key, $value := .Columns}}
	{{- if or (eq (index $value 4) "<-") (eq (index $value 4) "")}}
	{{$key}},
	{{- end}}
{{- end}}
	}

)

`
	buildTpl = `
{{- $tablename := .TableName}}
{{- $structName := .StructName}}
package dao

import (
	"context"
	"{{.ModulePath}}/{{.PackageName}}"
	"github.com/linbaozhong/gentity/pkg/ace"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

type {{.StructName}}Daoer interface {
	atype.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *{{.PackageName}}.{{.StructName}}, cols ...atype.Field) (int64, error)
	// InsertMulti 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...atype.Field) (int64, error)
	// UpdateMulti 批量更新多条数据
	// cols: 要更新的列名
	UpdateMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...atype.Field) (bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []atype.Field, cond ...atype.Condition) ([]*{{.PackageName}}.{{.StructName}}, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond ...atype.Condition) ([]*{{.PackageName}}.{{.StructName}}, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []atype.Field, cond ...atype.Condition) (*{{.PackageName}}.{{.StructName}}, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, args ...any) (*{{.PackageName}}.{{.StructName}}, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond ...atype.Condition) (*{{.PackageName}}.{{.StructName}}, error)
}

type {{.TableName}}Dao struct {
	db    atype.Executer
}

func {{.StructName}}(exec atype.Executer) {{$structName}}Daoer {
	return &{{.TableName}}Dao{db: exec}
}

// C Create {{ .TableName }}
func (p *{{.TableName}}Dao) C() *ace.Creator {
	return ace.NewCreate(p.db, {{.PackageName}}.{{$structName}}TableName)
}

// R Read {{ .TableName }}
func (p *{{.TableName}}Dao) R() *ace.Selector{
	return ace.NewSelect(p.db, {{.PackageName}}.{{$structName}}TableName)
}

// U Update {{ .TableName }}
func (p *{{.TableName}}Dao) U() *ace.Updater{
	return ace.NewUpdate(p.db, {{.PackageName}}.{{$structName}}TableName)
}

// D Delete {{ .TableName }}
func (p *{{.TableName}}Dao) D() *ace.Deleter{
	return ace.NewDelete(p.db, {{.PackageName}}.{{$structName}}TableName)
}


// Insert 返回 LastInsertId
func (p *{{.TableName}}Dao) Insert(ctx context.Context, sets ...atype.Setter) (int64, error) {
	if len(sets) == 0 {
		return 0, atype.ErrSetterEmpty
	}
	result, err := p.C().Set(sets...).Do(ctx)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertOne 返回 LastInsertId
// cols: 要插入的列名
func (p *{{.TableName}}Dao) InsertOne(ctx context.Context, bean *{{.PackageName}}.{{.StructName}}, cols ...atype.Field) (int64, error) {
	result, err := p.C().Cols(cols...).Struct(ctx, bean)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertMulti 批量插入,返回 RowsAffected
// cols: 要插入的列名
func (p *{{.TableName}}Dao) InsertMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...atype.Field) (int64, error) {
	lens := len(beans)
	if lens == 0 {
		return 0, atype.ErrBeanEmpty
	}
	args := make([]atype.Modeler, 0, lens)
	for _, bean := range beans {
		args = append(args, bean)
	}
	result, err := p.C().Cols(cols...).Struct(ctx, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}


// Update
func (p *{{.TableName}}Dao) Update(ctx context.Context, sets []atype.Setter, cond ...atype.Condition) (bool, error) {
	if len(sets) == 0 {
		return false, atype.ErrSetterEmpty
	}
	result, err := p.U().Where(cond...).Set(sets...).Do(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// UpdateMulti
// cols: 要更新的列名
func (p *{{.TableName}}Dao) UpdateMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...atype.Field) (bool, error) {
	lens := len(beans)
	if lens == 0 {
		return false, atype.ErrBeanEmpty
	}
	args := make([]atype.Modeler, 0, lens)
	for _, bean := range beans {
		args = append(args, bean)
	}
	result, err := p.U().Cols(cols...).Struct(ctx, args...)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Delete
func (p *{{.TableName}}Dao) Delete(ctx context.Context, cond ...atype.Condition) (bool, error) {
	result, err := p.D().Where(cond...).Do(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Get4Cols
func (p *{{.TableName}}Dao) Get4Cols(ctx context.Context, cols []atype.Field, cond ...atype.Condition) (*{{.PackageName}}.{{.StructName}}, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols({{.PackageName}}.{{$structName}}ReadableFields...)
	} else {
		c.Cols(cols...)
	}
	
	rows, err := c.Where(cond...).Limit(1).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	obj := {{.PackageName}}.New{{.StructName}}()
	defer obj.Free()

	objs, err := obj.Scan(rows, cols...)
	if err != nil {
		return nil, err
	}
	if len(objs) == 0 {
		return nil, atype.ErrNotFound
	}
	return objs[0], nil
}
//
// Find4Cols
func (p *{{.TableName}}Dao) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []atype.Field, cond ...atype.Condition) ([]*{{.PackageName}}.{{.StructName}}, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols({{.PackageName}}.{{$structName}}ReadableFields...)
	} else {
		c.Cols(cols...)
	}
	//
	if pageSize == 0 {
		pageSize = atype.PageSize
	}
	//
	rows, err := c.Where(cond...).Limit(pageSize, pageSize*pageIndex).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	obj := {{.PackageName}}.New{{.StructName}}()
	defer obj.Free()

	objs, err := obj.Scan(rows, cols...)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

// GetByID Read one {{.TableName}} By Primary Key value,
// Pass values in this order：{{ range $key,$value := .Keys}}{{if gt $key 0}},{{end}}{{$value}}{{ end}}
func (p *{{.TableName}}Dao) GetByID(ctx context.Context, args ...any) (*{{.PackageName}}.{{.StructName}}, error) {
	lens := len({{.PackageName}}.{{$structName}}PrimaryKeys)
	if lens != len(args) {
		return nil, atype.ErrArgsNotMatch
	}
	
	cond := make([]atype.Condition, 0, lens)
	for i, key := range {{.PackageName}}.{{$structName}}PrimaryKeys {
		cond = append(cond, key.Eq(args[i]))
	}
	return p.Get4Cols(ctx, []atype.Field{}, cond...)
}

// Get Read one {{.TableName}}
func (p *{{.TableName}}Dao) Get(ctx context.Context, cond ...atype.Condition) (*{{.PackageName}}.{{.StructName}}, error) {
	return p.Get4Cols(ctx, []atype.Field{}, cond...)
}

// Find
func (p *{{.TableName}}Dao) Find(ctx context.Context, pageIndex, pageSize uint, cond ...atype.Condition) ([]*{{.PackageName}}.{{.StructName}}, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []atype.Field{}, cond...)
}


// IDs
func (p *{{.TableName}}Dao) IDs(ctx context.Context, cond ...atype.Condition) ([]int64, error) {
	if len({{.PackageName}}.{{$structName}}PrimaryKeys) == 0 {
		return nil, atype.ErrPrimaryKeyNotMatch
	}
	c := p.R().Cols({{.PackageName}}.{{$structName}}PrimaryKeys[0])
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, atype.PageSize)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// Columns
func (p *{{.TableName}}Dao) Columns(ctx context.Context, col atype.Field, cond ...atype.Condition) ([]any, error) {
	c := p.R().Cols(col)
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]any, atype.PageSize)
	for rows.Next() {
		var v any
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		cols = append(cols, v)
	}

	return cols, nil
}

// Count
func (p *{{.TableName}}Dao) Count(ctx context.Context, cond ...atype.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *{{.TableName}}Dao) Sum(ctx context.Context, col atype.Field, cond ...atype.Condition) (int64, error) {
	return p.R().Sum(ctx, col, cond...)
}

// Exists
func (p *{{.TableName}}Dao) Exists(ctx context.Context, cond ...atype.Condition) (bool, error) {
	if len({{.PackageName}}.{{$structName}}PrimaryKeys) == 0 {
		return false, atype.ErrPrimaryKeyNotMatch
	}
	c := p.R().Cols({{.PackageName}}.{{$structName}}PrimaryKeys[0]).Where(cond...).Limit(1)
	rows, err := c.Query(ctx)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}


`
)
