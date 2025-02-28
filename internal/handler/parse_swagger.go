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
	"fmt"
	"github.com/vetcher/go-astra"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

func parseSwagger(filename, pkgPath string) {
	var _structFullName = filepath.Join(fullpath, filename)
	_file, e := astra.ParseFile(_structFullName,
		astra.IgnoreVariables|astra.IgnoreConstants|astra.IgnoreFunctions|
			astra.IgnoreInterfaces|astra.IgnoreTypes)
	if e != nil {
		showError(e)
		return
	}

	// 查找包含 models.User 的 Swagger 注释
	for _, decl := range _file.Methods {
		if decl.Docs != nil {
			for _, comment := range decl.Docs {
				fmt.Println(1, comment)
				if strings.Contains(comment, "dto.UserRegisterReq") {
					// 找到包含 models.User 的注释，现在查找结构体定义
					// structDef := findStructDefinition(pkg, "dto", "UserRegisterReq")
					// if structDef != nil {
					// 	printStructFields(structDef)
					// }
					for _, imp := range _file.Imports {
						if strings.Contains(imp.Package, "dto") {
							pkg := types.NewPackage(imp.Package, "")
							structDef := findStructDefinition(pkg, "dto", "UserRegisterReq")
							if structDef != nil {
								printStructFields(structDef)
							}
						}
					}
				}
			}
		}
	}
}

// 查找结构体定义
func findStructDefinition(pkg *types.Package, structPackage, structName string) *types.Struct {
	fmt.Println(1, pkg.Path(), structPackage)
	// 检查当前包
	// if pkg.Path() == structPackage {
	obj := pkg.Scope().Lookup(structName)
	if obj != nil && obj.Type().Underlying() != nil {
		if structType, ok := obj.Type().Underlying().(*types.Struct); ok {
			return structType
		}
	}
	// }

	// 检查导入的包
	for _, imp := range pkg.Imports() {
		fmt.Println(2, imp.Path(), structPackage)
		if imp.Path() == structPackage {
			obj := imp.Scope().Lookup(structName)
			if obj != nil && obj.Type().Underlying() != nil {
				if structType, ok := obj.Type().Underlying().(*types.Struct); ok {
					return structType
				}
			}
		}
	}

	return nil
}

// 打印结构体字段
func printStructFields(structType *types.Struct) {
	fmt.Println("结构体字段信息:")
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fmt.Printf("  字段名: %s, 类型: %s\n", field.Name(), field.Type())
	}
}

// 将 AST 节点转换为字符串
func astToString(node ast.Node) string {
	var buf strings.Builder
	_ = ast.Fprint(&buf, token.NewFileSet(), node, nil)
	return buf.String()
}
