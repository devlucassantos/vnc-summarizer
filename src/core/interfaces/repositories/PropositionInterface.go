package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
)

type Proposition interface {
	CreateProposition(proposition proposition.Proposition) (*uuid.UUID, error)
	GetPropositionsByCodes(codes []int) ([]proposition.Proposition, error)
}
