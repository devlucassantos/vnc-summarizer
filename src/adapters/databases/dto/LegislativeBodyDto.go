package dto

import (
	"github.com/google/uuid"
)

type LegislativeBody struct {
	Id      uuid.UUID `db:"legislative_body_id"`
	Code    int       `db:"legislative_body_code"`
	Name    string    `db:"legislative_body_name"`
	Acronym string    `db:"legislative_body_acronym"`
	*LegislativeBodyType
}
