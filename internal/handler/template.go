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

// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

{{- $tablename := .TableName}}
{{- $structName := .StructName}}
{{- $structNameLower := lower .StructName}}
{{- $keys := .Keys}}
package {{.PackageName}}

import (
	"database/sql"
	"sync"
	{{- if .HasTime }}
	"time"
	{{- end}}
	"{{.ModulePath}}/define/table/{{$structNameLower}}tbl"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

const {{$structName}}TableName = "{{.TableName}}"

var (
	{{lower .StructName}}Pool = sync.Pool{
		New: func() any {
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
	{{range $key, $value := .Columns}}p.{{index $value 0}} = {{getTypeValue $value}}
	{{end}}
	{{lower .StructName}}Pool.Put(p)
}


func (p *{{.StructName}}) TableName() string {
	return {{$structName}}TableName
}

func (p *{{.StructName}}) Scan(rows *sql.Rows, args ...dialect.Field) ([]*{{.StructName}}, bool, error) {
	defer rows.Close()
	{{.TableName}}s := make([]*{{.StructName}}, 0)

	if len(args) == 0 {
		args = {{$structNameLower}}tbl.ReadableFields
	}

	for rows.Next() {
		p := New{{.StructName}}()
		vals := make([]any, 0, len(args))
		for _, col := range args {
			switch col {
			{{- range $key, $value := .Columns}}
			case {{$structNameLower}}tbl.{{ index $value 0 }}:
				vals = append(vals, &p.{{ index $value 0 }})
			{{- end}}
			}
		}
		err := rows.Scan(vals...)
		if err != nil {
			return nil, false, err
		}
		{{.TableName}}s = append({{.TableName}}s, p)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	if len({{.TableName}}s) == 0 {
		return nil, false, sql.ErrNoRows
	}
	return {{.TableName}}s, true, nil
}

func (p *{{.StructName}})AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = {{$structNameLower}}tbl.WritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			{{- range $key, $value := .Columns}}
			case {{$structNameLower}}tbl.{{ index $value 0 }}:
				if p.{{ index $value 0 }}{{getZeroValue $value}} {
					continue
				}
				cols = append(cols, {{$structNameLower}}tbl.{{ index $value 0 }}.Quote())
				vals = append(vals, p.{{ index $value 0 }})
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
		case {{$structNameLower}}tbl.{{ index $value 0 }}:
			cols = append(cols, {{$structNameLower}}tbl.{{ index $value 0 }}.Quote())
			vals = append(vals, p.{{ index $value 0 }})
		{{- end}}
		}
	}
	return cols, vals
}

//
func (p *{{.StructName}}) AssignKeys() ([]dialect.Field, []any) {
	return {{$structNameLower}}tbl.PrimaryKeys, []any{
	{{- range $key,$value := .Keys}}
		p.{{$value}},
	{{- end}}
	}
}

`
var (
	//
	tableTpl = `

// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

{{- $tablename := .TableName}}
package {{ .TableName }}tbl

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

var (
{{- range $key, $value := .Columns}}
	{{ index $value 0 }} = dialect.Field{Name: "{{index $value 1}}",Table: "{{ $tablename }}",Type: "{{index $value 2}}"}
{{- end}}

{{- if .HasPrimaryKey}}
	// 主键
	PrimaryKeys = []dialect.Field{
	{{- range $key,$value := .Keys}}
		{{$value}},
	{{- end}}
	}
{{- end}}

	// 可写列
	WritableFields = []dialect.Field{
{{- range $key, $value := .Columns}}
	{{- if or (eq (index $value 4) "->") (eq (index $value 4) "")}}
	{{index $value 0}},
	{{- end}}
{{- end}}
	}
	// 可读列
	ReadableFields = []dialect.Field{
{{- range $key, $value := .Columns}}
	{{- if or (eq (index $value 4) "<-") (eq (index $value 4) "")}}
	{{index $value 0}},
	{{- end}}
{{- end}}
	}

)

`
	buildTpl = `

// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

{{- $tablename := .TableName}}
{{- $structName := .StructName}}
{{- $structNameLower := lower .StructName}}

package dao

import (
	"context"
	"{{.ModulePath}}/{{.PackageName}}"
	"{{.ModulePath}}/define/table/{{$structNameLower}}tbl"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

type {{.StructName}}Daoer interface {
	dialect.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *{{.PackageName}}.{{.StructName}}, cols ...dialect.Field) (int64, error)
	// InsertMulti 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...dialect.Field) (int64, error)
	// UpdateMulti 批量更新多条数据
	// cols: 要更新的列名
	UpdateMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...dialect.Field) (bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]*{{.PackageName}}.{{.StructName}}, bool, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]*{{.PackageName}}.{{.StructName}}, bool, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*{{.PackageName}}.{{.StructName}}, bool, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, args ...any) (*{{.PackageName}}.{{.StructName}}, bool, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond ...dialect.Condition) (*{{.PackageName}}.{{.StructName}}, bool, error)
	// GetSingle 按条件读取第一行的第一个字段
	GetSingle(ctx context.Context,col dialect.Field, cond ...dialect.Condition) (any, error)
	//
	IDs(ctx context.Context, cond ...dialect.Condition) ([]int64, error)
	//
	Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error)
}

type {{.TableName}}Dao struct {
	db    ace.Executer
}

func {{.StructName}}(exec ace.Executer) {{$structName}}Daoer {
	return &{{.TableName}}Dao{db: exec}
}

// C Create {{ .TableName }}
func (p *{{.TableName}}Dao) C() *ace.Creator {
	return p.db.Insert({{.PackageName}}.{{$structName}}TableName)
}

// R Read {{ .TableName }}
func (p *{{.TableName}}Dao) R() *ace.Selector{
	return p.db.Select({{.PackageName}}.{{$structName}}TableName)
}

// U Update {{ .TableName }}
func (p *{{.TableName}}Dao) U() *ace.Updater{
	return p.db.Update({{.PackageName}}.{{$structName}}TableName)
}

// D Delete {{ .TableName }}
func (p *{{.TableName}}Dao) D() *ace.Deleter{
	return p.db.Delete({{.PackageName}}.{{$structName}}TableName)
}


// Insert 返回 LastInsertId
func (p *{{.TableName}}Dao) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
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
func (p *{{.TableName}}Dao) InsertOne(ctx context.Context, bean *{{.PackageName}}.{{.StructName}}, cols ...dialect.Field) (int64, error) {
	result, err := p.C().Cols(cols...).Struct(ctx, bean)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertMulti 批量插入,返回 RowsAffected
// cols: 要插入的列名
func (p *{{.TableName}}Dao) InsertMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...dialect.Field) (int64, error) {
	lens := len(beans)
	if lens == 0 {
		return 0, atype.ErrBeanEmpty
	}
	args := make([]dialect.Modeler, 0, lens)
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
func (p *{{.TableName}}Dao) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
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
func (p *{{.TableName}}Dao) UpdateMulti(ctx context.Context, beans []*{{.PackageName}}.{{.StructName}}, cols ...dialect.Field) (bool, error) {
	lens := len(beans)
	if lens == 0 {
		return false, atype.ErrBeanEmpty
	}
	args := make([]dialect.Modeler, 0, lens)
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
func (p *{{.TableName}}Dao) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	result, err := p.D().Where(cond...).Do(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Get4Cols
func (p *{{.TableName}}Dao) Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*{{.PackageName}}.{{.StructName}}, bool, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols({{$structNameLower}}tbl.ReadableFields...)
	} else {
		c.Cols(cols...)
	}
	
	rows, err := c.Where(cond...).Limit(1).Query(ctx)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	obj := {{.PackageName}}.New{{.StructName}}()
	defer obj.Free()

	objs, has, err := obj.Scan(rows, cols...)
	if err != nil {
		return nil, false, err
	}
	if has {
		return objs[0], true, nil
	}

	return nil, false, atype.ErrNotFound
}
//
// Find4Cols
func (p *{{.TableName}}Dao) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]*{{.PackageName}}.{{.StructName}}, bool, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols({{$structNameLower}}tbl.ReadableFields...)
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
		return nil, false, err
	}
	defer rows.Close()

	obj := {{.PackageName}}.New{{.StructName}}()
	defer obj.Free()

	objs, has, err := obj.Scan(rows, cols...)
	if err != nil {
		return nil, false, err
	}
	return objs, has, nil
}

// GetByID Read one {{.TableName}} By Primary Key value,
// Pass values in this order：{{ range $key,$value := .Keys}}{{if gt $key 0}},{{end}}{{$value}}{{ end}}
func (p *{{.TableName}}Dao) GetByID(ctx context.Context, args ...any) (*{{.PackageName}}.{{.StructName}}, bool, error) {
	lens := len({{$structNameLower}}tbl.PrimaryKeys)
	if lens != len(args) {
		return nil, false, atype.ErrArgsNotMatch
	}
	
	cond := make([]dialect.Condition, 0, lens)
	for i, key := range {{$structNameLower}}tbl.PrimaryKeys {
		cond = append(cond, key.Eq(args[i]))
	}
	return p.Get4Cols(ctx, []dialect.Field{}, cond...)
}

// Get Read one {{.TableName}}
func (p *{{.TableName}}Dao) Get(ctx context.Context, cond ...dialect.Condition) (*{{.PackageName}}.{{.StructName}}, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond...)
}

// GetSingle Read the first column of the first row
func (p *{{.TableName}}Dao) GetSingle(ctx context.Context,col dialect.Field, cond ...dialect.Condition) (any, error) {
	c := p.R().Cols(col)
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var v any
	if rows.Next() {
		if err = rows.Scan(&v); err != nil {
			return nil, err
		}
	}

	return v, rows.Err()
}

// Find
func (p *{{.TableName}}Dao) Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]*{{.PackageName}}.{{.StructName}}, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond...)
}


// IDs
func (p *{{.TableName}}Dao) IDs(ctx context.Context, cond ...dialect.Condition) ([]int64, error) {
	if len({{$structNameLower}}tbl.PrimaryKeys) == 0 {
		return nil, atype.ErrPrimaryKeyNotMatch
	}
	c := p.R().Cols({{$structNameLower}}tbl.PrimaryKeys[0])
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, atype.PageSize)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

// Columns
func (p *{{.TableName}}Dao) Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error) {
	c := p.R().Cols(col)
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]any, atype.PageSize)
	for rows.Next() {
		var v any
		if err = rows.Scan(&v); err != nil {
			return nil, err
		}
		cols = append(cols, v)
	}
	return cols, rows.Err()
}

// Count
func (p *{{.TableName}}Dao) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *{{.TableName}}Dao) Sum(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (int64, error) {
	return p.R().Sum(ctx, col, cond...)
}

// Exists
func (p *{{.TableName}}Dao) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	if len({{$structNameLower}}tbl.PrimaryKeys) == 0 {
		return false, atype.ErrPrimaryKeyNotMatch
	}
	c := p.R().Cols({{$structNameLower}}tbl.PrimaryKeys[0]).Where(cond...).Limit(1)
	rows, err := c.Query(ctx)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), rows.Err()
}
`
)
