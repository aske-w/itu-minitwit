package controllers

// import (
// 	"aske-w/itu-minitwit/database"
// 	"aske-w/itu-minitwit/web/utils"

// 	"github.com/kataras/iris/v12"
// 	"github.com/kataras/iris/v12/mvc"
// 	"github.com/kataras/iris/v12/sessions"
// 	"golang.org/x/crypto/bcrypt"
// )

// type SignupController struct {
// 	Ctx iris.Context

// 	DB *database.SQLite
// 	// Session, binded using dependency injection from the main.go.
// 	Session *sessions.Session
// }

// func (c *SignupController) Post() mvc.Result {
// 	username := c.Ctx.FormValue("username")
// 	email := c.Ctx.FormValue("email")
// 	password := c.Ctx.FormValue("password")
// 	password2 := c.Ctx.FormValue("password2")

// 	error := ""

// 	if username == "" {
// 		error = "You have to enter a username"
// 	} else if email == "" {
// 		error = "You have to enter a valid email address"
// 	} else if password == "" {
// 		error = "You have to enter a password"
// 	} else if password2 != password {
// 		error = "The two passwords do not match"
// 	} else {

// 		_, err := utils.GetUserByUsername(username, c.DB, c.Ctx)

// 		if err != nil {
// 			error = "The username is already taken"
// 		} else {

// 			byteHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 			if err != nil {
// 				error = err.Error()
// 			} else {
// 				_, err := c.DB.db.Exec(`insert into user (username, email, pw_hash) values (?,?,?)`, username, email, string(byteHash))

// 				if err != nil {
// 					error = err.Error()
// 				} else {
// 					c.Ctx.Redirect("/login")
// 				}
// 			}

// 		}
// 	}
// 	return mvc.View{
// 		Name: "signup.html",
// 		Data: iris.Map{"Title": "Signup page", "error": error},
// 	}
// }

// func (c *SignupController) Get() mvc.Result {
// 	_, loggedIn := utils.GetUserIdFromSession(c.Session)
// 	if loggedIn {
// 		c.Ctx.Redirect("/")
// 	}
// 	return mvc.View{
// 		Name: "signup.html",
// 		Data: iris.Map{"Title": "Signup age"},
// 	}

// }
