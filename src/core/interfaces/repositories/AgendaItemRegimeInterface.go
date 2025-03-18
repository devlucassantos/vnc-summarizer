package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/agendaitemregime"
	"github.com/google/uuid"
)

type AgendaItemRegime interface {
	CreateAgendaItemRegime(agendaItemRegime agendaitemregime.AgendaItemRegime) (*uuid.UUID, error)
	GetAgendaItemRegimeByCode(code int) (*agendaitemregime.AgendaItemRegime, error)
}
