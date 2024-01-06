package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/newsletter"
)

type Newsletter interface {
	CreateNewsletter(newsletter newsletter.Newsletter) (*uuid.UUID, error)
}
