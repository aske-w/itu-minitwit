package controllers

import (
	"aske-w/itu-minitwit/database"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type LogoutController struct {
	Ctx iris.Context

	DB *database.SQLite
	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

func (c *LogoutController) Get() mvc.Result {

	c.Session.Set("user_id", nil)

	c.Ctx.Redirect("/")
	return mvc.View{
		Name: "index.html",
		Data: iris.Map{"Title": "Logout"},
	}
}
