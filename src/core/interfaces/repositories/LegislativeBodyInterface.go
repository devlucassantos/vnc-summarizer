package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/google/uuid"
)

type LegislativeBody interface {
	CreateLegislativeBody(legislativeBody legislativebody.LegislativeBody) (*uuid.UUID, error)
	GetLegislativeBodyByCode(code int) (*legislativebody.LegislativeBody, error)
	GetLegislativeBodiesByCodes(codes []int) ([]legislativebody.LegislativeBody, error)
}
