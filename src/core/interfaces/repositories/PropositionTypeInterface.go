package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
)

type PropositionType interface {
	GetPropositionTypeByCodeOrDefaultType(code string) (*propositiontype.PropositionType, error)
}
