package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
	"time"
	"vnc-summarizer/config/diconteiner"
)

func main() {
	if os.Getenv("APPLICATION_MODE") != "production" {
		err := godotenv.Load("config/.env")
		if err != nil {
			log.Fatal("Environment variables file not found: ", err.Error())
		}
	}

	backgroundDataService := diconteiner.GetBackgroundDataService()
	backgroundDataService.RegisterNewPropositions()

	for range time.NewTicker(time.Hour).C {
		backgroundDataService.RegisterNewPropositions()
		timeNow := time.Now()
		if timeNow.Hour() >= 18 {
			backgroundDataService.RegisterNewNewsletter(timeNow)
		} else if timeNow.Hour() < 6 {
			backgroundDataService.RegisterNewNewsletter(timeNow.AddDate(0, 0, -1))
		}
	}
}
