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
	"strconv"
	"strings"
	"text/template"
)

// TempData 表示生成template所需要的数据结构
type TempData struct {
	ParseTag    []string
	Module      string
	ModulePath  string
	FileName    string
	PackageName string
	Imports     map[string]struct{}
	StructName  string
	TableName   string
	CacheData   string // 数据缓存时长
	CacheList   string // list缓存时长
	CacheLimit  string // list缓存长度
	Columns     []Field
	PrimaryKey  Field
	Relation    Relation // 关系键
	// Relations     []Relation // ← 新增这一行
	HasPrimaryKey bool
	HasRef        bool // 有引用类型
	HasCache      bool
	// HasCustomType bool
	// HasTime       bool
	// VisitorName   string // 访问者字段名

	Path           string // @path 文档标签，自定义接口路径
	Entity         string // 前端实体前缀（由结构体名推导，如 user）
	Group          string // 接口分组前缀（--group）
	Platform       string // 前端目标平台：pc | h5 | mp
	SearchFormComp string // M5: 自动嵌套的 SearchForm 组件文件名（如 UserListReqSearchForm.vue），空表示不嵌套
}

// Field struct 字段
type Field struct {
	Name   string   // 字段名
	Col    string   // 数据库列名
	Json   jsonObj  // json名
	Type   string   // 类型
	Rw     string   // 数据库读写标志
	Valids []string // 数据校验规则
}

type jsonObj struct {
	Name      string
	OmitEmpty bool
	OmitZero  bool
}

// Relation 关系
type Relation struct {
	Name     string // 字段名
	Type     string // 类型
	Field    string // 参照字段
	Foreign  string // 外键字段
	Kind     string // slice，ptr，struct
	JoinType string // ← 新增这一行
	Alias    string // ← 新增这一行
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
	if v[0] == '*' {
		v = v[1:]
		switch v {
		case "string", "types.String":
			return "new(" + v + ")"
		case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64",
			"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
			"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.Float32",
			"types.Float64", "types.BigInt", "types.Money":
			return "new(" + v + ")"
		case "time.Time", "types.Time":
			return "new(" + v + ")" // `time.Time{}`
		case "bool", "types.Bool":
			return "new(" + v + ")"
		default:
			if v[:2] == "[]" {
				return "p." + t.Name + "[:0]"
			}
			return v + "{}"
		}
	}
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
		"hasPrefix":        strings.HasPrefix,
	}

	fileName = filepath.Join(fullpath, getBaseFilename(fileName)+"_do_"+d.StructName+".gen.go")

	return writeToFormatFile(fileName, _f, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/do_struct.tmpl")
		if e != nil {
			return e
		}
		return tmpl.ExecuteTemplate(ioWriter, "do_struct.tmpl", d)
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
		"split": strings.Split, // ← 新增这一行
		"sub": func(a, b int) int {
			return a - b
		},
		"getReturnValue": func(p, t string) string {
			return "[]" + p + "." + t
		},
		"upperFirst": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"toFieldName": func(colName string) string {
			parts := strings.Split(colName, "_")
			var result strings.Builder
			for _, part := range parts {
				if len(part) > 0 {
					result.WriteString(strings.ToUpper(part[:1]))
					result.WriteString(part[1:])
				}
			}
			return result.String()
		},
	}
	return writeToFormatFile(_fileName, _f, func(ioWriter io.Writer, funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(funcMap)
		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/dao.tmpl")
		if e != nil {
			// showError(e)
			return e
		}

		// // ↓↓↓ 新增：解析关联查询模板 ↓↓↓
		// _, e = tmpl.ParseFS(resources.TemplatesFS, "templates/dao_relation.tmpl")
		// if e != nil {
		//	return e
		// }
		// // ↑↑↑ 新增结束 ↑↑↑

		e = tmpl.ExecuteTemplate(ioWriter, "dao.tmpl", d)
		if e != nil {
			return e
		}

		// // ↓↓↓ 新增：执行关联查询模板 ↓↓↓
		// if len(d.Relations) > 0 {
		//	e = tmpl.ExecuteTemplate(ioWriter, "dao_relation.tmpl", d)
		// }
		// // ↑↑↑ 新增结束 ↑↑↑

		return e
	})
}

//
// func (d *TempData) writeDTO(parent string) error {
//	e := os.MkdirAll(parent, os.ModePerm)
//	if e != nil {
//		return e
//	}
//
//	if len(d.Relations) == 0 {
//		return nil
//	}
//
//	_fileName := filepath.Join(parent, getBaseFilename(d.FileName)+"_dto_"+d.StructName+"_relation.gen.go")
//
//	_f := template.FuncMap{
//		"lower":           strings.ToLower,
//		"getNotZeroValue": getNotZeroValue,
//	}
//
//	return writeToFormatFile(_fileName, _f, func(ioWriter io.Writer, funcMap template.FuncMap) error {
//		tmpl := template.New("").Funcs(funcMap)
//		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/dto_relation.tmpl")
//		if e != nil {
//			return e
//		}
//		return tmpl.ExecuteTemplate(ioWriter, "dto_relation.tmpl", d)
//	})
// }

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

// toJsType 将 Go 字段类型映射为 JSDoc 类型
func toJsType(t Field) string {
	v := t.Type
	if strings.HasPrefix(v, "*") {
		v = v[1:]
	}
	switch v {
	case "string", "types.String", "time.Time", "types.Time":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64",
		"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
		"types.BigInt", "types.Money",
		"float32", "float64", "types.Float32", "types.Float64":
		return "number"
	case "bool", "types.Bool":
		return "boolean"
	default:
		return "any"
	}
}

// entityStripSuffixes 推导实体前缀时需要剥离的后缀（动作词 + Req/Resp）
var entityStripSuffixes = []string{
	"Req", "Resp", "Request", "Response",
	"Register", "Login", "Create", "Update", "Edit", "Get", "List",
	"Search", "Delete", "Detail", "Info", "Query",
}

// toEntity 由结构体名推导实体前缀，如 UserRegisterReq -> user
func toEntity(structName string) string {
	s := structName
	for {
		matched := false
		for _, suf := range entityStripSuffixes {
			if strings.HasSuffix(s, suf) && len(s) > len(suf) {
				s = s[:len(s)-len(suf)]
				matched = true
				break
			}
		}
		if !matched {
			break
		}
	}
	return getFieldName(s)
}

// toPascal 下划线转大驼峰，如 user -> User
func toPascal(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return b.String()
}

// toValidRule 将单个 valid 标签转为前端校验规则对象字面量（JS）
func toValidRule(valid string, f Field, entity string) string {
	msg := "t('" + entity + "." + f.Json.Name + "_invalid')"
	if i := strings.Index(valid, "~"); i >= 0 {
		msg = "\"" + valid[i+1:] + "\""
		valid = valid[:i]
	}
	var sb strings.Builder
	sb.WriteString("{ ")
	if pos := strings.Index(valid, "("); pos > 0 {
		pos1 := strings.Index(valid, ")")
		if pos1 > pos {
			tag := valid[:pos]
			params := valid[pos+1 : pos1]
			switch tag {
			case "maxstringlength":
				sb.WriteString("max: " + params)
			case "minstringlength":
				sb.WriteString("min: " + params)
			case "stringlength", "length":
				if ps := strings.Split(params, "|"); len(ps) == 2 {
					sb.WriteString("min: " + ps[0] + ", max: " + ps[1])
				}
			case "range":
				if ps := strings.Split(params, "|"); len(ps) == 2 {
					sb.WriteString("min: " + ps[0] + ", max: " + ps[1])
				}
			case "min":
				sb.WriteString("min: " + params)
			case "max":
				sb.WriteString("max: " + params)
			case "in":
				sb.WriteString("oneOf: [")
				for i, p := range strings.Split(params, "|") {
					if i > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString("\"" + p + "\"")
				}
				sb.WriteString("]")
			}
		}
	} else {
		switch valid {
		case "required":
			sb.WriteString("required: true, message: " + msg)
		case "email":
			sb.WriteString("type: 'email'")
		case "mobile":
			sb.WriteString("pattern: /^1[3456789]\\d{9}$/")
		case "int":
			sb.WriteString("type: 'integer'")
		case "float":
			sb.WriteString("type: 'number'")
		case "alphanum":
			sb.WriteString("pattern: /^[a-zA-Z0-9]+$/")
		default:
			sb.WriteString("/* unsupported valid: " + valid + " */")
		}
	}
	sb.WriteString(" }")
	return sb.String()
}

// toEndpoint 推导接口路径（可被 @path 覆盖，M1 先用启发式）
func toEndpoint(structName, group string) string {
	s := structName
	for _, suf := range []string{"Req", "Resp", "Request", "Response"} {
		if strings.HasSuffix(s, suf) {
			s = s[:len(s)-len(suf)]
			break
		}
	}
	ep := toPascal(toEntity(structName))
	if strings.HasPrefix(s, ep) {
		s = s[len(ep):]
	}
	action := getFieldName(s)
	if group != "" {
		return "/" + group + "/" + action
	}
	return "/" + toEntity(structName) + "/" + action
}

// hasTag 判断 TempData 是否带某解析标签
func hasTag(td TempData, tag string) bool {
	for _, t := range td.ParseTag {
		if t == tag {
			return true
		}
	}
	return false
}

// toOptions 从 valid:"in(a|b|c)" 提取选项
func toOptions(f Field) []string {
	for _, v := range f.Valids {
		if strings.HasPrefix(v, "in(") && strings.HasSuffix(v, ")") {
			return strings.Split(v[len("in("):len(v)-1], "|")
		}
	}
	return nil
}

// toOptionsJS 将字段的 options 转成 JS 数组字面量（uni-app picker 的 range 用）
func toOptionsJS(f Field) string {
	opts := toOptions(f)
	if len(opts) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(opts))
	for _, o := range opts {
		parts = append(parts, strconv.Quote(o))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// mockValue 生成字段的 mock 初值
func mockValue(f Field) string {
	switch {
	case strings.Contains(f.Type, "string") || f.Type == "types.String":
		return "\"\""
	case strings.Contains(f.Type, "int") || strings.Contains(f.Type, "uint") ||
		strings.Contains(f.Type, "Money") || strings.Contains(f.Type, "Float"):
		return "0"
	case strings.Contains(f.Type, "bool"):
		return "false"
	default:
		return "null"
	}
}

// toMethod 推导 HTTP 方法：Search/List/Query/Get 命名 → get，其余 post
func toMethod(structName string) string {
	for _, kw := range []string{"Search", "List", "Query", "Get"} {
		if strings.Contains(structName, kw) {
			return "get"
		}
	}
	return "post"
}

// isListView 判断是否为列表型响应（生成 ListView 用）
func isListView(structName string) bool {
	for _, kw := range []string{"List", "Search", "Query", "Page"} {
		if strings.Contains(structName, kw) {
			return true
		}
	}
	return false
}

// isSlice 判断字段是否为切片（表格/卡片一般跳过）
func isSlice(f Field) bool {
	return strings.HasPrefix(f.Type, "[]")
}

// hasValid 判断字段是否带某 valid 标签（如 "email" / "mobile"）
func hasValid(f Field, v string) bool {
	for _, x := range f.Valids {
		if strings.HasPrefix(x, v) {
			return true
		}
	}
	return false
}

// toInputType 根据类型/校验推导表单控件类型
func toInputType(f Field) string {
	t := f.Type
	if strings.HasPrefix(t, "*") {
		t = t[1:]
	}
	switch {
	case toOptions(f) != nil:
		return "select"
	case strings.Contains(t, "bool"):
		return "checkbox"
	case strings.Contains(t, "int") || strings.Contains(t, "uint") ||
		strings.Contains(t, "Money") || strings.Contains(t, "Float"):
		return "number"
	case strings.Contains(t, "Time"):
		return "date"
	case hasValid(f, "email"):
		return "email"
	case hasValid(f, "mobile"):
		return "tel"
	case strings.EqualFold(f.Json.Name, "password"):
		return "password"
	default:
		return "text"
	}
}

// toPkName 返回主键 json 名，无主键回退 "id"（表格行 key 用）
func toPkName(td TempData) string {
	if td.PrimaryKey.Json.Name != "" {
		return td.PrimaryKey.Json.Name
	}
	return "id"
}
