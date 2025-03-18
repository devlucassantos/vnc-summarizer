package dto

import (
	"github.com/google/uuid"
)

type ArticleType struct {
	Id          uuid.UUID `db:"article_type_id"`
	Description string    `db:"article_type_description"`
	Codes       string    `db:"article_type_codes"`
	Color       string    `db:"article_type_color"`
}
