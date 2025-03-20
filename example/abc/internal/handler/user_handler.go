package handler

import (
	userService "abc/internal/service/user"
	"github.com/linbaozhong/gentity/pkg/api"
)

type user struct{}

func init() {
	api.RegisterRoute(&user{})
}

func (u *user) RegisterRoute(group api.Party) {
	_g := api.NewParty(group, "/user")

	_g.Post("/user_register", u.userRegister)
	_g.Get("/get", u.get)
}

func (u *user) userRegister(c api.Context) {
	api.Post(c, userService.UserRegister)
}

func (u *user) get(c api.Context) {
	api.Get(c, userService.GetUser)
}
