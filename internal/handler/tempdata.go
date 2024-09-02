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
	Columns     map[string][]string
	Keys        []string
	// PrimaryKey     []string
	// PrimaryKeyName string // struct pk属性名
	HasPrimaryKey bool
	HasState      bool
	HasCache      bool
	HasTime       bool
	HasString     bool
	HasConvert    bool
}

func getBaseFilename(filename string) string {
	f := filepath.Base(filename)
	pos := strings.LastIndex(f, ".")
	if pos == -1 {
		return f
	}
	return f[:pos]
}

// func (d *TempData) tableFilename(parent string) string {
// 	return filepath.Join(parent, getBaseFilename(d.FileName)+"_"+d.StructName+".go")
// }

func (d *TempData) writeToModel(fileName string) error {
	var buf bytes.Buffer
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"getTypeValue": func(t []string) interface{} {
			if len(t) < 3 {
				return `""`
			}
			var ret interface{}
			switch t[2] {
			case "string":
				ret = `""`
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64":
				ret = 0
			case "time.Time":
				ret = `time.Time{}`
			case "bool":
				ret = `false`
			default:
				ret = 0
			}
			return ret
		},
		"getZeroValue": func(t []string) interface{} {
			if len(t) < 3 {
				return `""`
			}
			var ret interface{}
			switch t[2] {
			case "string":
				ret = ` == ""`
			case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64":
				ret = ` == 0`
			case "float32", "float64":
				ret = ` == 0.0`
			case "time.Time":
				ret = `.IsZero()`
			case "bool":
				ret = ` == false`
			default:
				ret = ` == 0`
			}
			return ret
		},
		"getSqlValue": func(t []string) interface{} {
			switch t[2] {
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
		"getSqlType": func(t []string) interface{} {
			switch t[2] {
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

	e := template.Must(template.New("tableTpl").Funcs(funcMap).Parse(model_str)).Execute(&buf, d)
	if e != nil {
		showError(e)
		return e
	}

	// absPath, _ := filepath.Abs(fileName)
	fileName = fileName[:len(fileName)-3] + "_" + d.StructName + "_define.go"
	if fi, err := os.Stat(fileName); err == nil {
		if !fi.IsDir() {
			if err := os.Remove(fileName); err != nil {
				showError(err.Error())
				return err
			}
		}
	}

	f, e := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		showError(e.Error())
		return e
	}
	defer f.Close()

	_, e = f.Write(buf.Bytes())
	if e != nil {
		showError(e)
		return e
	}

	return nil
}
func (d *TempData) writeTo(w io.Writer) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return template.Must(template.New("tableTpl").Funcs(funcMap).Parse(tableTpl)).Execute(w, d)
}

// writeTable 将生成好的模块文件写到本地
func (d *TempData) writeTable(parent string) error {
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		showError(err)
		return err
	}

	fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_"+d.StructName+"_table.go") // d.tableFilename(parent)

	if fi, err := os.Stat(fileName); err == nil {
		if !fi.IsDir() {
			if err := os.Remove(fileName); err != nil {
				showError(err.Error())
				return err
			}
		}
	}

	f, e := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		showError(e.Error())
		return e
	}
	defer f.Close()

	var buf bytes.Buffer
	e = d.writeTo(&buf)
	if e != nil {
		showError(e.Error())
		return e
	}

	formatted, e := format.Source(buf.Bytes())
	if e != nil {
		showError(e.Error())
		return e
	}
	_, e = f.Write(formatted)
	if e != nil {
		showError(e.Error())
		return e
	}
	return e
}

func (d *TempData) writeToBuild(w io.Writer) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return template.Must(template.New("buildTpl").Funcs(funcMap).Parse(buildTpl)).Execute(w, d)
}

func (d *TempData) writeBuild(parent string) error {
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		showError(err)
		return err
	}

	fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_"+d.StructName+"_dao.go") // d.tableFilename(parent)

	if fi, err := os.Stat(fileName); err == nil {
		if !fi.IsDir() {
			if err := os.Remove(fileName); err != nil {
				showError(err.Error())
				return err
			}
		}
	}

	f, e := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		showError(e.Error())
		return e
	}
	defer f.Close()

	var buf bytes.Buffer
	e = d.writeToBuild(&buf)
	if e != nil {
		showError(e.Error())
		return e
	}

	formatted, e := format.Source(buf.Bytes())
	if e != nil {
		showError(e.Error())
		return e
	}
	_, e = f.Write(formatted)
	if e != nil {
		showError(e.Error())
		return e
	}
	return e
}
