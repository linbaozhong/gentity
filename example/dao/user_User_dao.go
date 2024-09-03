package dao

import (
	"context"
	"github.com/linbaozhong/gentity/example/model"
	"github.com/linbaozhong/gentity/pkg/ace"
	atype "github.com/linbaozhong/gentity/pkg/ace/types"
)

type UserDaoer interface {
	atype.Daoer
	ace.Cruder
	// InsertOne 插入一条数据，返回 LastInsertId
	// cols: 要插入的列名
	InsertOne(ctx context.Context, bean *model.User, cols ...atype.Field) (int64, error)
	// InsertMulti 批量插入多条数据,返回 RowsAffected
	// cols: 要插入的列名
	InsertMulti(ctx context.Context, beans []*model.User, cols ...atype.Field) (int64, error)
	// UpdateMulti 批量更新多条数据
	// cols: 要更新的列名
	UpdateMulti(ctx context.Context, beans []*model.User, cols ...atype.Field) (bool, error)
	// Find4Cols 分页查询指定列，返回一个slice
	Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []atype.Field, cond ...atype.Condition) ([]*model.User, error)
	// Find 分页查询，返回一个slice
	Find(ctx context.Context, pageIndex, pageSize uint, cond ...atype.Condition) ([]*model.User, error)
	// Get4Cols 读取一个对象的指定列
	Get4Cols(ctx context.Context, cols []atype.Field, cond ...atype.Condition) (*model.User, error)
	// GetByID 按主键查询，返回一个对象
	GetByID(ctx context.Context, args ...any) (*model.User, error)
	// Get 按条件读取一个对象
	Get(ctx context.Context, cond ...atype.Condition) (*model.User, error)
}

type userDao struct {
	db atype.Executer
}

func User(exec atype.Executer) UserDaoer {
	return &userDao{db: exec}
}

// C Create user
func (p *userDao) C() *ace.Creator {
	return ace.NewCreate(p.db, model.UserTableName)
}

// R Read user
func (p *userDao) R() *ace.Selector {
	return ace.NewSelect(p.db, model.UserTableName)
}

// U Update user
func (p *userDao) U() *ace.Updater {
	return ace.NewUpdate(p.db, model.UserTableName)
}

// D Delete user
func (p *userDao) D() *ace.Deleter {
	return ace.NewDelete(p.db, model.UserTableName)
}

// Insert 返回 LastInsertId
func (p *userDao) Insert(ctx context.Context, sets ...atype.Setter) (int64, error) {
	if len(sets) == 0 {
		return 0, atype.ErrSetterEmpty
	}
	result, err := p.C().Set(sets...).Do(ctx)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertOne 返回 LastInsertId
// cols: 要插入的列名
func (p *userDao) InsertOne(ctx context.Context, bean *model.User, cols ...atype.Field) (int64, error) {
	result, err := p.C().Cols(cols...).Struct(ctx, bean)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertMulti 批量插入,返回 RowsAffected
// cols: 要插入的列名
func (p *userDao) InsertMulti(ctx context.Context, beans []*model.User, cols ...atype.Field) (int64, error) {
	lens := len(beans)
	if lens == 0 {
		return 0, atype.ErrBeanEmpty
	}
	args := make([]atype.Modeler, 0, lens)
	for _, bean := range beans {
		args = append(args, bean)
	}
	result, err := p.C().Cols(cols...).Struct(ctx, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Update
func (p *userDao) Update(ctx context.Context, sets []atype.Setter, cond ...atype.Condition) (bool, error) {
	if len(sets) == 0 {
		return false, atype.ErrSetterEmpty
	}
	result, err := p.U().Where(cond...).Set(sets...).Do(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// UpdateMulti
// cols: 要更新的列名
func (p *userDao) UpdateMulti(ctx context.Context, beans []*model.User, cols ...atype.Field) (bool, error) {
	lens := len(beans)
	if lens == 0 {
		return false, atype.ErrBeanEmpty
	}
	args := make([]atype.Modeler, 0, lens)
	for _, bean := range beans {
		args = append(args, bean)
	}
	result, err := p.U().Cols(cols...).Struct(ctx, args...)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Delete
func (p *userDao) Delete(ctx context.Context, cond ...atype.Condition) (bool, error) {
	result, err := p.D().Where(cond...).Do(ctx)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	return n >= 0, err
}

// Get4Cols
func (p *userDao) Get4Cols(ctx context.Context, cols []atype.Field, cond ...atype.Condition) (*model.User, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols(model.UserReadableFields...)
	} else {
		c.Cols(cols...)
	}

	rows, err := c.Where(cond...).Limit(1).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	obj := model.NewUser()
	defer obj.Free()

	objs, err := obj.Scan(rows, cols...)
	if err != nil {
		return nil, err
	}
	if len(objs) == 0 {
		return nil, atype.ErrNotFound
	}
	return objs[0], nil
}

// Find4Cols
func (p *userDao) Find4Cols(ctx context.Context, pageIndex, pageSize uint, cols []atype.Field, cond ...atype.Condition) ([]*model.User, error) {
	c := p.R()
	if len(cols) == 0 {
		c.Cols(model.UserReadableFields...)
	} else {
		c.Cols(cols...)
	}
	//
	if pageSize == 0 {
		pageSize = atype.PageSize
	}
	//
	rows, err := c.Where(cond...).Limit(pageSize, pageSize*pageIndex).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	obj := model.NewUser()
	defer obj.Free()

	objs, err := obj.Scan(rows, cols...)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

// GetByID Read one user By Primary Key value,
// Pass values in this order：ID
func (p *userDao) GetByID(ctx context.Context, args ...any) (*model.User, error) {
	lens := len(model.UserPrimaryKeys)
	if lens != len(args) {
		return nil, atype.ErrArgsNotMatch
	}

	cond := make([]atype.Condition, 0, lens)
	for i, key := range model.UserPrimaryKeys {
		cond = append(cond, key.Eq(args[i]))
	}
	return p.Get4Cols(ctx, []atype.Field{}, cond...)
}

// Get Read one user
func (p *userDao) Get(ctx context.Context, cond ...atype.Condition) (*model.User, error) {
	return p.Get4Cols(ctx, []atype.Field{}, cond...)
}

// Find
func (p *userDao) Find(ctx context.Context, pageIndex, pageSize uint, cond ...atype.Condition) ([]*model.User, error) {
	return p.Find4Cols(ctx, pageIndex, pageSize, []atype.Field{}, cond...)
}

// IDs
func (p *userDao) IDs(ctx context.Context, cond ...atype.Condition) ([]int64, error) {
	if len(model.UserPrimaryKeys) == 0 {
		return nil, atype.ErrPrimaryKeyNotMatch
	}
	c := p.R().Cols(model.UserPrimaryKeys[0])
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, atype.PageSize)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// Columns
func (p *userDao) Columns(ctx context.Context, col atype.Field, cond ...atype.Condition) ([]any, error) {
	c := p.R().Cols(col)
	rows, err := c.Where(cond...).Limit(atype.MaxLimit).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]any, atype.PageSize)
	for rows.Next() {
		var v any
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		cols = append(cols, v)
	}

	return cols, nil
}

// Count
func (p *userDao) Count(ctx context.Context, cond ...atype.Condition) (int64, error) {
	return p.R().Count(ctx, cond...)
}

// Sum
func (p *userDao) Sum(ctx context.Context, col atype.Field, cond ...atype.Condition) (int64, error) {
	return p.R().Sum(ctx, col, cond...)
}

// Exists
func (p *userDao) Exists(ctx context.Context, cond ...atype.Condition) (bool, error) {
	if len(model.UserPrimaryKeys) == 0 {
		return false, atype.ErrPrimaryKeyNotMatch
	}
	c := p.R().Cols(model.UserPrimaryKeys[0]).Where(cond...).Limit(1)
	rows, err := c.Query(ctx)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
