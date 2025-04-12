package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type ExternalAuthor struct {
	connectionManager connectionManagerInterface
}

func NewExternalAuthorRepository(connectionManager connectionManagerInterface) *ExternalAuthor {
	return &ExternalAuthor{
		connectionManager: connectionManager,
	}
}

func (instance ExternalAuthor) CreateExternalAuthor(externalAuthor externalauthor.ExternalAuthor) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorId uuid.UUID
	externalAuthorType := externalAuthor.Type()
	err = postgresConnection.QueryRow(queries.ExternalAuthor().Insert(), externalAuthor.Name(),
		externalAuthorType.Id()).Scan(&externalAuthorId)
	if err != nil {
		log.Errorf("Error registering external author %s (Type code: %d): %s", externalAuthor.Name(),
			externalAuthorType.Code(), err.Error())
		return nil, err
	}

	log.Infof("External author %s (Type code: %d) successfully registered with ID %s", externalAuthor.Name(),
		externalAuthorType.Code(), externalAuthorId)
	return &externalAuthorId, nil
}

func (instance ExternalAuthor) GetExternalAuthorByNameAndTypeCode(name string, typeCode int) (
	*externalauthor.ExternalAuthor, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthor dto.ExternalAuthor
	err = postgresConnection.Get(&externalAuthor, queries.ExternalAuthor().Select().ByNameAndTypeCode(), name, typeCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("External author %s (Type code: %d) not found in database", name, typeCode)
			return nil, nil
		}
		log.Errorf("Error retrieving data for external author %s (Type code: %d) from the database: %s", name,
			typeCode, err.Error())
		return nil, err
	}

	externalAuthorType, err := externalauthortype.NewBuilder().
		Id(externalAuthor.ExternalAuthorType.Id).
		Code(externalAuthor.ExternalAuthorType.Code).
		Description(externalAuthor.ExternalAuthorType.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author type %s: %s",
			externalAuthor.ExternalAuthorType.Id, err.Error())
		return nil, err
	}

	externalAuthorDomain, err := externalauthor.NewBuilder().
		Id(externalAuthor.Id).
		Name(externalAuthor.Name).
		Type(*externalAuthorType).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author %s: %s", externalAuthor.Id, err.Error())
		return nil, err
	}

	return externalAuthorDomain, nil
}
