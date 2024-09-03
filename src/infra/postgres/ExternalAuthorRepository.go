package postgres

import (
	"database/sql"
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
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorId uuid.UUID
	err = postgresConnection.QueryRow(queries.ExternalAuthor().Insert(), externalAuthor.Name(),
		externalAuthor.Type()).Scan(&externalAuthorId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o autor externo %s: %s", externalAuthor.Name(), err.Error())
		return nil, err
	}

	log.Infof("Autor externo %s registrado com sucesso com o ID %s", externalAuthor.Name(), externalAuthorId)
	return &externalAuthorId, nil
}

func (instance ExternalAuthor) GetExternalAuthorByNameAndType(name string, _type string) (*external.ExternalAuthor, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorData dto.ExternalAuthor
	err = postgresConnection.Get(&externalAuthorData, queries.ExternalAuthor().Select().ByNameAndType(), name, _type)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("Autor externo %s não encontrado no banco de dados", name)
			return nil, nil
		}
		log.Errorf("Erro ao obter os dados do autor externo %s no banco de dados: %s", name, err.Error())
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
		log.Errorf("Erro ao validar os dados do autor externo %s: %s", name, err.Error())
		return nil, err
	}

	return externalAuthorDomain, nil
}
