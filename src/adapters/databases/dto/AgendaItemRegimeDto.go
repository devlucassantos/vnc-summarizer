package dto

import (
	"github.com/google/uuid"
)

type AgendaItemRegime struct {
	Id          uuid.UUID `db:"agenda_item_regime_id"`
	Code        int       `db:"agenda_item_regime_code"`
	Description string    `db:"agenda_item_regime_description"`
}
