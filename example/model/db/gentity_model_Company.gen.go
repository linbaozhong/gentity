// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package db

import (
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/define/table/companytbl"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"time"
)

const CompanyTableName = "company"

var (
	companyPool = ace.Pool{
		New: func() interface{} {
			return &Company{}
		},
	}
)

func NewCompany() *Company {
	obj := companyPool.Get().(*Company)
	return obj
}

// Free
func (p *Company) Free() {
	if p == nil {
		return
	}
	p.reset()
	companyPool.Put(p)
	ace.Dispose(p)
}

// reset
func (p *Company) reset() {
	p.Id = 0
	p.Platform = ""
	p.CorpId = ""
	p.CorpType = 0
	p.FullCorpName = ""
	p.CorpType2 = 0
	p.CorpName = ""
	p.Industry = ""
	p.IsAuthenticated = false
	p.LicenseCode = ""
	p.CorpLogoUrl = ""
	p.InviteUrl = ""
	p.InviteCode = ""
	p.IsEcologicalCorp = false
	p.AuthLevel = 0
	p.AuthChannel = ""
	p.AuthChannelType = ""
	p.State = 0
	p.StateTime = time.Time{}
	p.CreatedTime = time.Time{}

}

func (p *Company) TableName() string {
	return CompanyTableName
}

func (p *Company) AssignPtr(args ...dialect.Field) []any {
	if len(args) == 0 {
		args = companytbl.ReadableFields
	}

	vals := make([]any, 0, len(args))
	for _, col := range args {
		switch col {
		case companytbl.Id:
			vals = append(vals, &p.Id)
		case companytbl.Platform:
			vals = append(vals, &p.Platform)
		case companytbl.CorpId:
			vals = append(vals, &p.CorpId)
		case companytbl.CorpType:
			vals = append(vals, &p.CorpType)
		case companytbl.FullCorpName:
			vals = append(vals, &p.FullCorpName)
		case companytbl.CorpType2:
			vals = append(vals, &p.CorpType2)
		case companytbl.CorpName:
			vals = append(vals, &p.CorpName)
		case companytbl.Industry:
			vals = append(vals, &p.Industry)
		case companytbl.IsAuthenticated:
			vals = append(vals, &p.IsAuthenticated)
		case companytbl.LicenseCode:
			vals = append(vals, &p.LicenseCode)
		case companytbl.CorpLogoUrl:
			vals = append(vals, &p.CorpLogoUrl)
		case companytbl.InviteUrl:
			vals = append(vals, &p.InviteUrl)
		case companytbl.InviteCode:
			vals = append(vals, &p.InviteCode)
		case companytbl.IsEcologicalCorp:
			vals = append(vals, &p.IsEcologicalCorp)
		case companytbl.AuthLevel:
			vals = append(vals, &p.AuthLevel)
		case companytbl.AuthChannel:
			vals = append(vals, &p.AuthChannel)
		case companytbl.AuthChannelType:
			vals = append(vals, &p.AuthChannelType)
		case companytbl.State:
			vals = append(vals, &p.State)
		case companytbl.StateTime:
			vals = append(vals, &p.StateTime)
		case companytbl.CreatedTime:
			vals = append(vals, &p.CreatedTime)
		}
	}

	return vals
}

func (p *Company) Scan(rows *sql.Rows, args ...dialect.Field) ([]*Company, bool, error) {
	defer rows.Close()
	companys := make([]*Company, 0)

	if len(args) == 0 {
		args = companytbl.ReadableFields
	}

	for rows.Next() {
		p := NewCompany()
		vals := p.AssignPtr(args...)
		err := rows.Scan(vals...)
		if err != nil {
			return nil, false, err
		}
		companys = append(companys, p)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	if len(companys) == 0 {
		return nil, false, sql.ErrNoRows
	}
	return companys, true, nil
}

func (p *Company) AssignValues(args ...dialect.Field) ([]string, []any) {
	var (
		lens = len(args)
		cols []string
		vals []any
	)

	if len(args) == 0 {
		args = companytbl.WritableFields
		lens = len(args)
		cols = make([]string, 0, lens)
		vals = make([]any, 0, lens)
		for _, arg := range args {
			switch arg {
			case companytbl.Id:
				if p.Id == 0 {
					continue
				}
				cols = append(cols, companytbl.Id.Quote())
				vals = append(vals, p.Id)
			case companytbl.Platform:
				if p.Platform == "" {
					continue
				}
				cols = append(cols, companytbl.Platform.Quote())
				vals = append(vals, p.Platform)
			case companytbl.CorpId:
				if p.CorpId == "" {
					continue
				}
				cols = append(cols, companytbl.CorpId.Quote())
				vals = append(vals, p.CorpId)
			case companytbl.CorpType:
				if p.CorpType == 0 {
					continue
				}
				cols = append(cols, companytbl.CorpType.Quote())
				vals = append(vals, p.CorpType)
			case companytbl.FullCorpName:
				if p.FullCorpName == "" {
					continue
				}
				cols = append(cols, companytbl.FullCorpName.Quote())
				vals = append(vals, p.FullCorpName)
			case companytbl.CorpType2:
				if p.CorpType2 == 0 {
					continue
				}
				cols = append(cols, companytbl.CorpType2.Quote())
				vals = append(vals, p.CorpType2)
			case companytbl.CorpName:
				if p.CorpName == "" {
					continue
				}
				cols = append(cols, companytbl.CorpName.Quote())
				vals = append(vals, p.CorpName)
			case companytbl.Industry:
				if p.Industry == "" {
					continue
				}
				cols = append(cols, companytbl.Industry.Quote())
				vals = append(vals, p.Industry)
			case companytbl.IsAuthenticated:
				if p.IsAuthenticated == false {
					continue
				}
				cols = append(cols, companytbl.IsAuthenticated.Quote())
				vals = append(vals, p.IsAuthenticated)
			case companytbl.LicenseCode:
				if p.LicenseCode == "" {
					continue
				}
				cols = append(cols, companytbl.LicenseCode.Quote())
				vals = append(vals, p.LicenseCode)
			case companytbl.CorpLogoUrl:
				if p.CorpLogoUrl == "" {
					continue
				}
				cols = append(cols, companytbl.CorpLogoUrl.Quote())
				vals = append(vals, p.CorpLogoUrl)
			case companytbl.InviteUrl:
				if p.InviteUrl == "" {
					continue
				}
				cols = append(cols, companytbl.InviteUrl.Quote())
				vals = append(vals, p.InviteUrl)
			case companytbl.InviteCode:
				if p.InviteCode == "" {
					continue
				}
				cols = append(cols, companytbl.InviteCode.Quote())
				vals = append(vals, p.InviteCode)
			case companytbl.IsEcologicalCorp:
				if p.IsEcologicalCorp == false {
					continue
				}
				cols = append(cols, companytbl.IsEcologicalCorp.Quote())
				vals = append(vals, p.IsEcologicalCorp)
			case companytbl.AuthLevel:
				if p.AuthLevel == 0 {
					continue
				}
				cols = append(cols, companytbl.AuthLevel.Quote())
				vals = append(vals, p.AuthLevel)
			case companytbl.AuthChannel:
				if p.AuthChannel == "" {
					continue
				}
				cols = append(cols, companytbl.AuthChannel.Quote())
				vals = append(vals, p.AuthChannel)
			case companytbl.AuthChannelType:
				if p.AuthChannelType == "" {
					continue
				}
				cols = append(cols, companytbl.AuthChannelType.Quote())
				vals = append(vals, p.AuthChannelType)
			case companytbl.State:
				if p.State == 0 {
					continue
				}
				cols = append(cols, companytbl.State.Quote())
				vals = append(vals, p.State)
			case companytbl.StateTime:
				if p.StateTime.IsZero() {
					continue
				}
				cols = append(cols, companytbl.StateTime.Quote())
				vals = append(vals, p.StateTime)
			case companytbl.CreatedTime:
				if p.CreatedTime.IsZero() {
					continue
				}
				cols = append(cols, companytbl.CreatedTime.Quote())
				vals = append(vals, p.CreatedTime)
			}
		}
		return cols, vals
	}

	cols = make([]string, 0, lens)
	vals = make([]any, 0, lens)
	for _, arg := range args {
		switch arg {
		case companytbl.Id:
			cols = append(cols, companytbl.Id.Quote())
			vals = append(vals, p.Id)
		case companytbl.Platform:
			cols = append(cols, companytbl.Platform.Quote())
			vals = append(vals, p.Platform)
		case companytbl.CorpId:
			cols = append(cols, companytbl.CorpId.Quote())
			vals = append(vals, p.CorpId)
		case companytbl.CorpType:
			cols = append(cols, companytbl.CorpType.Quote())
			vals = append(vals, p.CorpType)
		case companytbl.FullCorpName:
			cols = append(cols, companytbl.FullCorpName.Quote())
			vals = append(vals, p.FullCorpName)
		case companytbl.CorpType2:
			cols = append(cols, companytbl.CorpType2.Quote())
			vals = append(vals, p.CorpType2)
		case companytbl.CorpName:
			cols = append(cols, companytbl.CorpName.Quote())
			vals = append(vals, p.CorpName)
		case companytbl.Industry:
			cols = append(cols, companytbl.Industry.Quote())
			vals = append(vals, p.Industry)
		case companytbl.IsAuthenticated:
			cols = append(cols, companytbl.IsAuthenticated.Quote())
			vals = append(vals, p.IsAuthenticated)
		case companytbl.LicenseCode:
			cols = append(cols, companytbl.LicenseCode.Quote())
			vals = append(vals, p.LicenseCode)
		case companytbl.CorpLogoUrl:
			cols = append(cols, companytbl.CorpLogoUrl.Quote())
			vals = append(vals, p.CorpLogoUrl)
		case companytbl.InviteUrl:
			cols = append(cols, companytbl.InviteUrl.Quote())
			vals = append(vals, p.InviteUrl)
		case companytbl.InviteCode:
			cols = append(cols, companytbl.InviteCode.Quote())
			vals = append(vals, p.InviteCode)
		case companytbl.IsEcologicalCorp:
			cols = append(cols, companytbl.IsEcologicalCorp.Quote())
			vals = append(vals, p.IsEcologicalCorp)
		case companytbl.AuthLevel:
			cols = append(cols, companytbl.AuthLevel.Quote())
			vals = append(vals, p.AuthLevel)
		case companytbl.AuthChannel:
			cols = append(cols, companytbl.AuthChannel.Quote())
			vals = append(vals, p.AuthChannel)
		case companytbl.AuthChannelType:
			cols = append(cols, companytbl.AuthChannelType.Quote())
			vals = append(vals, p.AuthChannelType)
		case companytbl.State:
			cols = append(cols, companytbl.State.Quote())
			vals = append(vals, p.State)
		case companytbl.StateTime:
			cols = append(cols, companytbl.StateTime.Quote())
			vals = append(vals, p.StateTime)
		case companytbl.CreatedTime:
			cols = append(cols, companytbl.CreatedTime.Quote())
			vals = append(vals, p.CreatedTime)
		}
	}
	return cols, vals
}

func (p *Company) AssignKeys() (dialect.Field, any) {
	return companytbl.PrimaryKey, p.Id
}

func (p *Company) AssignPrimaryKeyValues(result sql.Result) error {
	return nil
}
