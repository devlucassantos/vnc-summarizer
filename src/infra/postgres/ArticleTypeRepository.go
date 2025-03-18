package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type ArticleType struct {
	connectionManager connectionManagerInterface
}

func NewArticleTypeRepository(connectionManager connectionManagerInterface) *ArticleType {
	return &ArticleType{
		connectionManager: connectionManager,
	}
}

func (instance ArticleType) GetArticleTypeByCode(code string) (*articletype.ArticleType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articleType dto.ArticleType
	err = postgresConnection.Get(&articleType, queries.ArticleType().Select().ByCode(), code)
	if err != nil {
		log.Errorf("Error retrieving article type data with code %s from the database: %s", code, err.Error())
		return nil, err
	}

	articleTypeDomain, err := articletype.NewBuilder().
		Id(articleType.Id).
		Description(articleType.Description).
		Codes(articleType.Codes).
		Color(articleType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s: %s", articleType.Id, err.Error())
		return nil, err
	}

	return articleTypeDomain, nil
}
