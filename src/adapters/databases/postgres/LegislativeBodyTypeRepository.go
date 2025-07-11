package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type LegislativeBodyType struct {
	connectionManager connectionManagerInterface
}

func NewLegislativeBodyTypeRepository(connectionManager connectionManagerInterface) *LegislativeBodyType {
	return &LegislativeBodyType{
		connectionManager: connectionManager,
	}
}

func (instance LegislativeBodyType) CreateLegislativeBodyType(legislativeBodyType legislativebodytype.LegislativeBodyType) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var legislativeBodyTypeId uuid.UUID
	err = postgresConnection.QueryRow(queries.LegislativeBodyType().Insert(), legislativeBodyType.Code(),
		legislativeBodyType.Description()).Scan(&legislativeBodyTypeId)
	if err != nil {
		log.Errorf("Error registering legislative body type %d: %s", legislativeBodyType.Code(), err.Error())
		return nil, err
	}

	log.Infof("Legislative body type %d successfully registered with ID %s", legislativeBodyType.Code(),
		legislativeBodyTypeId)
	return &legislativeBodyTypeId, nil
}

func (instance LegislativeBodyType) GetLegislativeBodyTypeByCode(code int) (*legislativebodytype.LegislativeBodyType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var legislativeBodyType dto.LegislativeBodyType
	err = postgresConnection.Get(&legislativeBodyType, queries.LegislativeBodyType().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Legislative body type %d not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for legislative body type %d from the database: %s", code, err.Error())
		return nil, err
	}

	legislativeBodyTypeDomain, err := legislativebodytype.NewBuilder().
		Id(legislativeBodyType.Id).
		Code(legislativeBodyType.Code).
		Description(legislativeBodyType.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for legislative body type %d: %s", legislativeBodyType.Id, err.Error())
		return nil, err
	}

	return legislativeBodyTypeDomain, nil
}
