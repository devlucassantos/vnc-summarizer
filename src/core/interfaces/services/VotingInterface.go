package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
)

type Voting interface {
	RegisterNewVotes()
	RegisterNewVotingByCode(code string) (*uuid.UUID, error)
	GetVotesByCodes(codes []string) ([]voting.Voting, error)
}
