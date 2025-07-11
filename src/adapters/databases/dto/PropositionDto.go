package dto

import (
	"github.com/google/uuid"
	"time"
)

type Proposition struct {
	Id                   uuid.UUID `db:"proposition_id"`
	Code                 int       `db:"proposition_code"`
	OriginalTextUrl      string    `db:"proposition_original_text_url"`
	OriginalTextMimeType string    `db:"proposition_original_text_mime_type"`
	Title                string    `db:"proposition_title"`
	Content              string    `db:"proposition_content"`
	SubmittedAt          time.Time `db:"proposition_submitted_at"`
	ImageUrl             string    `db:"proposition_image_url"`
	ImageDescription     string    `db:"proposition_image_description"`
	*PropositionType
}
