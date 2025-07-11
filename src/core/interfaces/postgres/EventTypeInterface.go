package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
)

type EventType interface {
	GetEventTypeByCodeOrDefaultType(code string) (*eventtype.EventType, error)
}
