package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"time"
)

type Proposition interface {
	CreateProposition(proposition proposition.Proposition) error
	GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error)
	GetPropositionsByNewsletterId(newsletterId uuid.UUID) ([]proposition.Proposition, error)
	GetLatestPropositionCodes() ([]int, error)
}
