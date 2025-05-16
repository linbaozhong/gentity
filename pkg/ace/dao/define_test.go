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
	"github.com/linbaozhong/gentity/example/model/define/table/tblcompany"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"testing"
)

func TestDefine(t *testing.T) {
	_dai := DataAccessInterface{
		Name:   "GetCompany",
		Table:  do.CompanyTableName,
		Method: ace.Method_Get,
		Input:  ace.Where(tblcompany.Id.Eq(1), tblcompany.State.Eq(nil)).Order(tblcompany.Ctime),
		Output: do.Company{},
	}
	RegisterDpi(_dai)
	Run()
}
