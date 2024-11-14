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
	"github.com/linbaozhong/gentity/pkg/sqlparser"
	"github.com/linbaozhong/gentity/pkg/util"
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

	var buf bytes.Buffer
	buf.WriteString("// Code generated by gentity. DO NOT EDIT.\n\n")
	buf.WriteString("package " + packageName + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("	\"time\" \n")
	buf.WriteString("	atype \"github.com/linbaozhong/gentity/pkg/ace/types\" \n")
	// buf.WriteString("	\"sync/atomic\" \n")
	buf.WriteString(") \n\n")
	//
	// buf.WriteString("var nocopy_uint64 uint64\n\n")
	// buf.WriteString("// noCopy 用于防止嵌入它的结构体被复制。\n")
	// buf.WriteString("// 它实现了Locker接口的Lock和Unlock方法，但这些方法什么都不做。\n")
	// buf.WriteString("type noCopy struct {\n")
	// buf.WriteString("\tnocopy_uuid uint64 `json:\"-\"`")
	// buf.WriteString(" // 内部留用，禁止外部赋值\n")
	// buf.WriteString("}\n\n")
	// buf.WriteString("// UUID 只是为了实现接口方法，外部不建议调用。\n")
	// buf.WriteString("func (a *noCopy) UUID() uint64 {\n\tif a.nocopy_uuid == 0 {\n\t\ta.nocopy_uuid = atomic.AddUint64(&nocopy_uint64, 1)\n\t}\n\treturn a.nocopy_uuid\n}\n\nfunc (p *noCopy) Lock()   {}\nfunc (p *noCopy) Unlock() {}\n")

	for _, table := range schema {
		buf.WriteString("// tablename " + table.Name + "\n")
		buf.WriteString("// cache time.Minute time.Minute 1000 \n")
		buf.WriteString("type " + util.ParseField(table.Name) + " struct {\n")
		buf.WriteString("\tatype.AceModel\n")
		// buf.WriteString("\tnoCopy\n")
		for _, col := range table.ColumnsX {
			buf.WriteString("\t" + util.ParseField(col.Name) + "\t" + util.ParseFieldType(col.Type, col.Size, col.Unsigned))
			buf.WriteString("\t`json:\"" + col.Name + ",omitempty\" db:\"'" + col.Name + "'") //
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

func DB2Struct(tables map[string][]sqlparser.Column, packageName string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by gentity. DO NOT EDIT.\n\n")
	buf.WriteString("package " + packageName + "\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("	\"time\" \n")
	buf.WriteString("	atype \"github.com/linbaozhong/gentity/pkg/ace/types\" \n")
	// buf.WriteString("	\"sync/atomic\" \n")
	buf.WriteString(") \n\n")
	//
	// buf.WriteString("var nocopy_uint64 uint64\n\n")
	// buf.WriteString("// noCopy 用于防止嵌入它的结构体被复制。\n")
	// buf.WriteString("// 它实现了Locker接口的Lock和Unlock方法，但这些方法什么都不做。\n")
	// buf.WriteString("type noCopy struct {\n")
	// buf.WriteString("\tnocopy_uuid uint64 `json:\"-\"`")
	// buf.WriteString(" // 内部留用，禁止外部赋值\n")
	// buf.WriteString("}\n\n")
	// buf.WriteString("// UUID 只是为了实现接口方法，外部不建议调用。\n")
	// buf.WriteString("func (a *noCopy) UUID() uint64 {\n\tif a.nocopy_uuid == 0 {\n\t\ta.nocopy_uuid = atomic.AddUint64(&nocopy_uint64, 1)\n\t}\n\treturn a.nocopy_uuid\n}\n\nfunc (p *noCopy) Lock()   {}\nfunc (p *noCopy) Unlock() {}\n")

	for table, columns := range tables {
		buf.WriteString("// tablename " + table + "\n")
		buf.WriteString("// cache time.Minute time.Minute 1000 \n")
		buf.WriteString("type " + util.ParseField(table) + " struct {\n")
		buf.WriteString("\tatype.AceModel\n")
		// buf.WriteString("\tnoCopy\n")
		for _, col := range columns {
			buf.WriteString("\t" + util.ParseField(col.Name) + "\t" + util.ParseFieldType(col.Type, col.Size, col.Unsigned))
			buf.WriteString("\t`json:\"" + col.Name + ",omitempty\" db:\"'" + col.Name + "'") //
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
