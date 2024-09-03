package dto

import (
	"github.com/google/uuid"
	"time"
)

type Newsletter struct {
	Id            uuid.UUID `db:"newsletter_id"`
	ReferenceDate time.Time `db:"newsletter_reference_date"`
	Title         string    `db:"newsletter_title"`
	Description   string    `db:"newsletter_description"`
	CreatedAt     time.Time `db:"newsletter_created_at"`
	UpdatedAt     time.Time `db:"newsletter_updated_at"`
}
