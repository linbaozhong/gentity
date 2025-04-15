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

package dialect

//
// import "bytes"
//
// func Or(fns ...Condition) Condition {
// 	return func() (string, any) {
// 		if len(fns) == 0 {
// 			return "", nil
// 		}
// 		var (
// 			buf    bytes.Buffer
// 			params = make([]any, 0, len(fns))
// 		)
//
// 		buf.WriteString(Operator_or + "(")
// 		for i, fn := range fns {
// 			cond, val := fn()
// 			if i > 0 {
// 				buf.WriteString(Operator_and)
// 			}
// 			buf.WriteString(cond)
// 			if vals, ok := val.([]any); ok {
// 				params = append(params, vals...)
// 			} else {
// 				params = append(params, val)
// 			}
// 		}
// 		buf.WriteString(")")
//
// 		return buf.String(), params
// 	}
// }
//
// func And(fns ...Condition) Condition {
// 	return func() (string, any) {
// 		if len(fns) == 0 {
// 			return "", nil
// 		}
// 		var (
// 			buf    bytes.Buffer
// 			params = make([]any, 0, len(fns))
// 		)
//
// 		buf.WriteString(Operator_and + "(")
// 		for i, fn := range fns {
// 			cond, val := fn()
// 			if i > 0 {
// 				buf.WriteString(Operator_or)
// 			}
// 			buf.WriteString(cond)
// 			if vals, ok := val.([]any); ok {
// 				params = append(params, vals...)
// 			} else {
// 				params = append(params, val)
// 			}
// 		}
// 		buf.WriteString(")")
//
// 		return buf.String(), params
// 	}
// }
//
// // Asc 函数用于创建一个升序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// // 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "ASC" 和指定的字段列表。
// // 该函数可用于指定查询结果按指定字段进行升序排序。
// func Asc(fs ...Field) Order {
// 	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "ASC" 和字段列表
// 	return func() (string, []Field) {
// 		// 返回升序操作符
// 		return Operator_Asc, fs
// 	}
// }
//
// // Desc 函数用于创建一个降序排序的规则。它接收可变数量的 dialect.Field 类型的参数，
// // 返回一个实现了 dialect.Order 接口的函数，该函数会返回排序操作符 "DESC" 和指定的字段列表。
// // 该函数可用于指定查询结果按指定字段进行降序排序。
// func Desc(fs ...Field) Order {
// 	// 返回一个匿名函数，该函数实现了 dialect.Order 接口，返回排序操作符 "DESC" 和字段列表
// 	return func() (string, []Field) {
// 		// 返回降序操作符
// 		return Operator_Desc, fs
// 	}
// }
