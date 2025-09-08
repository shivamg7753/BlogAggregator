package main

import (
	"blogAggregator/internal/config"
	"blogAggregator/internal/database"
	"blogAggregator/internal/jobs"
	"blogAggregator/internal/server"
	"blogAggregator/docs"
	"fmt"
	"log"
	"time"
)

// @title           Blog Aggregator API
// @version         1.0
// @description     API for managing RSS feeds, users, and subscriptions
// @host            localhost:8080
// @BasePath        /
// @schemes         http

func main() {

	cfg := config.LoadConfig()

	// Swagger metadata
	docs.SwaggerInfo.Title = "Blog Aggregator API"
	docs.SwaggerInfo.Description = "API for managing RSS feeds, users, and subscriptions"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	database.ConnectDatabase(cfg.DBPath)

	r := server.NewRouter()
	go jobs.StartFeedUpdater(5 * time.Minute)
	fmt.Println("server is running :8080")
	err := r.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}
}
