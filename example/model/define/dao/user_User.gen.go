// Code generated by github.com/linbaozhong/gentity. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"
	"github.com/linbaozhong/gentity/example/model/db"
	"github.com/linbaozhong/gentity/example/model/define/table/usertbl"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/ace/dialect"
)

type UserDaoer interface {
	dialect.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *db.User, cols ...dialect.Field) (bool, error)
	// InsertBatch 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertBatch(ctx context.Context, beans []*db.User, cols ...dialect.Field) (int64, error)
	// UpdateById 按主键更新一条数据
	UpdateById(ctx context.Context, id uint64, sets ...dialect.Setter) (bool, error)
	// UpdateBatch 批量更新多条数据
	// cols: 要更新的列名
	UpdateBatch(ctx context.Context, beans []*db.User, cols ...dialect.Field) (bool, error)
	// DeleteById 按主键删除一条数据
	DeleteById(ctx context.Context, id uint64) (bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]*db.User, bool, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]*db.User, bool, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*db.User, bool, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, id uint64, cols ...dialect.Field) (*db.User, bool, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond ...dialect.Condition) (*db.User, bool, error)
	// GetFirstCell 按条件读取第一行的第一个字段
	GetFirstCell(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (any, bool, error)
	//
	IDs(ctx context.Context, cond ...dialect.Condition) ([]any, error)
	//
	Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error)
}

type userDao struct {
	db ace.Executer
}

func User(exec ace.Executer) UserDaoer {
	obj := &userDao{}
	obj.db = exec
	return obj
}

// C Create user
func (p *userDao) C() *ace.Creator {
	return p.db.C(db.UserTableName)
}

// R Read user
func (p *userDao) R() *ace.Selector {
	return p.db.R(db.UserTableName)
}

// U Update user
func (p *userDao) U() *ace.Updater {
	return p.db.U(db.UserTableName)
}

// D Delete user
func (p *userDao) D() *ace.Deleter {
	return p.db.D(db.UserTableName)
}

// Insert 返回 LastInsertId
func (p *userDao) Insert(ctx context.Context, sets ...dialect.Setter) (int64, error) {
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
func (p *userDao) InsertOne(ctx context.Context, bean *db.User, cols ...dialect.Field) (bool, error) {
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

// InsertBatch 批量插入,返回 RowsAffected
// cols: 要插入的列名
func (p *userDao) InsertBatch(ctx context.Context, beans []*db.User, cols ...dialect.Field) (int64, error) {
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
		Struct(ctx, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Update
func (p *userDao) Update(ctx context.Context, sets []dialect.Setter, cond ...dialect.Condition) (bool, error) {
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
func (p *userDao) UpdateById(ctx context.Context, id uint64, sets ...dialect.Setter) (bool, error) {
	return p.Update(ctx,
		sets,
		usertbl.PrimaryKey.Eq(id),
	)
}

// UpdateBatch
// cols: 要更新的列名
func (p *userDao) UpdateBatch(ctx context.Context, beans []*db.User, cols ...dialect.Field) (bool, error) {
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
		Struct(ctx, args...)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Delete
func (p *userDao) Delete(ctx context.Context, cond ...dialect.Condition) (bool, error) {
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
func (p *userDao) DeleteById(ctx context.Context, id uint64) (bool, error) {
	return p.Delete(ctx,
		usertbl.PrimaryKey.Eq(id),
	)
}

// Get4Cols 先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *userDao) Get4Cols(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (*db.User, bool, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols(usertbl.ReadableFields...)
	} else {
		c.Cols(cols...)
	}

	row, err := c.Where(cond...).
		QueryRow(ctx)
	if err != nil {
		return nil, false, err
	}

	obj := db.NewUser()

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

// Find4Cols 分页获取user} slice对象，先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *userDao) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []dialect.Field, cond ...dialect.Condition) ([]*db.User, bool, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols(usertbl.ReadableFields...)
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

	obj := db.NewUser()

	objs, has, err := obj.Scan(rows, cols...)
	if has {
		return objs, true, err
	}
	return nil, false, err
}

// GetByID 按主键读取一个user对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *userDao) GetByID(ctx context.Context, id uint64, cols ...dialect.Field) (*db.User, bool, error) {
	return p.Get4Cols(ctx, cols, usertbl.PrimaryKey.Eq(id))
}

// Get 按条件读取一个user对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *userDao) Get(ctx context.Context, cond ...dialect.Condition) (*db.User, bool, error) {
	return p.Get4Cols(ctx, []dialect.Field{}, cond...)
}

// GetFirstCell 按条件读取首行首列,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *userDao) GetFirstCell(ctx context.Context, col dialect.Field, cond ...dialect.Condition) (any, bool, error) {
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

// Find 按条件读取一个user slice对象,先判断第二返回值是否为true,再判断是否第三返回值为nil
func (p *userDao) Find(ctx context.Context, pageIndex, pageSize uint, cond ...dialect.Condition) ([]*db.User, bool, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []dialect.Field{}, cond...)
}

// IDs
func (p *userDao) IDs(ctx context.Context, cond ...dialect.Condition) ([]any, error) {
	c := p.R().Cols(usertbl.PrimaryKey)
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
func (p *userDao) Columns(ctx context.Context, col dialect.Field, cond ...dialect.Condition) ([]any, error) {
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
func (p *userDao) Count(ctx context.Context, cond ...dialect.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *userDao) Sum(ctx context.Context, cols []dialect.Field, cond ...dialect.Condition) (map[string]any, error) {
	return p.R().Sum(ctx, cols, cond...)
}

// Exists
func (p *userDao) Exists(ctx context.Context, cond ...dialect.Condition) (bool, error) {
	c := p.R().Cols(usertbl.PrimaryKey).Where(cond...)
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
func (p *userDao) onUpdate(ctx context.Context, ids ...uint64) error {
	return nil
}
