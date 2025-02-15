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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sync"
)

const dentityDTO = "gentity_dto.gen.go"

var (
	astOnce sync.Once
	dtoFile *os.File
)

func getAst(structFullName string) (*ast.File, error) {
	var src any
	_fset := token.NewFileSet()
	_astFile, e := parser.ParseFile(_fset, structFullName, src, parser.ParseComments)
	if e != nil {
		return nil, e
	}

	return _astFile, nil
}

// func generateCheck(filename string) error {
// 	if filename == dentityDTO {
// 		return nil
// 	}
//
// 	var structFullName = filepath.Join(fullpath, filename)
//
// 	astOnce.Do(func() {
// 		astFile, e := getAst(structFullName)
// 		if e != nil {
// 			showError(e)
// 			return
// 		}
// 		dtoFile, e = os.OpenFile(filepath.Join(fullpath, dentityDTO), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
// 		if e != nil {
// 			showError(e)
// 			return
// 		}
// 		var buf bytes.Buffer
// 		buf.WriteString("// Code generated by gentity. DO NOT EDIT.\n\n")
// 		buf.WriteString("package " + astFile.Name.Name + " \n\n")
// 		buf.WriteString("import (\n")
// 		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/validator\" \n")
// 		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/types\" \n")
// 		buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/conv\" \n")
// 		buf.WriteString("	\"net/http\" \n")
// 		buf.WriteString(") \n\n")
// 		_, e = dtoFile.Write(buf.Bytes())
// 		if e != nil {
// 			showError(e)
// 			return
// 		}
// 	})
//
// 	file, e := astra.ParseFile(structFullName,
// 		astra.IgnoreVariables|astra.IgnoreConstants|astra.IgnoreFunctions|
// 			astra.IgnoreInterfaces|astra.IgnoreTypes|astra.IgnoreMethods)
// 	if e != nil {
// 		return e
// 	}
// 	if len(file.Structures) == 0 {
// 		return nil
// 	}
//
// 	var buf bytes.Buffer
// 	for _, stru := range file.Structures {
// 		if !isChecker(stru.Docs) {
// 			continue
// 		}
// 		//
// 		receiver := strings.ToLower(stru.Name[0:1])
// 		// 生成 Initializer 接口方法
// 		writeInit(&buf, receiver, stru)
// 		// 生成 checker 接口方法
// 		writeCheck(&buf, receiver, stru)
// 	}
// 	_, e = dtoFile.Write(buf.Bytes())
// 	return e
// }

// func writeInit(buf *bytes.Buffer, receiver string, stru types.Struct) {
// 	buf.WriteString(fmt.Sprintf("func (%s *%s) Init() error {\n", receiver, stru.Name))
// 	//
// 	for _, field := range stru.Fields {
// 		for k, v := range field.Tags {
// 			if k != "valid" || len(v) == 0 || util.SliceContains(v, "-") {
// 				continue
// 			}
// 			vv := make([]string, len(v))
// 			for _, s := range v {
// 				tags := strings.Split(s, "~")
// 				if tags[0] == "required" {
// 					switch field.Type.String() {
// 					case "string", "types.String":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilString \n", receiver, field.Name))
// 					case "int", "types.Int":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilInt \n", receiver, field.Name))
// 					case "int8", "types.Int8":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilInt8 \n", receiver, field.Name))
// 					case "int16", "types.Int16":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilInt16 \n", receiver, field.Name))
// 					case "int32", "types.Int32":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilInt32 \n", receiver, field.Name))
// 					case "int64", "types.Int64":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilInt64 \n", receiver, field.Name))
// 					case "uint", "types.Uint":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilUint \n", receiver, field.Name))
// 					case "uint8", "types.Uint8":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilUint8 \n", receiver, field.Name))
// 					case "uint16", "types.Uint16":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilUint16 \n", receiver, field.Name))
// 					case "uint32", "types.Uint32":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilUint32 \n", receiver, field.Name))
// 					case "uint64", "types.Uint64":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilUint64 \n", receiver, field.Name))
// 					case "float32", "types.Float32":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilFloat32 \n", receiver, field.Name))
// 					case "float64", "types.Float64":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilFloat64 \n", receiver, field.Name))
// 					case "time.Time":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.NilTime \n", receiver, field.Name))
// 					case "types.Time":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.Time{} \n", receiver, field.Name))
// 					case "bool", "types.Bool":
// 						buf.WriteString(fmt.Sprintf("	%s.%s = types.Bool(-1) \n", receiver, field.Name))
// 					}
// 					continue
// 				}
// 				vv = append(vv, s)
// 			}
// 		}
// 	}
// 	buf.WriteString("	return nil\n")
// 	buf.WriteString("}\n")
// }
//
// func writeCheck(buf *bytes.Buffer, receiver string, stru types.Struct) {
// 	buf.WriteString(fmt.Sprintf("func (%s *%s) Check() error {\n", receiver, stru.Name))
// 	//
// 	for _, field := range stru.Fields {
// 		for k, v := range field.Tags {
// 			if k != "valid" || len(v) == 0 || util.SliceContains(v, "-") {
// 				continue
// 			}
// 			vv := make([]string, len(v))
// 			for _, s := range v {
// 				tags := strings.Split(s, "~")
// 				if tags[0] == "required" {
// 					writeRequired(tags, field, buf, receiver)
// 					continue
// 				}
// 				vv = append(vv, s)
// 			}
// 			for _, s := range vv {
// 				tags := strings.Split(s, "~")
// 				writeDefault(tags, field, buf, receiver)
// 			}
// 		}
// 	}
// 	buf.WriteString("	return nil\n")
// 	buf.WriteString("}\n")
// }

// func writeDefault(tags []string, field types.StructField, b *bytes.Buffer, receiver string) {
// 	if fn, ok := validator.TagMap[tags[0]]; ok {
// 		if len(tags) == 1 {
// 			tags = append(tags, "Wrong "+field.Name+" format")
// 		}
// 		b.WriteString(fmt.Sprintf("	if !validator.%s(", fn))
// 		if field.Type.String() == "string" {
// 			b.WriteString(fmt.Sprintf("%s.%s", receiver, field.Name))
// 		} else if strings.HasPrefix(field.Type.String(), "types.") {
// 			b.WriteString(fmt.Sprintf("%s.%s.String()", receiver, field.Name))
// 		} else {
// 			b.WriteString(fmt.Sprintf("conv.Any2String(%s.%s)", receiver, field.Name))
// 		}
// 		b.WriteString(") {\n")
// 		b.WriteString(fmt.Sprintf("		return types.NewError(http.StatusBadRequest, \"%s\")\n", tags[1]))
// 		b.WriteString("	}\n")
// 		return
// 	}
// 	//
// 	pos1 := strings.Index(tags[0], "(")
// 	if pos1 > 0 {
// 		tag := tags[0][:pos1]
// 		if len(tags) == 1 {
// 			tags = append(tags, "Wrong "+field.Name+" "+tag)
// 		}
// 		pos2 := strings.Index(tags[0], ")")
// 		if pos2 < 0 {
// 			return
// 		}
// 		params := strings.Split(tags[0][pos1+1:pos2], "|")
// 		if fn, ok := validator.ParamTagMap[tag]; ok {
// 			b.WriteString(fmt.Sprintf("	if !validator.%s(", fn))
// 			if field.Type.String() == "string" {
// 				b.WriteString(fmt.Sprintf("%s.%s", receiver, field.Name))
// 			} else if strings.HasPrefix(field.Type.String(), "types.") {
// 				b.WriteString(fmt.Sprintf("%s.%s.String()", receiver, field.Name))
// 			} else {
// 				b.WriteString(fmt.Sprintf("conv.Any2String(%s.%s)", receiver, field.Name))
// 			}
// 			for _, param := range params {
// 				b.WriteString(", \"" + param + "\"")
// 			}
// 			b.WriteString(") {\n")
// 			b.WriteString(fmt.Sprintf("		return types.NewError(http.StatusBadRequest, \"%s\")\n", tags[1]))
// 			b.WriteString("	}\n")
// 			return
// 		}
// 	}
// }
//
// func writeRequired(tags []string, field types.StructField, buf *bytes.Buffer, receiver string) {
// 	if len(tags) == 1 {
// 		tags = append(tags, "required")
// 	}
// 	switch field.Type.String() {
// 	case "string", "types.String":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilString {\n", receiver, field.Name))
// 	case "int", "types.Int":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilInt {\n", receiver, field.Name))
// 	case "int8", "types.Int8":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilInt8 {\n", receiver, field.Name))
// 	case "int16", "types.Int16":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilInt16 {\n", receiver, field.Name))
// 	case "int32", "types.Int32":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilInt32 {\n", receiver, field.Name))
// 	case "int64", "types.Int64":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilInt64 {\n", receiver, field.Name))
// 	case "uint", "types.Uint":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilUint {\n", receiver, field.Name))
// 	case "uint8", "types.Uint8":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilUint8 {\n", receiver, field.Name))
// 	case "uint16", "types.Uint16":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilUint16 {\n", receiver, field.Name))
// 	case "uint32", "types.Uint32":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilUint32 {\n", receiver, field.Name))
// 	case "uint64", "types.Uint64":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilUint64 {\n", receiver, field.Name))
// 	case "float32", "types.Float32":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilFloat32 {\n", receiver, field.Name))
// 	case "float64", "types.Float64":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.NilFloat64 {\n", receiver, field.Name))
// 	case "time.Time", "types.Time":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s.IsZero() {\n", receiver, field.Name))
// 	case "bool", "types.Bool":
// 		buf.WriteString(fmt.Sprintf("	if %s.%s == types.Bool(-1) {\n", receiver, field.Name))
// 	default:
// 		return
// 		// case "bool":
// 		// 	return
// 		// default:
// 		// 	buf.WriteString(fmt.Sprintf("	if util.IsStructZeroValue(%s.%s) {\n", receiver, field.Name))
// 	}
// 	buf.WriteString(fmt.Sprintf("		return types.NewError(http.StatusBadRequest, \"%s is %s\")\n", field.Name, tags[1]))
// 	buf.WriteString("	}\n")
// }
//
// func isChecker(docs []string) bool {
// 	for _, doc := range docs {
// 		doc = strings.TrimLeft(doc, " /")
// 		if strings.EqualFold(doc, "checker") {
// 			return true
// 		}
// 	}
// 	return false
// }
