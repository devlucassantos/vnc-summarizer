package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/google/uuid"
)

type Event interface {
	CreateEvent(event event.Event) (*uuid.UUID, error)
	UpdateEvent(event event.Event) error
	GetEventsByCodes(codes []int) ([]event.Event, error)
	GetEventsOccurringToday() ([]event.Event, error)
	GetEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished() ([]event.Event, error)
}
