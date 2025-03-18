package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
)

type Voting interface {
	CreateVoting(voting voting.Voting) (*uuid.UUID, error)
	GetVotingByCode(code string) (*voting.Voting, error)
	GetVotesByCodes(codes []string) ([]voting.Voting, error)
}
