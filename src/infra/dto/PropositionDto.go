package dto

import (
	"github.com/google/uuid"
	"time"
)

type Proposition struct {
	Id              uuid.UUID `db:"proposition_id"`
	Code            int       `db:"proposition_code"`
	OriginalTextUrl string    `db:"proposition_original_text_url"`
	Title           string    `db:"proposition_title"`
	Content         string    `db:"proposition_content"`
	SubmittedAt     time.Time `db:"proposition_submitted_at"`
	CreatedAt       time.Time `db:"proposition_created_at"`
	UpdatedAt       time.Time `db:"proposition_updated_at"`
}
