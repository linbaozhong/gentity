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
	"github.com/linbaozhong/gentity/pkg/sqlparser"
	"strconv"
	"strings"
	"unicode"
)

func ParseField(fieldName string) string {
	var (
		isUpper bool
		upper   strings.Builder
	)
	for i, ch := range fieldName {
		if i == 0 {
			upper.WriteRune(unicode.ToUpper(ch))
		} else {
			if ch == 95 {
				isUpper = true
				continue
			}
			if isUpper {
				upper.WriteRune(unicode.ToUpper(ch))
				isUpper = false
			} else {
				upper.WriteRune(ch)
			}
		}
	}
	return upper.String()
}

func ParseFieldType(col *sqlparser.Column) string {
	switch strings.ToUpper(col.Type) {
	case "BIGINT":
		if col.Unsigned {
			return "uint64"
		}
		return "int64"
	case "INT", "MEDIUMINT":
		if col.Unsigned {
			return "uint32"
		}
		return "int32"
	case "SMALLINT":
		if col.Unsigned {
			return "uint16"
		}
		return "int16"
	case "TINYINT":
		if col.Unsigned {
			return "uint8"
		}
		return "int8"
	case "VARCHAR", "LONGTEXT", "MEDIUMTEXT", "TEXT":
		return "string"
	case "BIT", "BOOLEAN":
		return "bool"
	case "FLOAT":
		return "float32"
	case "DOUBLE", "DECIMAL", "NUMERIC":
		return "float64"
	case "TIMESTAMP", "DATETIME", "DATE", "TIME":
		return "time.Time"
	default:
		return "any" // 对于未明确映射的类型，使用接口类型作为占位符
	}
}

func ParseFieldAceType(col *sqlparser.Column) string {
	switch strings.ToUpper(col.Type) {
	case "BIGINT":
		if col.Unsigned {
			return "types.BigInt"
		}
		return "types.Money"
	case "INT", "MEDIUMINT":
		if col.Unsigned {
			return "types.Uint32"
		}
		return "types.Int32"
	case "SMALLINT":
		if col.Unsigned {
			return "types.Uint16"
		}
		return "types.Int16"
	case "TINYINT":
		if col.Unsigned {
			return "types.Uint8"
		}
		return "types.Int8"
	case "VARCHAR", "LONGTEXT", "MEDIUMTEXT", "TEXT", "JSON":
		return "types.String"
	case "BIT", "BOOLEAN":
		return "types.Bool"
	case "FLOAT":
		return "types.Float32"
	case "DOUBLE", "DECIMAL", "NUMERIC":
		return "types.Float64"
	case "TIMESTAMP", "DATETIME", "DATE", "TIME":
		return "types.Time"
	default:
		return "any" // 对于未明确映射的类型，使用接口类型作为占位符
	}
}
func ParseFieldSize(col *sqlparser.Column) string {
	switch strings.ToUpper(col.Type) {
	case "VARCHAR", "LONGTEXT", "MEDIUMTEXT", "TEXT":
		return " size:" + strconv.Itoa(col.Size)
	// case "BIGINT":
	// 	return ""
	// case "INT", "MEDIUMINT":
	// 	return ""
	// case "SMALLINT":
	// 	return ""
	// case "TINYINT":
	// 	return ""
	// case "BIT":
	// 	return ""
	// case "FLOAT":
	// 	return ""
	// case "DOUBLE":
	// 	return ""
	// case "DECIMAL", "NUMERIC":
	// 	return ""
	// case "TIMESTAMP", "DATETIME", "DATE", "TIME":
	// 	return ""
	default:
		if col.Precision > 0 {
			if col.Scale > 0 {
				return " size:" + strconv.Itoa(col.Precision) + "|" + strconv.Itoa(col.Scale)
			}
			return " size:" + strconv.Itoa(col.Precision)
		}
		return "" // 对于未明确映射的类型，使用接口类型作为占位符
	}
}
