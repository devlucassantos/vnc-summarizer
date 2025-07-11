package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"os"
	"time"
	"vnc-summarizer/config/dicontainer"
	"vnc-summarizer/utils/datetime"
)

func main() {
	if os.Getenv("APPLICATION_MODE") != "production" {
		err := godotenv.Load("config/.env")
		if err != nil {
			log.Fatal("Environment variables file not found: ", err.Error())
		}
	}

	propositionService := dicontainer.GetPropositionService()
	newsletterService := dicontainer.GetNewsletterService()
	votingService := dicontainer.GetVotingService()
	eventService := dicontainer.GetEventService()

	for {
		startTime, err := datetime.GetCurrentDateTimeInBrazil()
		if err != nil {
			log.Fatal("datetime.GetCurrentDateTimeInBrazil(): ", err)
			return
		}

		propositionService.RegisterNewPropositions()
		votingService.RegisterNewVotes()
		eventService.UpdateEventsOccurringToday()
		eventService.RegisterNewEvents()

		if startTime.Hour() >= 18 {
			newsletterService.RegisterNewNewsletter(*startTime)
		} else if startTime.Hour() < 6 {
			newsletterService.RegisterNewNewsletter(startTime.AddDate(0, 0, -1))
		}

		if startTime.Hour() == 7 {
			eventService.UpdateEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished()
		}

		elapsedTime := time.Since(*startTime)
		sleepDuration := time.Hour - elapsedTime
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}
