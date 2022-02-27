package controllers

import (
	"aske-w/itu-minitwit/models"
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gorm.io/gorm"
)

type ApiController struct {
	Ctx iris.Context

	TimelineService *services.TimelineService
	MessageService  *services.MessageService
	UserService     *services.UserService
	AuthService     *services.AuthService

	DB *gorm.DB
	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

var LATEST = 0

func update_latest(c *ApiController) int {
	urlParams := c.Ctx.Params()
	tryLatest := urlParams.GetIntDefault("latest", -1)
	if tryLatest != -1 {
		LATEST = tryLatest
	}
	return LATEST
}

type MyResponse struct {
	status    int    `json:"status"`
	error_msg string `json:"error_msg"`
}

type FollowRequest struct {
	Follow   *string `json:"follow`
	Unfollow *string `json:"unfollow`
}

type FilteredMsg struct {
	Content  string `json:"content"`
	Pub_date string `json:"pub_date"`
	User     string `json:"user"`
}

type Message struct {
	Content string `json:"Content"`
}

type FilteredMsgs []iris.Map

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"pwd"`
	Email    string `json:"email"`
}

func not_req_from_simulator(ctx iris.Context) bool {
	auth := ctx.GetHeader("Authorization")

	if strings.Compare(auth, "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh") != 0 {
		ctx.StatusCode(403)
		ctx.JSON(iris.Map{"status": 403, "error_msg": "You are not authorized to use this resource!"})
		return false

	}

	return true

}

func getFilteredMsgs(rows *sql.Rows) FilteredMsgs {
	filtered_msgs := make(FilteredMsgs, 0)

	for rows.Next() {
		msg := &FilteredMsg{}
		err := rows.Scan(&msg.User, &msg.Content, &msg.Pub_date)
		utils.CheckError(err)
		temp := iris.Map{"user": msg.User, "content": msg.Content, "pub_date": msg.Pub_date}
		filtered_msgs = append(filtered_msgs, temp)
	}

	return filtered_msgs
}

func (c *ApiController) BeforeActivation(b mvc.BeforeActivation) {
	// Register endpoints on /api
	b.Handle("POST", "/register", "RegisterHandler") // Done
	b.Handle("GET", "/latest", "LatestHandler")      // Done √

	b.Handle("GET", "/msgs", "MsgHandler")                             // Done √
	b.Handle("GET", "/msgs/{username:string}", "UserMsgsGetHandler")   // Done √
	b.Handle("POST", "/msgs/{username:string}", "UserMsgsPostHandler") // Done √

	b.Handle("GET", "/fllws/{username:string}", "FollowersGetHandler") // Done
	b.Handle("POST", "/fllws/{username:string}", "FollowersPostHandler")

}

func (c *ApiController) RegisterHandler() {

	update_latest(c)

	registerUser := RegisterUser{}
	readBody(c, &registerUser)
	username := registerUser.Username
	email := registerUser.Email
	password := registerUser.Password

	var err error

	if username == "" {
		err = fmt.Errorf("you have to enter a username")
	} else if email == "" || !strings.Contains(email, "@") {
		err = fmt.Errorf("you have to enter a valid email address")
	} else if password == "" {
		err = fmt.Errorf("you have to enter a password")
	} else {

		exists, _ := c.UserService.CheckUsernameExists(username)

		if exists {
			err = fmt.Errorf("the username is already taken")
		} else {
			_, err := c.AuthService.CreateUser(username, email, password)
			if err == nil {

				c.Ctx.StatusCode(204)
				return
			}

		}

		// user, _ := c.UserService.FindByUsername(username)
		// count := c.UserService.CountUsers()
		// if user.ID != 0 && count > 0 {
		// 	err = fmt.Errorf("the username is already taken")
		// } else {
		// 	var byteHash []byte
		// 	byteHash, err = bcrypt.GenerateFromPassword([]byte(password), 10)
		// 	if err == nil {
		// 		err = c.DB.Exec(`insert into users (username, email, pw_hash) values (?,?,?)`, username, email, string(byteHash)).Error
		// 		if err == nil {
		// 			c.Ctx.StatusCode(204)
		// 			return
		// 		}
		// 	}
		// }
	}
	c.Ctx.StatusCode(400)
	c.Ctx.JSON(iris.Map{"status": 400, "error_msg": err})
}

func (c *ApiController) LatestHandler() {
	c.Ctx.JSON(iris.Map{"latest": LATEST})
}

func (c *ApiController) MsgHandler() {
	validToken := not_req_from_simulator(c.Ctx)

	if !validToken {
		return
	}

	no_msg := c.Ctx.Params().GetIntDefault("no", 100)
	msgs := []FilteredMsg{}

	c.DB.Model(&models.User{}).Select("users.username as user", "messages.text as content", "messages.pub_date").Joins(
		"INNER JOIN messages ON messages.author_id = users.id AND messages.flagged = 0",
	).Order("messages.pub_date DESC").Limit(no_msg).Scan(&msgs)

	c.Ctx.JSON(msgs)
}

func (c *ApiController) UserMsgsGetHandler(username string) {
	validToken := not_req_from_simulator(c.Ctx)

	if !validToken {
		return
	}

	no_msg := c.Ctx.Params().GetIntDefault("no", 100)
	profile_user_id, _ := c.UserService.UsernameToId(username)

	if profile_user_id == -1 {
		c.Ctx.StatusCode(404)
		return
	}

	msgs := []FilteredMsg{}

	c.DB.Table("messages, users").Select("users.username as User", "messages.text as Content", "messages.pub_date as Pub_date").Where(
		"messages.flagged = 0 AND users.id = messages.author_id AND users.id = ?", profile_user_id,
	).Order("messages.pub_date DESC").Limit(no_msg).Scan(&msgs)

	c.Ctx.JSON(msgs)
}

func (c *ApiController) UserMsgsPostHandler(username string) {
	update_latest(c)
	validToken := not_req_from_simulator(c.Ctx)

	if !validToken {
		return
	}

	userId, _ := c.UserService.UsernameToId(username)
	if userId == -1 {
		c.Ctx.StatusCode(404)
		return
	}

	msg := Message{}

	readBody(c, &msg)
	text := msg.Content
	if text != "" {
		c.MessageService.CreateMessage(userId, text)
	}
	c.Ctx.StatusCode(204)
}

func readBody(c *ApiController, v interface{}) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Ctx.Request().Body)
	json.Unmarshal(buf.Bytes(), &v)
}

func (c *ApiController) FollowersGetHandler(username string) {
	update_latest(c)

	validToken := not_req_from_simulator(c.Ctx)
	if !validToken {
		return
	}

	num_followers := c.Ctx.Params().GetIntDefault("no", 100)

	follower_names := c.UserService.GetFollowersByUsername(username, num_followers)

	c.Ctx.StatusCode(200)
	c.Ctx.JSON(iris.Map{"follows": follower_names})
}

func (c *ApiController) FollowersPostHandler(username string) {
	validToken := not_req_from_simulator(c.Ctx)
	if !validToken {
		return
	}

	body := FollowRequest{}
	readBody(c, &body)

	userId, _ := c.UserService.UsernameToId(username)

	if body.Follow != nil && body.Unfollow == nil {
		// follow
		followerId, _ := c.UserService.UsernameToId(*body.Follow)
		c.UserService.FollowUser(userId, followerId)
	} else if body.Follow == nil && body.Unfollow != nil {
		// un follow
		followerId, _ := c.UserService.UsernameToId(*body.Unfollow)
		c.UserService.FollowUser(userId, followerId)
	}

	c.Ctx.StatusCode(204)

}
