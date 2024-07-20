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
{{- $primary := .PrimaryKeyName}}
package {{.PackageName}}

import (
	"database/sql"
	"sync"
	{{- range $key,$value := .Imports}}
	"{{ $value }}"
	{{- end}}
	"{{.ModulePath}}/table/{{.TableName}}"
	"github.com/linbaozhong/gentity/pkg/ace"
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

func (p *{{.StructName}}) Scan(rows *sql.Rows, args ...atype.Field) ([]atype.Modeler, error) {
	defer rows.Close()
	{{.TableName}}s := make([]atype.Modeler, 0)
	lens := len(args)
	if lens == 0 {
		for rows.Next() {
			p := New{{.StructName}}()
			err := rows.Scan(
				{{- range $key, $value := .Columns}}
				&p.{{$key}},
				{{- end}}
			)
			if err != nil {
				return nil, err
			}
			{{.TableName}}s = append({{.TableName}}s, p)
		}
		return {{.TableName}}s, nil
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


func (p *{{.StructName}})AssignColumns(args ...atype.Field) []string {
	var lens = len(args)
	if lens == 0 {
		return []string{
		{{- range $key, $value := .Columns}}
		{{- if eq $primary $key}} {{continue}} {{end}}
			{{$tablename}}.{{ $key }}.Quote(),
		{{- end}}
		}
	}
	var (
		cols = make([]string, 0, lens)
	)
	for _, arg := range args {
		switch arg {
		{{- range $key, $value := .Columns}}
		{{- if eq $primary $key}} {{continue}} {{end}}
		case {{$tablename}}.{{ $key }}:
			cols = append(cols, {{$tablename}}.{{ $key }}.Quote())
		{{- end}}
		}
	}
	return cols
}

func (p *{{.StructName}})AssignValues(args ...atype.Field) []any {
	var lens = len(args)
	if lens == 0 {
		return []any{
		{{- range $key, $value := .Columns}}
			p.{{ $key }},
		{{- end}}
		}
	}
	var vals = make([]any, 0, lens)
	for _, arg := range args {
		switch arg {
		{{- range $key, $value := .Columns}}
		case {{$tablename}}.{{ $key }}:
			vals = append(vals, p.{{ $key }})
		{{- end}}
		}
	}
	return vals
}
`
var (
	//
	tableTpl = `
{{- $tablename := .TableName}}
package {{ .TableName }}

import (
	"github.com/linbaozhong/gentity/pkg/ace"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

var (
	PrimaryKey = atype.Field{Name: "{{index .PrimaryKey 0}}",Table: "{{ $tablename }}"}
{{- range $key, $value := .Columns}}
	{{ $key }} = atype.Field{Name: "{{index $value 0}}",Table: "{{ $tablename }}"}
{{- end}}
)

`
	buildTpl = `
package {{ .TableName }}

import (
	"{{.ModulePath}}/db"
	"{{.ModulePath}}/table/{{.TableName}}"
	"github.com/linbaozhong/gentity/pkg/ace"
)

// CreateX 新增 {{ .TableName }}
func CreateX(exec ace.Executer) *ace.Creator {
	return ace.NewCreate(exec, &db.{{.StructName}}{})
}

// UpdateX 修改 {{ .TableName }}
func UpdateX(exec ace.Executer) *ace.Updater{
	return ace.NewUpdate(exec, &db.{{.StructName}}{})
}


// DeleteX 删除 {{ .TableName }}
func DeleteX(exec ace.Executer) *ace.Deleter{
	return ace.NewDelete(exec, &db.{{.StructName}}{})
}


// SelectX 查询 {{ .TableName }}
func SelectX(exec ace.Executer) *ace.Selector{
	return ace.NewSelect(exec, &db.{{.StructName}}{})
}
`
)
