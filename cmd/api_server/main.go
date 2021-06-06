package main

import (
	"database/sql"
	"log"

	admin_delivery "github.com/forum-api-back/internal/pkg/admin/handler"
	admin_repo "github.com/forum-api-back/internal/pkg/admin/repository"
	admin_usecase "github.com/forum-api-back/internal/pkg/admin/usecase"
	forum_delivery "github.com/forum-api-back/internal/pkg/forum/handler"
	forum_repo "github.com/forum-api-back/internal/pkg/forum/repository"
	forum_usecase "github.com/forum-api-back/internal/pkg/forum/usecase"
	post_delivery "github.com/forum-api-back/internal/pkg/post/handler"
	post_repo "github.com/forum-api-back/internal/pkg/post/repository"
	post_usecase "github.com/forum-api-back/internal/pkg/post/usecase"
	thread_delivery "github.com/forum-api-back/internal/pkg/thread/handler"
	thread_repo "github.com/forum-api-back/internal/pkg/thread/repository"
	thread_usecase "github.com/forum-api-back/internal/pkg/thread/usecase"
	user_delivery "github.com/forum-api-back/internal/pkg/user/handler"
	user_repo "github.com/forum-api-back/internal/pkg/user/repository"
	user_usecase "github.com/forum-api-back/internal/pkg/user/usecase"

	"github.com/fasthttp/router"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
)

func main() {
	// Connect to postgreSql db
	postgreSqlConn, err := sql.Open(
		"postgres",
		"user=postgres "+
			"password=postgres "+
			"dbname=forum_db "+
			"host=localhost "+
			"port=5432 ",
	)
	postgreSqlConn.SetMaxIdleConns(100)
	postgreSqlConn.SetMaxIdleConns(100)

	if err != nil {
		log.Fatal(err)
	}
	defer postgreSqlConn.Close()
	if err := postgreSqlConn.Ping(); err != nil {
		log.Fatal(err)
	}

	userRepo := user_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	forumRepo := forum_repo.NewSessionPostgresqlRepository(postgreSqlConn)

	userUCase := user_usecase.NewUseCase(userRepo, forumRepo)
	userHandler := user_delivery.NewHandler(userUCase)

	forumUCase := forum_usecase.NewUseCase(forumRepo, userRepo)
	forumHandler := forum_delivery.NewHandler(forumUCase)

	threadRepo := thread_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	threadUCase := thread_usecase.NewUseCase(threadRepo, forumRepo)
	threadHandler := thread_delivery.NewHandler(threadUCase)

	postRepo := post_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	postUCase := post_usecase.NewUseCase(postRepo, threadRepo, forumRepo, userRepo)
	postHandler := post_delivery.NewHandler(postUCase)

	adminRepo := admin_repo.NewSessionPostgresqlRepository(postgreSqlConn)
	adminUCase := admin_usecase.NewUseCase(adminRepo)
	adminHandler := admin_delivery.NewHandler(adminUCase)

	mainRouter := router.New()
	mainRouter.POST("/api/forum/create", forumHandler.CreateNewForum)
	mainRouter.GET("/api/forum/{slug}/details", forumHandler.GetForumDetails)
	mainRouter.POST("/api/forum/{slug}/create", threadHandler.CreateNewThread)
	mainRouter.GET("/api/forum/{slug}/users", userHandler.GetUsersByForum)
	mainRouter.GET("/api/forum/{slug}/threads", threadHandler.GetThreadsByForum)
	mainRouter.GET("/api/post/{id}/details", postHandler.GetPostDetails)
	mainRouter.POST("/api/post/{id}/details", postHandler.UpdatePostDetails)
	mainRouter.POST("/api/service/clear", adminHandler.ClearBase)
	mainRouter.GET("/api/service/status", adminHandler.GetBaseDetails)
	mainRouter.POST("/api/thread/{slug_or_id}/create", postHandler.CreateNewPosts)
	mainRouter.GET("/api/thread/{slug_or_id}/details", threadHandler.GetThreadDetails)
	mainRouter.POST("/api/thread/{slug_or_id}/details", threadHandler.UpdateThreadDetails)
	mainRouter.GET("/api/thread/{slug_or_id}/posts", postHandler.GetPostsByThread)
	mainRouter.POST("/api/thread/{slug_or_id}/vote", threadHandler.UpdateThreadVote)
	mainRouter.POST("/api/user/{nickname}/create", userHandler.CreateNewUser)
	mainRouter.GET("/api/user/{nickname}/profile", userHandler.GetUserProfile)
	mainRouter.POST("/api/user/{nickname}/profile", userHandler.UpdateUserProfile)

	if err := fasthttp.ListenAndServe(":5000", mainRouter.Handler); err != nil {
		log.Fatal(err)
	}
}
