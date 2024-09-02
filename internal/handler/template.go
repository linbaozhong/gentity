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
{{- $keys := .Keys}}
package {{.PackageName}}

import (
	"database/sql"
	"sync"
	// {{- range $key,$value := .Imports}}
	// "{{ $value }}"
	// {{- end}}
	{{- if .HasTime }}
	"time"
	{{- end}}
	"{{.ModulePath}}/table/{{.TableName}}"
	// "github.com/linbaozhong/gentity/pkg/ace"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

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
	return "{{.TableName}}"
}

func (p *{{.StructName}}) Scan(rows *sql.Rows, args ...atype.Field) ([]*{{.StructName}}, error) {
	defer rows.Close()
	{{.TableName}}s := make([]*{{.StructName}}, 0)

	if len(args) == 0 {
		args = {{$tablename}}.ReadableFields
	}

	for rows.Next() {
		p := New{{.StructName}}()
		vals := make([]any, 0, len(args))
		for _, col := range args {
			switch col {
			{{- range $key, $value := .Columns}}
			case {{$tablename}}.{{ $key }}:
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
//
// func (p *{{.StructName}})AssignColumns(args ...atype.Field) []string {
// 	var lens = len(args)
// 	if lens > 0 {
// 		cols := make([]string, 0, lens)
// 		for _, arg := range args {
// 			switch arg {
// 			{{- range $key, $value := .Columns}}
// 			case {{$tablename}}.{{ $key }}:
// 				cols = append(cols, {{$tablename}}.{{ $key }}.Quote())
// 			{{- end}}
// 			}
// 		}
// 		return cols
// 	}
//
// 	cols := make([]string, 0, len({{$tablename}}.WritableFields))
// 	for _, col := range {{$tablename}}.WritableFields {
// 		cols = append(cols, col.Quote())
// 	}
// 	return cols
// }

func (p *{{.StructName}})AssignValues(args ...atype.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = {{$tablename}}.WritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			{{- range $key, $value := .Columns}}
			case {{$tablename}}.{{ $key }}:
				if p.{{ $key }}{{getZeroValue $value}} {
					continue
				}
				cols = append(cols, {{$tablename}}.{{ $key }}.Quote())
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
		case {{$tablename}}.{{ $key }}:
			cols = append(cols, {{$tablename}}.{{ $key }}.Quote())
			vals = append(vals, p.{{ $key }})
		{{- end}}
		}
	}
	return cols, vals
}

//
func (p *{{.StructName}}) AssignKeys() ([]atype.Field, []any) {
	return {{$tablename}}.PrimaryKeys, []any{
	{{- range $key,$value := .Keys}}
		p.{{$value}},
	{{- end}}
	}
}


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
	WritableFields = []atype.Field {
{{- range $key, $value := .Columns}}
	{{- if or (eq (index $value 4) "->") (eq (index $value 4) "")}}
	{{$key}},
	{{- end}}
{{- end}}
	}
	// 可读列
	ReadableFields = []atype.Field {
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
package dao

import (
	"context"
	"{{.ModulePath}}/db"
	"{{.ModulePath}}/table/{{.TableName}}"
	"github.com/linbaozhong/gentity/pkg/ace"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

type {{.TableName}}Dao struct {
	db    ace.Executer
	table string
}

func {{.StructName}}(exec ace.Executer) *{{.TableName}}Dao {
	return &{{.TableName}}Dao{db: exec, table: "{{.TableName}}"}
}

// C Create {{ .TableName }}
func (p *{{.TableName}}Dao) C() *ace.Creator {
	return ace.NewCreate(p.db, p.table)
}

// R Read {{ .TableName }}
func (p *{{.TableName}}Dao) R() *ace.Selector{
	return ace.NewSelect(p.db, p.table)
}

// U Update {{ .TableName }}
func (p *{{.TableName}}Dao) U() *ace.Updater{
	return ace.NewUpdate(p.db, p.table)
}

// D Delete {{ .TableName }}
func (p *{{.TableName}}Dao) D() *ace.Deleter{
	return ace.NewDelete(p.db, p.table)
}


// Insert 返回 LastInsertId
func (p *{{.TableName}}Dao) Insert(ctx context.Context, sets []atype.Setter) (int64, error) {
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

// Single4Cols
func (p *{{.TableName}}Dao) Single4Cols(ctx context.Context, cols []atype.Field, cond ...atype.Condition) (*{{.PackageName}}.{{.StructName}}, error) {
	c := p.R().Cols(cols...).Where(cond...).Limit(1)
	rows, err := c.Query(ctx)
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
// Multi4Cols
func (p *{{.TableName}}Dao) Multi4Cols(ctx context.Context, cols []atype.Field, cond ...atype.Condition) ([]*{{.PackageName}}.{{.StructName}}, error) {
	c := p.R().Cols(cols...).Where(cond...).Limit(1000)
	rows, err := c.Query(ctx)
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

// Get Read a {{.TableName}} By Primary Key value,
// Pass values in this order：{{ range $key,$value := .Keys}}{{$value}},{{ end}}
func (p *{{.TableName}}Dao) Get(ctx context.Context, args ...any) (*{{.PackageName}}.{{.StructName}}, error) {
	lens := len({{.TableName}}.PrimaryKeys)
	if lens != len(args) {
		return nil, atype.ErrArgsNotMatch
	}
	
	cond := make([]atype.Condition, 0, lens)
	for i, key := range {{.TableName}}.PrimaryKeys {
		cond = append(cond, key.Eq(args[i]))
	}
	return p.Single4Cols(ctx, []atype.Field{}, cond...)
}

// Single
func (p *{{.TableName}}Dao) Single(ctx context.Context, cond ...atype.Condition) (*{{.PackageName}}.{{.StructName}}, error) {
	return p.Single4Cols(ctx, []atype.Field{}, cond...)
}

// Multi
func (p *{{.TableName}}Dao) Multi(ctx context.Context, cond ...atype.Condition) ([]*{{.PackageName}}.{{.StructName}}, error) {
	return p.Multi4Cols(ctx, []atype.Field{}, cond...)
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
	c := p.R().Cols({{.TableName}}.PrimaryKeys...).Where(cond...).Limit(1)
	rows, err := c.Query(ctx)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}


`
)
