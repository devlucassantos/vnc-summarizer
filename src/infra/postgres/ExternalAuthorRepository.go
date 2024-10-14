package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type ExternalAuthor struct {
	connectionManager connectionManagerInterface
}

func NewExternalAuthorRepository(connectionManager connectionManagerInterface) *ExternalAuthor {
	return &ExternalAuthor{
		connectionManager: connectionManager,
	}
}

func (instance ExternalAuthor) CreateExternalAuthor(externalAuthor external.ExternalAuthor) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorId uuid.UUID
	err = postgresConnection.QueryRow(queries.ExternalAuthor().Insert(), externalAuthor.Name(),
		externalAuthor.Type()).Scan(&externalAuthorId)
	if err != nil {
		log.Errorf("Error registering external author %s - %s: %s", externalAuthor.Name(),
			externalAuthor.Type(), err.Error())
		return nil, err
	}

	log.Infof("External author %s - %s successfully registered with ID %s", externalAuthor.Name(),
		externalAuthor.Type(), externalAuthorId)
	return &externalAuthorId, nil
}

func (instance ExternalAuthor) GetExternalAuthorByNameAndType(name string, _type string) (*external.ExternalAuthor, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorData dto.ExternalAuthor
	err = postgresConnection.Get(&externalAuthorData, queries.ExternalAuthor().Select().ByNameAndType(), name, _type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("External author %s - %s not found in database", name, _type)
			return nil, nil
		}
		log.Errorf("Error retrieving data for external author %s - %s from the database: %s", name, _type, err.Error())
		return nil, err
	}

	externalAuthorDomain, err := external.NewBuilder().
		Id(externalAuthorData.Id).
		Name(externalAuthorData.Name).
		Type(externalAuthorData.Type).
		CreatedAt(externalAuthorData.CreatedAt).
		UpdatedAt(externalAuthorData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author %s - %s: %s", name, _type, err.Error())
		return nil, err
	}

	return externalAuthorDomain, nil
}
