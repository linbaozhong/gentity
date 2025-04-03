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

package ace

// Create
func Create(x ...Executer) CreateBuilder {
	if len(x) > 0 {
		return newCreate(x[0])
	}
	return newCreate(GetDB())
}

// Stmt
func Stmt() StmtBuilder {
	return newStmt()
}

// Select
func Select(x ...Executer) SelectBuilder {
	if len(x) > 0 {
		return newSelect(x[0])
	}
	return newSelect(GetDB())
}

// Update
func Update(x ...Executer) UpdateBuilder {
	if len(x) > 0 {
		return newUpdate(x[0])
	}
	return newUpdate(GetDB())
}

// Delete
func Delete(x ...Executer) DeleteBuilder {
	if len(x) > 0 {
		return newDelete(x[0])
	}
	return newDelete(GetDB())
}
