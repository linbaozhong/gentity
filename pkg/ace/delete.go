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
	"context"
	"database/sql"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace/types"
	"strings"
	"sync"
)

type (
	Deleter struct {
		db          Executer
		object      types.Modeler
		where       strings.Builder
		whereParams []interface{}
		command     strings.Builder
	}
)

var (
	deletePool = sync.Pool{
		New: func() interface{} {
			obj := &Deleter{}
			return obj
		},
	}
)

// Deleter
func NewDelete(db Executer, mod types.Modeler) *Deleter {
	if db == nil || mod == nil {
		panic("db or table is nil")
		return nil
	}
	obj := deletePool.Get().(*Deleter)
	obj.db = db
	obj.object = mod
	obj.command.Reset()
	return obj

}

func (d *Deleter) Free() {
	d.object = nil
	d.where.Reset()
	d.whereParams = d.whereParams[:]
	deletePool.Put(d)
}

func (d *Deleter) String() string {
	return fmt.Sprintf("%s  %v", d.command.String(), d.whereParams)
}

// Where
func (d *Deleter) Where(fns ...types.Condition) *Deleter {
	if len(fns) == 0 {
		return d
	}
	if d.where.Len() == 0 {
		d.where.WriteString("(")
	} else {
		d.where.WriteString(types.Operator_and + "(")
	}
	for i, fn := range fns {
		if i > 0 {
			d.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		d.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			d.whereParams = append(d.whereParams, vals...)
		} else {
			d.whereParams = append(d.whereParams, val)
		}
	}
	d.where.WriteString(")")

	return d
}

// And
func (d *Deleter) And(fns ...types.Condition) *Deleter {
	if len(fns) == 0 {
		return d
	}

	if d.where.Len() == 0 {
		d.where.WriteString("(")
	} else {
		d.where.WriteString(types.Operator_and + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			d.where.WriteString(types.Operator_or)
		}
		cond, val := fn()
		d.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			d.whereParams = append(d.whereParams, vals...)
		} else {
			d.whereParams = append(d.whereParams, val)
		}
	}
	d.where.WriteString(")")
	return d
}

// Or
func (d *Deleter) Or(fns ...types.Condition) *Deleter {
	if len(fns) == 0 {
		return d
	}

	if d.where.Len() == 0 {
		d.where.WriteString("(")
	} else {
		d.where.WriteString(types.Operator_or + "(")
	}

	for i, fn := range fns {
		if i > 0 {
			d.where.WriteString(types.Operator_and)
		}
		cond, val := fn()
		d.where.WriteString(cond)
		if vals, ok := val.([]any); ok {
			d.whereParams = append(d.whereParams, vals...)
		} else {
			d.whereParams = append(d.whereParams, val)
		}
	}
	d.where.WriteString(")")
	return d
}

// Do
func (d *Deleter) Do(ctx context.Context) (sql.Result, error) {
	defer d.Free()

	d.command.WriteString("DELETE FROM " + types.Quote_Char + d.object.TableName() + types.Quote_Char)
	// WHERE
	if d.where.Len() > 0 {
		d.command.WriteString(" WHERE " + d.where.String())
	}

	return d.db.ExecContext(ctx, d.command.String(), d.whereParams...)
}
