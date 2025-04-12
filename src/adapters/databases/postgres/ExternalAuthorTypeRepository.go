package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type ExternalAuthorType struct {
	connectionManager connectionManagerInterface
}

func NewExternalAuthorTypeRepository(connectionManager connectionManagerInterface) *ExternalAuthorType {
	return &ExternalAuthorType{
		connectionManager: connectionManager,
	}
}

func (instance ExternalAuthorType) CreateExternalAuthorType(externalAuthorType externalauthortype.ExternalAuthorType) (
	*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorTypeId uuid.UUID
	err = postgresConnection.QueryRow(queries.ExternalAuthorType().Insert(), externalAuthorType.Code(),
		externalAuthorType.Description()).Scan(&externalAuthorTypeId)
	if err != nil {
		log.Errorf("Error registering external author type %d: %s", externalAuthorType.Code(), err.Error())
		return nil, err
	}

	log.Infof("External author type %d successfully registered with ID %s", externalAuthorType.Code(),
		externalAuthorTypeId)
	return &externalAuthorTypeId, nil
}

func (instance ExternalAuthorType) GetExternalAuthorTypeByCode(code int) (*externalauthortype.ExternalAuthorType,
	error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorType dto.ExternalAuthorType
	err = postgresConnection.Get(&externalAuthorType, queries.ExternalAuthorType().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("External author type %d not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for external author type %d from the database: %s", code, err.Error())
		return nil, err
	}

	externalAuthorTypeDomain, err := externalauthortype.NewBuilder().
		Id(externalAuthorType.Id).
		Code(externalAuthorType.Code).
		Description(externalAuthorType.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author type %s: %s", externalAuthorType.Id, err.Error())
		return nil, err
	}

	return externalAuthorTypeDomain, nil
}
