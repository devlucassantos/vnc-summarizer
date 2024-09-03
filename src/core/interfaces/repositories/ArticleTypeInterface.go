package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
)

type ArticleType interface {
	GetArticleTypeByCodeOrDefaultType(code string) (*articletype.ArticleType, error)
}
