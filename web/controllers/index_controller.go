package controllers

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/entity"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type IndexController struct {
	Ctx iris.Context

	DB *database.SQLite
	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
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

type Timeline struct {
	UserId          int
	Username        string
	Email           string
	Pw_hash         string
	Message_id      int
	Author_id       int
	Text            string
	Pub_date        int
	Flagged         int
	Gravatar_Url    func(email string, size int) string
	Format_Datetime func(timestamp int) string
}

type Timelines []*Timeline

// type Message struct {
// 	Message_id int
// 	Author_id  int
// 	Text       string
// 	Pub_date   int
// 	Flagged    int
// }

func getMessages() []entity.Message {

	db, err := database.ConnectSqlite()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	rows, err := db.Conn.Query(`SELECT * FROM message`)

	if err != nil {
		log.Fatalf("2: error selecting all messages: %v", err)
	}

	messages := make([]entity.Message, 0)

	for rows.Next() {
		message := entity.Message{}
		err = rows.Scan(&message.Message_id, &message.Author_id, &message.Text, &message.Pub_date, &message.Flagged)
		if err != nil {
			log.Fatalf("error scanning rows %v", err)
		}

		messages = append(messages, message)
	}

	return messages
}

func getUserByUsername(username string) (*User, error) {
	db, err := database.ConnectSqlite()
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
	db, err := database.ConnectSqlite()
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

func public_timeline(c *IndexController) []*Timeline {
	// var timeline Timeline

	rows, err := c.DB.Conn.Query(" SELECT * FROM user INNER JOIN message ON message.author_id = user.user_id AND message.flagged = 0 ORDER BY message.pub_date DESC LIMIT ?", 10)
	checkError(err)
	defer rows.Close()

	timeline := make(Timelines, 0)

	for rows.Next() {
		group := &Timeline{
			Gravatar_Url:    gravatar_url,
			Format_Datetime: format_datetime,
		}
		// user := entity.User{}
		err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
		checkError(err)
		// group.Gravatar = gravatar_url
		timeline = append(timeline, group)
	}

	return timeline
}

// https://docs.iris-go.com/iris/contents/sessions
// func beforeRequest(){}

func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {

}

func (c *IndexController) Get() mvc.Result {

	// var messages entity.Messages

	// err := c.DB.Select(c.Ctx, &messages, "SELECT * from message desc limit ?", 10)
	// checkError(err)

	// for i := 0; i < len(messages); i++ {
	// 	fmt.Println(messages[i])

	// }
	messages := public_timeline(c)

	// messages := getMessages()

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
