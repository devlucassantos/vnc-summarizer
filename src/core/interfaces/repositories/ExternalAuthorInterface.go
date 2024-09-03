package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/google/uuid"
)

type ExternalAuthor interface {
	CreateExternalAuthor(externalAuthor external.ExternalAuthor) (*uuid.UUID, error)
	GetExternalAuthorByNameAndType(name string, _type string) (*external.ExternalAuthor, error)
}
