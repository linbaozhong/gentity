package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linbaozhong/gentity/pkg/ace"
	"github.com/linbaozhong/gentity/pkg/types"
	"tui/internal/model/define/dao"
	"tui/internal/model/define/table/tbluser"
	"tui/internal/model/do"
	"tui/internal/model/dto"
)

func UserRegister(c context.Context, in *dto.UserRegisterReq, out *dto.UserRegisterResp) error {
	// todo: 在这里做用户注册，返回用户信息
	fmt.Println("UserRegister:", in)

	out.UserID = 12345678
	out.UserName = "哈利蔺特"

	return nil
}

func GetUser(c context.Context, in *dto.GetUserReq, out *dto.GetUserResp) error {
	user := do.NewUser()
	// 第一种方法
	e := db.Table(do.UserTableName).
		Where(tbluser.Id.Eq(in.UserID)).Select().Get(c, user)

	// 第二种方法
	e = ace.Table(do.UserTableName).
		Where(tbluser.Id.Eq(in.UserID)).Select().Get(c, user)

	// 第三种方法
	user, has, e := dao.User(db).GetByID(c, types.BigInt(*in.UserID))

	// 第四种方法
	user, has, e = dao.User(db).Get(c, ace.Where(tbluser.Id.Eq(in.UserID)))

	if e != nil {
		return e
	}
	if !has {
		return errors.New("user not found")
	}

	// 上面四种方法都可以，根据个人喜好选择一种

	// 这里是将do.User转换为dto.GetUserResp
	out.UserID = user.Id.Uint64()
	out.Email = user.Email.String()
	out.UserName = user.Name.String()

	return nil
}
