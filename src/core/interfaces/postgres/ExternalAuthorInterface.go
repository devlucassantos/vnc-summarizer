package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/google/uuid"
)

type ExternalAuthor interface {
	CreateExternalAuthor(externalAuthor externalauthor.ExternalAuthor) (*uuid.UUID, error)
	GetExternalAuthorByNameAndTypeCode(name string, typeCode int) (*externalauthor.ExternalAuthor, error)
}
