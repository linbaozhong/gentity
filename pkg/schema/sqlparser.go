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

package schema

import (
	"bytes"
	"github.com/linbaozhong/gentity/pkg/util"
	"github.com/linbaozhong/sqlparser"
	"io"
	"os"
)

func Sql2Struct(buf []byte) ([]byte, error) {
	return reader2Struct(bytes.NewReader(buf))
}

func SqlFile2Struct(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return reader2Struct(f)
}

func reader2Struct(r io.Reader) ([]byte, error) {
	parser := sqlparser.NewParser(r)
	schema, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for _, table := range schema {
		buf.WriteString("type " + util.ParseField(table.Name) + " struct {\n")
		for _, col := range table.Columns {
			buf.WriteString("\t" + util.ParseField(col.Name) + "\t" + util.ParseFieldType(col.Type, col.Size))
			buf.WriteString("\t`json:\"" + col.Name + "\" db:\"'" + col.Name + "'") //
			if util.SliceContains(table.PrimaryKey, col.Name) {
				buf.WriteString(" pk")
			}
			if col.AutoIncr {
				buf.WriteString(" auto")
			}
			buf.WriteString("\"` \n")
		}
		buf.WriteString("}")
	}
	return buf.Bytes(), nil
}
