package controllers

import (
	"aske-w/itu-minitwit/models"
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gorm.io/gorm"
)

type IndexController struct {
	Ctx iris.Context

	DB *gorm.DB

	TimelineService *services.TimelineService
	MessageService  *services.MessageService
	UserService     *services.UserService
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
	timeline := make(Timelines, 0)
	c.DB.Model(&models.User{}).Find("id as userId,Username,Email,Pw_hash,Message_id,Author_id,Text,Pub_date,Flagged").Joins("INNER JOIN message ON message.author_id = user.user_id AND message.flagged = 0").Order("message.pub_date DESC").Limit(30).Scan(&timeline)
	// rows, err := c.DB.db.Query(" SELECT * FROM user INNER JOIN message ON message.author_id = user.user_id AND message.flagged = 0 ORDER BY message.pub_date DESC LIMIT ?", 30)

	// utils.CheckError(err)
	// defer rows.Close()

	// for rows.Next() {
	// 	group := &Timeline{
	// 		Gravatar_Url:    gravatar_url,
	// 		Format_Datetime: format_datetime,
	// 	}
	// 	err = rows.Scan(&group.UserId, &group.Username, &group.Email, &group.Pw_hash, &group.Message_id, &group.Author_id, &group.Text, &group.Pub_date, &group.Flagged)
	// 	utils.CheckError(err)

	// 	timeline = append(timeline, group)
	// }

	return timeline
}

func private_timeline(c *IndexController, userId string) []*Timeline {
	rows, err := c.DB.Raw(`
	select  user.*, message.* from message, user
	where message.flagged = 0 and message.author_id = user.user_id and (
		user.user_id = ? or
		user.user_id in (select whom_id from follower
								where who_id = ?))
	order by message.pub_date desc limit ?`, userId, userId, 10).Rows()

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

func user_timeline(c *IndexController, userId uint) []*Timeline {
	rows, err := c.DB.Raw(`
	select  user.*, message.* from message, user where
	user.user_id = message.author_id and user.user_id = ?
	order by message.pub_date desc limit ?`, userId, 30).Rows()

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
func (c *IndexController) User() (*models.User, error) {
	var userId = c.UserId()
	user := &models.User{}

	c.DB.First(user, userId)
	return user, nil
}
func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {

	b.Handle("GET", "/{username:string}", "UserTimelineHandler")
	b.Handle("GET", "/{username:string}/follow", "FollowHandler")
	b.Handle("GET", "/{username:string}/unfollow", "UnfollowHandler")
	b.Handle("POST", "/add_message", "AddMessageHandler")
}

func (c *IndexController) get_user_id(username string) string {

	user := &models.User{}
	// c.DB.Get(c.Ctx, &userId, "select user_id from user where username = ?", username)
	c.DB.First(user).Where("username = ?", username)
	return string(user.ID)
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
		"delete from follower where who_id=? and whom_id=?",
		userId, whomId,
	)
	c.Ctx.Redirect("/" + username)
	return mvc.View{}

}

func (c *IndexController) FollowHandler(username string) mvc.View {
	// """Adds the current user as follower of the given user."""
	userId := c.getUserId()
	whom, err := c.UserService.FindByUsername(username)
	fmt.Println(userId, whom, err)
	if userId != -1 || err != nil {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}
	c.UserService.FollowUser(userId, int(whom.ID))

	c.Ctx.Redirect("/" + username)
	return mvc.View{}
}

func (c *IndexController) AddMessageHandler() mvc.Result {
	user := c.getUser()
	if user == nil {
		return c.errorPage("You need to be logged in")
	}

	text := c.Ctx.FormValue("text")
	if text != "" {
		c.MessageService.CreateMessage(int(user.ID), text)
	}
	c.Ctx.Redirect("/")
	return mvc.View{}
}

func (c *IndexController) UserTimelineHandler(username string) mvc.View {

	profile_user, err := c.UserService.FindByUsername(username)

	if err != nil {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}
	var followed bool
	if c.isLoggedIn() {

		followed = c.UserService.UserIsFollowing(c.getUserId(), int(profile_user.ID))

	}
	messages, err := c.TimelineService.GetUserTimeline(int(profile_user.ID))

	utils.CheckError(err)
	return mvc.View{
		Name: "timeline.html",
		Data: iris.Map{
			"Title":       profile_user.Username + "'s timeline",
			"User":        c.getUser(),
			"LoggedIn":    c.isLoggedIn(),
			"Messages":    messages,
			"ProfileUser": profile_user,
			"Endpoint":    "user_timeline",
			"Followed":    followed,
		},
	}
}

func (c *IndexController) GetPublic() mvc.Result {
	messages, err := c.TimelineService.GetPublicTimeline()
	if err != nil {
		return c.errorPage(err.Error())
	}
	return mvc.View{
		Name: "timeline.html",
		Data: iris.Map{
			"Title":    "Public timeline",
			"Messages": messages,
			"User":     c.getUser(),
			"LoggedIn": c.isLoggedIn(),
			"Endpoint": "timeline",
		},
	}

}

func (c *IndexController) Get() mvc.Result {

	loggedIn := c.isLoggedIn()
	if !loggedIn {
		c.Ctx.Redirect("/public")
		return mvc.View{}
	}
	var messages []*Timeline
	// messages = []Timeline{} //private_timeline(c, userId)

	user := c.getUser()

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

func (c *IndexController) errorPage(message string) mvc.Result {
	return mvc.View{
		Data: iris.Map{"Message": message},
		Code: 404,
	}
}

/*
	Returns -1 if not found
*/
func (c *IndexController) getUserId() int {
	userId, _ := utils.GetUserIdFromSession(c.Session)
	return userId
}

func (c *IndexController) isLoggedIn() bool {
	_, loggedIn := utils.GetUserIdFromSession(c.Session)
	return loggedIn
}

func (c *IndexController) getUser() *models.User {
	userId, loggedIn := utils.GetUserIdFromSession(c.Session)
	if !loggedIn {
		return nil
	}
	user, err := c.UserService.GetById(userId)
	utils.CheckError(err)
	return user
}
