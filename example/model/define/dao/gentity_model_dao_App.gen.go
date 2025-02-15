// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/define/table/tblapp"
	"github.com/linbaozhong/gentity/example/model/do"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

type apper interface {
	dialect.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *do.App, cols ...dialect.Field) (bool, error)
	// InsertBatch 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertBatch(ctx context.Context, beans []*do.App, cols ...dialect.Field) (int64, error)
	// UpdateById 按主键更新一条数据
	UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error)
	// UpdateBatch 批量更新多条数据
	// cols: 要更新的列名
	UpdateBatch(ctx context.Context, beans []*do.App, cols ...dialect.Field) (bool, error)
	// DeleteById 按主键删除一条数据
	DeleteById(ctx context.Context, id types.BigInt) (bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]do.App, bool, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]do.App, bool, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*do.App, bool, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*do.App, bool, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond ...dialect.Condition) (*do.App, bool, error)
	// GetFirstCell 按条件读取第一行的第一个字段
	GetFirstCell(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (any, bool, error)
	//
	IDs(ctx context.Context, cond ...dialect.Condition) ([]any, error)
	//
	Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error)
}

type daoApp struct {
	db ace.Executer
}

func App(exec ace.Executer) apper {
	_obj := &daoApp{}
	_obj.db = exec
	return _obj
}

// C Create app
func (p *daoApp) C() *ace.Creator {
	return p.db.C(do.AppTableName)
}

// R Read app
func (p *daoApp) R() *ace.Selector {
	return p.db.R(do.AppTableName)
}

// U Update app
func (p *daoApp) U() *ace.Updater {
	return p.db.U(do.AppTableName)
}

// D Delete app
func (p *daoApp) D() *ace.Deleter {
	return p.db.D(do.AppTableName)
}

// Insert 返回 LastInsertId
func (p *daoApp) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
	if len(sets) == 0 {
		return 0, dialect.ErrSetterEmpty
	}
	_result, e := p.C().
		Set(sets...).
		Exec(ctx)
	if e != nil {
		log.Error(e)
		return 0, e
	}
	return _result.LastInsertId()
}

// InsertOne 返回 LastInsertId
// cols: 要插入的列名
func (p *daoApp) InsertOne(ctx context.Context, bean *do.App, cols ...dialect.Field) (bool, error) {
	_result, e := p.C().
		Cols(cols...).
		Struct(ctx, bean)
	if e != nil {
		log.Error(e)
		return false, e
	}

	bean.AssignPrimaryKeyValues(_result)

	_n, e := _result.RowsAffected()
	return _n > 0, e
}

// InsertBatch 批量插入,返回 RowsAffected。禁止在事务中使用
// cols: 要插入的列名
func (p *daoApp) InsertBatch(ctx context.Context, beans []*do.App, cols ...dialect.Field) (int64, error) {
	_lens := len(beans)
	if _lens == 0 {
		return 0, dialect.ErrBeanEmpty
	}
	_args := make([]dialect.Modeler, 0, _lens)
	for _, _bean := range beans {
		_args = append(_args, _bean)
	}
	_result, e := p.C().
		Cols(cols...).
		StructBatch(ctx, _args...)
	if e != nil {
		log.Error(e)
		return 0, e
	}

	return _result.RowsAffected()
}

// Update
func (p *daoApp) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
	if len(sets) == 0 {
		return false, dialect.ErrSetterEmpty
	}
	_result, e := p.U().
		Where(cond...).
		Set(sets...).
		Exec(ctx)
	if e != nil {
		log.Error(e)
		return false, e
	}
	_n, e := _result.RowsAffected()
	return _n >= 0, e
}

// UpdateById
func (p *daoApp) UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error) {
	return p.Update(ctx,
		sets,
		tblapp.PrimaryKey.Eq(id),
	)
}

// UpdateBatch 批量更新,禁止在事务中使用
// cols: 要更新的列名
func (p *daoApp) UpdateBatch(ctx context.Context, beans []*do.App, cols ...dialect.Field) (bool, error) {
	_lens := len(beans)
	if _lens == 0 {
		return false, dialect.ErrBeanEmpty
	}
	_args := make([]dialect.Modeler, 0, _lens)
	for _, _bean := range beans {
		_args = append(_args, _bean)
	}
	_result, e := p.U().
		Cols(cols...).
		StructBatch(ctx, _args...)
	if e != nil {
		log.Error(e)
		return false, e
	}
	_n, e := _result.RowsAffected()
	return _n >= 0, e
}

// Delete
func (p *daoApp) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	_result, e := p.D().
		Where(cond...).
		Exec(ctx)
	if e != nil {
		log.Error(e)
		return false, e
	}
	_n, e := _result.RowsAffected()
	return _n >= 0, e
}

// DeleteById
func (p *daoApp) DeleteById(ctx context.Context, id types.BigInt) (bool, error) {
	return p.Delete(ctx,
		tblapp.PrimaryKey.Eq(id),
	)
}

// Get4Cols 先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoApp) Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*do.App, bool, error) {
	_c := p.R()
	if len(cols) == 0 {
		_c.Cols(tblapp.ReadableFields...)
	} else {
		_c.Cols(cols...)
	}

	_row, e := _c.Where(cond...).
		QueryRow(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}

	_obj := do.NewApp()

	e = _row.Scan(_obj.AssignPtr(cols...)...)
	switch e {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return _obj, true, nil
	default:
		log.Error(e)
		return nil, false, e
	}
}

// Find4Cols 分页获取app slice对象，先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoApp) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]do.App, bool, error) {
	_c := p.R()
	if len(cols) == 0 {
		_c.Cols(tblapp.ReadableFields...)
	} else {
		_c.Cols(cols...)
	}
	//
	if pageSize == 0 {
		pageSize = dialect.PageSize
	}
	//
	_rows, e := _c.Where(cond...).
		Limit(pageSize, pageSize*pageIndex).
		Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}
	defer _rows.Close()

	_obj := do.NewApp()

	_objs, has, e := _obj.Scan(_rows, cols...)
	if has {
		return _objs, true, nil
	}
	log.Error(e)
	return nil, false, e
}

// GetByID 按主键读取一个app对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoApp) GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*do.App, bool, error) {
	return p.Get4Cols(ctx, cols, tblapp.PrimaryKey.Eq(id))
}

// Get 按条件读取一个app对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoApp) Get(ctx context.Context, cond ...dialect.Condition) (*do.App, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond...)
}

// GetFirstCell 按条件读取首行首列,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoApp) GetFirstCell(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (any, bool, error) {
	_c := p.R().Cols(col)
	_row, e := _c.Where(cond...).QueryRow(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}

	var _v any
	e = _row.Scan(&_v)
	switch e {
	case sql.ErrNoRows:
		return nil, false, nil
	case nil:
		return _v, true, nil
	default:
		log.Error(e)
		return nil, false, e
	}
}

// Find 按条件读取一个app slice对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoApp) Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]do.App, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond...)
}

// IDs
func (p *daoApp) IDs(ctx context.Context, cond ...dialect.Condition) ([]any, error) {
	_c := p.R().Cols(tblapp.PrimaryKey)
	_rows, e := _c.Where(cond...).
		Limit(dialect.MaxLimit).
		Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, e
	}
	defer _rows.Close()

	_ids := make([]any, 0, dialect.PageSize)
	for _rows.Next() {
		var id types.BigInt
		if e = _rows.Scan(&id); e != nil {
			log.Error(e)
			return nil, e
		}
		_ids = append(_ids, id)
	}

	return _ids, _rows.Err()
}

// Columns
func (p *daoApp) Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error) {
	_c := p.R().Cols(col)
	_rows, e := _c.Where(cond...).
		Limit(dialect.MaxLimit).
		Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, e
	}
	defer _rows.Close()

	_cols := make([]any, 0, dialect.PageSize)
	for _rows.Next() {
		var _v any
		if e = _rows.Scan(&_v); e != nil {
			log.Error(e)
			return nil, e
		}
		_cols = append(_cols, _v)
	}
	return _cols, _rows.Err()
}

// Count
func (p *daoApp) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *daoApp) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	return p.R().Sum(ctx, cols, cond...)
}

// Exists
func (p *daoApp) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	_c := p.R().Cols(tblapp.PrimaryKey).Where(cond...)
	_row, e := _c.QueryRow(ctx)
	if e != nil {
		log.Error(e)
		return false, e
	}

	var id types.BigInt
	e = _row.Scan(&id)
	switch e {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, e
	}
}
