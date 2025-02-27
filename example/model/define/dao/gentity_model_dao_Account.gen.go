// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/table/tblaccount"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

type accounter interface {
	dialect.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *db.Account, cols ...dialect.Field) (bool, error)
	// InsertBatch 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertBatch(ctx context.Context, beans []*db.Account, cols ...dialect.Field) (int64, error)
	// UpdateById 按主键更新一条数据
	UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error)
	// UpdateBatch 批量更新多条数据
	// cols: 要更新的列名
	UpdateBatch(ctx context.Context, beans []*db.Account, cols ...dialect.Field) (bool, error)
	// DeleteById 按主键删除一条数据
	DeleteById(ctx context.Context, id types.BigInt) (bool, error)
	// SelectAll 读取所有数据
	SelectAll(ctx context.Context, s *ace.Selector) ([]db.Account, bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]db.Account, bool, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond []dialect.Condition, sort ...dialect.Order) ([]db.Account, bool, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (*db.Account, bool, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*db.Account, bool, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) (*db.Account, bool, error)
	// GetFirstCell 按条件读取第一行的第一个字段
	GetFirstCell(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (any, bool, error)
	//
	IDs(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) ([]any, error)
	//
	Columns(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]any, error)
}

type daoAccount struct {
	db ace.Executer
}

func Account(exec ...ace.Executer) accounter {
	_obj := &daoAccount{}
	if len(exec) > 0 {
		_obj.db = exec[0]
	} else {
		_obj.db = ace.GetDB()
	}
	return _obj
}

// C Create account
func (p *daoAccount) C() *ace.Creator {
	return p.db.C(db.AccountTableName)
}

// R Read account
func (p *daoAccount) R() *ace.Selector {
	return p.db.R(db.AccountTableName)
}

// U Update account
func (p *daoAccount) U() *ace.Updater {
	return p.db.U(db.AccountTableName)
}

// D Delete account
func (p *daoAccount) D() *ace.Deleter {
	return p.db.D(db.AccountTableName)
}

// Insert 返回 LastInsertId
func (p *daoAccount) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
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
func (p *daoAccount) InsertOne(ctx context.Context, bean *db.Account, cols ...dialect.Field) (bool, error) {
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
func (p *daoAccount) InsertBatch(ctx context.Context, beans []*db.Account, cols ...dialect.Field) (int64, error) {
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
func (p *daoAccount) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
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
func (p *daoAccount) UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error) {
	return p.Update(ctx,
		sets,
		tblaccount.PrimaryKey.Eq(id),
	)
}

// UpdateBatch 批量更新,禁止在事务中使用
// cols: 要更新的列名
func (p *daoAccount) UpdateBatch(ctx context.Context, beans []*db.Account, cols ...dialect.Field) (bool, error) {
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
func (p *daoAccount) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
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
func (p *daoAccount) DeleteById(ctx context.Context, id types.BigInt) (bool, error) {
	return p.Delete(ctx,
		tblaccount.PrimaryKey.Eq(id),
	)
}

// SelectAll 查询所有
func (p *daoAccount) SelectAll(ctx context.Context, s *ace.Selector) ([]db.Account, bool, error) {
	if len(s.GetTableName()) == 0 {
		s.SetTableName(db.AccountTableName)
	}

	_cols := s.GetCols()
	_rows, e := s.Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}
	defer _rows.Close()

	_obj := db.NewAccount()
	_objs, has, e := _obj.Scan(_rows, _cols...)
	if has {
		return _objs, true, nil
	}
	if e == nil || e == sql.ErrNoRows {
		return nil, false, nil
	}
	log.Error(e)
	return nil, false, e
}

// Get4Cols 先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoAccount) Get4Cols(ctx context.Context, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (*db.Account, bool, error) {
	_c := p.R()
	if len(cols) == 0 {
		_c.Cols(tblaccount.ReadableFields...)
	} else {
		_c.Cols(cols...)
	}

	_row, e := _c.Where(cond...).
		OrderFunc(sort...).
		QueryRow(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}

	_obj := db.NewAccount()

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

// Find4Cols 分页获取account slice对象，先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoAccount) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]db.Account, bool, error) {
	_c := p.R()
	if len(cols) == 0 {
		_c.Cols(tblaccount.ReadableFields...)
	} else {
		_c.Cols(cols...)
	}
	//
	_rows, e := _c.Where(cond...).
		OrderFunc(sort...).
		Page(pageIndex, pageSize).
		Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}
	defer _rows.Close()

	_obj := db.NewAccount()

	_objs, has, e := _obj.Scan(_rows, cols...)
	if has {
		return _objs, true, nil
	}
	if e == nil || e == sql.ErrNoRows {
		return nil, false, nil
	}
	log.Error(e)
	return nil, false, e
}

// GetByID 按主键读取一个account对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoAccount) GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*db.Account, bool, error) {
	return p.Get4Cols(ctx, cols, []dialect.Condition{tblaccount.PrimaryKey.Eq(id)})
}

// Get 按条件读取一个account对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoAccount) Get(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) (*db.Account, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond, sort...)
}

// GetFirstCell 按条件读取首行首列,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoAccount) GetFirstCell(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (any, bool, error) {
	_c := p.R().Cols(col)
	_row, e := _c.Where(cond...).
		OrderFunc(sort...).
		QueryRow(ctx)
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

// Find 按条件读取一个account slice对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoAccount) Find(ctx context.Context, pageIndex, pageSize uint, cond []dialect.Condition, sort ...dialect.Order) ([]db.Account, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond, sort...)
}

// IDs
func (p *daoAccount) IDs(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) ([]any, error) {
	_c := p.R().Cols(tblaccount.PrimaryKey)
	_rows, e := _c.Where(cond...).
		OrderFunc(sort...).
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
func (p *daoAccount) Columns(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]any, error) {
	_c := p.R().Cols(col)
	_rows, e := _c.Where(cond...).
		Limit(dialect.MaxLimit).
		OrderFunc(sort...).
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
func (p *daoAccount) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *daoAccount) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	return p.R().Sum(ctx, cols, cond...)
}

// Exists
func (p *daoAccount) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	_c := p.R().Cols(tblaccount.PrimaryKey).Where(cond...)
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
