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

func TestRegister(t *testing.T) {
	app := NewApp("development")

	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "user",
		"password": "123123",
		"email":    "user@gmail.com",
	}

	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusNoContent)
}

func TestRegisterWithExistingCredentials(t *testing.T) {
	app := NewApp("development")

	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "user",
		"password": "123123",
		"email":    "user@gmail.com",
	}

	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusNoContent)
	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusBadRequest)
}

func TestSigninWithNonexistingCredentials(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	form := map[string]interface{}{
		"username": "nonexistinguser",
		"password": "123123",
	}

	e.POST("/api/signin").WithJSON(form).Expect().Status(httptest.StatusBadRequest).JSON().Object().Raw()
}

func TestRegisterAndSignin(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	form1 := map[string]interface{}{
		"username": "user",
		"password": "123123",
		"email":    "user@gmail.com",
	}

	e.POST("/api/register").WithJSON(form1).Expect().Status(httptest.StatusNoContent)

	form2 := map[string]interface{}{
		"username": "user",
		"password": "123123",
	}

	e.POST("/api/signin").WithJSON(form2).Expect().Status(httptest.StatusOK)
}

func TestRegisterAndSigninAndPostMessage(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	// Register
	form := map[string]interface{}{
		"username": "user",
		"password": "123123",
		"email":    "user@gmail.com",
	}
	e.POST("/api/register").WithJSON(form).Expect().Status(httptest.StatusNoContent)

	//Sign in
	form2 := map[string]interface{}{
		"username": "user",
		"password": "123123",
	}

	token := e.POST("/api/signin").WithJSON(form2).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	//Post message
	form3 := map[string]interface{}{
		"content": "Test post message",
	}

	e.POST("/api/tweets").WithJSON(form3).WithHeader("Authorization", "Bearer "+token).Expect().Status(httptest.StatusNoContent)
}

func TestRegisterAndSigninAndPostMessageWithWrongBearerToken(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	//Invalid token
	invalidToken := "invalidToken"

	//Post message
	form3 := map[string]interface{}{
		"content": "Test post message",
	}

	e.POST("/api/tweets").WithJSON(form3).WithHeader("Authorization", "Bearer "+invalidToken).Expect().Status(httptest.StatusUnauthorized)
}

func TestFollowAndUnfollow(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	userFollowingForm := map[string]interface{}{
		"username": "user1",
		"email":    "user1@email.com",
		"password": "123123",
	}

	userToFollowForm := map[string]interface{}{
		"username": "user2",
		"email":    "user2@email.com",
		"password": "123123",
	}

	userFollowingSigninForm := map[string]interface{}{
		"username": userFollowingForm["username"],
		"password": userFollowingForm["password"],
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

//Public timeline
func TestPublicTimeline(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	//Register user1
	user1 := map[string]interface{}{
		"username": "user1",
		"email":    "user1@email.com",
		"password": "123123",
	}
	e.POST("/api/register").WithJSON(user1).Expect().Status(httptest.StatusNoContent)
	//Register user2
	user2 := map[string]interface{}{
		"username": "user2",
		"email":    "user2@email.com",
		"password": "123123",
	}
	e.POST("/api/register").WithJSON(user2).Expect().Status(httptest.StatusNoContent)

	//User1 posts message1
	user1Token := e.POST("/api/signin").WithJSON(user1).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	message1 := map[string]interface{}{
		"content": "Message1",
	}
	e.POST("/api/tweets").WithJSON(message1).WithHeader("Authorization", "Bearer "+user1Token).Expect().Status(httptest.StatusNoContent)

	//User7 posts message2
	user2Token := e.POST("/api/signin").WithJSON(user2).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	message2 := map[string]interface{}{
		"content": "Message2",
	}
	e.POST("/api/tweets").WithJSON(message2).WithHeader("Authorization", "Bearer "+user2Token).Expect().Status(httptest.StatusNoContent)

	//User6 posts message3
	message3 := map[string]interface{}{
		"content": "Message3",
	}
	e.POST("/api/tweets").WithJSON(message3).WithHeader("Authorization", "Bearer "+user1Token).Expect().Status(httptest.StatusNoContent)

	//Check if messages are in public timeline
	tweets := e.GET("/api/tweets").Expect().JSON().Array()

	// Check ordering
	tweets.Element(0).Object().ValueEqual("Text", message3["content"])
	tweets.Element(1).Object().ValueEqual("Text", message2["content"])
	tweets.Element(2).Object().ValueEqual("Text", message1["content"])
}

//Private timeline
func TestPrivateTimeline(t *testing.T) {
	app := NewApp("development")
	e := httptest.New(t, app)

	// Register user 1
	user1 := map[string]interface{}{
		"username": "user1",
		"email":    "user1@email.com",
		"password": "123123",
	}

	e.POST("/api/register").WithJSON(user1).Expect().Status(httptest.StatusNoContent)

	// register user 2
	user2 := map[string]interface{}{
		"username": "user2",
		"email":    "user2@email.com",
		"password": "123123",
	}

	e.POST("/api/register").WithJSON(user2).Expect().Status(httptest.StatusNoContent)

	//Sign in to user 1 to grab bearer token
	signedInUser1 := map[string]interface{}{
		"username": user1["username"],
		"password": user1["password"],
	}

	user1token := e.POST("/api/signin").WithJSON(signedInUser1).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	//Sign in to user 1 to grab bearer token
	signedInUser2 := map[string]interface{}{
		"username": user2["username"],
		"password": user2["password"],
	}

	user2token := e.POST("/api/signin").WithJSON(signedInUser2).Expect().Status(httptest.StatusOK).JSON().Object().Raw()["access_token"].(string)

	// user 1 follow user 2
	username2 := user2["username"].(string)

	//Intially user does not follow
	e.GET("/api/users/"+username2+"/isfollowing").WithHeader("Authorization", "Bearer "+user1token).Expect().Status(httptest.StatusOK).JSON().Object().ValueEqual("isFollowing", false)

	//User follows
	e.POST("/api/users/"+username2+"/follow").WithHeader("Authorization", "Bearer "+user1token).Expect().Status(httptest.StatusOK)
	e.GET("/api/users/"+username2+"/isfollowing").WithHeader("Authorization", "Bearer "+user1token).Expect().Status(httptest.StatusOK).JSON().Object().ValueEqual("isFollowing", true)

	// user 1 posts msg
	message1 := map[string]interface{}{
		"content": "Message from user1",
	}
	e.POST("/api/tweets").WithJSON(message1).WithHeader("Authorization", "Bearer "+user1token).Expect().Status(httptest.StatusNoContent)

	// user 2 posts msg
	message2 := map[string]interface{}{
		"content": "Message from user2",
	}
	e.POST("/api/tweets").WithJSON(message2).WithHeader("Authorization", "Bearer "+user2token).Expect().Status(httptest.StatusNoContent)

	// user 1 checks private timeline to find the msgs
	tweets := e.GET("/api/timeline").WithHeader("Authorization", "Bearer "+user1token).Expect().JSON().Array()
	tweets.Element(0).Object().ValueEqual("Text", message2["content"])
	tweets.Element(1).Object().ValueEqual("Text", message1["content"])
}
