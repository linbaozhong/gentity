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

package cmd

var model_str = `
package {{.PackageName}}

import (
	"fmt"
	"sync"
	"database/sql"
	{{if .HasTime}}"time"{{end}}
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
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		{{range $key, $value := .Columns}}case "{{index $value 0}}":
			values[i] = new({{getSqlValue $value}})
		{{end}}
		default:
			values[i] = new(any)
		}
	}
	return values, nil
}

func (p *{{.StructName}})assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		{{range $key, $value := .Columns}}case "{{index $value 0}}":
			value,ok := values[i].({{getSqlValue $value}})
			if !ok {
				return fmt.Errorf("unexpected type %T for field {{index $value 0}}", value)
			}
			{{- $v := index $value 2}}
			{{- if or (eq $v "string") (eq $v "int64") (eq $v "bool") (eq $v "float64") (eq $v "time.Time")}}
			p.{{$key}} = value.{{getSqlType $value}}
			{{else}}
			p.{{$key}} = {{index $value 2}}(value.{{getSqlType $value}})
			{{- end}}
		{{end}}
		}
	}
	return nil
}
`
var (
	//
	tableName = "TableName"
	//
	tableTpl = `
		package table

		type _{{.StructName}} struct {
			TableName string
			ColumnNames        []string          //可读列名
			WriteColumnNames   []string          //可写列名
			ColumnName2Comment map[string]string //列名和列描述映射
			ColumnName2Json    map[string]string //列名和JSON Key映射
			{{if .HasPrimaryKey}}PrimaryKey TableField{{end}}
		{{range $key, $value := .Columns}} {{ $key }} TableField
		{{end}}
		}

		var (
			{{.StructName}}  _{{.StructName}}
		)

		func init() {
			{{.StructName}}.TableName = "{{lower .TableName}}"
			{{ $.StructName}}.ColumnNames = make([]string,0,{{len .Columns}})
			{{ $.StructName}}.WriteColumnNames = make([]string,0,{{len .Columns}})
			{{.StructName}}.ColumnName2Json = make(map[string]string,{{len .Columns}})
			{{.StructName}}.ColumnName2Comment = make(map[string]string,{{len .Columns}})

		{{if .HasPrimaryKey}}
			{{.StructName}}.PrimaryKey = TableField{
				Name: "{{index .PrimaryKey 0}}",
				Json: "{{index .PrimaryKey 1}}",
				Table: {{$.StructName}}.TableName,
			}
		{{end}}
		{{range $key, $value := .Columns}}
		{{ $.StructName}}.{{$key}} = TableField{
			Name: "{{index $value 0}}",
			Json: "{{index $value 1}}",
			Comment: "{{index $value 3}}",
			Table: {{$.StructName}}.TableName,
		}
		{{ $rw := index $value 4 }}{{if eq $rw "<-"}}
		{{ $.StructName}}.ColumnNames = append({{ $.StructName}}.ColumnNames,"{{index $value 0}}"){{else if eq $rw "->"}}
		{{ $.StructName}}.WriteColumnNames = append({{ $.StructName}}.WriteColumnNames,"{{index $value 0}}"){{else if eq $rw ""}}
		{{ $.StructName}}.ColumnNames = append({{ $.StructName}}.ColumnNames,"{{index $value 0}}")
		{{ $.StructName}}.WriteColumnNames = append({{ $.StructName}}.WriteColumnNames,"{{index $value 0}}"){{end}}
		{{ $.StructName}}.ColumnName2Json["{{index $value 0}}"] = "{{index $value 1}}"
		{{ $.StructName}}.ColumnName2Comment["{{index $value 0}}"] = "{{index $value 3}}"
		{{end}}
		}
		`
)
