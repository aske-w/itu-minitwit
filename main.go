package main

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/web/controllers"
	"fmt"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
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

func initDatabase() {
	fmt.Println("INIT DATABASE")

	db, err := database.ConnectSqlite()
	if err != nil {
		log.Fatalf("No database found: %v", err)
	}

	// db.Conn.Exec(`
	// 	create table users (
	// 		id integer primary key autoincrement,
	// 		username string not null,
	// 		email string not null,
	// 		pw_hash string not null
	// 	);
	// `)

	// db.Conn.Exec(`
	// 	create table followers (
	// 		who_id integer,
	// 		whom_id integer
	//   	);
	// `)

	// db.Conn.Exec(`
	// 	create table messages (
	// 		message_id integer primary key autoincrement,
	// 		author_id integer not null,
	// 		text string not null,
	// 		pub_date integer,
	// 		flagged integer
	//   	);
	// `)

	db.Conn.Exec(`INSERT INTO users (username, email, pw_hash) values ('christian', 'cger@itu.dk', 'secret')`)

	db.Conn.Exec(`INSERT INTO messages (author_id, text, pub_date, flagged) values (1, 'Tweet 1', 1000, false)`)
	db.Conn.Exec(`INSERT INTO messages (author_id, text, pub_date, flagged) values (1, 'Tweet 1', 1000, false)`)
	// db.Conn.Exec(`INSERT INTO products (product_id, product_name, product_price) values (2, 'Havestol', 1000)`)
}

func before(ctx iris.Context) {
	shareInformation := "this is a sharable information between handlers"

	requestPath := ctx.Path()
	println("Before the indexHandler or contactHandler: " + requestPath)

	// if ctx.Session {
	// }

	// ctx.SetUser()
	ctx.Values().Set("info", shareInformation)
	ctx.Next()
}

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug") // more logging

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

	app.UseGlobal(before)
	app.Listen(":8080")
}
