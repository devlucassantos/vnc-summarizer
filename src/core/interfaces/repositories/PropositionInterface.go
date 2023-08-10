package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/proposition"
)

type Proposition interface {
	CreateProposition(proposition proposition.Proposition) (*uuid.UUID, error)
	GetLatestPropositionCodes() ([]int, error)
}
