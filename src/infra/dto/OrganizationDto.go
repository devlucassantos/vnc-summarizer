package dto

import (
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	Id        uuid.UUID `db:"organization_id"`
	Code      int       `db:"organization_code"`
	Name      string    `db:"organization_name"`
	Acronym   string    `db:"organization_acronym"`
	Nickname  string    `db:"organization_nickname"`
	Active    bool      `db:"organization_active"`
	CreatedAt time.Time `db:"organization_created_at"`
	UpdatedAt time.Time `db:"organization_updated_at"`
}
