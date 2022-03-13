package main

import (
	"aske-w/itu-minitwit/database"
	"aske-w/itu-minitwit/environment"
	"aske-w/itu-minitwit/middleware"
	"aske-w/itu-minitwit/models"
	"aske-w/itu-minitwit/services"
	"aske-w/itu-minitwit/web/controllers"
	"log"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	app := iris.New()
	// app.Logger().SetLevel("debug") // more logging
	app.Favicon("./web/public/favicon.ico")
	// Load env's
	environment.InitEnv()

	app.Use(logger.New())  // logs request
	app.Use(recover.New()) // handles panics (shows 404)

	// Register middleware
	app.Use(middleware.InitMiddleware)

	// Configure sessions manager.
	sess := sessions.New(sessions.Config{
		Cookie:                      "itu-minitwit-cookie",
		AllowReclaim:                true,
		DisableSubdomainPersistence: true,
	})
	app.Use(sess.Handler())

	// Add html files
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "./web/public")

	// Register default error view
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "Error occured"))
		ctx.View("shared/error.html")
	})

	db, err := database.ConnectMySql()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
	userService := services.NewUserService(db)
	authService := services.NewAuthService(db)
	timelineService := services.NewTimelineService(db)
	messageService := services.NewMessageService(db)

	// I cant figure out how to have global DI, when using MVC pattern?
	index := mvc.New(app.Party("/"))
	index.Register(timelineService)
	index.Register(messageService)
	index.Register(userService)
	index.Handle(new(controllers.IndexController))

	auth := mvc.New(app.Party("/"))
	auth.Register(userService)
	auth.Register(authService)
	auth.Handle(new(controllers.AuthController))

	// Setup prometheus for monitoring
	var count int64 = 0
	var avgFollowers float64 = 0
	usersCount := promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "minitwit",
		Name:      "total_users_count",
		Help:      "The total amount of users in the database",
	})
	avgFollowersCount := promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "minitwit",
		Name:      "average_followers_count",
		Help:      "The total amount of users in the database",
	})
	//run non-middleware metrics data collection for in separate thread.
	// middleware data is collected in ./middleware/prometheusMiddleware.go
	go func() {
		for {
			db.Model(&models.User{}).Count(&count)
			db.Raw("select ((select count(follower_id) from followers) / (select count(*) from users));").Scan(&avgFollowers)
			usersCount.Set(float64(count))
			avgFollowersCount.Set(avgFollowers)
			time.Sleep(60 * time.Second)
		}
	}()

	app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	// make sure the latest row is in the database
	db.FirstOrCreate(&models.Latest{
		// id is always 0
		ID: 0,
	})
	api := mvc.New(app.Party("/api"))
	api.Register(db)
	api.Register(timelineService)
	api.Register(messageService)
	api.Register(userService)
	api.Register(authService)
	api.Handle(new(controllers.ApiController))

	tweetsAPI := app.Party("/tweets")
	{
		tweetsAPI.Get("/", func(ctx iris.Context) {
			tweets := []services.Tweet{}
			err := db.Model(&models.User{}).Select("users.id as UserId", "users.Username", "users.Email", "messages.id", "messages.Author_id", "messages.Text", "messages.Pub_date", "messages.Flagged").Joins("INNER JOIN messages ON messages.author_id = users.id AND messages.flagged = 0").Order("messages.pub_date DESC").Limit(30).Scan(&tweets).Error

			if err != nil {
				// return nil, err
			}

			services.AddAvatarAndDates(&tweets)

			ctx.JSON(tweets)
		})
	}

	app.Listen(":8080", iris.WithOptimizations)
}
