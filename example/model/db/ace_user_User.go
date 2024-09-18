// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package db

import (
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/define/table/usertbl"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"sync"
	"time"
)

const UserTableName = "user"

var (
	userPool = sync.Pool{
		New: func() interface{} {
			return &User{}
		},
	}
)

func NewUser() *User {
	return userPool.Get().(*User)
}

// Free
func (p *User) Free() {
	if p == nil {
		return
	}
	p.ID = 0
	p.Name = ""
	p.Avatar = ""
	p.Nickname = ""
	p.Status = 0
	p.IsAllow = false
	p.CreatedTime = time.Time{}

	userPool.Put(p)
}

func (p *User) TableName() string {
	return UserTableName
}

func (p *User) AssignPtr(args ...dialect.Field) []any {
	if len(args) == 0 {
		args = usertbl.ReadableFields
	}

	vals := make([]any, 0, len(args))
	for _, col := range args {
		switch col {
		case usertbl.ID:
			vals = append(vals, &p.ID)
		case usertbl.Name:
			vals = append(vals, &p.Name)
		case usertbl.Avatar:
			vals = append(vals, &p.Avatar)
		case usertbl.Nickname:
			vals = append(vals, &p.Nickname)
		case usertbl.Status:
			vals = append(vals, &p.Status)
		case usertbl.IsAllow:
			vals = append(vals, &p.IsAllow)
		case usertbl.CreatedTime:
			vals = append(vals, &p.CreatedTime)
		}
	}

	return vals
}

func (p *User) Scan(rows *sql.Rows, args ...dialect.Field) ([]*User, bool, error) {
	defer rows.Close()
	users := make([]*User, 0)

	if len(args) == 0 {
		args = usertbl.ReadableFields
	}

	for rows.Next() {
		p := NewUser()
		vals := p.AssignPtr(args...)
		err := rows.Scan(vals...)
		if err != nil {
			return nil, false, err
		}
		users = append(users, p)
	}
	if len(users) == 0 {
		return nil, false, sql.ErrNoRows
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	return users, true, nil
}

func (p *User) AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = usertbl.WritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			case usertbl.ID:
				if p.ID == 0 {
					continue
				}
				cols = append(cols, usertbl.ID.Quote())
				vals = append(vals, p.ID)
			case usertbl.Name:
				if p.Name == "" {
					continue
				}
				cols = append(cols, usertbl.Name.Quote())
				vals = append(vals, p.Name)
			case usertbl.Avatar:
				if p.Avatar == "" {
					continue
				}
				cols = append(cols, usertbl.Avatar.Quote())
				vals = append(vals, p.Avatar)
			case usertbl.Nickname:
				if p.Nickname == "" {
					continue
				}
				cols = append(cols, usertbl.Nickname.Quote())
				vals = append(vals, p.Nickname)
			case usertbl.Status:
				if p.Status == 0 {
					continue
				}
				cols = append(cols, usertbl.Status.Quote())
				vals = append(vals, p.Status)
			case usertbl.IsAllow:
				if p.IsAllow == false {
					continue
				}
				cols = append(cols, usertbl.IsAllow.Quote())
				vals = append(vals, p.IsAllow)
			case usertbl.CreatedTime:
				if p.CreatedTime.IsZero() {
					continue
				}
				cols = append(cols, usertbl.CreatedTime.Quote())
				vals = append(vals, p.CreatedTime)
			}
		}
		return cols, vals
	}

	cols = make([]string, 0, lens)
	vals = make([]any, 0, lens)
	for _, arg := range args {
		switch arg {
		case usertbl.ID:
			cols = append(cols, usertbl.ID.Quote())
			vals = append(vals, p.ID)
		case usertbl.Name:
			cols = append(cols, usertbl.Name.Quote())
			vals = append(vals, p.Name)
		case usertbl.Avatar:
			cols = append(cols, usertbl.Avatar.Quote())
			vals = append(vals, p.Avatar)
		case usertbl.Nickname:
			cols = append(cols, usertbl.Nickname.Quote())
			vals = append(vals, p.Nickname)
		case usertbl.Status:
			cols = append(cols, usertbl.Status.Quote())
			vals = append(vals, p.Status)
		case usertbl.IsAllow:
			cols = append(cols, usertbl.IsAllow.Quote())
			vals = append(vals, p.IsAllow)
		case usertbl.CreatedTime:
			cols = append(cols, usertbl.CreatedTime.Quote())
			vals = append(vals, p.CreatedTime)
		}
	}
	return cols, vals
}

func (p *User) AssignKeys() ([]dialect.Field, []any) {
	return usertbl.PrimaryKeys, []any{
		p.ID,
	}
}
