package main

import (
	"database/sql"
	user_delivery "github.com/forum-api-back/internal/pkg/user/handler"
	user_repo "github.com/forum-api-back/internal/pkg/user/repository"
	user_usecase "github.com/forum-api-back/internal/pkg/user/usecase"
	"log"

	"github.com/fasthttp/router"
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

	mainRouter := router.New()
	mainRouter.POST("/api/user/{nickname}/create", userHandler.CreateNewUser)
	mainRouter.GET("/api/user/{nickname}/profile", userHandler.GetUserProfile)
	mainRouter.POST("/api/user/{nickname}/profile", userHandler.UpdateUserProfile)

	if err := fasthttp.ListenAndServe(":8080", mainRouter.Handler); err != nil {
		log.Fatal(err)
	}
}
