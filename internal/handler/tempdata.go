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
	Columns     [][]string
	// Keys        [][]string
	PrimaryKey []string
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

func writeDaoBase(parent string) error {
	err := os.MkdirAll(parent, os.ModePerm)
	if err != nil {
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
			return err
		}
		return tmpl.ExecuteTemplate(ioWriter, "dao_base.tmpl", nil)
	})
}

func (d *TempData) writeToModel(fileName string) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"getTypeValue": func(t []string) any {
			if len(t) < 3 {
				return `""`
			}
			var ret any
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
		"getZeroValue": func(t []string) any {
			if len(t) < 3 {
				return `""`
			}
			var ret any
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
		"getSqlValue": func(t []string) any {
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
		"getSqlType": func(t []string) any {
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

	fileName = filepath.Join(fullpath, getBaseFilename(fileName)+"_"+d.StructName+".gen.go")

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

	fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_"+d.StructName+".gen.go") // d.tableFilename(parent)
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
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

	fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_"+d.StructName+".gen.go")

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
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
				// showError(err)
				return err
			}
		}
	}

	f, err := os.OpenFile(fullFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		// showError(err.Error())
		return err
	}
	defer f.Close()
	var buf bytes.Buffer
	err = fn(&buf, funcMap)
	if err != nil {
		// showError(err.Error())
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// showError(err.Error())
		return err
	}
	_, err = f.Write(formatted)
	// _, err = f.Write(buf.Bytes())
	if err != nil {
		// showError(err.Error())
		return err
	}
	return err
}
