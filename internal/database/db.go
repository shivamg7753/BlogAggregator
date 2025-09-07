package database

import (
	"blogAggregator/internal/models"
	"fmt"
	"log"
   "gorm.io/driver/postgres"
	"gorm.io/gorm"
)






var DB *gorm.DB

func ConnectDatabase(dsn string){
	var err error

	DB,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("failed to connect database :",err)
	}
	err= DB.AutoMigrate(&models.User{},&models.Feed{},&models.Post{},&models.Subscription{})
	if err!=nil{
		log.Fatal("failed to migrate database:",err)
	}
	fmt.Println("database connected & migrated successfully")
}