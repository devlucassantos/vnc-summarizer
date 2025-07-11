package dto

import (
	"github.com/google/uuid"
)

type Article struct {
	Id uuid.UUID `db:"article_id"`
	*ArticleType
	*Proposition
	*Newsletter
	*Voting
	*Event
}
