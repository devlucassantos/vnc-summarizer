package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"time"
)

type Proposition interface {
	CreateProposition(proposition proposition.Proposition) error
	GetLatestPropositionCodes() ([]int, error)
	GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error)
}
