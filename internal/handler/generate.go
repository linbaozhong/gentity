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
	"github.com/linbaozhong/gentity/pkg/util"
	"github.com/linbaozhong/gentity/pkg/validator"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func initDTOFile(prefix string) {
	astOnce.Do(func() {
		var structFullName = filepath.Join(fullpath, prefix)
		_astFile, e := getAst(structFullName)
		if e != nil {
			showError(e)
			return
		}
		dtoFile, e = os.OpenFile(filepath.Join(fullpath, dentityDTO), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if e != nil {
			showError(e)
			return
		}
		var buf bytes.Buffer
		buf.WriteString("// Code generated by gentity. DO NOT EDIT.\n\n")
		buf.WriteString("package " + _astFile.Name.Name + " \n\n")
		buf.WriteString("import (\n")
		buf.WriteString("	\"bytes\" \n")
		buf.WriteString("	\"errors\" \n")
		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/gjson\" \n")
		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/log\" \n")
		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/validator\" \n")
		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/types\" \n")
		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/conv\" \n")
		buf.WriteString("	\"net/http\" \n")
		buf.WriteString(") \n\n")
		buf.WriteString("var _ = conv.Any2String(\"\")\n")
		buf.WriteString("var _ = bytes.NewBuffer(nil)\n")
		_, e = dtoFile.Write(buf.Bytes())
		if e != nil {
			showError(e)
			return
		}
	})
}

func generateDao(tds []TempData, prefix string) error {
	for _, td := range tds {
		// 写table文件
		e := td.writeTable(filepath.Join(tablePath, "tbl"+strings.ToLower(td.StructName)))
		// e := td.writeTable(filepath.Join(daoPath, strings.ToLower(td.StructName)))
		if e != nil {
			showError(e.Error())
			return e
		}

		// 写model文件
		e = td.writeToModel(prefix)
		if e != nil {
			showError(e)
			return e
		}

		// 写dal文件
		e = td.writeBuild(filepath.Join(daoPath, "dao"+strings.ToLower(td.StructName)))
		// e = td.writeBuild(daoPath)
		if e != nil {
			showError(e.Error())
			return e
		}
	}
	return nil
}

func generateDTO(tds []TempData, prefix string) error {
	initDTOFile(prefix)

	for _, td := range tds {
		// 写DTO文件
		e := writeToDTO(td)
		if e != nil {
			showError(e)
			return e
		}
	}
	return nil
}

func writeToDTO(d TempData) error {
	_funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"sub": func(a, b int) int {
			return a - b
		},
		"getType":           getType,
		"getTypeValue":      getTypeValue,
		"getZeroValue":      getZeroValue,
		"getNotZeroValue":   getNotZeroValue,
		"getSqlValue":       getSqlValue,
		"getSqlType":        getSqlType,
		"getUnmarshalJSON":  getUnmarshalJSON,
		"getTypeNil":        getTypeNil,
		"getValidFunc":      getValidFunc,
		"getValidParamFunc": getValidParamFunc,
		"getFieldString":    getFieldString,
		"contains":          contains,
	}
	var fn func(ioWriter io.Writer, _funcMap template.FuncMap) error
	fn = func(ioWriter io.Writer, _funcMap template.FuncMap) error {
		tmpl := template.New("").Funcs(_funcMap)
		_, e := tmpl.ParseFS(resources.TemplatesFS, "templates/dto_struct.tmpl")
		if e != nil {
			showError(e)
			return e
		}
		return tmpl.ExecuteTemplate(ioWriter, "dto_struct.tmpl", d)
	}

	var buf = new(bytes.Buffer)
	e := fn(buf, _funcMap)
	if e != nil {
		showError(e)
		return e
	}
	_formatted, _ := format.Source(buf.Bytes())
	if _formatted == nil {
		_, e = dtoFile.Write(buf.Bytes())
	} else {
		_, e = dtoFile.Write(_formatted)
	}
	if e != nil {
		showError(e.Error())
	}
	return e
}
func getTypeNil(t Field) any {
	switch t.Type {
	case "string", "types.String":
		return "types.NilString"
	case "int", "types.Int":
		return "types.NilInt"
	case "int8", "types.Int8":
		return "types.NilInt8"
	case "int16", "types.Int16":
		return "types.NilInt16"
	case "int32", "types.Int32":
		return "types.NilInt32"
	case "int64", "types.Int64", "types.Money":
		return "types.NilInt64"
	case "uint", "types.Uint":
		return "types.NilUint"
	case "uint8", "types.Uint8":
		return "types.NilUint8"
	case "uint16", "types.Uint16":
		return "types.NilUint16"
	case "uint32", "types.Uint32":
		return "types.NilUint32"
	case "uint64", "types.Uint64", "types.BigInt":
		return "types.NilUint64"
	case "float32", "types.Float32":
		return "types.NilFloat32"
	case "float64", "types.Float64":
		return "types.NilFloat64"
	case "time.Time":
		return "types.NilTime"
	case "types.Time":
		return "types.Time{}"
	case "bool", "types.Bool":
		return "types.Bool(-1)"
	default:
		switch t.Type[0] {
		case '*', '[', 'm':
			return "nil"
		default:
			return `"请将字段类型改为指针/切片/映射等引用类型"`
		}
	}
}
func getFieldString(t Field) string {
	switch t.Type {
	case "string":
		return `p.` + t.Name
	case "types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64",
		"types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64",
		"types.BigInt", "types.Money", "types.Float32", "types.Float64", "types.Time",
		"types.Bool", "types.String":
		return `p.` + t.Name + ".String()"
	default:
		return `conv.Any2String(p.` + t.Name + `)`
	}
}

// getValidFunc
// n: 字段名
func getValidFunc(t, n string) (fo funcObj) {
	fo = funcObj{}
	pos := strings.Index(t, "~")
	if pos > 0 {
		fo.Err = t[pos+1:]
	} else {
		fo.Err = "Wrong " + n + " format"
	}
	var tag string
	if pos < 0 {
		tag = t
	} else {
		tag = t[:pos]
	}

	if fn, ok := validator.TagMap[tag]; ok {
		fo.Func = fn
		return
	}
	return
}

type funcObj struct {
	Func  string
	Param string
	Err   string
}

// getValidParamFunc
// n: 字段名
func getValidParamFunc(t, n string) (fo funcObj) {
	fo = funcObj{}
	pos := strings.Index(t, "(")
	if pos < 1 {
		return
	}
	pos1 := strings.Index(t, ")")
	if pos1 < 0 || pos1 <= pos {
		return
	}
	_params := strings.Split(t[pos+1:pos1], "|")

	var tag string
	if pos < 0 {
		tag = t
	} else {
		tag = t[:pos]
	}

	pos = strings.Index(t, "~")
	if pos > 0 {
		fo.Err = t[pos+1:]
	} else {
		fo.Err = "Wrong " + n + " " + tag
	}

	if fn, ok := validator.ParamTagMap[tag]; ok {
		fo.Func = fn
		fo.Param = strconv.Quote(strings.Join(_params, `","`))
		return
	}
	return
}

func contains(s []string, e string) bool {
	before, after, found := strings.Cut(e, "|")
	if util.SliceContains(s, before) {
		return true
	}
	if found {
		return util.SliceContains(s, after)
	}
	return false
}
