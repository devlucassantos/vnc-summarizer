package dto

import "github.com/google/uuid"

type News struct {
	Id    uuid.UUID `db:"news_id"`
	Views int       `db:"news_views"`
	*Proposition
	*Newsletter
}
