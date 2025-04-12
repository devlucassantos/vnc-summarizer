package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
)

type EventSituation interface {
	GetEventSituationByCodeOrDefaultSituation(code string) (*eventsituation.EventSituation, error)
}
