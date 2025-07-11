package dto

import (
	"github.com/google/uuid"
)

type ExternalAuthorType struct {
	Id          uuid.UUID `db:"external_author_type_id"`
	Code        int       `db:"external_author_type_code"`
	Description string    `db:"external_author_type_description"`
}
