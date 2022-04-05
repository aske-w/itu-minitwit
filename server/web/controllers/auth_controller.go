package controllers

import (
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type AuthController struct {
	Ctx iris.Context

	AuthService *services.AuthService
	UserService *services.UserService

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

// func (c *AuthController) BeforeActivation(b mvc.BeforeActivation) {
// 	// Register endpoints on /api
// 	// b.Handle("GET", "/signup", "SignupViewHandler")

// }

func (c *AuthController) GetSignup() mvc.Result {
	_, loggedIn := utils.GetUserIdFromSession(c.Session)
	if loggedIn {
		c.Ctx.Redirect("/")
	}
	return mvc.View{
		Name: "signup.html",
		Data: iris.Map{"Title": "Signup age"},
	}
}
func (c *AuthController) PostSignup() mvc.Result {
	username := c.Ctx.FormValue("username")
	email := c.Ctx.FormValue("email")
	password := c.Ctx.FormValue("password")
	password2 := c.Ctx.FormValue("password2")

	error := ""

	if username == "" {
		error = "You have to enter a username"
	} else if email == "" {
		error = "You have to enter a valid email address"
	} else if password == "" {
		error = "You have to enter a password"
	} else if password2 != password {
		error = "The two passwords do not match"
	} else {

		user, _ := c.UserService.FindByUsername(username)

		if user != nil {
			error = "The username is already taken"
		} else {
			_, err := c.AuthService.CreateUser(username, email, password)
			if err != nil {
				error = err.Error()
			} else {
				c.Ctx.Redirect("/login")
			}

		}
	}
	return mvc.View{
		Name: "signup.html",
		Data: iris.Map{"Title": "Signup page", "error": error},
	}
}

func (c *AuthController) GetLogin() mvc.Result {
	_, loggedIn := utils.GetUserIdFromSession(c.Session)
	if loggedIn {
		c.Ctx.Redirect("/")
	}
	return mvc.View{
		Name: "login.html",
		Data: iris.Map{"Title": "Login page"},
	}
}

func (c *AuthController) PostLogin() mvc.Result {
	username := c.Ctx.FormValue("username")
	password := c.Ctx.FormValue("password")
	error := ""
	user, _ := c.UserService.FindByUsername(username)
	if user == nil {
		error = "Invalid username"
	} else {
		passwordMatch := c.AuthService.CheckPassword(user, password)

		if !passwordMatch {
			error = "Invalid password"
		} else {
			c.Session.Set("user_id", int(user.ID))
			c.Ctx.Redirect("/")
		}
	}

	return mvc.View{
		Name: "login.html",
		Data: iris.Map{"Title": "Login page", "error": error},
	}
}

func (c *AuthController) GetLogout() mvc.Result {

	c.Session.Set("user_id", nil)

	c.Ctx.Redirect("/")
	return mvc.View{
		Name: "layout.html",
		Data: iris.Map{"Title": "Logout"},
	}
}
