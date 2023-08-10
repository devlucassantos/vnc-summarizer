package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-write-api/core/domains/keyword"
	"vnc-write-api/infra/dto"
	"vnc-write-api/infra/postgres/queries"
)

type Keyword struct {
	connectionManager ConnectionManagerInterface
}

func NewKeywordRepository(connectionManager ConnectionManagerInterface) *Keyword {
	return &Keyword{
		connectionManager: connectionManager,
	}
}

func (instance Keyword) CreateKeyword(keyword keyword.Keyword) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var keywordId uuid.UUID
	err = postgresConnection.QueryRow(queries.Keyword().Insert(), keyword.Keyword()).Scan(&keywordId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a palavra-chave %s: %s", keyword.Keyword(), err.Error())
		return nil, err
	}

	log.Infof("Palavra-chave %s registrada com sucesso com o ID %s", keyword.Keyword(), keywordId)
	return &keywordId, nil
}

func (instance Keyword) GetKeywordByKeyword(key string) (*keyword.Keyword, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var keywordData dto.Keyword
	err = postgresConnection.Get(&keywordData, queries.Keyword().Select().ByKeyword(), key)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("Palavra-chave %s não encontrada no banco de dados", key)
			return nil, nil
		}
		log.Errorf("Erro ao obter os dados da palavra-chave %s no banco de dados: %s", key, err.Error())
		return nil, err
	}

	keywordDomain, err := keyword.NewBuilder().
		Id(keywordData.Id).
		Keyword(keywordData.Keyword).
		Active(keywordData.Active).
		CreatedAt(keywordData.CreatedAt).
		UpdatedAt(keywordData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados da palavra-chave %s: %s", keywordData.Keyword, err.Error())
		return nil, err
	}

	return keywordDomain, nil
}
