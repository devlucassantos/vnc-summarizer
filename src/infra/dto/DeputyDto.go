package dto

import (
	"github.com/google/uuid"
	"time"
)

type Deputy struct {
	Id            uuid.UUID `db:"deputy_id"`
	Code          int       `db:"deputy_code"`
	Cpf           string    `db:"deputy_cpf"`
	Name          string    `db:"deputy_name"`
	ElectoralName string    `db:"deputy_electoral_name"`
	ImageUrl      string    `db:"deputy_image_url"`
	Active        bool      `db:"deputy_active"`
	CreatedAt     time.Time `db:"deputy_created_at"`
	UpdatedAt     time.Time `db:"deputy_updated_at"`
	*Party
}
