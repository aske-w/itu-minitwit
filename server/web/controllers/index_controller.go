package controllers

import (
	"aske-w/itu-minitwit/models"
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/utils"
	"fmt"

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

func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {

	b.Handle("POST", "/add_message", "AddMessageHandler")
	b.Handle("GET", "/{username:string}", "UserTimelineHandler")
	b.Handle("GET", "/{username:string}/follow", "FollowHandler")
	b.Handle("GET", "/{username:string}/unfollow", "UnfollowHandler")
}

func (c *IndexController) UnfollowHandler(username string) mvc.View {
	userId := c.getUserId()
	follower, err := c.UserService.FindByUsername(username)
	if err != nil {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}

	c.UserService.UnfollowUser(userId, int(follower.ID))
	c.Ctx.Redirect("/" + username)
	return mvc.View{}

}

func (c *IndexController) FollowHandler(username string) mvc.View {
	// """Adds the current user as follower of the given user."""
	userId := c.getUserId()
	whom, err := c.UserService.FindByUsername(username)

	if userId == -1 || err != nil {
		return mvc.View{
			Data: iris.Map{"Message": "User not found"},
			Code: 404,
		}
	}
	c.UserService.FollowUser(userId, int(whom.ID))

	c.Ctx.Redirect("/" + username)
	return mvc.View{}
}

type AddMessageForm struct {
	Text string `form:"text"` // or just "colors".
}

func (c *IndexController) AddMessageHandler() mvc.Result {
	fmt.Println("add message handler")
	user := c.getUser()
	if user == nil {
		fmt.Println("user nil")
		return c.errorPage("You need to be logged in")
	}

	var form AddMessageForm
	c.Ctx.ReadForm(&form)
	fmt.Println("text", form)
	if form.Text != "" {
		c.MessageService.CreateMessage(int(user.ID), form.Text)
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
	user := c.getUser()
	messages, _ := c.TimelineService.GetPrivateTimeline(int(user.ID))

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

// https://github.com/kataras/iris/issues/1704#issuecomment-761177806
// fix add message error
func (c *IndexController) HandleError(ctx iris.Context, err error) {
	if iris.IsErrPath(err) {
		// to ignore any "schema: invalid path" you can check the error type
		// and don't stop the execution.
		return // continue.
	}

	ctx.StopExecution()
}
