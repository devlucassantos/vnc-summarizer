package postgres

import (
	"database/sql"
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
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articleType dto.ArticleType
	err = postgresConnection.Get(&articleType, queries.ArticleType().Select().ByCode(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			err = postgresConnection.Get(&articleType, queries.ArticleType().Select().DefaultOption())
			if err != nil {
				log.Error("Erro ao obter os dados do tipo de matéria padrão para casos onde o código do tipo de "+
					"matéria procurado não foi encontrado no banco de dados: ", err.Error())
				return nil, err
			}
		} else {
			log.Error("Erro ao obter os dados do tipo de matéria com código %s no banco de dados: ", code, err.Error())
			return nil, err
		}
	}

	articleTypeData, err := articletype.NewBuilder().
		Id(articleType.Id).
		Description(articleType.Description).
		Color(articleType.Color).
		SortOrder(articleType.SortOrder).
		CreatedAt(articleType.CreatedAt).
		UpdatedAt(articleType.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do tipo de matéria %s: %s", articleType.Id,
			err.Error())
		return nil, err
	}

	return articleTypeData, nil
}
