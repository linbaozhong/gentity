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

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"strconv"
)

// Limit 设置查询结果的限制条件
//
//	size 大小
//	start 开始位置
func (o *orm) Limit(size uint, start ...uint) Builder {
	if size == 0 {
		o.limit = ""
		return o
	}
	if len(start) > 0 {
		o.limit = " LIMIT " + strconv.Itoa(int(size)) + " OFFSET " + strconv.Itoa(int(start[0]))
	} else {
		o.limit = " LIMIT " + strconv.Itoa(int(size))
	}

	return o
}

// Page 分页查询
//
//	pageIndex 页码, 从1开始
//	pageSize 页大小
func (o *orm) Page(pageIndex, pageSize uint) Builder {
	if pageSize == 0 {
		return o.Limit(0)
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return o.Limit(pageSize, (pageIndex-1)*pageSize)
}

// LimitByBookmark 按上页最后一条记录的主键值作为书签查询下一页数据
//
//	size 页大小
//	bm 书签条件，如果是正序，书签条件为大于，如果是倒序，书签条件为小于
//
// 例如：
//
//	正序： 	orm.PageByBookmark(10, Id.Gt(lastId))
//	倒序： 	orm.PageByBookmark(10, Id.Lt(lastId))
func (o *orm) PageByBookmark(size uint, bm dialect.Condition) Builder {
	return o.Where(bm).Limit(size)
}
