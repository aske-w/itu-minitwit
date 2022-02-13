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

	rows, err := c.DB.Conn.Query(" SELECT * FROM user INNER JOIN message ON message.author_id = user.user_id AND message.flagged = 0 ORDER BY message.pub_date DESC LIMIT ?", 10)

	utils.CheckError(err)
	defer rows.Close()

	timeline := make(Timelines, 0)

	for rows.Next() {
		group := &Timeline{
			Gravatar_Url:    gravatar_url,
			Format_Datetime: format_datetime,
		}
		err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
		utils.CheckError(err)
		timeline = append(timeline, group)
	}

	return timeline
}

func private_timeline(c *IndexController, userId string) []*Timeline {
	rows, err := c.DB.Conn.Query(`
	select  user.*, message.* from message, user
	where message.flagged = 0 and message.author_id = user.user_id and (
		user.user_id = ? or
		user.user_id in (select whom_id from follower
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

		err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
		utils.CheckError(err)
		timeline = append(timeline, group)
	}

	return timeline
}

func user_timeline(c *IndexController, userId int) []*Timeline {
	rows, err := c.DB.Conn.Query(`
	select  user.*, message.* from message, user where
	user.user_id = message.author_id and user.user_id = ?
	order by message.pub_date desc limit ?`, userId, 30)

	utils.CheckError(err)
	defer rows.Close()

	timeline := make(Timelines, 0)

	for rows.Next() {
		group := &Timeline{
			Gravatar_Url:    gravatar_url,
			Format_Datetime: format_datetime,
		}

		err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
		utils.CheckError(err)
		timeline = append(timeline, group)
	}

	return timeline
}

func (c *IndexController) UserId() string {
	return c.Session.GetString("user_id")
}
func (c *IndexController) User() (entity.User, error) {
	return utils.GetUserById(c.UserId(), c.DB, c.Ctx)
}
func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {

	b.Handle("GET", "/{username:string}", "UserTimelineHandler")
	b.Handle("GET", "/{username:string}/follow", "FollowHandler")
	b.Handle("GET", "/{username:string}/unfollow", "UnfollowHandler")
	b.Handle("POST", "/add_message", "AddMessageHandler")
}

func (c *IndexController) get_user_id(username string) string {
	var userId string
	c.DB.Get(c.Ctx, &userId, "select user_id from user where username = ?", username)
	return userId
}

func (c *IndexController) UnfollowHandler(username string) mvc.View {
	userId := c.UserId()
	whomId := c.get_user_id(username)
	if userId == "" || whomId == "" {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}

	c.DB.Exec(
		c.Ctx,
		"delete from follower where who_id=? and whom_id=?",
		userId, whomId,
	)
	c.Ctx.Redirect("/" + username)
	return mvc.View{}

}

func (c *IndexController) FollowHandler(username string) mvc.View {
	// """Adds the current user as follower of the given user."""
	userId := c.UserId()
	whomId := c.get_user_id(username)
	if userId == "" || whomId == "" {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}
	c.DB.Exec(
		c.Ctx,
		"insert into follower (who_id, whom_id) values (?, ?)",
		userId, whomId,
	)
	c.Ctx.Redirect("/" + username)
	return mvc.View{}
}

func (c *IndexController) AddMessageHandler() mvc.View {

	userId := c.UserId()
	if userId == "" {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}

	text := c.Ctx.FormValue("text")
	if text != "" {
		c.DB.Exec(
			c.Ctx,
			"insert into message (author_id, text, pub_date, flagged)	values (?, ?, ?, 0)",
			userId,
			text,
			time.Now().Unix(),
		)
	}
	c.Ctx.Redirect("/")
	return mvc.View{}
}

func (c *IndexController) UserTimelineHandler(username string) mvc.View {
	user, _ := c.User()
	profile_user, err := utils.GetUserByUsername(username, c.DB, c.Ctx)
	if err != nil {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}
	var followed bool
	c.DB.Get(c.Ctx, &followed, `
	select 1 from follower where
            follower.who_id = ? and follower.whom_id = ?
	`, user.User_id, profile_user.User_id)

	messages := user_timeline(c, profile_user.User_id)

	return mvc.View{
		Name: "timeline.html",
		Data: iris.Map{
			"Title":       profile_user.Username + "'s timeline",
			"User":        user,
			"LoggedIn":    user.User_id > 0,
			"Messages":    messages,
			"ProfileUser": profile_user,
			"Endpoint":    "user_timeline",
			"Followed":    followed,
		},
	}
}

func (c *IndexController) GetPublic() mvc.Result {
	user, _ := c.User()
	messages := public_timeline(c)
	return mvc.View{
		Name: "timeline.html",
		Data: iris.Map{
			"Title":    "Public timeline",
			"Messages": messages,
			"User":     user,
			"LoggedIn": user.User_id > 0,
			"Endpoint": "timeline",
		},
	}

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
		Name: "timeline.html",
		Data: iris.Map{
			"Title":    "My timeline",
			"Messages": messages,
			"User":     user,
			"LoggedIn": loggedIn,
			"Endpoint": "timeline",
		},
	}

}
