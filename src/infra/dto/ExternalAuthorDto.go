package dto

import (
	"github.com/google/uuid"
	"time"
)

type ExternalAuthor struct {
	Id        uuid.UUID `db:"external_author_id"`
	Name      string    `db:"external_author_name"`
	Type      string    `db:"external_author_type"`
	CreatedAt time.Time `db:"external_author_created_at"`
	UpdatedAt time.Time `db:"external_author_updated_at"`
}
