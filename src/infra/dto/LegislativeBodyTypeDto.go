package dto

import (
	"github.com/google/uuid"
)

type LegislativeBodyType struct {
	Id          uuid.UUID `db:"legislative_body_type_id"`
	Code        int       `db:"legislative_body_type_code"`
	Description string    `db:"legislative_body_type_description"`
}
