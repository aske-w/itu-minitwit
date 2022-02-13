package controllers

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/entity"
	"aske-w/itu-minitwit/web/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	Ctx iris.Context

	DB *database.SQLite
	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

func (c *LoginController) Post() mvc.Result {
	username := c.Ctx.FormValue("username")
	password := c.Ctx.FormValue("password")
	error := ""
	var user entity.User
	err := c.DB.Get(c.Ctx, &user, "select * from user where username = ?", username)
	if err != nil {
		error = "Invalid username"
	} else {

		pwErr := bcrypt.CompareHashAndPassword([]byte(user.Pw_Hash), []byte(password))

		if pwErr != nil {
			error = "Invalid password"
		} else {
			c.Session.Set("user_id", user.User_id)
			c.Ctx.Redirect("/")
		}
	}

	return mvc.View{
		Name: "login.html",
		Data: iris.Map{"Title": "Login page", "error": error},
	}

}
func (c *LoginController) Get() mvc.Result {
	_, loggedIn := utils.GetUserIdFromSession(c.Session)
	if loggedIn {
		c.Ctx.Redirect("/")
	}
	return mvc.View{
		Name: "login.html",
		Data: iris.Map{"Title": "Login page"},
	}

}
