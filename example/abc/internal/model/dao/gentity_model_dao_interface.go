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

// // AccountDao 账号信息DAO
// // @Namespace account
// type AccountDao interface {
// 	// FindByID 根据ID查询账号信息
// 	// @Statement SELECT * FROM account WHERE id = :id
// 	FindByID(ctx context.Context, id int) (*do.Account, error)
// 	// FindByPhone 根据手机号查询账号信息
// 	// @Statement SELECT * FROM account WHERE login_name = :phone
// 	FindByPhone(ctx context.Context, phone string) (*do.Account, error)
// 	// FindByIDs 根据ID列表查询账号信息
// 	// @Statement SELECT * FROM account WHERE id IN (:ids)
// 	FindByIDs(ctx context.Context, ids []int) ([]*do.Account, error)
// }

// CompanyDao 公司信息DAO
// @Namespace company
type CompanyDao interface {
	// FindByID 根据ID查询公司信息
	// @Statement SELECT company.id,account.id FROM company as c,account as a WHERE a.id = :do.Account.ID
	FindByID(ctx context.Context, id int) (comp *do.Company, err error)
	// Create 创建公司信息
	// @Statement INSERT INTO company (long_name, short_name, address, email, contact_name, contact_telephone, contact_mobile, contact_email, legal_name, creator, state, status, ctime, utime) values (:long_name, :short_name, :address, :email, :contact_name, :contact_telephone, :contact_mobile, :contact_email, :legal_name, :creator, :state, :status, :ctime, :utime)
	Create(ctx context.Context, company *do.Company) (int64, error)
	// Update 更新公司信息
	// @Statement Update company Set long_name = :company.LongName,short_name = :company.ShortName WHERE id = :company.ID
	Update(ctx context.Context, company *do.Company) (int64, error)
}
