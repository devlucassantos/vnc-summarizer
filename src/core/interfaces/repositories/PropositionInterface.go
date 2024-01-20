package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"time"
)

type Proposition interface {
	CreateProposition(proposition proposition.Proposition) (*uuid.UUID, error)
	GetLatestPropositionCodes() ([]int, error)
	GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error)
}
