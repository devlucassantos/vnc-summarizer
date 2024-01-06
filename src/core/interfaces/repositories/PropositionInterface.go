package repositories

import (
	"github.com/google/uuid"
	"time"
	"vnc-write-api/core/domains/proposition"
)

type Proposition interface {
	CreateProposition(proposition proposition.Proposition) (*uuid.UUID, error)
	GetLatestPropositionCodes() ([]int, error)
	GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error)
}
