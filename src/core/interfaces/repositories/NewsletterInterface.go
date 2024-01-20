package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
)

type Newsletter interface {
	CreateNewsletter(newsletter newsletter.Newsletter) (*uuid.UUID, error)
}
