// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/table/tbldocumenttemplate"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
	"github.com/linbaozhong/gentity/pkg/log"
	"github.com/linbaozhong/gentity/pkg/types"
)

//type document_templateer interface {
//	dialect.Daoer
//	// InsertOne 插入一条数据，返回 LastInsertId
//	// cols: 要插入的列名
//	InsertOne(ctx context.Context, bean *db.DocumentTemplate, cols ...dialect.Field) (bool, error)
//	// InsertBatch 批量插入,返回 RowsAffected。禁止在事务中使用
//	// cols: 要插入的列名，如果为空，则插入结构体字段对应所有列
//	InsertBatch(ctx context.Context, beans []*db.DocumentTemplate, cols ...dialect.Field) (int64, error)
//  // UpdateById 按主键更新一条数据
//	UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error)
//	// UpdateBatch 批量更新,禁止在事务中使用
//	// cols: 要更新的列名，如果为空，则更新结构体所有字段对应列，包含零值字段
//	UpdateBatch(ctx context.Context, beans []*db.DocumentTemplate, cols ...dialect.Field) (bool, error)
//	// DeleteById 按主键删除一条数据
//  DeleteById(ctx context.Context,id types.BigInt) (bool, error)
//	// SelectAll 读取所有数据
//	SelectAll(ctx context.Context, s ace.SelectBuilder) ([]db.DocumentTemplate, bool, error)
//	// Find4Cols 分页查询指定列，返回一个slice
//	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]db.DocumentTemplate, bool, error)
//	// Find 分页查询，返回一个slice
//	Find(ctx context.Context, pageIndex, pageSize uint, cond []dialect.Condition, sort ...dialect.Order) ([]db.DocumentTemplate, bool, error)
//	// Get4Cols 读取一个对象的指定列
//	Get4Cols(ctx context.Context, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (*db.DocumentTemplate, bool, error)
//	// GetByID 按主键查询，返回一个对象
//	GetByID(ctx context.Context,id types.BigInt, cols ...dialect.Field) (*db.DocumentTemplate, bool, error)
//	// Get 按条件读取一个对象
//	Get(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) (*db.DocumentTemplate, bool, error)
//
//	// GetFirstCell 按条件读取第一行的第一个字段
//	GetFirstCell(ctx context.Context,col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (any, bool, error)
//	//
//	IDs(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) ([]any, error)
//	//
//	Columns(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]any, error)
//}

type daoDocumentTemplate struct {
	db ace.Executer
}

func DocumentTemplate(exec ...ace.Executer) *daoDocumentTemplate {
	_obj := &daoDocumentTemplate{}
	_obj.db = ace.GetExec(exec...)
	return _obj
}

// Create document_template
func DocumentTemplateCreate(exec ...ace.Executer) ace.CreateBuilder {
	return ace.Create(exec...).Table(db.DocumentTemplateTableName)
}

// Select document_template
func DocumentTemplateSelect(exec ...ace.Executer) ace.SelectBuilder {
	return ace.Select(exec...).Table(db.DocumentTemplateTableName)
}

// Update document_template
func DocumentTemplateUpdate(exec ...ace.Executer) ace.UpdateBuilder {
	return ace.Update(exec...).Table(db.DocumentTemplateTableName)
}

// Delete document_template
func DocumentTemplateDelete(exec ...ace.Executer) ace.DeleteBuilder {
	return ace.Delete(exec...).Table(db.DocumentTemplateTableName)
}

// Insert 返回 LastInsertId
func (p *daoDocumentTemplate) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
	if len(sets) == 0 {
		return 0, dialect.ErrSetterEmpty
	}
	_result, e := DocumentTemplateCreate(p.db).
		Sets(sets...).
		Exec(ctx)
	if e != nil {
		log.Error(e)
		return 0, e
	}
	return _result.LastInsertId()
}

// InsertOne 返回 LastInsertId
// cols: 要插入的列名
func (p *daoDocumentTemplate) InsertOne(ctx context.Context, bean *db.DocumentTemplate, cols ...dialect.Field) (bool, error) {
	_result, e := DocumentTemplateCreate(p.db).
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
// cols: 要插入的列名，如果为空，则插入结构体字段对应所有列
func (p *daoDocumentTemplate) InsertBatch(ctx context.Context, beans []*db.DocumentTemplate, cols ...dialect.Field) (int64, error) {
	_lens := len(beans)
	if _lens == 0 {
		return 0, dialect.ErrBeanEmpty
	}
	_args := make([]dialect.Modeler, 0, _lens)
	for _, _bean := range beans {
		_args = append(_args, _bean)
	}
	_result, e := DocumentTemplateCreate(p.db).
		Cols(cols...).
		BatchStruct(ctx, _args...)
	if e != nil {
		log.Error(e)
		return 0, e
	}

	return _result.RowsAffected()
}

// Update
func (p *daoDocumentTemplate) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
	if len(sets) == 0 {
		return false, dialect.ErrSetterEmpty
	}
	_result, e := DocumentTemplateUpdate(p.db).
		Where(cond...).
		Sets(sets...).
		Exec(ctx)
	if e != nil {
		log.Error(e)
		return false, e
	}
	_n, e := _result.RowsAffected()
	return _n >= 0, e
}

// UpdateById
func (p *daoDocumentTemplate) UpdateById(ctx context.Context, id types.BigInt, sets ...dialect.Setter) (bool, error) {
	return p.Update(ctx,
		sets,
		tbldocumenttemplate.PrimaryKey.Eq(id),
	)
}

// UpdateBatch 批量更新,禁止在事务中使用
// cols: 要更新的列名，如果为空，则更新结构体所有字段对应列，包含零值字段
func (p *daoDocumentTemplate) UpdateBatch(ctx context.Context, beans []*db.DocumentTemplate, cols ...dialect.Field) (bool, error) {
	_lens := len(beans)
	if _lens == 0 {
		return false, dialect.ErrBeanEmpty
	}
	_args := make([]dialect.Modeler, 0, _lens)
	for _, _bean := range beans {
		_args = append(_args, _bean)
	}
	_result, e := DocumentTemplateUpdate(p.db).
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
func (p *daoDocumentTemplate) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	_result, e := DocumentTemplateDelete(p.db).
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
func (p *daoDocumentTemplate) DeleteById(ctx context.Context, id types.BigInt) (bool, error) {
	return p.Delete(ctx,
		tbldocumenttemplate.PrimaryKey.Eq(id),
	)
}

// SelectAll 查询所有
func (p *daoDocumentTemplate) SelectAll(ctx context.Context, s ace.SelectBuilder) ([]db.DocumentTemplate, bool, error) {
	if len(s.GetTableName()) == 0 {
		s.SetTableName(db.DocumentTemplateTableName)
	}

	_cols := s.GetCols()
	if len(_cols) == 0 {
		_cols = tbldocumenttemplate.ReadableFields
		s.Cols(_cols...)
	}

	_rows, e := s.Query(ctx)
	if e != nil {
		log.Error(e)
		return nil, false, e
	}
	defer _rows.Close()

	_obj := db.NewDocumentTemplate()
	_objs, has, e := _obj.Scan(_rows, _cols...)
	if has {
		return _objs, true, nil
	}
	if e == nil || e == sql.ErrNoRows {
		return _objs, false, nil
	}
	log.Error(e)
	return _objs, false, e
}

// Get4Cols 先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoDocumentTemplate) Get4Cols(ctx context.Context, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (*db.DocumentTemplate, bool, error) {
	_c := DocumentTemplateSelect(p.db)
	if len(cols) == 0 {
		_c.Cols(tbldocumenttemplate.ReadableFields...)
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

	_obj := db.NewDocumentTemplate()

	e = _row.Scan(_obj.AssignPtr(cols...)...)
	switch e {
	case nil:
		return _obj, true, nil
	case sql.ErrNoRows:
		return _obj, false, nil
	default:
		log.Error(e)
		return _obj, false, e
	}
}

// Find4Cols 分页获取document_template slice对象，先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoDocumentTemplate) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]db.DocumentTemplate, bool, error) {
	_c := DocumentTemplateSelect(p.db)
	if len(cols) == 0 {
		_c.Cols(tbldocumenttemplate.ReadableFields...)
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

	_obj := db.NewDocumentTemplate()

	_objs, has, e := _obj.Scan(_rows, cols...)
	if has {
		return _objs, true, nil
	}
	if e == nil || e == sql.ErrNoRows {
		return _objs, false, nil
	}
	log.Error(e)
	return _objs, false, e
}

// GetByID 按主键读取一个document_template对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoDocumentTemplate) GetByID(ctx context.Context, id types.BigInt, cols ...dialect.Field) (*db.DocumentTemplate, bool, error) {
	return p.Get4Cols(ctx, cols, []dialect.Condition{tbldocumenttemplate.PrimaryKey.Eq(id)})
}

// Get 按条件读取一个document_template对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoDocumentTemplate) Get(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) (*db.DocumentTemplate, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond, sort...)
}

// GetFirstCell 按条件读取首行首列,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoDocumentTemplate) GetFirstCell(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) (any, bool, error) {
	_c := DocumentTemplateSelect(p.db).Cols(col)
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
	case nil:
		return _v, true, nil
	case sql.ErrNoRows:
		return _v, false, nil
	default:
		log.Error(e)
		return _v, false, e
	}
}

// Find 按条件读取一个document_template slice对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *daoDocumentTemplate) Find(ctx context.Context, pageIndex, pageSize uint, cond []dialect.Condition, sort ...dialect.Order) ([]db.DocumentTemplate, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond, sort...)
}

// IDs
func (p *daoDocumentTemplate) IDs(ctx context.Context, cond []dialect.Condition, sort ...dialect.Order) ([]any, error) {
	_c := DocumentTemplateSelect(p.db).Cols(tbldocumenttemplate.PrimaryKey)
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
func (p *daoDocumentTemplate) Columns(ctx context.Context, col dialect.Field, cond []dialect.Condition, sort ...dialect.Order) ([]any, error) {
	_c := DocumentTemplateSelect(p.db).Cols(col)
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
func (p *daoDocumentTemplate) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return DocumentTemplateSelect(p.db).Count(ctx, cond...)
}

// Sum
func (p *daoDocumentTemplate) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	return DocumentTemplateSelect(p.db).Sum(ctx, cols, cond...)
}

// Exists
func (p *daoDocumentTemplate) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	_c := DocumentTemplateSelect(p.db).Cols(tbldocumenttemplate.PrimaryKey).Where(cond...)
	_row, e := _c.QueryRow(ctx)
	if e != nil {
		log.Error(e)
		return false, e
	}

	var id types.BigInt
	e = _row.Scan(&id)
	switch e {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, e
	}
}
