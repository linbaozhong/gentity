package handler

import (
	"github.com/linbaozhong/gentity/pkg/api"
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
	api.Get(c, userService.GetUser)
}
