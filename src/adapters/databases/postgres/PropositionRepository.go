package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type Proposition struct {
	connectionManager connectionManagerInterface
}

func NewPropositionRepository(connectionManager connectionManagerInterface) *Proposition {
	return &Proposition{
		connectionManager: connectionManager,
	}
}

func (instance Proposition) CreateProposition(proposition proposition.Proposition) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to register the proposition %d: %s", proposition.Code(),
			err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var articleId uuid.UUID
	propositionArticle := proposition.Article()
	articleType := propositionArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert(), articleType.Id(), propositionArticle.ReferenceDateTime()).
		Scan(&articleId)
	if err != nil {
		log.Errorf("Error registering proposition %d as article:  %s", proposition.Code(), err.Error())
		return nil, err
	}

	var propositionImageUrl, imageDescription *string
	if proposition.ImageUrl() != "" {
		imageUrl := proposition.ImageUrl()
		propositionImageUrl = &imageUrl
	}
	if proposition.ImageDescription() != "" {
		description := proposition.ImageDescription()
		imageDescription = &description
	}

	var propositionId uuid.UUID
	propositionType := proposition.Type()
	err = transaction.QueryRow(queries.Proposition().Insert(), proposition.Code(), proposition.OriginalTextUrl(),
		proposition.OriginalTextMimeType(), proposition.Title(), proposition.Content(), proposition.SubmittedAt(),
		propositionImageUrl, imageDescription, proposition.SpecificType(), propositionType.Id(), articleId).Scan(
		&propositionId)
	if err != nil {
		log.Errorf("Error registering the proposition %d: %s", proposition.Code(), err.Error())
		return nil, err
	}

	for _, deputyData := range proposition.Deputies() {
		var propositionAuthorId uuid.UUID
		deputyParty := deputyData.Party()
		err = transaction.QueryRow(queries.PropositionAuthor().Insert().Deputy(), propositionId, deputyData.Id(),
			deputyParty.Id(), deputyData.FederatedUnit()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Error registering deputy %s as the author of the proposition %d: %s", deputyData.Id(),
				proposition.Code(), err.Error())
			return nil, err
		}
		log.Infof("Deputy %s registered as author of the proposition %d with ID %s", deputyData.Id(),
			proposition.Code(), propositionAuthorId)
	}

	for _, externalAuthorData := range proposition.ExternalAuthors() {
		var propositionAuthorId uuid.UUID
		err = transaction.QueryRow(queries.PropositionAuthor().Insert().ExternalAuthor(), propositionId,
			externalAuthorData.Id()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Error registering external author %s as the author of the proposition %d: %s",
				externalAuthorData.Id(), proposition.Code(), err.Error())
			return nil, err
		}
		log.Infof("External author %s registered as author of the proposition %d with ID %s",
			externalAuthorData.Id(), proposition.Code(), propositionAuthorId)
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to register proposition %d: %s", proposition.Code(),
			err.Error())
		return nil, err
	}

	log.Infof("Proposition %d successfully registered with ID %s (Article ID: %s)", proposition.Code(),
		propositionId, articleId)
	return &propositionId, err
}

func (instance Proposition) GetPropositionsByCodes(codes []int) ([]proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionCodes []interface{}
	for _, code := range codes {
		propositionCodes = append(propositionCodes, code)
	}

	var propositions []dto.Proposition
	err = postgresConnection.Select(&propositions, queries.Proposition().Select().ByCodes(len(propositionCodes)),
		propositionCodes...)
	if err != nil {
		log.Error("Error retrieving the proposition data by codes from the database: ", err.Error())
		return nil, err
	}

	var propositionData []proposition.Proposition
	for _, propositionDto := range propositions {
		propositionDomain, err := proposition.NewBuilder().
			Id(propositionDto.Id).
			Code(propositionDto.Code).
			OriginalTextUrl(propositionDto.OriginalTextUrl).
			OriginalTextMimeType(propositionDto.OriginalTextMimeType).
			Title(propositionDto.Title).
			Content(propositionDto.Content).
			SubmittedAt(propositionDto.SubmittedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s: %s", propositionDto.Id, err.Error())
			return nil, err
		}
		propositionData = append(propositionData, *propositionDomain)
	}

	return propositionData, nil
}
