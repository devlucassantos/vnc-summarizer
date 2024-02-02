package dto

import (
	"github.com/google/uuid"
	"time"
)

type Newsletter struct {
	Id            uuid.UUID `db:"newsletter_id"`
	Title         string    `db:"newsletter_title"`
	Content       string    `db:"newsletter_content"`
	ReferenceDate time.Time `db:"newsletter_reference_date"`
	Active        bool      `db:"newsletter_active"`
	CreatedAt     time.Time `db:"newsletter_created_at"`
	UpdatedAt     time.Time `db:"newsletter_updated_at"`
}
