package config

import (
	"github.com/joho/godotenv"
	"log"
)

func loadEnvFile() {
	err := godotenv.Load("api/config/.env")
	if err != nil {
		log.Fatal("Environment variables file not found: ", err)
	}
}
