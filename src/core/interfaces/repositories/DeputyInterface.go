package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/google/uuid"
)

type Deputy interface {
	CreateDeputy(deputy.Deputy) (*uuid.UUID, error)
	UpdateDeputy(deputy deputy.Deputy) error
	GetDeputyByCode(code int) (*deputy.Deputy, error)
}
