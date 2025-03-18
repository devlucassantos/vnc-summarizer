package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type EventSituation struct {
	connectionManager connectionManagerInterface
}

func NewEventSituationRepository(connectionManager connectionManagerInterface) *EventSituation {
	return &EventSituation{
		connectionManager: connectionManager,
	}
}

func (instance EventSituation) GetEventSituationByCodeOrDefaultSituation(code string) (*eventsituation.EventSituation,
	error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var eventSituation dto.EventSituation
	err = postgresConnection.Get(&eventSituation, queries.EventSituation().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = postgresConnection.Get(&eventSituation, queries.EventSituation().Select().DefaultOption())
			if err != nil {
				log.Error("Error retrieving the default event situation data for cases where the event situation "+
					"code searched was not found in the database: ", err.Error())
				return nil, err
			}
		} else {
			log.Errorf("Error retrieving event situation data with code %s from the database: %s", code,
				err.Error())
			return nil, err
		}
	}

	eventSituationDomain, err := eventsituation.NewBuilder().
		Id(eventSituation.Id).
		Description(eventSituation.Description).
		Codes(eventSituation.Codes).
		Color(eventSituation.Color).
		IsFinished(eventSituation.IsFinished).
		Build()
	if err != nil {
		log.Errorf("Error validating data for event situation %s: %s", eventSituation.Id, err.Error())
		return nil, err
	}

	return eventSituationDomain, nil
}
