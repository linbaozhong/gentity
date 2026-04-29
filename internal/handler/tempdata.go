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
	ParseTag      []string
	Module        string
	ModulePath    string
	FileName      string
	PackageName   string
	Imports       []string
	StructName    string
	TableName     string
	CacheData     string // 数据缓存时长
	CacheList     string // list缓存时长
	CacheLimit    string // list缓存长度
	Columns       []Field
	PrimaryKey    Field
	RelationX     Relation // 关系键
	HasPrimaryKey bool
	HasRef        bool // 有引用类型
	HasCache      bool
	HasCustomType bool
	HasTime       bool
	VisitorName   string // 访问者字段名
}

// Field struct 字段
type Field struct {
	Name string // 字段名
	// info   typeInfo // 字段长度(占用字节数)
	Col    string   // 数据库列名
	Json   jsonObj  // json名
	Type   string   // 类型
	Rw     string   // 数据库读写标志
	Valids []string // 数据校验规则
	// idx    int      // 索引
}

type jsonObj struct {
	Name      string
	OmitEmpty bool
	OmitZero  bool
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
	_f := filepath.Base(filename)
	pos := strings.LastIndex(_f, ".")
	if pos == -1 {
		return _f
	}
	return _f[:pos]
}

func writeDaoBase(parent string, tds []TempData) error {
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		showError(e)
		return e
	}

	_fileName := filepath.Join(parent, "gentity_model.gen.go")

	_funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return writeToFormatFile(_fileName, _funcMap, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		_tmpl := template.New("").Funcs(funcMap)
		_, e := _tmpl.ParseFS(resources.TemplatesFS, "templates/dao_base.tmpl")
		if e != nil {
			showError(e)
			return e
		}

		e = _tmpl.ExecuteTemplate(ioWriter, "dao_base.tmpl", tds)
		if e != nil {
			showError(e)
		}
		return e
	})
}

func getType(t Field) string {
	_v := t.Type
	switch _v {
	case "string", "types.String", "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64",
		"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
		"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.Float32",
		"types.Float64", "types.BigInt", "types.Money", "time.Time", "types.Time", "bool", "types.Bool":
		return _v
	default:
		return "any"
	}
}

func getUnmarshalJSON(t Field) string {
	prefix := "p." + t.Name + " = "
	switch t.Type {
	case "string":
		return prefix + `value.Str`
	case "int64":
		return prefix + `value.Int()`
	case "uint64":
		return prefix + `value.Uint()`
	case "int", "int8", "int16", "int32",
		"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64":
		return prefix + t.Type + "(value.Int())"
	case "uint", "uint8", "uint16", "uint32",
		"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
		"types.BigInt":
		return prefix + t.Type + "(value.Uint())"
	case "types.Money":
		return `e = types.Unmarshal(value, &p.` + t.Name + `)`
	case "float32", "types.Float32", "types.Float64":
		return prefix + t.Type + "(value.Float())"
	case "float64":
		return prefix + "value.Float()"
	case "types.String":
		return prefix + `types.String(value.Str)`
	case "time.Time":
		return prefix + "value.Time()"
	case "types.Time":
		return prefix + "types.Time{Time: value.Time()}"
	case "bool":
		return prefix + "value.Bool()"
	case "types.Bool":
		return prefix + "util.IIF(value.Bool(),types.Bool(1),types.Bool(0))"
	default:
		fn := `func(value gjson.Result) ` + t.Type + ` {
					var obj ` + t.Type + `
					e := types.Unmarshal(value, &obj)
					if e != nil {
						panic(e)
					}
					return obj
				}(value)`
		return `e = types.Unmarshal(value, &p.` + t.Name + `, ` + fn + `)`
	}
}

func getTypeValue(t Field) any {
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
}

func getZeroValue(t Field) any {
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
}

func getNotZeroValue(t Field) any {
	v := t.Type
	switch v {
	case "string", "types.String":
		return `p.` + t.Name + ` != ""`
	case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64",
		"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
		"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.BigInt", "types.Money":
		return `p.` + t.Name + ` != 0`
	case "float32", "float64", "types.Float32", "types.Float64":
		return `p.` + t.Name + ` != 0.0`
	case "time.Time", "types.Time":
		return `!p.` + t.Name + `.IsZero()`
	case "bool", "types.Bool":
		return `p.` + t.Name + ` != false`
	default:
		return `p.` + t.Name + ` != nil`
	}
}

func getSqlValue(t Field) any {
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
}
func getSqlType(t Field) any {
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
}

func (d *TempData) writeToModel(fileName string) error {
	_f := template.FuncMap{
		"lower": strings.ToLower,
		"sub": func(a, b int) int {
			return a - b
		},
		"getType":          getType,
		"getTypeValue":     getTypeValue,
		"getZeroValue":     getZeroValue,
		"getNotZeroValue":  getNotZeroValue,
		"getSqlValue":      getSqlValue,
		"getSqlType":       getSqlType,
		"getUnmarshalJSON": getUnmarshalJSON,
	}

	fileName = filepath.Join(fullpath, getBaseFilename(fileName)+"_do_"+d.StructName+".gen.go")

	return writeToFormatFile(fileName, _f, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/struct.tmpl")
		if e != nil {
			return e
		}
		return tmpl.ExecuteTemplate(ioWriter, "struct.tmpl", d)
	})
}

// writeTable 将生成好的模块文件写到本地
func (d *TempData) writeTable(parent string) error {
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		// showError(e)
		return e
	}

	_fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_tbl_"+d.StructName+".gen.go") // d.tableFilename(parent)
	_f := template.FuncMap{
		"lower":   strings.ToLower,
		"getType": getType,
	}
	return writeToFormatFile(_fileName, _f, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/table.tmpl")
		if e != nil {
			return e
		}
		return tmpl.ExecuteTemplate(ioWriter, "table.tmpl", d)
	})

}

func (d *TempData) writeBuild(parent string) error {
	e := os.MkdirAll(parent, os.ModePerm)
	if e != nil {
		// showError(e)
		return e
	}

	_fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_dao_"+d.StructName+".gen.go")

	_f := template.FuncMap{
		"lower": strings.ToLower,
		"getReturnValue": func(p, t string) string {
			return "[]" + p + "." + t
		},
	}
	return writeToFormatFile(_fileName, _f, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/dao.tmpl")
		if e != nil {
			// showError(e)
			return e
		}
		return tmpl.ExecuteTemplate(ioWriter, "dao.tmpl", d)
	})
}

func writeToFormatFile(fullFilename string, funcMap template.FuncMap, fn func(ioWriter io.Writer, funcMap template.FuncMap) error) error {
	if fi, e := os.Stat(fullFilename); e == nil {
		if !fi.IsDir() {
			if e := os.Remove(fullFilename); e != nil {
				showError(e)
				return e
			}
		}
	}

	_f, e := os.OpenFile(fullFilename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if e != nil {
		showError(e.Error())
		return e
	}
	defer _f.Close()
	var buf bytes.Buffer
	e = fn(&buf, funcMap)
	if e != nil {
		showError(e.Error())
		return e
	}

	formatted, _ := format.Source(buf.Bytes())
	if formatted == nil {
		_, e = _f.Write(buf.Bytes())
	} else {
		_, e = _f.Write(formatted)
	}
	if e != nil {
		showError(e.Error())
		return e
	}
	return e
}

//
// // typeInfo 记录类型的对齐系数和占用大小 (基于 64 位系统)
// type typeInfo struct {
// 	size  int
// 	align int
// }

// // 内置类型的大小和对齐系数映射表
// var builtinTypes = map[string]typeInfo{
// 	// 布尔
// 	"bool": {1, 1},
// 	// 整数
// 	"int8": {1, 1}, "uint8": {1, 1}, "byte": {1, 1},
// 	"int16": {2, 2}, "uint16": {2, 2},
// 	"int32": {4, 4}, "uint32": {4, 4}, "float32": {4, 4}, "rune": {4, 4},
// 	"int64": {8, 8}, "uint64": {8, 8}, "float64": {8, 8},
// 	"int": {8, 8}, "uint": {8, 8}, "uintptr": {8, 8},
// 	// 复数
// 	"complex64": {8, 4}, "complex128": {16, 8},
// 	// 字符串 (64位系统下，指针8+长度8)
// 	"string": {16, 8},
//
// 	// ==========================================
// 	// 接口类型标识符
// 	// ==========================================
// 	// any 是 interface{} 的别名，底层是 eface (type指针 + data指针)，16字节，8对齐
// 	"any": {16, 8},
// 	// error 是 interface 的特例，底层是 iface (itab指针 + data指针)，16字节，8对齐
// 	"error": {16, 8},
//
// 	// ==========================================
// 	// 常见跨包类型 (精确硬编码)
// 	// ==========================================
// 	"time.Time":       {24, 8},
// 	"sync.Mutex":      {8, 8},
// 	"decimal.Decimal": {16, 8},
// 	"types.Money":     {8, 8},
// 	"types.BigInt":    {8, 8},
// 	"types.String":    {16, 8},
// 	"types.Int":       {8, 8},
// 	"types.Int8":      {1, 1},
// 	"types.Int16":     {2, 2},
// 	"types.Int32":     {4, 4},
// 	"types.Int64":     {8, 8},
// 	"types.Uint":      {8, 8},
// 	"types.Uint8":     {1, 1},
// 	"types.Uint16":    {2, 2},
// 	"types.Uint32":    {4, 4},
// 	"types.Uint64":    {8, 8},
// 	"types.Float32":   {4, 4},
// 	"types.Float64":   {8, 8},
// 	"types.Time":      {8, 8},
// 	"types.Bool":      {1, 1},
// 	"types.Point":     {8, 8},
// }
//
// // getExprInfo 根据 AST 表达式推断类型信息
// func getExprInfo(typeName string) typeInfo {
// 	// 基础类型 (int, string 等)
// 	if info, ok := builtinTypes[typeName]; ok {
// 		return info
// 	}
// 	// 同包内的未知自定义类型, 安全起见按 8 字节对齐处理
// 	return typeInfo{8, 8}
// }
