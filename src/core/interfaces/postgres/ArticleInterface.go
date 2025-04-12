package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/google/uuid"
	"time"
)

type Article interface {
	GetArticlesByReferenceDate(date time.Time) ([]article.Article, error)
	GetNewsletterArticlesByNewsletterId(newsletterId uuid.UUID) ([]article.Article, error)
}
