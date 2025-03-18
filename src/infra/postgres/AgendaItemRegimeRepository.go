package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/agendaitemregime"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type AgendaItemRegime struct {
	connectionManager connectionManagerInterface
}

func NewAgendaItemRegimeRepository(connectionManager connectionManagerInterface) *AgendaItemRegime {
	return &AgendaItemRegime{
		connectionManager: connectionManager,
	}
}

func (instance AgendaItemRegime) CreateAgendaItemRegime(agendaItemRegime agendaitemregime.AgendaItemRegime) (*uuid.UUID,
	error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var agendaItemRegimeId uuid.UUID
	err = postgresConnection.QueryRow(queries.AgendaItemRegime().Insert(), agendaItemRegime.Code(),
		agendaItemRegime.Description()).Scan(&agendaItemRegimeId)
	if err != nil {
		log.Errorf("Error registering agenda item regime %d: %s", agendaItemRegime.Code(), err.Error())
		return nil, err
	}

	log.Infof("Agenda item regime %d successfully registered with ID %s", agendaItemRegime.Code(),
		agendaItemRegimeId)
	return &agendaItemRegimeId, nil
}

func (instance AgendaItemRegime) GetAgendaItemRegimeByCode(code int) (*agendaitemregime.AgendaItemRegime, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var agendaItemRegime dto.AgendaItemRegime
	err = postgresConnection.Get(&agendaItemRegime, queries.AgendaItemRegime().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Agenda item regime %d not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for agenda item regime %d from the database: %s", code, err.Error())
		return nil, err
	}

	agendaItemRegimeDomain, err := agendaitemregime.NewBuilder().
		Id(agendaItemRegime.Id).
		Code(agendaItemRegime.Code).
		Description(agendaItemRegime.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for agenda item regime %d: %s", agendaItemRegime.Id, err.Error())
		return nil, err
	}

	return agendaItemRegimeDomain, nil
}
