package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
)

type PropositionType interface {
	GetPropositionTypeByCode(code string) (*proptype.PropositionType, error)
}
