package controllers

import (
	"aske-w/itu-minitwit/database/sqlite"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type IndexController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// // Our UserService, it's an interface which
	// // is binded from the main application.
	// Service services.UserService

	// // Session, binded using dependency injection from the main.go.
	// Session *sessions.Session
}

type User struct {
	Id       int
	Username string
	Email    string
	Pw_hash  string
}

type Follower struct {
	Who_id  int
	Whom_id int
}

type Message struct {
	Message_id int
	Author_id  int
	Text       string
	Pub_date   int
	Flagged    int
}

func getMessages() []Message {

	db, err := sqlite.ConnectSqlite()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	rows, err := db.Conn.Query(`SELECT * FROM messages`)

	if err != nil {
		log.Fatalf("2: error selecting all messages: %v", err)
	}

	messages := make([]Message, 0)

	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.Message_id, &message.Author_id, &message.Text, &message.Pub_date, &message.Flagged)
		if err != nil {
			log.Fatalf("error scanning rows %v", err)
		}

		messages = append(messages, message)
	}

	return messages
}

func getUserByUsername(username string) (*User, error) {
	db, err := sqlite.ConnectSqlite()
	checkError(err)

	rows, err := db.Conn.Query(`select id, username, email, pw_hash from user where username = ?`, username)
	defer rows.Close()
	checkError(err)

	if rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Pw_hash)
		checkError(err)

		return &user, nil
	}

	return nil, errors.New("Can't find user with given username")
}

func getUserById(id int) (*User, error) {
	db, err := sqlite.ConnectSqlite()
	checkError(err)

	rows, err := db.Conn.Query(`select id, username, email, pw_hash from user where id = ?`, id)
	defer rows.Close()
	checkError(err)

	if rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.Pw_hash)
		checkError(err)

		return &user, nil
	}

	return nil, errors.New("Can't find user with given id")
}

//     """Format a timestamp for display."""
func format_datetime(timestamp int) string {
	unix := time.Unix(int64(timestamp), 0)
	return unix.Format("2006-01-02T15:04:05Z07:00")
}

//     """Return the gravatar image for the given email address."""
func gravatar_url(email string, size int) string {
	stripped := strings.Trim(email, "")
	lowered := strings.ToLower(stripped)
	valid := strings.ToValidUTF8(lowered, "")

	hasher := md5.New()
	data := []byte(valid)
	hasher.Write(data)
	md5Email := hex.EncodeToString(hasher.Sum(nil))

	url := fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", md5Email, size)
	return url
}

// https://docs.iris-go.com/iris/contents/sessions
// func beforeRequest(){}

func (c *IndexController) Get() mvc.Result {

	messages := getMessages()

	return mvc.View{
		Name: "index.html",
		Data: iris.Map{"Title": "Index page", "Messages": messages},
	}

}

func checkError(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}
