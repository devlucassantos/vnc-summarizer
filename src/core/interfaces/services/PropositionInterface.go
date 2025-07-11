package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
)

type Proposition interface {
	RegisterNewPropositions()
	RegisterNewPropositionByCode(code int) (*uuid.UUID, error)
	GetPropositionsByCodes(codes []int) ([]proposition.Proposition, error)
}
