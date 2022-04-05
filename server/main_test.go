package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12/httptest"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// type registerForm struct {
// 	username string `form:"username"`
// 	pwd      string `form:"pwd"`
// 	email    string `form:"email"`
// }

func login(e *httptest.Expect, username, password string) *httpexpect.Request {
	return e.POST("/login").WithFormField("username", username).WithFormField("password", password).WithHeader("Content-Type", "application/x-www-form-urlencoded")
}

func logout(e *httptest.Expect) *httpexpect.Request {
	return e.GET("/logout")
}

func postMessage(e *httptest.Expect, text string) *httpexpect.Request {
	return e.POST("/add_message").WithFormField("text", text).WithHeader("Content-Type", "application/x-www-form-urlencoded")
}

func follow(e *httptest.Expect, username string) *httpexpect.Request {
	return e.GET("/" + username + "/follow")
}

func unfollow(e *httptest.Expect, username string) *httpexpect.Request {
	return e.GET("/" + username + "/unfollow")
}

func getTimeLine(e *httptest.Expect, username string) *httpexpect.Request {
	return e.GET("/" + username)
}

func getCookie(cookies []*http.Cookie, name string) *string {
	for _, c := range cookies {
		if c.Name == name {
			return &c.Value
		}
	}
	return nil
}

func TestRegister(t *testing.T) {
	app := NewApp("development")

	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "user2",
		"pwd":      "123123",
		"email":    "user2@gmail.com",
	}

	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusNoContent)
}

func TestSignin(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "user2",
		"pwd":      "123123",
	}

	e.POST("/api/signin").WithJSON(form).Expect().Status(httptest.StatusNoContent)
	// TODO change the contains method to something more specific/unique
	// Doesnt take Sessions into account, optimally it should check for "My timeline"
	// login := login(e, "user1", "123").Expect().Status(httptest.StatusOK).Body().Contains("My timeline")
	// login(e, "user1", "123").Expect().Status(httptest.StatusOK)
}

// func TestLogout(t *testing.T) {
// 	app := NewApp("development")
// 	e := httptest.New(t, app)

// 	// logout(e).Expect().Status(httptest.StatusOK).Body().Contains("Public timeline")
// 	logout(e).Expect().Status(httptest.StatusOK)
// }

// func TestSignupAndLogin(t *testing.T) {
// 	app := NewApp("development")
// 	e := httptest.New(t, app)

// 	register(e, "user1", "user1@example.com", "123", "123").Expect().Status(httptest.StatusOK)
// 	login(e, "user1", "123").Expect().Status(httptest.StatusOK)
// }

// func TestSignupAndLoginAndPostMessage(t *testing.T) {
// 	app := NewApp("development")
// 	e := httptest.New(t, app)

// 	// db, err := database.ConnectMySql("development")
// 	// if err != nil {
// 	// 	log.Fatalf("error connecting to the database: %v", err)
// 	// }

// 	register(e, "user1", "user1@example.com", "123", "123").Expect().Status(httptest.StatusOK)
// 	login(e, "user1", "123").Expect().Status(httptest.StatusOK)
// 	postMessage(e, "Should be the same").Expect().Status(httptest.StatusOK)

// 	// messageService := services.NewMessageService(db)
// 	// message, err := messageService.CreateMessage(1, "test")
// 	// assert.Equal(t, message.Author_id, 1, "Should be the same")
// 	// assert.Equal(t, message.Text, "test", "Should be the same")
// }

// func TestFollowAndUnfollow(t *testing.T) {
// 	app := NewApp("development")
// 	e := httptest.New(t, app)

// 	register(e, "user1", "user1@example.com", "123", "123").Expect().Status(httptest.StatusOK)
// 	register(e, "user2", "user2@example.com", "123", "123").Expect().Status(httptest.StatusOK)
// 	login(e, "user1", "123").Expect().Status(httptest.StatusOK)
// 	follow(e, "user2").Expect().Status(httptest.StatusOK)
// 	unfollow(e, "user2").Expect().Status(httptest.StatusOK)
// }

// // func AddMessage(t *testing.T) {
// // 	app := NewApp("development")
// // 	e := httptest.New(t, app)

// // }

// // func TestMessageRecording(t *testing.T) {

// // }

// // func TestTimelines(t *testing.T) {

// // }
