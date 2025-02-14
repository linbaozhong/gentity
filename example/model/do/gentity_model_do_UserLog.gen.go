// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package do

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/linbaozhong/gentity/example/model/define/table/tbluserlog"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/ace/pool"
	"github.com/linbaozhong/gentity/pkg/gjson"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

const UserLogTableName = "user_log"

var (
	userlogPool = pool.New(ace.Context, func() any {
		obj := &UserLog{}
		obj.UUID()
		return obj
	})
)

func NewUserLog() *UserLog {
	obj := userlogPool.Get().(*UserLog)
	return obj
}

// MarshalJSON
func (p *UserLog) MarshalJSON() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	buf.WriteByte('{')
	if p.Id != 0 {
		buf.WriteString(`"id":` + types.Marshal(p.Id) + `,`)
	}
	if p.UserId != 0 {
		buf.WriteString(`"user_id":` + types.Marshal(p.UserId) + `,`)
	}
	if !p.LoginTime.IsZero() {
		buf.WriteString(`"login_time":` + types.Marshal(p.LoginTime) + `,`)
	}
	if p.Device != "" {
		buf.WriteString(`"device":` + types.Marshal(p.Device) + `,`)
	}
	if p.Os != "" {
		buf.WriteString(`"os":` + types.Marshal(p.Os) + `,`)
	}
	if p.OsVersion != "" {
		buf.WriteString(`"os_version":` + types.Marshal(p.OsVersion) + `,`)
	}
	if p.AppName != "" {
		buf.WriteString(`"app_name":` + types.Marshal(p.AppName) + `,`)
	}
	if p.AppVersion != "" {
		buf.WriteString(`"app_version":` + types.Marshal(p.AppVersion) + `,`)
	}
	if p.Ip != "" {
		buf.WriteString(`"ip":` + types.Marshal(p.Ip) + `,`)
	}
	if p.User != nil {
		buf.WriteString(`"user":` + types.Marshal(p.User) + `,`)
	}
	if l := buf.Len(); l > 1 {
		buf.Truncate(l - 1)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON
func (p *UserLog) UnmarshalJSON(data []byte) error {

	ok := gjson.ValidBytes(data)
	if !ok {
		return errors.New("invalid json")
	}
	result := gjson.ParseBytes(data)
	result.ForEach(func(key, value gjson.Result) bool {
		var e error
		switch key.Str {
		case "id":
			e = types.Unmarshal(value, &p.Id, types.BigInt(value.Uint()))
		case "user_id":
			e = types.Unmarshal(value, &p.UserId, types.BigInt(value.Uint()))
		case "login_time":
			e = types.Unmarshal(value, &p.LoginTime, types.Time{Time: value.Time()})
		case "device":
			e = types.Unmarshal(value, &p.Device, types.String(value.Str))
		case "os":
			e = types.Unmarshal(value, &p.Os, types.String(value.Str))
		case "os_version":
			e = types.Unmarshal(value, &p.OsVersion, types.String(value.Str))
		case "app_name":
			e = types.Unmarshal(value, &p.AppName, types.String(value.Str))
		case "app_version":
			e = types.Unmarshal(value, &p.AppVersion, types.String(value.Str))
		case "ip":
			e = types.Unmarshal(value, &p.Ip, types.String(value.Str))
		case "user":
			e = types.Unmarshal(value, &p.User, func(value gjson.Result) User {
				var obj User
				e := types.Unmarshal(value, &obj)
				if e != nil {
					panic(e)
				}
				return obj
			}(value))
		}
		if e != nil {
			log.Error(e)
			return false
		}
		return true
	})
	return nil
}

// Free
func (p *UserLog) Free() {
	if p == nil {
		return
	}

	userlogPool.Put(p)
}

// Reset
func (p *UserLog) Reset() {
	p.Id = 0
	p.UserId = 0
	p.LoginTime = types.Time{}
	p.Device = ""
	p.Os = ""
	p.OsVersion = ""
	p.AppName = ""
	p.AppVersion = ""
	p.Ip = ""
	p.User = User{}

}

func (p *UserLog) TableName() string {
	return UserLogTableName
}

func (p *UserLog) AssignPtr(args ...dialect.Field) []any {
	if len(args) == 0 {
		args = tbluserlog.ReadableFields
	}

	vals := make([]any, 0, len(args))
	for _, col := range args {
		switch col {
		case tbluserlog.Id:
			vals = append(vals, &p.Id)
		case tbluserlog.UserId:
			vals = append(vals, &p.UserId)
		case tbluserlog.LoginTime:
			vals = append(vals, &p.LoginTime)
		case tbluserlog.Device:
			vals = append(vals, &p.Device)
		case tbluserlog.Os:
			vals = append(vals, &p.Os)
		case tbluserlog.OsVersion:
			vals = append(vals, &p.OsVersion)
		case tbluserlog.AppName:
			vals = append(vals, &p.AppName)
		case tbluserlog.AppVersion:
			vals = append(vals, &p.AppVersion)
		case tbluserlog.Ip:
			vals = append(vals, &p.Ip)
		case tbluserlog.User:
			vals = append(vals, &p.User)
		}
	}

	return vals
}

func (p *UserLog) Scan(rows *sql.Rows, args ...dialect.Field) ([]UserLog, bool, error) {
	defer rows.Close()
	user_logs := make([]UserLog, 0)

	if len(args) == 0 {
		args = tbluserlog.ReadableFields
	}

	for rows.Next() {
		p := NewUserLog()
		vals := p.AssignPtr(args...)
		err := rows.Scan(vals...)
		if err != nil {
			log.Error(err)
			return nil, false, err
		}
		user_logs = append(user_logs, *p)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		return nil, false, err
	}
	if len(user_logs) == 0 {
		return nil, false, sql.ErrNoRows
	}
	return user_logs, true, nil
}

func (p *UserLog) AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = tbluserlog.WritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			case tbluserlog.Id:
				if p.Id == 0 {
					continue
				}
				cols = append(cols, tbluserlog.Id.Quote())
				vals = append(vals, p.Id)
			case tbluserlog.UserId:
				if p.UserId == 0 {
					continue
				}
				cols = append(cols, tbluserlog.UserId.Quote())
				vals = append(vals, p.UserId)
			case tbluserlog.LoginTime:
				if p.LoginTime.IsZero() {
					continue
				}
				cols = append(cols, tbluserlog.LoginTime.Quote())
				vals = append(vals, p.LoginTime)
			case tbluserlog.Device:
				if p.Device == "" {
					continue
				}
				cols = append(cols, tbluserlog.Device.Quote())
				vals = append(vals, p.Device)
			case tbluserlog.Os:
				if p.Os == "" {
					continue
				}
				cols = append(cols, tbluserlog.Os.Quote())
				vals = append(vals, p.Os)
			case tbluserlog.OsVersion:
				if p.OsVersion == "" {
					continue
				}
				cols = append(cols, tbluserlog.OsVersion.Quote())
				vals = append(vals, p.OsVersion)
			case tbluserlog.AppName:
				if p.AppName == "" {
					continue
				}
				cols = append(cols, tbluserlog.AppName.Quote())
				vals = append(vals, p.AppName)
			case tbluserlog.AppVersion:
				if p.AppVersion == "" {
					continue
				}
				cols = append(cols, tbluserlog.AppVersion.Quote())
				vals = append(vals, p.AppVersion)
			case tbluserlog.Ip:
				if p.Ip == "" {
					continue
				}
				cols = append(cols, tbluserlog.Ip.Quote())
				vals = append(vals, p.Ip)
			}
		}
		return cols, vals
	}

	cols = make([]string, 0, lens)
	vals = make([]any, 0, lens)
	for _, arg := range args {
		switch arg {
		case tbluserlog.Id:
			cols = append(cols, tbluserlog.Id.Quote())
			vals = append(vals, p.Id)
		case tbluserlog.UserId:
			cols = append(cols, tbluserlog.UserId.Quote())
			vals = append(vals, p.UserId)
		case tbluserlog.LoginTime:
			cols = append(cols, tbluserlog.LoginTime.Quote())
			vals = append(vals, p.LoginTime)
		case tbluserlog.Device:
			cols = append(cols, tbluserlog.Device.Quote())
			vals = append(vals, p.Device)
		case tbluserlog.Os:
			cols = append(cols, tbluserlog.Os.Quote())
			vals = append(vals, p.Os)
		case tbluserlog.OsVersion:
			cols = append(cols, tbluserlog.OsVersion.Quote())
			vals = append(vals, p.OsVersion)
		case tbluserlog.AppName:
			cols = append(cols, tbluserlog.AppName.Quote())
			vals = append(vals, p.AppName)
		case tbluserlog.AppVersion:
			cols = append(cols, tbluserlog.AppVersion.Quote())
			vals = append(vals, p.AppVersion)
		case tbluserlog.Ip:
			cols = append(cols, tbluserlog.Ip.Quote())
			vals = append(vals, p.Ip)
		}
	}
	return cols, vals
}

func (p *UserLog) AssignKeys() (dialect.Field, any) {
	return tbluserlog.PrimaryKey, p.Id
}

func (p *UserLog) AssignPrimaryKeyValues(result sql.Result) error {
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = types.BigInt(id)
	return nil
}
