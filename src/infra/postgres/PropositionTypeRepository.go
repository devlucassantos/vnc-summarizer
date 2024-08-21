package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type PropositionType struct {
	connectionManager connectionManagerInterface
}

func NewPropositionTypeRepository(connectionManager connectionManagerInterface) *PropositionType {
	return &PropositionType{
		connectionManager: connectionManager,
	}
}

func (instance PropositionType) GetPropositionTypeByCode(code string) (*proptype.PropositionType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionType dto.PropositionType
	err = postgresConnection.Get(&propositionType, queries.PropositionType().Select().ByCode(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			err = postgresConnection.Get(&propositionType, queries.PropositionType().Select().OtherOption())
			if err != nil {
				log.Error("Erro ao obter os dados do tipo de proposição 'Outras Proposições' no banco de dados: ", err.Error())
				return nil, err
			}
		} else {
			log.Error("Erro ao obter os dados dos tipos de proposição no banco de dados: ", err.Error())
			return nil, err
		}
	}

	propositionTypeData, err := proptype.NewBuilder().
		Id(propositionType.Id).
		Description(propositionType.Description).
		Codes(propositionType.Codes).
		CreatedAt(propositionType.CreatedAt).
		UpdatedAt(propositionType.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do tipo de proposição %s: %s", propositionType.Id,
			err.Error())
		return nil, err
	}

	return propositionTypeData, nil
}
