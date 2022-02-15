package main

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/environment"
	"aske-w/itu-minitwit/web/controllers"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func main() {

	app := iris.New()
	// app.Logger().SetLevel("debug") // more logging

	// Load env's
	environment.InitEnv()

	app.Use(logger.New())  // logs request
	app.Use(recover.New()) // handles panics (shows 404)

	// Configure sessions manager.
	sess := sessions.New(sessions.Config{Cookie: "itu-minitwit-cookie"})
	app.Use(sess.Handler())

	// Add html files
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "./web/public")

	// Register default error view
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "Error occured"))
		ctx.View("shared/error.html")
	})

	db, err := database.ConnectSqlite()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	// I cant figure out how to have global DI, when using MVC pattern?
	index := mvc.New(app.Party("/"))
	// register db in dependecy injection container
	index.Register(db)

	index.Handle(new(controllers.IndexController))

	login := mvc.New(app.Party("/login"))
	// register db in dependecy injection container
	login.Register(db)
	login.Handle(new(controllers.LoginController))

	logout := mvc.New(app.Party("/logout"))
	logout.Register(db)
	logout.Handle(new(controllers.LogoutController))

	signup := mvc.New(app.Party("/signup"))
	signup.Register(db)
	signup.Handle(new(controllers.SignupController))

	app.Listen(":8080")
}
