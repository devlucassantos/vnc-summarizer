package dto

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	Id                uuid.UUID `db:"article_id"`
	Views             int       `db:"article_views"`
	AverageRating     float64   `db:"article_average_rating"`
	NumberOfRatings   int       `db:"article_number_of_ratings"`
	ReferenceDateTime time.Time `db:"article_reference_date_time"`
	CreatedAt         time.Time `db:"article_created_at"`
	UpdatedAt         time.Time `db:"article_updated_at"`
	*ArticleType
	*Proposition
	*Newsletter
}
