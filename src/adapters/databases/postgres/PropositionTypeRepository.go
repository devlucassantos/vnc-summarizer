package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type PropositionType struct {
	connectionManager connectionManagerInterface
}

func NewPropositionTypeRepository(connectionManager connectionManagerInterface) *PropositionType {
	return &PropositionType{
		connectionManager: connectionManager,
	}
}

func (instance PropositionType) GetPropositionTypeByCodeOrDefaultType(code string) (*propositiontype.PropositionType,
	error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionType dto.PropositionType
	err = postgresConnection.Get(&propositionType, queries.PropositionType().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = postgresConnection.Get(&propositionType, queries.PropositionType().Select().DefaultOption())
			if err != nil {
				log.Error("Error retrieving the default proposition type data for cases where the proposition type "+
					"code searched was not found in the database: ", err.Error())
				return nil, err
			}
		} else {
			log.Errorf("Error retrieving proposition type data with code %s from the database: %s", code,
				err.Error())
			return nil, err
		}
	}

	propositionTypeDomain, err := propositiontype.NewBuilder().
		Id(propositionType.Id).
		Description(propositionType.Description).
		Codes(propositionType.Codes).
		Color(propositionType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for proposition type %s: %s", propositionType.Id, err.Error())
		return nil, err
	}

	return propositionTypeDomain, nil
}
