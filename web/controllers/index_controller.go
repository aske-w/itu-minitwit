package controllers

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/entity"
	"aske-w/itu-minitwit/web/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

	rows, err := c.DB.Conn.Query(" SELECT * FROM users INNER JOIN message ON message.author_id = users.id AND message.flagged = 0 ORDER BY message.pub_date DESC LIMIT ?", 10)
	utils.CheckError(err)
	defer rows.Close()

	timeline := make(Timelines, 0)

	for rows.Next() {
		group := &Timeline{
			Gravatar_Url:    gravatar_url,
			Format_Datetime: format_datetime,
		}
		// user := entity.User{}
		err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
		utils.CheckError(err)
		// group.Gravatar = gravatar_url
		timeline = append(timeline, group)
	}

	return timeline
}

func (c *IndexController) GetPublic() mvc.Result {
	messages := public_timeline(c)
	return mvc.View{
		Name: "index.html",
		Data: iris.Map{"Title": "Index page", "Messages": messages, "User": nil, "LoggedIn": false},
	}

}

func private_timeline(c *IndexController, userId string) []*Timeline {
	rows, err := c.DB.Conn.Query(`
	select message.*, users.* from message, users
	where message.flagged = 0 and message.author_id = users.id and (
		users.id = ? or
		users.id in (select whom_id from follower
								where who_id = ?))
	order by message.pub_date desc limit ?`, userId, userId, 10)
	utils.CheckError(err)
	defer rows.Close()

	timeline := make(Timelines, 0)

	for rows.Next() {
		group := &Timeline{
			Gravatar_Url:    gravatar_url,
			Format_Datetime: format_datetime,
		}
		// user := entity.User{}
		err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
		utils.CheckError(err)
		// group.Gravatar = gravatar_url
		timeline = append(timeline, group)
	}

	return timeline
}

func (c *IndexController) Get() mvc.Result {

	var messages []*Timeline

	userId, loggedIn := utils.GetUserIdFromSession(c.Session)
	var user entity.User
	if loggedIn {
		_user, err := utils.GetUserById(userId, c.DB, c.Ctx)
		utils.CheckError(err)
		messages = private_timeline(c, userId)
		user = _user
	} else {
		c.Ctx.Redirect("/public")
	}

	return mvc.View{
		Name: "index.html",
		Data: iris.Map{"Title": "Index page", "Messages": messages, "User": user, "LoggedIn": loggedIn},
	}

}
