package dto

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id                uuid.UUID `db:"event_id"`
	Code              int       `db:"event_code"`
	Title             string    `db:"event_title"`
	Description       string    `db:"event_description"`
	StartsAt          time.Time `db:"event_starts_at"`
	EndsAt            time.Time `db:"event_ends_at"`
	Location          string    `db:"event_location"`
	IsInternal        bool      `db:"event_is_internal"`
	VideoUrl          string    `db:"event_video_url"`
	SpecificType      string    `db:"event_specific_type"`
	SpecificSituation string    `db:"event_specific_situation"`
	*EventType
	*EventSituation
}
