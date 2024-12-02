package err

import "github.com/linbaozhong/gentity/pkg/types"

var (
	ErrUserName = types.NewError(501, "用户名不能为空")
	ErrPassword = types.NewError(502, "密码不能为空")
	ErrUserID   = types.NewError(503, "用户ID不能为空")

	Err_AuthToken_NotFound    = types.NewError(504, "token不存在")
	Err_Authorization_Limited = types.NewError(505, "无权限")
)
