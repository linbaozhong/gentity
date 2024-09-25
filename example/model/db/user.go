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

package db

import "time"

// tablename user
type User struct {
	ID          uint64    `json:"id" db:"'id' pk auto"`
	Name        string    `json:"name" db:"name"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Status      int8      `json:"status" db:"status"`
	IsAllow     bool      `json:"is_allow" db:"is_allow"`
	CreatedTime time.Time `json:"created_time" db:"created_time <-"`
}
