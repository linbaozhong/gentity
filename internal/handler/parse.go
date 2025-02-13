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

const (
	jsonTag  = "json"
	formTag  = "form"
	urlTag   = "url"
	paramTag = "param"
	dbTag    = "db"
	validTag = "valid"
)

func parseType(fieldType string) (kind, typeName string) {
	if strings.HasPrefix(fieldType, "[]") {
		return "slice", fieldType[2:]
	} else if strings.HasPrefix(fieldType, "*") {
		return "ptr", fieldType[1:]
	} else {
		return "struct", fieldType
	}
}

func parseFile(filename, pkgPath string, tags ...string) ([]TempData, error) {
	var structFullName = filepath.Join(fullpath, filename)
	astFile, err := getAst(structFullName)
	if err != nil {
		showError(err)
		return nil, err
	}

	file, err := astra.ParseFile(structFullName,
		astra.IgnoreVariables|astra.IgnoreConstants|astra.IgnoreFunctions|
			astra.IgnoreInterfaces|astra.IgnoreTypes|astra.IgnoreMethods)
	if err != nil {
		showError(err)
		return nil, err
	}
	if len(file.Structures) == 0 {
		return nil, nil
	}

	imps := make([]string, 0, len(file.Imports))
	for _, im := range file.Imports {
		imps = append(imps, im.Package)
	}
	// 文件中全部符合条件的struct模板
	tplsData := make([]TempData, 0, len(file.Structures))
	// 遍历struct
	for _, stru := range file.Structures {
		// 如果struct名称为空,或者是一个私有struct,或者field为空,返回
		if len(stru.Name) == 0 || len(stru.Fields) == 0 /*|| unicode.IsLower(rune(stru.Name[0]))*/ {
			continue
		}

		tempData := TempData{
			ParseTag:    make([]string, 0, 2),
			ModulePath:  pkgPath,
			PackageName: astFile.Name.Name,
			Imports:     imps,
			StructName:  stru.Name,
			CacheData:   "time.Minute",
			CacheList:   "time.Minute",
			CacheLimit:  "1000",
			FileName:    structFullName,
			PrimaryKey:  Field{},
			RelationX:   Relation{},
			Columns:     make([]Field, 0, 20),
		}
		// 解析struct文档
		parseDocs(&tempData, stru.Docs, tags...)
		if len(tempData.ParseTag) == 0 {
			continue
		}

		for _, field := range stru.Fields {
			if len(field.Tags) == 0 {
				continue
			}
			var (
				pk  string
				rw  string // 禁止读写 -，只读<-，只写->
				ref string // 关系键
			)
			var _namejson = Field{}
			for k, v := range field.Tags {
				switch k {
				case jsonTag, urlTag, formTag, paramTag:
					_namejson.Json = parseJson(v) // json_name
				case dbTag:
					_namejson.Col, pk, rw, ref = parseTagsForDB(v) // column_name
					if len(ref) == 0 {
						continue
					}
					_ref := strings.Split(ref, "|")
					if len(_ref) != 2 {
						continue
					}
					tempData.RelationX.Kind, tempData.RelationX.Type = parseType(field.Type.String())
					tempData.RelationX.Name = field.Name
					tempData.RelationX.Field = _ref[0]
					tempData.RelationX.Foreign = _ref[1]
				case validTag:
					_namejson.Valids = moveToFront(v, "required")
				}
			}
			_namejson.Name = field.Name
			_namejson.Type = field.Type.String()
			_namejson.Rw = rw

			if strings.HasPrefix(_namejson.Type, "types") {
				tempData.HasCustomType = true
			}
			if _namejson.Type == "time.Time" {
				tempData.HasTime = true
			}

			if _namejson.Col == "" {
				if _namejson.Json.Name == "" {
					_namejson.Col = getFieldName(field.Name)
				} else {
					_namejson.Col = _namejson.Json.Name
				}
			}

			tempData.Columns = append(tempData.Columns, _namejson)
			if pk != "" {
				tempData.HasPrimaryKey = true
				// tempData.Keys = append(tempData.Keys, _namejson)
				tempData.PrimaryKey = _namejson
				// // 主键是否是自增
				// if tempData.PrimaryKey[4] == "<-" {
				// 	tempData.HasCustomType = true
				// }
			}
			if r := _namejson.Type[:1]; r == "*" || r == "[" || r == "m" {
				tempData.HasRef = true
			}
		}
		tplsData = append(tplsData, tempData)
	}

	return tplsData, err
}

func parseJson(keys []string) jsonObj {
	var json jsonObj
	for _, key := range keys {
		switch key {
		case "omitempty":
			json.OmitEmpty = true
		case "omitzero":
			json.OmitZero = true
		default:
			json.Name = key
		}
	}
	return json
}

// moveToFront 函数用于将指定字符串移到切片的第一个位置
func moveToFront(slice []string, target string) []string {
	// 遍历切片，查找目标字符串的索引
	for i, str := range slice {
		if str == target {
			// 如果找到目标字符串，将其移动到切片的第一个位置
			if i > 0 {
				// 把目标字符串暂存起来
				temp := slice[i]
				// 将目标字符串之前的元素依次向后移动一位
				copy(slice[1:i+1], slice[0:i])
				// 将目标字符串放到切片的第一个位置
				slice[0] = temp
			}
			break
		}
	}
	return slice
}
func parseDocs(tmp *TempData, docs []string, tags ...string) {
	for _, doc := range docs {
		doc = strings.TrimLeft(doc, " /")

		for _, tag := range tags {
			if strings.Contains(doc, tag) {
				tmp.ParseTag = append(tmp.ParseTag, tag)
				if tag == "tablename" {
					tmp.TableName = strings.TrimSpace(strings.TrimLeft(doc, tag))
					break
				}
			}
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

// columnName 列名
// key 主键
// rw 读写标志
// ref 关系键
// fk 关系外键
func parseTagsForDB(matchs []string) (columnName, key, rw, ref string) {
	s := strings.Split(matchs[0], " ")
	if len(s) == 1 {
		if strings.HasPrefix(s[0], "ref:") {
			ref = s[0][4:]
		} else if s[0] == "-" || s[0] == "->" || s[0] == "<-" {
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
		if strings.HasPrefix(v, "ref:") {
			ref = v[4:]
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
