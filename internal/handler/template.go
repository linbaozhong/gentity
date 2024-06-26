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
package {{.PackageName}}

import (
	"fmt"
	"sync"
	{{if .HasTime}}"time"{{end}}
	"{{.ModulePath}}/table/{{.TableName}}"
	"github.com/linbaozhong/gentity/pkg/orm/sql"
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

//TableName
func (*{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

func (p *{{.StructName}})scanValues(columns []string) ([]any, error) {
	{{- $tablename := .TableName}}
	values := make([]any, len(columns))
	for i,column := range columns {
		switch column {
		{{- range $key, $value := .Columns}}
		case {{$tablename}}.{{ $key }}.Name:
			values[i] = new({{getSqlValue $value}})
		{{- end}}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

func (p *{{.StructName}})assignValues(columns []string, values []any) error {
	{{- $tablename := .TableName}}
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i,column := range columns {
		switch column {
		{{- range $key, $value := .Columns}}
		case {{$tablename}}.{{ $key }}.Name:
			value,ok := values[i].({{getSqlValue $value}})
			if !ok {
				return fmt.Errorf("unexpected type %T for field {{index $value 0}}", value)
			}
			{{- $v := index $value 2}}
			{{- if or (eq $v "string") (eq $v "int64") (eq $v "bool") (eq $v "float64") (eq $v "time.Time")}}
			p.{{$key}} = value.{{getSqlType $value}}
			{{- else}}
			p.{{$key}} = {{index $value 2}}(value.{{getSqlType $value}})
			{{- end}}
		{{- end}}
		}
	}
	return nil
}
`
var (
	//
	tableTpl = `
package {{ .TableName }}

import (
	"github.com/linbaozhong/gentity/pkg/orm/sql"
)
{{- $tablename := .TableName}}
const (
	TableName = "{{ .TableName }}"
)

var (
{{- range $key, $value := .Columns}}
	{{ $key }} = sql.Field{Name: "{{index $value 0}}",table: "{{ $tablename }}"}
{{- end}}
)
`
	buildTpl = `
package table

import (
	{{if .HasTime}}"time"{{end}}
	"{{.ModulePath}}/table/{{.TableName}}"
	"github.com/linbaozhong/gentity/pkg/orm/sql"
)

// {{.StructName}}Create 新增 {{ .TableName }}
func {{.StructName}}Create() *sql.Creator {
	return sql.NewCreate({{ .TableName }}.TableName)
}

// {{.StructName}}Update 修改 {{ .TableName }}
func {{.StructName}}Update() *sql.Updater{
	return sql.NewUpdate({{ .TableName }}.TableName)
}


// {{.StructName}}Delete 删除 {{ .TableName }}
func {{.StructName}}Delete() *sql.Deleter{
	return sql.NewDelete({{ .TableName }}.TableName)
}


// {{.StructName}}Query 查询 {{ .TableName }}，返回 []{{.StructName}}
func {{.StructName}}Query() *sql.Selector{
	return sql.NewSelect({{ .TableName }}.TableName)
}
`
)
