package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file not found, using defaults")
	}

	return Config{
		Port:      getEnv("PORT"),
		DBPath:    getEnv("dsn"),
		JWTSecret: getEnv("JWT_SECRET"),
	}
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatalf("required enviornment variable %s is not set",key)
	return ""
}
