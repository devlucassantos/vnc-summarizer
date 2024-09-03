package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type Proposition struct {
	connectionManager connectionManagerInterface
}

func NewPropositionRepository(connectionManager connectionManagerInterface) *Proposition {
	return &Proposition{
		connectionManager: connectionManager,
	}
}

func (instance Proposition) CreateProposition(proposition proposition.Proposition) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro da proposição %d: %s", proposition.Code(),
			err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var propositionId uuid.UUID
	err = transaction.QueryRow(queries.Proposition().Insert(), proposition.Code(), proposition.OriginalTextUrl(),
		proposition.Title(), proposition.Content(), proposition.SubmittedAt(), proposition.ImageUrl(),
		proposition.SpecificType()).Scan(&propositionId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a proposição %d: %s", proposition.Code(), err.Error())
		return err
	}

	for _, deputyData := range proposition.Deputies() {
		var propositionAuthorId uuid.UUID
		deputyParty := deputyData.Party()
		err = transaction.QueryRow(queries.PropositionAuthor().Insert().Deputy(), propositionId, deputyData.Id(),
			deputyParty.Id()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Erro ao cadastrar deputado(a) %s como autor(a) da proposição %d: %s", deputyData.Id(),
				proposition.Code(), err.Error())
			continue
		}

		log.Infof("Deputado(a) %s cadastrado(a) como autor(a) da proposição %d com o ID %s", deputyData.Id(),
			proposition.Code(), propositionAuthorId)
	}

	for _, externalAuthorData := range proposition.ExternalAuthors() {
		var propositionAuthorId uuid.UUID
		err = transaction.QueryRow(queries.PropositionAuthor().Insert().ExternalAuthor(), propositionId,
			externalAuthorData.Id()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Erro ao cadastrar autor externo %s como autor da proposição %d: %s", externalAuthorData.Id(),
				proposition.Code(), err.Error())
			continue
		}

		log.Infof("Autor externo %s cadastrado como autor da proposição %d com o ID %s", externalAuthorData.Id(),
			proposition.Code(), propositionAuthorId)
	}

	var articleId uuid.UUID
	propositionArticle := proposition.Article()
	articleType := propositionArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert().Proposition(), propositionId, articleType.Id()).Scan(&articleId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a proposição %d como matéria: %s", proposition.Code(), err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Erro ao confirmar transação para o cadastro da proposição %d: %s", proposition.Code(),
			err.Error())
		return err
	}

	log.Infof("Proposição %d registrada com sucesso com o ID %s (ID da Matéria: %s)",
		proposition.Code(), propositionId, articleId)
	return nil
}

func (instance Proposition) GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

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
			CreatedAt(propositionDetails.CreatedAt).
			UpdatedAt(propositionDetails.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados da proposição %s: %s", propositionDetails,
				err.Error())
			return nil, err
		}

		propositionData = append(propositionData, *propositionDomain)
	}

	return propositionData, nil
}

func (instance Proposition) GetPropositionsByNewsletterId(newsletterId uuid.UUID) ([]proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterPropositions []dto.Proposition
	err = postgresConnection.Select(&newsletterPropositions, queries.NewsletterProposition().Select().ByNewsletterId(),
		newsletterId)
	if err != nil {
		log.Errorf("Erro ao obter os dados das proposições do boletim %s no banco de dados: %s",
			newsletterId, err.Error())
		return nil, err
	}

	var propositions []proposition.Proposition
	for _, propositionData := range newsletterPropositions {
		propositionDomain, err := proposition.NewBuilder().
			Id(propositionData.Id).
			Code(propositionData.Code).
			OriginalTextUrl(propositionData.OriginalTextUrl).
			Title(propositionData.Title).
			Content(propositionData.Content).
			SubmittedAt(propositionData.SubmittedAt).
			CreatedAt(propositionData.CreatedAt).
			UpdatedAt(propositionData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados da proposição %s do boletim %s: %s",
				propositionData.Id, newsletterId, err.Error())
			continue
		}

		propositions = append(propositions, *propositionDomain)
	}

	return propositions, nil
}

func (instance Proposition) GetLatestPropositionCodes() ([]int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionCodes []int
	err = postgresConnection.Select(&propositionCodes, queries.Proposition().Select().LatestPropositionsCodes())
	if err != nil {
		log.Error("Erro ao obter os últimos códigos das proposições inseridas no banco de dados: ", err.Error())
		return nil, err
	}

	return propositionCodes, nil
}
