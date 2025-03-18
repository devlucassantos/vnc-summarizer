package services

import "github.com/google/uuid"

type Event interface {
	RegisterNewEvents()
	UpdateEventsOccurringToday()
	UpdateEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished()
	RegisterNewEventByCode(code int) (*uuid.UUID, error)
}
