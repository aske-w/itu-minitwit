package controllers

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/web/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"golang.org/x/crypto/bcrypt"
)

type ApiController struct {
	Ctx iris.Context

	DB *database.SQLite
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
	content  string `json:"content"`
	pub_date string `json:"pub_date"`
	user     string `json:"user"`
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
		err := rows.Scan(&msg.user, &msg.content, &msg.pub_date)
		utils.CheckError(err)
		temp := iris.Map{"user": msg.user, "content": msg.content, "pub_date": msg.pub_date}
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
		user, _ := utils.GetUserByUsername(username, c.DB, c.Ctx)
		count := utils.CountEntries("user", c.DB)
		if user.User_id != 0 && count > 0 {
			err = fmt.Errorf("the username is already taken")
		} else {
			var byteHash []byte
			byteHash, err = bcrypt.GenerateFromPassword([]byte(password), 10)
			if err == nil {
				_, err = c.DB.db.Exec(`insert into user (username, email, pw_hash) values (?,?,?)`, username, email, string(byteHash))
				if err == nil {
					c.Ctx.StatusCode(204)
					return
				}
			}
		}
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

	rows, err := c.DB.db.Query("SELECT user.username, message.text, message.pub_date  FROM user INNER JOIN message ON message.author_id = user.user_id AND message.flagged = 0 ORDER BY message.pub_date DESC LIMIT ?", no_msg)
	utils.CheckError(err)
	defer rows.Close()

	msgs := getFilteredMsgs(rows)

	// for msg := range msgs {
	// 	fmt.Println("msg:" + msgs[msg].user)
	// }

	c.Ctx.JSON(msgs)
}

func (c *ApiController) UserMsgsGetHandler(username string) {
	validToken := not_req_from_simulator(c.Ctx)

	if !validToken {
		return
	}

	no_msg := c.Ctx.Params().GetIntDefault("no", 100)
	profile_user, err := utils.GetUserByUsername(username, c.DB, c.Ctx)

	if err != nil {
		c.Ctx.StatusCode(404)
		return
	}

	query := `SELECT user.username, message.text, message.pub_date FROM message, user WHERE message.flagged = 0 AND user.user_id = message.author_id AND user.user_id = ? ORDER BY message.pub_date DESC LIMIT ?`

	rows, err := c.DB.db.Query(query, profile_user.User_id, no_msg)
	utils.CheckError(err)
	defer rows.Close()

	c.Ctx.JSON(getFilteredMsgs(rows))
}

func (c *ApiController) UserMsgsPostHandler(username string) {
	update_latest(c)
	validToken := not_req_from_simulator(c.Ctx)

	if !validToken {
		return
	}

	user, err := utils.GetUserByUsername(username, c.DB, c.Ctx)
	if err != nil {
		c.Ctx.StatusCode(404)
		return
	}

	userId := user.User_id
	msg := Message{}

	readBody(c, &msg)
	text := msg.Content
	if text != "" {
		_, err := c.DB.Exec(
			c.Ctx,
			"insert into message (author_id, text, pub_date, flagged)	values (?, ?, ?, 0)",
			userId,
			text,
			time.Now().Unix(),
		)

		utils.CheckError(err)
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

	user, err := utils.GetUserByUsername(username, c.DB, c.Ctx)

	if err != nil {
		c.Ctx.StatusCode(404)
		return
	}

	no_followers := c.Ctx.Params().GetIntDefault("no", 100)
	rows, err := c.DB.db.Query("SELECT user.username FROM user INNER JOIN follower ON follower.whom_id=user.user_id WHERE follower.who_id = ? LIMIT ?", user.User_id, no_followers)
	utils.CheckError(err)
	defer rows.Close()

	follower_names := make([]string, 0)

	for rows.Next() {
		var follower_name string
		err := rows.Scan(&follower_name)

		utils.CheckError(err)

		follower_names = append(follower_names, follower_name)
	}

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

	myUser, err := utils.GetUserByUsername(username, c.DB, c.Ctx)

	if err != nil {
		c.Ctx.StatusCode(404)
		return
	}

	if body.Follow != nil && body.Unfollow == nil {

		whomUser, err := utils.GetUserByUsername(*body.Follow, c.DB, c.Ctx)

		if err != nil {
			c.Ctx.StatusCode(404)
			return
		}

		c.DB.Exec(
			c.Ctx,
			"insert into follower (who_id, whom_id) values (?, ?)",
			myUser.User_id, whomUser.User_id,
		)

		c.Ctx.StatusCode(204)
		return
	} else if body.Follow == nil && body.Unfollow != nil {

		whomUser, err := utils.GetUserByUsername(*body.Unfollow, c.DB, c.Ctx)

		if err != nil {
			c.Ctx.StatusCode(404)
			return
		}

		c.DB.Exec(
			c.Ctx,
			"delete from follower where who_id=? and whom_id=?",
			myUser.User_id, whomUser.User_id,
		)

		c.Ctx.StatusCode(204)
	} else {
		panic("Something is very wrong")
	}

}
