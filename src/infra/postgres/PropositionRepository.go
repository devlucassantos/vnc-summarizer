package postgres

import (
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-write-api/core/domains/proposition"
	"vnc-write-api/infra/postgres/queries"
)

type Proposition struct {
	connectionManager ConnectionManagerInterface
}

func NewPropositionRepository(connectionManager ConnectionManagerInterface) *Proposition {
	return &Proposition{
		connectionManager: connectionManager,
	}
}

func (instance Proposition) CreateProposition(proposition proposition.Proposition) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro da proposição %d: %s", proposition.Code(), err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var propositionId uuid.UUID
	err = transaction.QueryRow(queries.Proposition().Insert(), proposition.Code(), proposition.OriginalTextUrl(),
		proposition.Title(), proposition.Summary(), proposition.SubmittedAt()).Scan(&propositionId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a proposição %d: %s", proposition.Code(), err.Error())
		return nil, err
	}

	for _, deputyData := range proposition.Deputies() {
		var propositionAuthorId uuid.UUID
		currentParty := deputyData.CurrentParty()
		err = transaction.QueryRow(queries.PropositionAuthor().InsertDeputy(), propositionId, deputyData.Id(),
			currentParty.Id()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Erro ao cadastrar deputado(a) %s como autor(a) da proposição %d: %s", deputyData.Id(),
				proposition.Code(), err.Error())
			continue
		}

		log.Infof("Deputado(a) %s cadastrado(a) como autor(a) da proposição %d com o ID %s", deputyData.Id(),
			proposition.Code(), propositionAuthorId)
	}

	for _, organizationData := range proposition.Organizations() {
		var propositionAuthorId uuid.UUID
		err = transaction.QueryRow(queries.PropositionAuthor().InsertOrganization(), propositionId, organizationData.Id()).
			Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Erro ao cadastrar organização %s como autora da proposição %d: %s", organizationData.Id(),
				proposition.Code(), err.Error())
			continue
		}

		log.Infof("Organização %s cadastrada como autora da proposição %d com o ID %s", organizationData.Id(),
			proposition.Code(), propositionAuthorId)
	}

	for _, keywordData := range proposition.Keywords() {
		err = transaction.QueryRow(queries.PropositionKeyword().Insert(), propositionId, keywordData.Id()).Err()
		if err != nil {
			log.Errorf("Erro ao cadastrar palavra-chave %s para a proposição %d: %s", keywordData.Id(),
				proposition.Code(), err.Error())
			continue
		}

		log.Infof("Palavra-chave %s cadastrada para a proposição %d", keywordData.Id(), proposition.Code())
	}

	err = transaction.Commit()
	if err != nil {
		log.Error("Erro ao confirmar transação para o cadastro da proposição %d: %s", proposition.Code(), err.Error())
		return nil, err
	}

	log.Infof("Proposição %d registrada com sucesso com o ID %s", proposition.Code(), propositionId)
	return &propositionId, nil
}

func (instance Proposition) GetLatestPropositionCodes() ([]int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var propositionCodes []int
	err = postgresConnection.Select(&propositionCodes, queries.Proposition().Select().LatestPropositionsCodes())
	if err != nil {
		log.Error("Erro ao obter os últimos códigos das proposições inseridas no banco de dados: ", err.Error())
		return nil, err
	}

	return propositionCodes, nil
}
