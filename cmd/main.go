package main

import (
	"blogAggregator/internal/config"
	"blogAggregator/internal/database"
	"blogAggregator/internal/jobs"
	"blogAggregator/internal/server"
	"fmt"
	"log"
	"time"
)

func main() {

	cfg := config.LoadConfig()

	database.ConnectDatabase(cfg.DBPath)

	r := server.NewRouter()
	go jobs.StartFeedUpdater(5 * time.Minute)
	fmt.Println("server is running :8080")
	err := r.Run(":" + cfg.Port)
	if err != nil {
		log.Fatal(err)
	}
}
