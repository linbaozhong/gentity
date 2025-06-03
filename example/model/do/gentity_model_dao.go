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

package do

import "context"

// AccountDao 账号信息DAO
// @tablename account
type AccountDao interface {
	// FindByID 根据ID查询账号信息
	// @Condition
	FindByID(ctx context.Context, id int) (Account, error)
	FindByIDs(ctx context.Context, ids []int) ([]Account, error)
	FindByPhone(ctx context.Context, phone string) (Account, error)
}
