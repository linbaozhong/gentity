package handler

import (
	"github.com/linbaozhong/gentity/pkg/api"
	"github.com/linbaozhong/gentity/pkg/serverpush"
	"github.com/linbaozhong/gentity/pkg/token"
	_ "reader/internal/model/dto"
	userService "reader/internal/service/user"
)

type user struct{}

func init() {
	api.RegisterRoute(&user{})
}

func (u *user) RegisterRoute(group api.Party) {
	g := api.NewParty(group, "/user")

	g.Post("/user_register", u.userRegister)
	g.Get("/get", u.get)
	g.Get("/sse", u.sse)
}

// @Summary 增加用工企业
// @Tags employmentCompany 用工企业
// @Accept  mpfd
// @Produce  json
// @Security ApiKeyAuth
// @Param user body dto.UserRegisterReq true "用户注册"
// @Success 200 {object} string ""
// @Router /v1/employmentCompany/employmentCompanyAdd [Post]
func (u *user) userRegister(c api.Context) {
	api.Post(c, userService.UserRegister)
}

func (u *user) get(c api.Context) {
	if api.ReadCache(c) {
		return
	}
	api.Get(c, userService.GetUser)
}

func (u *user) sse(c api.Context) {
	var _clientId string
	values := c.Request().URL.Query()

	_tk := values.Get("token")
	if _tk != "" {
		_clientId, _, _ = token.GetIDAndTokenFromCipher(_tk)
	}

	_lastEventId := values.Get("last_event_id")
	serverpush.ServeHTTP(c, _clientId, _lastEventId)
}
