package main

import (
	"aske-w/itu-minitwit/web/controllers"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type (
	request struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	response struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug") // more logging

	// Add html files
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "./web/public")
	// app.HandleDir("/public", iris.Dir("./web/public"))

	// Register default error view
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "Error occured"))
		ctx.View("shared/error.html")
	})

	index := mvc.New(app.Party("/"))
	index.Handle(new(controllers.IndexController))

	user := mvc.New(app.Party("/user"))
	user.Handle(new(controllers.UserController))

	// app.Handle("GET", "/", indexHandler)
	app.Listen(":8080")
}

func indexHandler(ctx iris.Context) {
	resp := response{
		ID:      "1",
		Message: "updated successfully",
	}
	ctx.JSON(resp)
}
