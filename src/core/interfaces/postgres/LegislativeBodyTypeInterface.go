package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/google/uuid"
)

type LegislativeBodyType interface {
	CreateLegislativeBodyType(legislativeBodyType legislativebodytype.LegislativeBodyType) (*uuid.UUID, error)
	GetLegislativeBodyTypeByCode(code int) (*legislativebodytype.LegislativeBodyType, error)
}
