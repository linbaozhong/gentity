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
	"context"
	"github.com/linbaozhong/gentity/example/abc/internal/model/do"
)

// AccountDao 账号信息DAO
// @Namespace account
type AccountDao interface {
	// FindByID 根据ID查询账号信息
	// @Statement SELECT * FROM account WHERE id = :id
	FindByID(ctx context.Context, id int) (*do.Account, error)
	// FindByPhone 根据手机号查询账号信息
	// @Statement SELECT * FROM account WHERE login_name = :phone
	FindByPhone(ctx context.Context, phone string) (*do.Account, error)
	// FindByIDs 根据ID列表查询账号信息
	// @Statement SELECT * FROM account WHERE id IN (:ids)
	FindByIDs(ctx context.Context, ids []int) ([]*do.Account, error)
}

// CompanyDao 公司信息DAO
// @Namespace company
type CompanyDao interface {
	// FindByID 根据ID查询公司信息
	// @Statement SELECT company.id,account.id FROM company as c,account as a WHERE a.id = :id
	FindByID(ctx context.Context, id int) (*do.Company, error)
}
