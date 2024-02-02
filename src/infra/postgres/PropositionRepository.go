package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-write-api/infra/dto"
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

func (instance Proposition) CreateProposition(proposition proposition.Proposition) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro da proposição %d: %s", proposition.Code(),
			err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var propositionId uuid.UUID
	err = transaction.QueryRow(queries.Proposition().Insert(), proposition.Code(), proposition.OriginalTextUrl(),
		proposition.Title(), proposition.Content(), proposition.SubmittedAt()).Scan(&propositionId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a proposição %d: %s", proposition.Code(), err.Error())
		return err
	}

	for _, deputyData := range proposition.Deputies() {
		var propositionAuthorId uuid.UUID
		deputyParty := deputyData.Party()
		err = transaction.QueryRow(queries.PropositionAuthor().InsertDeputy(), propositionId, deputyData.Id(),
			deputyParty.Id()).Scan(&propositionAuthorId)
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
		err = transaction.QueryRow(queries.PropositionAuthor().InsertOrganization(), propositionId,
			organizationData.Id()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Erro ao cadastrar organização %s como autora da proposição %d: %s", organizationData.Id(),
				proposition.Code(), err.Error())
			continue
		}

		log.Infof("Organização %s cadastrada como autora da proposição %d com o ID %s", organizationData.Id(),
			proposition.Code(), propositionAuthorId)
	}

	var newsId uuid.UUID
	err = transaction.QueryRow(queries.News().InsertProposition(), propositionId).Scan(&newsId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a proposição %d como matéria: %s", proposition.Code(), err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Error("Erro ao confirmar transação para o cadastro da proposição %d: %s", proposition.Code(),
			err.Error())
		return err
	}

	log.Infof("Proposição %d registrada com sucesso com o ID %s (ID da Matéria: %s)",
		proposition.Code(), propositionId, newsId)
	return nil
}

func (instance Proposition) GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var propositions []dto.Proposition
	err = postgresConnection.Select(&propositions, queries.Proposition().Select().ByDate(), date)
	if err != nil {
		log.Error("Erro ao obter os dados das proposições por data no banco de dados: ", err.Error())
		return nil, err
	}

	var propositionData []proposition.Proposition
	for _, propositionDetails := range propositions {
		propositionDomain, err := proposition.NewBuilder().
			Id(propositionDetails.Id).
			Code(propositionDetails.Code).
			OriginalTextUrl(propositionDetails.OriginalTextUrl).
			Title(propositionDetails.Title).
			Content(propositionDetails.Content).
			SubmittedAt(propositionDetails.SubmittedAt).
			Active(propositionDetails.Active).
			CreatedAt(propositionDetails.CreatedAt).
			UpdatedAt(propositionDetails.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro durante a construção da estrutura de dados da proposição %s: %s", propositionDetails,
				err.Error())
			return nil, err
		}

		propositionData = append(propositionData, *propositionDomain)
	}

	return propositionData, nil
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
