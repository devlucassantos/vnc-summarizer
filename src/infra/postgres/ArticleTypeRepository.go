package postgres

import (
	"database/sql"
	"errors"
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

func (instance ArticleType) GetArticleTypeByCodeOrDefaultType(code string) (*articletype.ArticleType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articleType dto.ArticleType
	err = postgresConnection.Get(&articleType, queries.ArticleType().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = postgresConnection.Get(&articleType, queries.ArticleType().Select().DefaultOption())
			if err != nil {
				log.Error("Error retrieving the default article type data for cases where the article type code "+
					"searched was not found in the database: ", err.Error())
				return nil, err
			}
		} else {
			log.Error("Error retrieving article type data with code %s from the database: ", code, err.Error())
			return nil, err
		}
	}

	articleTypeData, err := articletype.NewBuilder().
		Id(articleType.Id).
		Description(articleType.Description).
		Codes(articleType.Codes).
		Color(articleType.Color).
		SortOrder(articleType.SortOrder).
		CreatedAt(articleType.CreatedAt).
		UpdatedAt(articleType.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s: %s", articleType.Id, err.Error())
		return nil, err
	}

	return articleTypeData, nil
}
