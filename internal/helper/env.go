package helper

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	return os.Getenv(key)
}
