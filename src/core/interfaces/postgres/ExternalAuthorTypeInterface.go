package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/google/uuid"
)

type ExternalAuthorType interface {
	CreateExternalAuthorType(externalAuthorType externalauthortype.ExternalAuthorType) (*uuid.UUID, error)
	GetExternalAuthorTypeByCode(code int) (*externalauthortype.ExternalAuthorType, error)
}
