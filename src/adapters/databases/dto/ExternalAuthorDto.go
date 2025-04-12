package dto

import (
	"github.com/google/uuid"
)

type ExternalAuthor struct {
	Id   uuid.UUID `db:"external_author_id"`
	Name string    `db:"external_author_name"`
	*ExternalAuthorType
}
