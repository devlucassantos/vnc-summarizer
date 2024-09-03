package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"time"
)

type Newsletter interface {
	CreateNewsletter(newsletter newsletter.Newsletter, propositions []proposition.Proposition) error
	UpdateNewsletter(newsletter newsletter.Newsletter, newPropositions []proposition.Proposition) error
	GetNewsletterByReferenceDate(referenceDate time.Time) (*newsletter.Newsletter, error)
}
