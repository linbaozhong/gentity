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
	"github.com/vetcher/go-astra"
	"path/filepath"
	"strings"
	"unicode"
)

func parseFile(filename, pkgPath string) error {
	tempData := new(TempData)
	tempData.ModulePath = pkgPath

	var structFullName = filepath.Join(fullpath, filename)
	astFile, err := getAst(structFullName)
	if err != nil {
		showError(err)
		return err
	}

	tempData.PackageName = astFile.Name.Name

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

	for _, im := range file.Imports {
		tempData.Imports = append(tempData.Imports, im.Package)
	}

	for _, stru := range file.Structures {
		tempData.TableName = ""
		tempData.HasTime = false
		tempData.HasString = false
		tempData.HasConvert = false
		tempData.HasCache = false
		tempData.HasPrimaryKey = false
		tempData.HasState = false
		tempData.CacheData = "time.Minute"
		tempData.CacheList = "time.Minute"
		tempData.CacheLimit = "1000"
		tempData.PrimaryKey = nil
		// tempData.PrimaryKeyName = ""
		// tempData.Keys = make([][]string, 0, 1)
		tempData.Columns = make([][]string, 0, 20)
		tempData.FileName = structFullName
		tempData.StructName = stru.Name
		// 解析struct文档
		parseDocs(tempData, stru.Docs)
		if tempData.TableName == "" {
			continue
		}

		for _, field := range stru.Fields {
			if len(field.Tags) == 0 {
				continue
			}
			var (
				pk string
				rw string // 禁止读写 -，只读<-，只写->
			)
			var _namejson = make([]string, 5)
			for k, v := range field.Tags {
				if k == "json" {
					_namejson[2] = v[0] // json_name
				} else if k == "db" {
					_namejson[1], pk, rw = parseTagsForDB(v) // column_name
					// } else if k == "default" {
					// 	_namejson[3] = v[0]
				}
			}
			_namejson[0] = field.Name
			_namejson[4] = rw
			_namejson[3] = field.Type.String()
			switch _namejson[3] {
			case "types.Time":
				tempData.HasTime = true
			case "time.Time":
				tempData.HasTime = true
			case "string":
				tempData.HasString = true
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"float32", "float64", "bool", "types.Money":
				tempData.HasConvert = true
			}

			if _namejson[1] == "" {
				if _namejson[2] == "" {
					_namejson[1] = getFieldName(field.Name)
				} else {
					_namejson[1] = _namejson[2]
				}
			}

			tempData.Columns = append(tempData.Columns, _namejson)
			if pk != "" {
				tempData.HasPrimaryKey = true
				// tempData.Keys = append(tempData.Keys, _namejson)
				tempData.PrimaryKey = _namejson
				// 主键是否是自增
				if tempData.PrimaryKey[4] == "<-" {
					tempData.HasTime = true
				}
			}
			if _namejson[1] == "state" {
				tempData.HasState = true
			}
		}
		// 如果struct名称为空,或者是一个私有struct,或者field为空,返回
		if len(tempData.StructName) == 0 ||
			tempData.StructName[:1] == strings.ToLower(tempData.StructName[:1]) ||
			len(tempData.Columns) == 0 {
			return nil
		}

		// 写table文件
		err = tempData.writeTable(filepath.Join(tablePath, "tbl"+tempData.TableName))
		if err != nil {
			showError(err.Error())
			return err
		}

		// 写model文件
		err = tempData.writeToModel(filename)
		if err != nil {
			showError(err)
			return err
		}

		// 写dal文件
		err = tempData.writeBuild(daoPath)
		if err != nil {
			showError(err.Error())
			return err
		}
	}

	return err
}

func parseDocs(tmp *TempData, docs []string) {
	for _, doc := range docs {
		doc = strings.TrimLeft(doc, " /")
		if strings.Contains(doc, "tablename") {
			tmp.TableName = strings.TrimSpace(strings.TrimLeft(doc, "tablename"))
			continue
		}
		if strings.HasPrefix(doc, "cache ") {
			tmp.HasCache = true
			cache := strings.Replace(strings.TrimSpace(strings.TrimLeft(doc, "cache")), "  ", " ", -1)
			caches := strings.Split(cache, " ")
			if len(caches) >= 3 {
				tmp.CacheData = caches[0]
				tmp.CacheList = caches[1]
				tmp.CacheLimit = caches[2]
			}
		}
	}
}

func parseTagsForDB(matchs []string) (columnName string, key string, rw string) {
	s := strings.Split(strings.ToLower(matchs[0]), " ")
	if len(s) == 1 {
		if s[0] == "-" || s[0] == "->" || s[0] == "<-" {
			rw = s[0]
		} else {
			columnName = strings.Replace(s[0], "'", "", -1)
		}
		return
	}
	col := &columnName
	k := new(string)
	for _, v := range s {
		if v == "" {
			continue
		}
		if v[:1] == "'" {
			*col = strings.Replace(v, "'", "", -1)
			continue
		}
		if v == "pk" {
			k = col
			continue
		}
		if v == "-" || v == "->" || v == "<-" {
			rw = v
		}

		if v == "auto" {
			rw = "<-"
		}
	}
	key = *k
	return
}

func getFieldName(name string) string {
	bs := bytes.NewBuffer([]byte{})

	pre_lower := true // 前一个字母是小写
	for i, s := range []rune(name) {
		// 如果是大写字母
		if unicode.IsUpper(s) {
			r := unicode.ToLower(s)
			if i == 0 {
				bs.WriteRune(r)
			} else {
				if pre_lower {
					bs.WriteByte(byte(95)) // 写下划线
				}
				bs.WriteRune(r)
			}
			pre_lower = false
			continue
		}
		pre_lower = true
		bs.WriteRune(s)
	}
	return bs.String()
}
