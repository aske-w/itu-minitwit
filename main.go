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

func main() {
	// initDatabase()
	// connect to db

	// db, err := sqlite.ConnectSqlite()
	// if err != nil {
	// 	log.Fatalf("error connecting to the MySQL database: %v", err)
	// }
	// query := `CREATE TABLE IF NOT EXISTS product(product_id int primary key auto_increment, product_name text,
	//     product_price int, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)`
	// println(db.Conn.Exec(query))

	app := iris.New()
	// app.Logger().SetLevel("debug") // more logging

	app.Use(logger.New())  // logs request
	app.Use(recover.New()) // handles panics (shows 404)

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

	index := mvc.New(app.Party("/"))

	db, err := database.ConnectSqlite()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	index.Register(db)
	index.Handle(new(controllers.IndexController))

	app.Listen(":8080")
}
