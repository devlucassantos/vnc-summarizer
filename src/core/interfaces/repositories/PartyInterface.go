package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/party"
)

type Party interface {
	CreateParty(party party.Party) (*uuid.UUID, error)
	UpdateParty(party party.Party) error
	GetPartyByCode(code int) (*party.Party, error)
}
