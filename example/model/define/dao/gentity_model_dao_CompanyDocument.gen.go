// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/table/tblcompanydocument"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

type company_documenter interface {
	dialect.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *db.CompanyDocument, cols ...dialect.Field) (bool, error)
	// InsertBatch 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertBatch(ctx context.Context, beans []*db.CompanyDocument, cols ...dialect.Field) (int64, error)
	// UpdateById 按主键更新一条数据
	UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error)
	// UpdateBatch 批量更新多条数据
	// cols: 要更新的列名
	UpdateBatch(ctx context.Context, beans []*db.CompanyDocument, cols ...dialect.Field) (bool, error)
	// DeleteById 按主键删除一条数据
	DeleteById(ctx context.Context, id types.BigInt) (bool, error)
	// SelectAll 读取所有数据
	SelectAll(ctx context.Context, s *ace.Selector) ([]db.CompanyDocument, bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]db.CompanyDocument, bool, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond []dialect.Condition, sort ...dialect.Order) ([]db.CompanyDocument, bool, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (*db.CompanyDocument, bool, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*db.CompanyDocument, bool, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) (*db.CompanyDocument, bool, error)
	// GetFirstCell 按条件读取第一行的第一个字段
	GetFirstCell(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (any, bool, error)
	//
	IDs(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) ([]any, error)
	//
	Columns(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]any, error)
}

type daoCompanyDocument struct {
	db ace.Executer
}

func CompanyDocument(exec ...ace.Executer) company_documenter {
	_obj := &daoCompanyDocument{}
	if len(exec) > 0 {
		_obj.db = exec[0]
	} else {
		_obj.db = ace.GetDB()
	}
	return _obj
}

// C Create company_document
func (p *daoCompanyDocument) C() *ace.Creator {
	return p.db.C(db.CompanyDocumentTableName)
}

// R Read company_document
func (p *daoCompanyDocument) R() *ace.Selector {
	return p.db.R(db.CompanyDocumentTableName)
}

// U Update company_document
func (p *daoCompanyDocument) U() *ace.Updater {
	return p.db.U(db.CompanyDocumentTableName)
}

// D Delete company_document
func (p *daoCompanyDocument) D() *ace.Deleter {
	return p.db.D(db.CompanyDocumentTableName)
}

// Insert 返回 LastInsertId
func (p *daoCompanyDocument) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
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
func (p *daoCompanyDocument) InsertOne(ctx context.Context, bean *db.CompanyDocument, cols ...dialect.Field) (bool, error) {
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
func (p *daoCompanyDocument) InsertBatch(ctx context.Context, beans []*db.CompanyDocument, cols ...dialect.Field) (int64, error) {
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
func (p *daoCompanyDocument) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
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
func (p *daoCompanyDocument) UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error) {
	return p.Update(ctx,
		sets,
		tblcompanydocument.PrimaryKey.Eq(id),
	)
}

// UpdateBatch 批量更新,禁止在事务中使用
// cols: 要更新的列名
func (p *daoCompanyDocument) UpdateBatch(ctx context.Context, beans []*db.CompanyDocument, cols ...dialect.Field) (bool, error) {
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
func (p *daoCompanyDocument) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
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
func (p *daoCompanyDocument) DeleteById(ctx context.Context, id types.BigInt) (bool, error) {
	return p.Delete(ctx,
		tblcompanydocument.PrimaryKey.Eq(id),
	)
}

// SelectAll 查询所有
func (p *daoCompanyDocument) SelectAll(ctx context.Context, s *ace.Selector) ([]db.CompanyDocument, bool, error) {
	if len(s.GetTableName()) == 0 {
		s.SetTableName(db.CompanyDocumentTableName)
	}

	_cols := s.GetCols()
	_rows, e := s.Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}
	defer _rows.Close()

	_obj := db.NewCompanyDocument()
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
func (p *daoCompanyDocument) Get4Cols(ctx context.Context, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (*db.CompanyDocument, bool, error) {
	_c := p.R()
	if len(cols) == 0 {
		_c.Cols(tblcompanydocument.ReadableFields...)
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

	_obj := db.NewCompanyDocument()

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

// Find4Cols 分页获取company_document slice对象，先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoCompanyDocument) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]db.CompanyDocument, bool, error) {
	_c := p.R()
	if len(cols) == 0 {
		_c.Cols(tblcompanydocument.ReadableFields...)
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

	_obj := db.NewCompanyDocument()

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

// GetByID 按主键读取一个company_document对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoCompanyDocument) GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*db.CompanyDocument, bool, error) {
	return p.Get4Cols(ctx, cols, []dialect.Condition{tblcompanydocument.PrimaryKey.Eq(id)})
}

// Get 按条件读取一个company_document对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoCompanyDocument) Get(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) (*db.CompanyDocument, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond, sort...)
}

// GetFirstCell 按条件读取首行首列,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoCompanyDocument) GetFirstCell(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (any, bool, error) {
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

// Find 按条件读取一个company_document slice对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoCompanyDocument) Find(ctx context.Context, pageIndex, pageSize uint, cond []dialect.Condition, sort ...dialect.Order) ([]db.CompanyDocument, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond, sort...)
}

// IDs
func (p *daoCompanyDocument) IDs(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) ([]any, error) {
	_c := p.R().Cols(tblcompanydocument.PrimaryKey)
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
func (p *daoCompanyDocument) Columns(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]any, error) {
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
func (p *daoCompanyDocument) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *daoCompanyDocument) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	return p.R().Sum(ctx, cols, cond...)
}

// Exists
func (p *daoCompanyDocument) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	_c := p.R().Cols(tblcompanydocument.PrimaryKey).Where(cond...)
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
