package dto

import (
	"github.com/google/uuid"
)

type EventSituation struct {
	Id          uuid.UUID `db:"event_situation_id"`
	Description string    `db:"event_situation_description"`
	Codes       string    `db:"event_situation_codes"`
	Color       string    `db:"event_situation_color"`
	IsFinished  bool      `db:"event_situation_is_finished"`
}
