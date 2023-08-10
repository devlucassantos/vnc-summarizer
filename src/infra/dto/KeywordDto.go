package dto

import (
	"github.com/google/uuid"
	"time"
)

type Keyword struct {
	Id        uuid.UUID `db:"keyword_id"`
	Keyword   string    `db:"keyword_keyword"`
	Active    bool      `db:"keyword_active"`
	CreatedAt time.Time `db:"keyword_created_at"`
	UpdatedAt time.Time `db:"keyword_updated_at"`
}
