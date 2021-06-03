package main

import (
	"database/sql"
	forum_delivery "github.com/forum-api-back/internal/pkg/forum/handler"
	forum_repo "github.com/forum-api-back/internal/pkg/forum/repository"
	forum_usecase "github.com/forum-api-back/internal/pkg/forum/usecase"
	"log"

	"github.com/fasthttp/router"
	user_delivery "github.com/forum-api-back/internal/pkg/user/handler"
	user_repo "github.com/forum-api-back/internal/pkg/user/repository"
	user_usecase "github.com/forum-api-back/internal/pkg/user/usecase"
	"github.com/valyala/fasthttp"
)

func main() {
	// Connect to postgreSql db
	postgreSqlConn, err := sql.Open(
		"postgres",
		"user=docker "+
			"password=docker "+
			"dbname=docker "+
			"host=localhost "+
			"port=5432 ",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()
	if err := postgreSqlConn.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := user_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	userUCase := user_usecase.NewUseCase(userRepo)
	userHandler := user_delivery.NewHandler(userUCase)

	forumRepo := forum_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	forumUCase := forum_usecase.NewUseCase(forumRepo, userRepo)
	forumHandler := forum_delivery.NewHandler(forumUCase)

	mainRouter := router.New()
	mainRouter.POST("/api/forum/{slug}/create", forumHandler.CreateNewForum)
	mainRouter.GET("/api/forum/{slug}/details", forumHandler.GetForumDetails)
	mainRouter.POST("/api/user/{nickname}/create", userHandler.CreateNewUser)
	mainRouter.GET("/api/user/{nickname}/profile", userHandler.GetUserProfile)
	mainRouter.POST("/api/user/{nickname}/profile", userHandler.UpdateUserProfile)

	if err := fasthttp.ListenAndServe(":8080", mainRouter.Handler); err != nil {
		log.Fatal(err)
	}
}
