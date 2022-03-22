package main

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/environment"
	"aske-w/itu-minitwit/models"
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/controllers"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	sigKey = []byte("secret")
	encKey = []byte("GCM_AES_256_secret_shared_key_32")
)

type UserClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

var userService services.UserService
var authService services.AuthService
var timelineService services.TimelineService
var messageService services.MessageService

func main() {
	app := iris.Default()

	environment.InitEnv()

	db, err := database.ConnectMySql()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	userService = *services.NewUserService(db)
	authService = *services.NewAuthService(db)
	timelineService = *services.NewTimelineService(db)
	messageService = *services.NewMessageService(db)

	signer := jwt.NewSigner(jwt.HS256, sigKey, 60*time.Minute)

	verifier := jwt.NewVerifier(jwt.HS256, sigKey)

	verifier.WithDefaultBlocklist()

	authMiddleware := verifier.Verify(func() interface{} {
		return new(UserClaims)
	})

	app.Post("/api/signin", signinHandler(signer, db))

	app.Get("/api/tweets", indexHandler(db))
	app.Post("/api/tweets", storeTweetHandler()).Use(authMiddleware)

	app.Get("/api/users/{username}", userHandler(db))
	app.Get("/api/users/{username}/tweets", userTweets(db))

	app.Post("/api/users/{username}/follow", followHandler()).Use(authMiddleware)
	app.Get("/api/users/{username}/isfollowing", isFollowingHandler()).Use(authMiddleware)

	app.Get("/api/timeline", timeline(db)).Use(authMiddleware)

	// Simulator endpoints
	app.Post("/api/msgs/{username}", simulatorStoreTweetHandler())
	app.Post("/api/register", signupHandler(db))
	app.Post("/api/fllws/{username}", simulatorFollowHandler())

	app.Listen(":8080")
}

func isFollowingHandler() iris.Handler {
	return func(ctx iris.Context) {
		claims := jwt.Get(ctx).(*UserClaims)
		username := ctx.Params().Get("username")

		followee, err := userService.FindByUsername(username)

		if err != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "Cant find user"})

			return
		}

		isFollowing := userService.UserIsFollowing(claims.Id, followee.ID)

		ctx.StatusCode(200)
		ctx.JSON(iris.Map{"isFollowing": isFollowing})

	}
}

func followHandler() iris.Handler {
	return func(ctx iris.Context) {
		claims := jwt.Get(ctx).(*UserClaims)
		username := ctx.Params().Get("username")

		followee, err := userService.FindByUsername(username)

		if err != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "Cant find user"})

			return
		}

		isFollowingAlready := userService.UserIsFollowing(claims.Id, followee.ID)

		if isFollowingAlready {
			// Unfollow
			_, err := userService.UnfollowUser(claims.Id, followee.ID)

			if err != nil {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{"error": "Cant unfollow user"})

				return
			}
		} else {
			// Follow
			_, err := userService.FollowUser(claims.Id, followee.ID)

			if err != nil {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{"error": "Cant follow user"})

				return
			}
		}

		ctx.StatusCode(200)
		ctx.JSON(iris.Map{"success": true})
	}
}

type FollowUserRequest struct {
	Follow   string `json:"follow"`
	Unfollow string `json:"unfollow"`
}

func simulatorFollowHandler() iris.Handler {
	return func(ctx iris.Context) {
		request := FollowUserRequest{}
		err := ctx.ReadJSON(&request)

		if err != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": err.Error()})

			return
		}

		username := ""

		if len(request.Follow) != 0 {
			username = request.Follow
		} else {
			username = request.Unfollow
		}

		auth, authErr := userService.FindByUsername(ctx.Params().Get("username"))

		if authErr != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "Cant find user"})

			return
		}

		userToFollow, err := userService.FindByUsername(username)

		if err != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "Cant find user"})

			return
		}

		isFollowingAlready := userService.UserIsFollowing(auth.ID, userToFollow.ID)

		if isFollowingAlready {
			// Unfollow
			_, err := userService.UnfollowUser(auth.ID, userToFollow.ID)

			if err != nil {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{"error": "Cant unfollow user"})

				return
			}
		} else {
			// Follow
			_, err := userService.FollowUser(auth.ID, userToFollow.ID)

			if err != nil {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{"error": "Cant follow user"})

				return
			}
		}

		ctx.StatusCode(204)
		ctx.JSON(iris.Map{"success": true})
	}
}

func userTweets(db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		tweets := []services.Tweet{}
		err := db.Raw(`
			SELECT
				users.id as UserId,
				users.Username,
				users.Email,
				messages.id as Message_id,
				messages.Author_id,
				messages.Text,
				messages.Pub_date,
				messages.Flagged
			from users, messages
			where
				messages.flagged = 0 and
				messages.author_id = users.id and
				(
					users.username = ?
				)
			order by messages.pub_date DESC
			limit ?
		`, ctx.Params().Get("username"), 30).Scan(&tweets).Error

		if err != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "Tweets not found"})

			return
		}

		services.AddAvatarAndDates(&tweets)

		ctx.JSON(tweets)
	}
}

func userHandler(db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		var user models.User
		result := db.First(&user, "username = ?", ctx.Params().Get("username"))

		if result.Error != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "User not found"})

			return
		}

		ctx.JSON(iris.Map{
			"id":       user.ID,
			"username": user.Username,
		})
	}
}

func indexHandler(db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		tweets := []services.Tweet{}
		err := db.Model(&models.User{}).Select("users.id as UserId", "users.Username", "users.Email", "messages.id as Message_id", "messages.Author_id", "messages.Text", "messages.Pub_date", "messages.Flagged").Joins("INNER JOIN messages ON messages.author_id = users.id AND messages.flagged = 0").Order("messages.pub_date DESC").Limit(30).Scan(&tweets).Error

		if err != nil {
			// return nil, err
		}

		services.AddAvatarAndDates(&tweets)

		ctx.JSON(tweets)
	}
}

type StoreTweetRequest struct {
	Text string `json:"content"`
}

func storeTweetHandler() iris.Handler {
	return func(ctx iris.Context) {
		tweetRequest := StoreTweetRequest{}
		err := ctx.ReadJSON(&tweetRequest)

		if err != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": err.Error()})

			return
		}

		claims := jwt.Get(ctx).(*UserClaims)

		fmt.Println(claims.Id, tweetRequest.Text)

		messageErr := messageService.CreateMessage(int(claims.Id), tweetRequest.Text)

		if messageErr != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": messageErr})

			return
		}

		ctx.StatusCode(204)
	}
}

func simulatorStoreTweetHandler() iris.Handler {
	return func(ctx iris.Context) {
		tweetRequest := StoreTweetRequest{}
		err := ctx.ReadJSON(&tweetRequest)

		if err != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": err.Error()})

			return
		}

		username := ctx.Params().Get("username")

		user, userErr := userService.FindByUsername(username)

		if userErr != nil {
			ctx.StatusCode(404)
			ctx.JSON(iris.Map{"error": "Cant find user"})

			return
		}

		messageErr := messageService.CreateMessage(int(user.ID), tweetRequest.Text)

		if messageErr != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": messageErr})

			return
		}

		ctx.StatusCode(204)
	}
}

func signupHandler(db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		user := controllers.RegisterUser{}
		err := ctx.ReadJSON(&user)

		errors := make([]string, 0)

		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{"error": err.Error()})

			return
		}

		if user.Username == "" {
			errors = append(errors, "You have to enter a username")
		}

		if user.Email == "" || !strings.Contains(user.Email, "@") {
			errors = append(errors, "You have to enter a valid email address")
		}

		if user.Password == "" {
			errors = append(errors, "You have to enter a password")
		}

		if len(errors) > 0 {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"errors": errors})

			return
		}

		exists, _ := userService.CheckUsernameExists(user.Username)

		if exists {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"errors": [1]string{"The username is already taken"}})

			return
		}

		err = authService.CreateUser(user.Username, user.Email, user.Password)

		if err == nil {
			// update_latest(c)
			ctx.StatusCode(204)

			return
		}
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func signinHandler(signer *jwt.Signer, db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		// Sign in logic
		loginRequest := LoginRequest{}
		err := ctx.ReadJSON(&loginRequest)

		if err != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": err.Error()})

			return
		}

		user := &models.User{}
		result := db.First(&user, "username = ?", loginRequest.Username)

		if result.Error != nil {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"errors": "Invalid username and/or password"})

			return
		}

		passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Pw_Hash), []byte(loginRequest.Password))

		if passwordErr != nil {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"errors": "Invalid username and/or password"})

			return
		}

		claims := UserClaims{
			Id:       user.ID,
			Username: user.Username,
		}

		token, err := signer.Sign(claims)
		if err != nil {
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		ctx.StatusCode(200)
		ctx.JSON(iris.Map{
			"access_token": string(token[:]),
			"username":     user.Username,
		})
	}
}

func timeline(db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		claims := jwt.Get(ctx).(*UserClaims)

		tweets := []services.Tweet{}
		err := db.Raw(`
			SELECT
				users.id as UserId,
				users.Username,
				users.Email,
				messages.id as Message_id,
				messages.Author_id,
				messages.Text,
				messages.Pub_date,
				messages.Flagged
			from users, messages
			where
				messages.flagged = 0 and
				messages.author_id = users.id and
				(
					users.id = ? or
					users.id in (select follower_id from followers where user_id = ?)
				)
			order by messages.pub_date DESC
			limit ?
		`, claims.Id, claims.Id, 30).Scan(&tweets).Error

		if err != nil {

		}

		services.AddAvatarAndDates(&tweets)

		ctx.JSON(tweets)
	}
}
