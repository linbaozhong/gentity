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

package dao

import (
	"fmt"
	"github.com/linbaozhong/gentity/example/model/define/table/tblcompanydocument"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"testing"
)

func TestDefine(t *testing.T) {
	_dai := DataAccessInterface{
		NameSpace: do.CompanyDocumentTableName,
		Children: []DataAccessInterface{
			{
				Name:   "GetById",
				Method: ace.Method_Get,
				Input:  ace.Where(tblcompanydocument.Id.Eq(1), tblcompanydocument.State.Eq(nil)).Order(tblcompanydocument.Ctime),
				Output: &do.Company{},
			},
			{
				Name:   "ListByCompanyId",
				Method: ace.Method_List,
				Input:  ace.Where(tblcompanydocument.Company.Eq(1), tblcompanydocument.State.Eq(nil)).Order(tblcompanydocument.Ctime),
				Output: []do.Company{},
			},
		},
	}
	Run(_dai)
}

func Run(dai DataAccessInterface) {
	if dai.NameSpace == "" && dai.Table == "" {
		return
	}
	if dai.NameSpace == "" {
		dai.NameSpace = dai.Table
	} else if dai.Table == "" {
		dai.Table = dai.NameSpace
	}

	if len(dai.Children) > 0 {
		for _, _dai := range dai.Children {
			if _dai.NameSpace == "" {
				_dai.NameSpace = dai.NameSpace
			} else {
				_dai.NameSpace = dai.NameSpace + "." + _dai.NameSpace
			}
			Run(_dai)
		}
		return
	}

	fmt.Println(`// `, dai.NameSpace+"."+dai.Name)
	if dai.Title != "" {
		fmt.Println(`// @Title `, dai.Title)
	}
	if dai.Description != "" {
		fmt.Println(`// @Description `, dai.Description)
	}
	fmt.Println(`func `, dai.Name, `(ctx context.Context,`)
}
