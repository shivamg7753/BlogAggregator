package server

import (
	"blogAggregator/internal/handlers"
	"blogAggregator/internal/middleware"
	"net/http"
  swaggerFiles "github.com/swaggo/files"
  ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	_ "blogAggregator/docs"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth
	r.POST("/login", handlers.Login)
	//protected routes
	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())

	//users
	r.POST("/users/register", handlers.RegisterUser)
	r.POST("/users", handlers.CreateUser)
	authRoutes.POST("/subscriptions", handlers.SubscribeFeed)
	authRoutes.GET("/users/:id/feed", handlers.GetUserFeed)

	//feeds
	r.POST("/feeds", handlers.CreateFeed)
	r.GET("/feeds", handlers.ListFeeds)
	r.POST("/feeds/refresh", handlers.RefreshFeed)

	//post
	r.GET("/posts", handlers.ListPosts)

	return r
}
