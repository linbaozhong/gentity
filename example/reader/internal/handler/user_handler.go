package handler

import (
	"github.com/linbaozhong/gentity/pkg/api"
	userService "reader/internal/service/user"
)

type user struct{}

func init() {
	api.Instances = append(api.Instances, &user{})
}

func (u *user) RegisterRouter(group api.Party) {
	g := api.NewParty(group, "/user")

	g.Post("/user_register", u.userRegister)
	g.Get("/get", u.get)
}

func (u *user) userRegister(c api.Context) {
	api.Post(c, userService.UserRegister)
}

func (u *user) get(c api.Context) {
	api.Get(c, userService.GetUser)
}
