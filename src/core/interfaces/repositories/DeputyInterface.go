package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/deputy"
)

type Deputy interface {
	CreateDeputy(deputy deputy.Deputy) (*uuid.UUID, error)
	UpdateDeputy(deputy deputy.Deputy) error
	GetDeputyByCode(code int) (*deputy.Deputy, error)
}
