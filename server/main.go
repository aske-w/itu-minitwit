package main

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/environment"
	"aske-w/itu-minitwit/models"
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/controllers"
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

func main() {
	app := iris.Default()

	// app.Logger().SetLevel("debug") // more logging
	// app.Favicon("./web/public/favicon.ico")
	// Load env's
	environment.InitEnv()

	// app.Use(logger.New())  // logs request
	// app.Use(recover.New()) // handles panics (shows 404)

	// Register middleware
	// app.Use(middleware.InitMiddleware)

	// Configure sessions manager.
	// sess := sessions.New(sessions.Config{
	// 	Cookie:                      "itu-minitwit-cookie",
	// 	AllowReclaim:                true,
	// 	DisableSubdomainPersistence: true,
	// })
	// app.Use(sess.Handler())

	// Add html files
	// tmpl := iris.HTML("./web/views", ".html").
	// 	Layout("shared/layout.html").
	// 	Reload(true)
	// app.RegisterView(tmpl)
	// app.HandleDir("/public", "./web/public")

	// // Register default error view
	// app.OnAnyErrorCode(func(ctx iris.Context) {
	// 	ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "Error occured"))
	// 	ctx.View("shared/error.html")
	// })

	db, err := database.ConnectMySql()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	userService := services.NewUserService(db)
	authService := services.NewAuthService(db)
	// timelineService := services.NewTimelineService(db)
	// messageService := services.NewMessageService(db)

	// // I cant figure out how to have global DI, when using MVC pattern?
	// index := mvc.New(app.Party("/"))
	// index.Register(timelineService)
	// index.Register(messageService)
	// index.Register(userService)
	// index.Handle(new(controllers.IndexController))

	// auth := mvc.New(app.Party("/"))
	// auth.Register(userService)
	// auth.Register(authService)
	// auth.Handle(new(controllers.AuthController))

	// // Setup prometheus for monitoring
	// var count int64 = 0
	// var avgFollowers float64 = 0
	// usersCount := promauto.NewGauge(prometheus.GaugeOpts{
	// 	Subsystem: "minitwit",
	// 	Name:      "total_users_count",
	// 	Help:      "The total amount of users in the database",
	// })
	// avgFollowersCount := promauto.NewGauge(prometheus.GaugeOpts{
	// 	Subsystem: "minitwit",
	// 	Name:      "average_followers_count",
	// 	Help:      "The total amount of users in the database",
	// })
	// //run non-middleware metrics data collection for in separate thread.
	// // middleware data is collected in ./middleware/prometheusMiddleware.go
	// go func() {
	// 	for {
	// 		db.Model(&models.User{}).Count(&count)
	// 		db.Raw("select ((select count(follower_id) from followers) / (select count(*) from users));").Scan(&avgFollowers)
	// 		usersCount.Set(float64(count))
	// 		avgFollowersCount.Set(avgFollowers)
	// 		time.Sleep(60 * time.Second)
	// 	}
	// }()

	// app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	// // make sure the latest row is in the database
	// db.FirstOrCreate(&models.Latest{
	// 	// id is always 0
	// 	ID: 0,
	// })
	// api := mvc.New(app.Party("/api"))
	// api.Register(db)
	// api.Register(timelineService)
	// api.Register(messageService)
	// api.Register(userService)
	// api.Register(authService)
	// api.Handle(new(controllers.ApiController))

	signer := jwt.NewSigner(jwt.HS256, sigKey, 60*time.Minute)

	verifier := jwt.NewVerifier(jwt.HS256, sigKey)

	verifier.WithDefaultBlocklist()

	authMiddleware := verifier.Verify(func() interface{} {
		return new(UserClaims)
	})

	app.Post("/api/signup", signupHandler(db, userService, authService))
	app.Post("/api/signin", signinHandler(signer, db))

	app.Get("/api/tweets", indexHandler(db))
	app.Post("/api/tweets", storeTweetHandler(db)).Use(authMiddleware)

	app.Get("/api/users/{username}", userHandler(db))
	app.Get("/api/users/{username}/tweets", userTweets(db))

	// protectedAPI := app.Party("/api/protected")
	// protectedAPI.Use(authMiddleware)
	app.Get("/api/timeline", timeline(db)).Use(authMiddleware)

	app.Listen(":8080")
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
	Text string `json:"text"`
}

func storeTweetHandler(db *gorm.DB) iris.Handler {
	return func(ctx iris.Context) {
		tweetRequest := StoreTweetRequest{}
		err := ctx.ReadJSON(&tweetRequest)

		if err != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": err.Error()})

			return
		}

		claims := jwt.Get(ctx).(*UserClaims)

		message := models.Message{
			Author_id: int(claims.Id),
			Text:      tweetRequest.Text,
			Pub_date:  int(time.Now().Unix()),
		}

		result := db.Create(&message)

		if result.Error != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": result.Error.Error()})

			return
		}

		ctx.StatusCode(200)
	}
}

func signupHandler(db *gorm.DB, userService *services.UserService, authService *services.AuthService) iris.Handler {
	return func(ctx iris.Context) {
		user := controllers.RegisterUser{}
		err := ctx.ReadJSON(&user)

		if err != nil {
			ctx.StatusCode(422)
			ctx.JSON(iris.Map{"error": err.Error()})
			// ctx.JSON(err)
			// ctx.WriteString(err.Error())

			return
		}

		if user.Username == "" {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"status": 400, "error_msg": "you have to enter a username"})

			return
		}

		if user.Email == "" || !strings.Contains(user.Email, "@") {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"status": 400, "error_msg": "you have to enter a valid email address"})

			return
		}

		if user.Password == "" {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"status": 400, "error_msg": "you have to enter a password"})

			return
		}

		exists, _ := userService.CheckUsernameExists(user.Username)

		if exists {
			ctx.StatusCode(400)
			ctx.JSON(iris.Map{"status": 400, "error_msg": "the username is already taken"})

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
