package main

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/environment"
	"aske-w/itu-minitwit/services"
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
	sess := sessions.New(sessions.Config{
		Cookie:                      "itu-minitwit-cookie",
		AllowReclaim:                true,
		DisableSubdomainPersistence: true,
	})
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
	userService := services.NewUserService(db)
	authService := services.NewAuthService(db)
	timelineService := services.NewTimelineService(db)
	messageService := services.NewMessageService(db)

	// I cant figure out how to have global DI, when using MVC pattern?
	index := mvc.New(app.Party("/"))
	index.Register(db)
	index.Register(timelineService)
	index.Register(messageService)
	index.Register(userService)
	index.Handle(new(controllers.IndexController))

	Auth := mvc.New(app.Party("/"))
	Auth.Register(userService)
	Auth.Register(authService)
	Auth.Handle(new(controllers.AuthController))

	// api := mvc.New(app.Party("/api"))
	// api.Register(db)
	// api.Handle(new(controllers.ApiController))

	app.Listen(":8080")
}
