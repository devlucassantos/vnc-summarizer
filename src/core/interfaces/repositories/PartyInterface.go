package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
)

type Party interface {
	CreateParty(party party.Party) (*uuid.UUID, error)
	UpdateParty(party party.Party) error
	GetPartyByCode(code int) (*party.Party, error)
}
