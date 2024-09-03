package model

import (
	"database/sql"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
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
	p.Avatar = ""
	p.CreatedTime = time.Time{}
	p.ID = 0
	p.IsAllow = false
	p.Name = ""
	p.Nickname = ""
	p.Status = 0

	userPool.Put(p)
}

func (p *User) TableName() string {
	return UserTableName
}

func (p *User) Scan(rows *sql.Rows, args ...atype.Field) ([]*User, error) {
	defer rows.Close()
	users := make([]*User, 0)

	if len(args) == 0 {
		args = UserReadableFields
	}

	for rows.Next() {
		p := NewUser()
		vals := make([]any, 0, len(args))
		for _, col := range args {
			switch col {
			case UserTbl.Avatar:
				vals = append(vals, &p.Avatar)
			case UserTbl.CreatedTime:
				vals = append(vals, &p.CreatedTime)
			case UserTbl.ID:
				vals = append(vals, &p.ID)
			case UserTbl.IsAllow:
				vals = append(vals, &p.IsAllow)
			case UserTbl.Name:
				vals = append(vals, &p.Name)
			case UserTbl.Nickname:
				vals = append(vals, &p.Nickname)
			case UserTbl.Status:
				vals = append(vals, &p.Status)
			}
		}
		err := rows.Scan(vals...)
		if err != nil {
			return nil, err
		}
		users = append(users, p)
	}
	return users, nil
}

func (p *User) AssignValues(args ...atype.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = UserWritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			case UserTbl.Avatar:
				if p.Avatar == "" {
					continue
				}
				cols = append(cols, UserTbl.Avatar.Quote())
				vals = append(vals, p.Avatar)
			case UserTbl.CreatedTime:
				if p.CreatedTime.IsZero() {
					continue
				}
				cols = append(cols, UserTbl.CreatedTime.Quote())
				vals = append(vals, p.CreatedTime)
			case UserTbl.ID:
				if p.ID == 0 {
					continue
				}
				cols = append(cols, UserTbl.ID.Quote())
				vals = append(vals, p.ID)
			case UserTbl.IsAllow:
				if p.IsAllow == false {
					continue
				}
				cols = append(cols, UserTbl.IsAllow.Quote())
				vals = append(vals, p.IsAllow)
			case UserTbl.Name:
				if p.Name == "" {
					continue
				}
				cols = append(cols, UserTbl.Name.Quote())
				vals = append(vals, p.Name)
			case UserTbl.Nickname:
				if p.Nickname == "" {
					continue
				}
				cols = append(cols, UserTbl.Nickname.Quote())
				vals = append(vals, p.Nickname)
			case UserTbl.Status:
				if p.Status == 0 {
					continue
				}
				cols = append(cols, UserTbl.Status.Quote())
				vals = append(vals, p.Status)
			}
		}
		return cols, vals
	}

	cols = make([]string, 0, lens)
	vals = make([]any, 0, lens)
	for _, arg := range args {
		switch arg {
		case UserTbl.Avatar:
			cols = append(cols, UserTbl.Avatar.Quote())
			vals = append(vals, p.Avatar)
		case UserTbl.CreatedTime:
			cols = append(cols, UserTbl.CreatedTime.Quote())
			vals = append(vals, p.CreatedTime)
		case UserTbl.ID:
			cols = append(cols, UserTbl.ID.Quote())
			vals = append(vals, p.ID)
		case UserTbl.IsAllow:
			cols = append(cols, UserTbl.IsAllow.Quote())
			vals = append(vals, p.IsAllow)
		case UserTbl.Name:
			cols = append(cols, UserTbl.Name.Quote())
			vals = append(vals, p.Name)
		case UserTbl.Nickname:
			cols = append(cols, UserTbl.Nickname.Quote())
			vals = append(vals, p.Nickname)
		case UserTbl.Status:
			cols = append(cols, UserTbl.Status.Quote())
			vals = append(vals, p.Status)
		}
	}
	return cols, vals
}

func (p *User) AssignKeys() ([]atype.Field, []any) {
	return UserPrimaryKeys, []any{
		p.ID,
	}
}

var (
	UserTbl = struct {
		Avatar      atype.Field
		CreatedTime atype.Field
		ID          atype.Field
		IsAllow     atype.Field
		Name        atype.Field
		Nickname    atype.Field
		Status      atype.Field
	}{
		Avatar:      atype.Field{Name: "avatar", Table: UserTableName},
		CreatedTime: atype.Field{Name: "created_time", Table: UserTableName},
		ID:          atype.Field{Name: "id", Table: UserTableName},
		IsAllow:     atype.Field{Name: "is_allow", Table: UserTableName},
		Name:        atype.Field{Name: "name", Table: UserTableName},
		Nickname:    atype.Field{Name: "nickname", Table: UserTableName},
		Status:      atype.Field{Name: "status", Table: UserTableName},
	}
	// 主键
	UserPrimaryKeys = []atype.Field{
		UserTbl.ID,
	}

	// 可写列
	UserWritableFields = []atype.Field{
		UserTbl.Avatar,
		UserTbl.IsAllow,
		UserTbl.Name,
		UserTbl.Nickname,
		UserTbl.Status,
	}
	// 可读列
	UserReadableFields = []atype.Field{
		UserTbl.Avatar,
		UserTbl.CreatedTime,
		UserTbl.ID,
		UserTbl.IsAllow,
		UserTbl.Name,
		UserTbl.Nickname,
		UserTbl.Status,
	}
)
