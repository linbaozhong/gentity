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

import (
	"bytes"
	"fmt"
	"github.com/linbaozhong/gentity/internal/resources"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TempData 表示生成template所需要的数据结构
type TempData struct {
	Module      string
	ModulePath  string
	FileName    string
	PackageName string
	Imports     []string
	StructName  string
	TableName   string
	CacheData   string // 数据缓存时长
	CacheList   string // list缓存时长
	CacheLimit  string // list缓存长度
	Columns     []Field
	// Keys        [][]string
	PrimaryKey    Field
	RelationX     Relation // 关系键
	HasPrimaryKey bool
	HasState      bool
	HasCache      bool
	HasTime       bool
	HasString     bool
	HasConvert    bool
}

// Field struct 字段
type Field struct {
	Name string // 字段名
	Col  string // 数据库列名
	Json string // json名
	Type string // 类型
	Rw   string // 数据库读写标志
}

// Relation 关系
type Relation struct {
	Name    string // 字段名
	Type    string // 类型
	Field   string // 参照字段
	Foreign string // 外键字段
	Kind    string // slice，ptr，struct
}

func getBaseFilename(filename string) string {
	f := filepath.Base(filename)
	pos := strings.LastIndex(f, ".")
	if pos == -1 {
		return f
	}
	return f[:pos]
}

func writeDaoBase(parent string) error {
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		fmt.Println("----")
		showError(err)
		return err
	}

	fileName := filepath.Join(parent, "gentity_model.gen.go")

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return writeToFormatFile(fileName, funcMap, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, err := tmpl.ParseFS(resources.TemplatesFS, "templates/dao_base.tmpl")
		if err != nil {
			showError(err)
			return err
		}

		err = tmpl.ExecuteTemplate(ioWriter, "dao_base.tmpl", nil)
		if err != nil {
			showError(err)
		}
		return err
	})
}

func getType(t Field) string {
	v := t.Type
	switch v {
	case "string", "types.String", "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64",
		"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
		"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.Float32",
		"types.Float64", "types.BigInt", "types.Money", "time.Time", "types.Time", "bool", "types.Bool":
		return v
	default:
		return "any"
	}
}
func (d *TempData) writeToModel(fileName string) error {
	funcMap := template.FuncMap{
		"lower":   strings.ToLower,
		"getType": getType,
		"getTypeValue": func(t Field) any {
			v := t.Type
			switch v {
			case "string", "types.String":
				return `""`
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64",
				"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
				"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.Float32",
				"types.Float64", "types.BigInt", "types.Money":
				return 0
			case "time.Time", "types.Time":
				return `types.Time{}` // `time.Time{}`
			case "bool", "types.Bool":
				return `false`
			default:
				if v[:2] == "[]" {
					return "p." + t.Name + "[:0]"
				}
				return v + "{}"
			}
		},
		"getZeroValue": func(t Field) any {
			v := t.Type
			switch v {
			case "string", "types.String":
				return ` == ""`
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
				"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
				"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.BigInt", "types.Money":
				return ` == 0`
			case "float32", "float64", "types.Float32", "types.Float64":
				return ` == 0.0`
			case "time.Time", "types.Time":
				return `.IsZero()`
			case "bool", "types.Bool":
				return ` == false`
			default:
				return ` == nil`
			}
		},
		"getSqlValue": func(t Field) any {
			switch t.Type {
			case "string":
				return "sql.NullString"
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64":
				return "sql.NullInt64"
			case "float32", "float64":
				return "sql.NullFloat64"
			case "time.Time":
				return "sql.NullTime"
			case "bool":
				return "sql.NullBool"
			default:
				return "sql.NullInt64"
			}
		},
		"getSqlType": func(t Field) any {
			switch t.Type {
			case "string":
				return "String"
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64":
				return "Int64"
			case "float32", "float64":
				return "Float64"
			case "time.Time":
				return "Time"
			case "bool":
				return "Bool"
			default:
				return "Int64"
			}
		},
		// "marshal": JSONValue,
	}

	fileName = filepath.Join(fullpath, getBaseFilename(fileName)+"_do_"+d.StructName+".gen.go")

	return writeToFormatFile(fileName, funcMap, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, err := tmpl.ParseFS(resources.TemplatesFS, "templates/struct.tmpl")
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(ioWriter, "struct.tmpl", d)
	})
}

// writeTable 将生成好的模块文件写到本地
func (d *TempData) writeTable(parent string) error {
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		// showError(err)
		return err
	}

	fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_tbl_"+d.StructName+".gen.go") // d.tableFilename(parent)
	funcMap := template.FuncMap{
		"lower":   strings.ToLower,
		"getType": getType,
	}
	return writeToFormatFile(fileName, funcMap, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, err := tmpl.ParseFS(resources.TemplatesFS, "templates/table.tmpl")
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(ioWriter, "table.tmpl", d)
	})

}

func (d *TempData) writeBuild(parent string) error {
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		// showError(err)
		return err
	}

	fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_dao_"+d.StructName+".gen.go")

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"getReturnValue": func(p, t string) string {
			return "[]" + p + "." + t
		},
	}
	return writeToFormatFile(fileName, funcMap, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, err := tmpl.ParseFS(resources.TemplatesFS, "templates/dao.tmpl")
		if err != nil {
			// showError(err)
			return err
		}
		return tmpl.ExecuteTemplate(ioWriter, "dao.tmpl", d)
	})
}

func writeToFormatFile(fullFilename string, funcMap template.FuncMap, fn func(ioWriter io.Writer, funcMap template.FuncMap) error) error {
	if fi, err := os.Stat(fullFilename); err == nil {
		if !fi.IsDir() {
			if err := os.Remove(fullFilename); err != nil {
				showError(err)
				return err
			}
		}
	}

	f, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		showError(err.Error())
		return err
	}
	defer f.Close()
	var buf bytes.Buffer
	err = fn(&buf, funcMap)
	if err != nil {
		showError(err.Error())
		return err
	}

	formatted, _ := format.Source(buf.Bytes())
	if formatted == nil {
		_, err = f.Write(buf.Bytes())
	} else {
		_, err = f.Write(formatted)
	}
	if err != nil {
		showError(err.Error())
		return err
	}
	return err
}
