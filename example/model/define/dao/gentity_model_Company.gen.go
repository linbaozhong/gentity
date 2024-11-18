// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/table/companytbl"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/cachego"
	"github.com/linbaozhong/gentity/pkg/conv"
	"golang.org/x/sync/singleflight"
	"time"
)

type CompanyDaoer interface {
	dialect.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *db.Company, cols ...dialect.Field) (bool, error)
	// InsertBatch 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertBatch(ctx context.Context, beans []*db.Company, cols ...dialect.Field) (int64, error)
	// UpdateById 按主键更新一条数据
	UpdateById(ctx context.Context, id uint64, sets ...dialect.Setter) (bool, error)
	// UpdateBatch 批量更新多条数据
	// cols: 要更新的列名
	UpdateBatch(ctx context.Context, beans []*db.Company, cols ...dialect.Field) (bool, error)
	// DeleteById 按主键删除一条数据
	DeleteById(ctx context.Context, id uint64) (bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]*db.Company, bool, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]*db.Company, bool, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*db.Company, bool, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, id uint64, cols ...dialect.Field) (*db.Company, bool, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond ...dialect.Condition) (*db.Company, bool, error)
	// GetFirstCell 按条件读取第一行的第一个字段
	GetFirstCell(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (any, bool, error)
	//
	IDs(ctx context.Context, cond ...dialect.Condition) ([]any, error)
	//
	Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error)
}

type companyDao struct {
	db    ace.Executer
	cache cachego.Cache
	sg    singleflight.Group
}

func Company(exec ace.Executer) CompanyDaoer {
	obj := &companyDao{}
	obj.db = exec
	obj.cache = exec.Cache(db.CompanyTableName)
	return obj
}

// C Create company
func (p *companyDao) C() *ace.Creator {
	return p.db.C(db.CompanyTableName)
}

// R Read company
func (p *companyDao) R() *ace.Selector {
	return p.db.R(db.CompanyTableName)
}

// U Update company
func (p *companyDao) U() *ace.Updater {
	return p.db.U(db.CompanyTableName)
}

// D Delete company
func (p *companyDao) D() *ace.Deleter {
	return p.db.D(db.CompanyTableName)
}

// Insert 返回 LastInsertId
func (p *companyDao) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
	if len(sets) == 0 {
		return 0, dialect.ErrSetterEmpty
	}
	result, err := p.C().
		Set(sets...).
		Exec(ctx)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertOne 返回 LastInsertId
// cols: 要插入的列名
func (p *companyDao) InsertOne(ctx context.Context, bean *db.Company, cols ...dialect.Field) (bool, error) {
	result, err := p.C().
		Cols(cols...).
		Struct(ctx, bean)
	if err != nil {
		return false, err
	}

	bean.AssignPrimaryKeyValues(result)

	n, err := result.RowsAffected()
	return n > 0, err
}

// InsertBatch 批量插入,返回 RowsAffected。禁止在事务中使用
// cols: 要插入的列名
func (p *companyDao) InsertBatch(ctx context.Context, beans []*db.Company, cols ...dialect.Field) (int64, error) {
	lens := len(beans)
	if lens == 0 {
		return 0, dialect.ErrBeanEmpty
	}
	args := make([]dialect.Modeler, 0, lens)
	for _, bean := range beans {
		args = append(args, bean)
	}
	result, err := p.C().
		Cols(cols...).
		StructBatch(ctx, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Update
func (p *companyDao) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
	if len(sets) == 0 {
		return false, dialect.ErrSetterEmpty
	}
	result, err := p.U().
		Where(cond...).
		Set(sets...).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// UpdateById
func (p *companyDao) UpdateById(ctx context.Context, id uint64, sets ...dialect.Setter) (bool, error) {
	return p.Update(ctx,
		sets,
		companytbl.PrimaryKey.Eq(id),
	)
}

// UpdateBatch 批量更新,禁止在事务中使用
// cols: 要更新的列名
func (p *companyDao) UpdateBatch(ctx context.Context, beans []*db.Company, cols ...dialect.Field) (bool, error) {
	lens := len(beans)
	if lens == 0 {
		return false, dialect.ErrBeanEmpty
	}
	args := make([]dialect.Modeler, 0, lens)
	for _, bean := range beans {
		args = append(args, bean)
	}
	result, err := p.U().
		Cols(cols...).
		StructBatch(ctx, args...)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Delete
func (p *companyDao) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	result, err := p.D().
		Where(cond...).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// DeleteById
func (p *companyDao) DeleteById(ctx context.Context, id uint64) (bool, error) {
	return p.Delete(ctx,
		companytbl.PrimaryKey.Eq(id),
	)
}

// Get4Cols 先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *companyDao) Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*db.Company, bool, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols(companytbl.ReadableFields...)
	} else {
		c.Cols(cols...)
	}

	row, err := c.Where(cond...).
		QueryRow(ctx)
	if err != nil {
		return nil, false, err
	}

	obj := db.NewCompany()

	err = row.Scan(obj.AssignPtr(cols...)...)
	switch err {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return obj, true, nil
	default:
		return nil, false, err
	}
}

// Find4Cols 分页获取company} slice对象，先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *companyDao) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]*db.Company, bool, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols(companytbl.ReadableFields...)
	} else {
		c.Cols(cols...)
	}
	//
	if pageSize == 0 {
		pageSize = dialect.PageSize
	}
	//
	rows, err := c.Where(cond...).
		Limit(pageSize, pageSize*pageIndex).
		Query(ctx)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	obj := db.NewCompany()

	objs, has, err := obj.Scan(rows, cols...)
	if has {
		return objs, true, err
	}
	return nil, false, err
}

// GetByID 按主键读取一个company对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *companyDao) GetByID(ctx context.Context, id uint64, cols ...dialect.Field) (*db.Company, bool, error) {
	obj, has, e := p.getCache(ctx, id)
	if has {
		return obj, has, nil
	}

	v, e, _ := p.sg.Do(conv.Any2String(id), func() (any, error) {
		obj, has, e = p.Get4Cols(ctx, cols, companytbl.PrimaryKey.Eq(id))
		if has {
			e = p.setCache(ctx, obj)
		}
		return obj, e
	})
	if v != nil {
		return v.(*db.Company), true, e
	}

	return nil, false, e
}

// Get 按条件读取一个company对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *companyDao) Get(ctx context.Context, cond ...dialect.Condition) (*db.Company, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond...)
}

// GetFirstCell 按条件读取首行首列,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *companyDao) GetFirstCell(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (any, bool, error) {
	c := p.R().Cols(col)
	row, err := c.Where(cond...).QueryRow(ctx)
	if err != nil {
		return nil, false, err
	}

	var v any
	err = row.Scan(&v)
	switch err {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return v, true, nil
	default:
		return nil, false, err
	}
}

// Find 按条件读取一个company slice对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *companyDao) Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]*db.Company, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond...)
}

// IDs
func (p *companyDao) IDs(ctx context.Context, cond ...dialect.Condition) ([]any, error) {
	c := p.R().Cols(companytbl.PrimaryKey)
	rows, err := c.Where(cond...).
		Limit(dialect.MaxLimit).
		Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]any, 0, dialect.PageSize)
	for rows.Next() {
		var id uint64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

// Columns
func (p *companyDao) Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error) {
	c := p.R().Cols(col)
	rows, err := c.Where(cond...).
		Limit(dialect.MaxLimit).
		Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]any, 0, dialect.PageSize)
	for rows.Next() {
		var v any
		if err = rows.Scan(&v); err != nil {
			return nil, err
		}
		cols = append(cols, v)
	}
	return cols, rows.Err()
}

// Count
func (p *companyDao) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *companyDao) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	return p.R().Sum(ctx, cols, cond...)
}

// Exists
func (p *companyDao) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	c := p.R().Cols(companytbl.PrimaryKey).Where(cond...)
	row, err := c.QueryRow(ctx)
	if err != nil {
		return false, err
	}

	var id uint64
	err = row.Scan(&id)
	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// onUpdate
func (p *companyDao) onUpdate(ctx context.Context, ids ...uint64) error {
	for _, id := range ids {
		if err := p.cache.Delete(ctx, cachego.GetIdHashKey(conv.Any2String(id))); err != nil {
			return err
		}
	}

	return p.cache.PrefixDelete(ctx, "s:")
}

// getCache
func (p *companyDao) getCache(ctx context.Context, id uint64) (*db.Company, bool, error) {
	s, err := p.cache.Fetch(ctx, cachego.GetIdHashKey(conv.Any2String(id)))
	if err != nil {
		return nil, false, err
	}
	if len(s) == 0 {
		return nil, false, nil
	}
	obj := db.NewCompany()
	err = json.Unmarshal(s, obj)
	if err != nil {
		return nil, false, err
	}
	return obj, true, nil
}

// setCache
func (p *companyDao) setCache(ctx context.Context, obj *db.Company) error {
	s, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return p.cache.Save(ctx, cachego.GetIdHashKey(conv.Any2String(obj.Id)), string(s), time.Minute)
}
