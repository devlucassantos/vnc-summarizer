package dto

import (
	"github.com/google/uuid"
	"time"
)

type Party struct {
	Id        uuid.UUID `db:"party_id"`
	Code      int       `db:"party_code"`
	Name      string    `db:"party_name"`
	Acronym   string    `db:"party_acronym"`
	ImageUrl  string    `db:"party_image_url"`
	Active    bool      `db:"party_active"`
	CreatedAt time.Time `db:"party_created_at"`
	UpdatedAt time.Time `db:"party_updated_at"`
}
