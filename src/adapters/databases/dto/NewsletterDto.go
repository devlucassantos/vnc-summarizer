package dto

import (
	"github.com/google/uuid"
	"time"
)

type Newsletter struct {
	Id            uuid.UUID `db:"newsletter_id"`
	ReferenceDate time.Time `db:"newsletter_reference_date"`
	Description   string    `db:"newsletter_description"`
}
