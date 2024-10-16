package dto

import (
	"github.com/google/uuid"
	"time"
)

type ArticleType struct {
	Id          uuid.UUID `db:"article_type_id"`
	Description string    `db:"article_type_description"`
	Codes       string    `db:"article_type_codes"`
	Color       string    `db:"article_type_color"`
	SortOrder   int       `db:"article_type_sort_order"`
	CreatedAt   time.Time `db:"article_type_created_at"`
	UpdatedAt   time.Time `db:"article_type_updated_at"`
}
