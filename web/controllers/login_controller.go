package controllers

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/entity"

	"fmt"

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
	err := c.DB.Get(c.Ctx, &user, "select * from users where username = ?", username)
	if err != nil {
		error = "Invalid username"
	} else {

		fmt.Println(user)
		pwErr := bcrypt.CompareHashAndPassword([]byte(user.Pw_Hash), []byte(password))

		if pwErr != nil {
			error = "Invalid password"
		} else {
			c.Session.Set("user_id", user.ID)
			c.Ctx.Redirect("/")
		}
	}

	return mvc.View{
		Name: "login.html",
		Data: iris.Map{"Title": "Login page", "error": error},
	}

}
func (c *LoginController) Get() mvc.Result {
	_, loggedIn := getUserFromSession(c.Session)
	if loggedIn {
		c.Ctx.Redirect("/")
	}
	return mvc.View{
		Name: "login.html",
		Data: iris.Map{"Title": "Login page"},
	}

}

/*
Returns the user_id from the session as the first in the tuple and if it was succesful in the second part
*/
func getUserFromSession(session *sessions.Session) (string, bool) {
	user_id := session.GetString("user_id")
	if len(user_id) == 0 {
		return "", false
	}
	return user_id, true
}
