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
	"github.com/linbaozhong/gentity/pkg/util"
	"github.com/linbaozhong/gentity/pkg/validator"
	"github.com/vetcher/go-astra"
	"github.com/vetcher/go-astra/types"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const dentityDTO = "gentity_dto.go"

func generateCheck(filename string, pkgPath string) error {
	if filename == dentityDTO {
		return nil
	}

	fset := token.NewFileSet()
	var src any
	var structFullName = filepath.Join(fullpath, filename)

	f, err := parser.ParseFile(fset, structFullName, src, parser.ParseComments)
	if err != nil {
		showError(err)
		return err
	}

	file, err := astra.ParseFile(structFullName,
		astra.IgnoreVariables|astra.IgnoreConstants|astra.IgnoreFunctions|
			astra.IgnoreInterfaces|astra.IgnoreTypes|astra.IgnoreMethods)
	if err != nil {
		showError(err)
		return err
	}
	if len(file.Structures) == 0 {
		return nil
	}
	//
	dtoFile, err := os.OpenFile(filepath.Join(fullpath, dentityDTO), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer dtoFile.Close()
	//
	var buf bytes.Buffer
	buf.WriteString("package " + f.Name.Name + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/validator\" \n")
	buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/types\" \n")
	buf.WriteString("	\"github.com/linbaozhong/gentity/pkg/conv\" \n")
	buf.WriteString(") \n\n")

	for _, stru := range file.Structures {
		if !isChecker(stru.Docs) {
			continue
		}
		//
		receiver := strings.ToLower(stru.Name[0:1])
		//
		buf.WriteString(fmt.Sprintf("func (%s *%s) Check() error {\n", receiver, stru.Name))
		//
		for _, field := range stru.Fields {
			for k, v := range field.Tags {
				if k != "valid" {
					continue
				}
				if util.SliceContains(v, "-") {
					continue
				}
				if len(v) == 0 {
					continue
				}
				//
				vv := make([]string, len(v))
				for _, s := range v {
					tags := strings.Split(s, "~")
					if tags[0] == "required" {
						writeRequired(tags, field, &buf, receiver)
						continue
					}
					vv = append(vv, s)
				}
				for _, s := range vv {
					tags := strings.Split(s, "~")
					writeDefault(tags, field, &buf, receiver)
				}
			}
		}
		buf.WriteString("	return nil\n")
		buf.WriteString("}\n")
	}
	_, err = dtoFile.Write(buf.Bytes())
	if err != nil {
		showError(err)
		return err
	}

	return nil
}

func writeDefault(tags []string, field types.StructField, b *bytes.Buffer, receiver string) {
	if fn, ok := validator.TagMap[tags[0]]; ok {
		if len(tags) == 1 {
			tags = append(tags, "Wrong "+tags[0]+" format")
		}
		b.WriteString(fmt.Sprintf("	if !validator.%s(", fn))
		if field.Type.String() == "string" {
			b.WriteString(fmt.Sprintf("%s.%s", receiver, field.Name))
		} else if strings.HasPrefix(field.Type.String(), "types.") {
			b.WriteString(fmt.Sprintf("%s.%s.String()", receiver, field.Name))
		} else {
			b.WriteString(fmt.Sprintf("conv.Any2String(%s.%s)", receiver, field.Name))
		}
		b.WriteString(") {\n")
		b.WriteString(fmt.Sprintf("		return types.NewError(30001, \"%s\")\n", tags[1]))
		b.WriteString("	}\n")
		return
	}
	//
	pos1 := strings.Index(tags[0], "(")
	if pos1 > 0 {
		tag := tags[0][:pos1]
		if len(tags) == 1 {
			tags = append(tags, "Wrong "+field.Name+" "+tag)
		}
		pos2 := strings.Index(tags[0], ")")
		if pos2 < 0 {
			return
		}
		params := strings.Split(tags[0][pos1+1:pos2], "|")
		if fn, ok := validator.ParamTagMap[tag]; ok {
			b.WriteString(fmt.Sprintf("	if !validator.%s(", fn))
			if field.Type.String() == "string" {
				b.WriteString(fmt.Sprintf("%s.%s", receiver, field.Name))
			} else if strings.HasPrefix(field.Type.String(), "types.") {
				b.WriteString(fmt.Sprintf("%s.%s.String()", receiver, field.Name))
			} else {
				b.WriteString(fmt.Sprintf("conv.Any2String(%s.%s)", receiver, field.Name))
			}
			for _, param := range params {
				b.WriteString(", \"" + param + "\"")
			}
			b.WriteString(") {\n")
			b.WriteString(fmt.Sprintf("		return types.NewError(30001, \"%s\")\n", tags[1]))
			b.WriteString("	}\n")
			return
		}
	}
}

func writeRequired(tags []string, field types.StructField, buf *bytes.Buffer, receiver string) {
	if len(tags) == 1 {
		tags = append(tags, "required")
	}
	switch field.Type.String() {
	case "string", "types.String":
		buf.WriteString(fmt.Sprintf("	if %s.%s == \"\" {\n", receiver, field.Name))
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		fallthrough
	case "types.Int", "types.Int8", "types.Int16", "types.Int32", "types.Int64", "types.Uint", "types.Uint8", "types.Uint16", "types.Uint32", "types.Uint64", "types.Float32", "types.Float64":
		buf.WriteString(fmt.Sprintf("	if %s.%s == 0 {\n", receiver, field.Name))
	default:
		fmt.Println(field.Type.String())
		buf.WriteString(fmt.Sprintf("	if %s.%s == nil {\n", receiver, field.Name))
	}
	buf.WriteString(fmt.Sprintf("		return types.NewError(30001, \"%s is %s\")\n", field.Name, tags[1]))
	buf.WriteString("	}\n")
}

func isChecker(docs []string) bool {
	for _, doc := range docs {
		doc = strings.TrimLeft(doc, " /")
		if strings.EqualFold(doc, "checker") {
			return true
		}
	}
	return false
}
