package main

import (
	"os"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// func postMessage(e *httptest.Expect, text string) *httpexpect.Request {
// 	return e.POST("/add_message").WithFormField("text", text).WithHeader("Content-Type", "application/x-www-form-urlencoded")
// }

// func follow(e *httptest.Expect, username string) *httpexpect.Request {
// 	return e.GET("/" + username + "/follow")
// }

// func unfollow(e *httptest.Expect, username string) *httpexpect.Request {
// 	return e.GET("/" + username + "/unfollow")
// }

// func getTimeLine(e *httptest.Expect, username string) *httpexpect.Request {
// 	return e.GET("/" + username)
// }

// func getCookie(cookies []*http.Cookie, name string) *string {
// 	for _, c := range cookies {
// 		if c.Name == name {
// 			return &c.Value
// 		}
// 	}
// 	return nil
// }

func TestRegister(t *testing.T) {
	app := NewApp("development")

	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "user2",
		"password": "123123",
		"email":    "user2@gmail.com",
	}

	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusNoContent)
}

func TestSignin(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "user2",
		"password": "123123",
	}

	e.POST("/api/signin").WithJSON(form).Expect().Status(httptest.StatusOK).JSON().Object().Raw()
}

func TestRegisterAndSignin(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	form1 := map[string]interface{}{
		"username": "user2",
		"password": "123123",
		"email":    "user2@gmail.com",
	}

	e.POST("/api/register").WithJSON(form1).Expect().Status(httptest.StatusNoContent)

	form2 := map[string]interface{}{
		"username": "user2",
		"password": "123123",
	}

	e.POST("/api/signin").WithJSON(form2).Expect().Status(httptest.StatusNoContent)
}

func TestRegisterAndSigninAndPostMessage(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	// Register
	form := map[string]interface{}{
		"username": "user4",
		"password": "123123",
		"email":    "user4@gmail.com",
	}
	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusNoContent)

	//Sign in
	form2 := map[string]interface{}{
		"username": "user4",
		"password": "123123",
	}

	token := e.POST("/api/signin").WithJSON(form2).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	//Post message
	form3 := map[string]interface{}{
		"content": "Test post message",
	}

	e.POST("/api/tweets").WithJSON(form3).WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusNoContent)
}

func TestFollowAndUnfollow(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	userFollowingForm := map[string]interface{}{
		"username": "user15",
		"email":    "user15@email.com",
		"password": "123123",
	}

	userToFollowForm := map[string]interface{}{
		"username": "user16",
		"email":    "user16@email.com",
		"password": "123123",
	}

	userFollowingSigninForm := map[string]interface{}{
		"username": "user15",
		"password": "123123",
	}

	e.POST("/api/register").WithJSON(userFollowingForm).Expect().Status(httptest.StatusNoContent)
	e.POST("/api/register").WithJSON(userToFollowForm).Expect().Status(httptest.StatusNoContent)

	token := e.POST("/api/signin").WithJSON(userFollowingSigninForm).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	username := userToFollowForm["username"].(string)

	//Intially user does not follow
	e.GET("/api/users/"+username+"/isfollowing").WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusOK).JSON().Object().ValueEqual("isFollowing", false)

	//User follows
	e.POST("/api/users/"+username+"/follow").WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusOK)
	e.GET("/api/users/"+username+"/isfollowing").WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusOK).JSON().Object().ValueEqual("isFollowing", true)

	//User unfollows
	e.POST("/api/users/"+username+"/follow").WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusOK)
	e.GET("/api/users/"+username+"/isfollowing").WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusOK).JSON().Object().ValueEqual("isFollowing", false)
}
