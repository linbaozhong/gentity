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
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/util"
	"github.com/linbaozhong/sqlparser"
	"io"
	"os"
	"strings"
)

func Sql2Struct(buf []byte, packageName string) ([]byte, error) {
	return reader2Struct(bytes.NewReader(buf), packageName)
}

func SqlFile2Struct(filePath, packageName string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return reader2Struct(f, packageName)
}

func reader2Struct(r io.Reader, packageName string) ([]byte, error) {
	parser := sqlparser.NewParser(r)
	schema, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	//
	//for _, x := range schema["company"].ColumnsX {
	//	fmt.Printf("%+v \n", x)
	//}

	var buf bytes.Buffer
	buf.WriteString("package " + packageName + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("	\"time\" \n")
	buf.WriteString(") \n\n")
	for _, table := range schema {
		buf.WriteString("// tablename " + table.Name + "\n")
		buf.WriteString("type " + util.ParseField(table.Name) + " struct {\n")
		for _, col := range table.ColumnsX {
			buf.WriteString("\t" + util.ParseField(col.Name) + "\t" + util.ParseFieldType(col.Type, col.Size))
			buf.WriteString("\t`json:\"" + col.Name + "\" db:\"'" + col.Name + "'") //
			if strings.ToUpper(col.Key) == dialect.PrimaryKey {
				buf.WriteString(" pk")
			}
			if col.AutoIncr {
				buf.WriteString(" auto")
			}
			buf.WriteString(fmt.Sprintf("\"`	// %s\n", strings.TrimSpace(col.Comment)))
		}
		buf.WriteString("} \n\n")
	}
	return buf.Bytes(), nil
}

func DB2Struct(tables map[string][]dialect.Column, packageName string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("package " + packageName + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("	\"time\" \n")
	buf.WriteString(") \n\n")

	for table, columns := range tables {
		buf.WriteString("// tablename " + table + "\n")
		buf.WriteString("type " + util.ParseField(table) + " struct {\n")
		for _, col := range columns {
			buf.WriteString("\t" + util.ParseField(col.Name) + "\t" + util.ParseFieldType(col.Type, col.Size))
			buf.WriteString("\t`json:\"" + col.Name + "\" db:\"'" + col.Name + "'") //
			if strings.ToUpper(col.Key) == dialect.PrimaryKey {
				buf.WriteString(" pk")
			}
			if col.AutoIncr {
				buf.WriteString(" auto")
			}
			buf.WriteString(fmt.Sprintf("\"`	// %s\n", strings.ReplaceAll(col.Comment, "\n", "")))
		}
		buf.WriteString("} \n\n")
	}
	return buf.Bytes(), nil
}
