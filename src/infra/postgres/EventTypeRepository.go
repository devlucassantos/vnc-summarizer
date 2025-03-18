package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type EventType struct {
	connectionManager connectionManagerInterface
}

func NewEventTypeRepository(connectionManager connectionManagerInterface) *EventType {
	return &EventType{
		connectionManager: connectionManager,
	}
}

func (instance EventType) GetEventTypeByCodeOrDefaultType(code string) (*eventtype.EventType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var eventType dto.EventType
	err = postgresConnection.Get(&eventType, queries.EventType().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = postgresConnection.Get(&eventType, queries.EventType().Select().DefaultOption())
			if err != nil {
				log.Error("Error retrieving the default event type data for cases where the event type code "+
					"searched was not found in the database: ", err.Error())
				return nil, err
			}
		} else {
			log.Errorf("Error retrieving event type data with code %s from the database: %s", code, err.Error())
			return nil, err
		}
	}

	eventTypeDomain, err := eventtype.NewBuilder().
		Id(eventType.Id).
		Description(eventType.Description).
		Codes(eventType.Codes).
		Color(eventType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for event type %s: %s", eventType.Id, err.Error())
		return nil, err
	}

	return eventTypeDomain, nil
}
