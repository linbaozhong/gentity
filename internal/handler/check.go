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
