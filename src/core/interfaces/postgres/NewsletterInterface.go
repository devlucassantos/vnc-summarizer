package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"time"
)

type Newsletter interface {
	CreateNewsletter(newsletter newsletter.Newsletter) (*uuid.UUID, error)
	UpdateNewsletter(newsletter newsletter.Newsletter, newArticles []article.Article) error
	GetNewsletterByReferenceDate(referenceDate time.Time) (*newsletter.Newsletter, error)
}
