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

package orm

import (
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

type Option func(*orm)

// WithWhere 配置where条件
func WithWhere(args ...dialect.Condition) Option {
	return func(o *orm) {
		o.Where(args...)
	}
}

// WithOrderBy 配置order by条件
func WithOrderBy(args ...dialect.Order) Option {
	return func(o *orm) {
		o.OrderFunc(args...)
	}
}

// WithGroupBy 配置group by条件
func WithGroupBy(args ...dialect.Field) Option {
	return func(o *orm) {
		o.Group(args...)
	}
}

// WithHaving 配置having条件
func WithHaving(args ...dialect.Condition) Option {
	return func(o *orm) {
		o.Having(args...)
	}
}

// WithLimit 配置limit条件
func WithLimit(size uint, start ...uint) Option {
	return func(o *orm) {
		o.Limit(size, start...)
	}
}

// WithJoin 配置join条件
func WithJoin(joinType dialect.JoinType, left, right dialect.Field, fns ...dialect.Condition) Option {
	return func(o *orm) {
		o.Join(joinType, left, right, fns...)
	}
}

// WithSet 为字段赋值
func WithSet(args ...dialect.Setter) Option {
	return func(o *orm) {
		o.Set(args...)
	}
}

// WithExpr 用表达式为字段赋值
func WithExpr(args ...dialect.ExprSetter) Option {
	return func(o *orm) {
		o.SetExpr(args...)
	}
}
