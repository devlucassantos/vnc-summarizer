package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
)

type ArticleType interface {
	GetArticleTypeByCode(code string) (*articletype.ArticleType, error)
}
