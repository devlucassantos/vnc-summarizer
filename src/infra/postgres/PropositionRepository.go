package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
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
		log.Error("connectionManager.createConnection(): ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to register the proposition %d: %s", proposition.Code(), err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var propositionImageUrl *string
	if proposition.ImageUrl() != "" {
		imageUrl := proposition.ImageUrl()
		propositionImageUrl = &imageUrl
	}

	var propositionId uuid.UUID
	err = transaction.QueryRow(queries.Proposition().Insert(), proposition.Code(), proposition.OriginalTextUrl(),
		proposition.Title(), proposition.Content(), proposition.SubmittedAt(), propositionImageUrl,
		proposition.SpecificType()).Scan(&propositionId)
	if err != nil {
		log.Errorf("Error registering the proposition %d: %s", proposition.Code(), err.Error())
		return err
	}

	for _, deputyData := range proposition.Deputies() {
		var propositionAuthorId uuid.UUID
		deputyParty := deputyData.Party()
		err = transaction.QueryRow(queries.PropositionAuthor().Insert().Deputy(), propositionId, deputyData.Id(),
			deputyParty.Id()).Scan(&propositionAuthorId)
		if err != nil {
			log.Errorf("Error registering deputy %s as the author of the proposition %d: %s", deputyData.Id(),
				proposition.Code(), err.Error())
			continue
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
			continue
		}

		log.Infof("External author %s registered as author of the proposition %d with ID %s",
			externalAuthorData.Id(), proposition.Code(), propositionAuthorId)
	}

	var articleId uuid.UUID
	propositionArticle := proposition.Article()
	articleType := propositionArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert().Proposition(), propositionId, articleType.Id()).Scan(&articleId)
	if err != nil {
		log.Errorf("Error registering proposition %d as article:  %s", proposition.Code(), err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to register proposition %d: %s", proposition.Code(),
			err.Error())
		return err
	}

	log.Infof("Proposition %d successfully registered with ID %s (Article ID: %s)", proposition.Code(),
		propositionId, articleId)
	return nil
}

func (instance Proposition) GetPropositionsByDate(date time.Time) ([]proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionArticles []dto.Article
	err = postgresConnection.Select(&propositionArticles, queries.Proposition().Select().ByDate(), date)
	if err != nil {
		log.Error("Error retrieving the proposition data by date from the database: ", err.Error())
		return nil, err
	}

	var propositionData []proposition.Proposition
	for _, propositionArticle := range propositionArticles {
		articleType, err := articletype.NewBuilder().
			Id(propositionArticle.ArticleType.Id).
			Description(propositionArticle.ArticleType.Description).
			Codes(propositionArticle.ArticleType.Codes).
			Color(propositionArticle.ArticleType.Color).
			SortOrder(propositionArticle.ArticleType.SortOrder).
			CreatedAt(propositionArticle.ArticleType.CreatedAt).
			UpdatedAt(propositionArticle.ArticleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s",
				propositionArticle.ArticleType.Id, propositionArticle.Id, err.Error())
			return nil, err
		}

		articleData, err := article.NewBuilder().Type(*articleType).Build()
		if err != nil {
			log.Errorf("Error validating data for article %s of proprosition %s: %s", propositionArticle.Id,
				propositionArticle.Proposition.Id, err.Error())
			return nil, err
		}

		propositionDomain, err := proposition.NewBuilder().
			Id(propositionArticle.Proposition.Id).
			Code(propositionArticle.Proposition.Code).
			OriginalTextUrl(propositionArticle.Proposition.OriginalTextUrl).
			Title(propositionArticle.Proposition.Title).
			Content(propositionArticle.Proposition.Content).
			SubmittedAt(propositionArticle.Proposition.SubmittedAt).
			Article(*articleData).
			CreatedAt(propositionArticle.Proposition.CreatedAt).
			UpdatedAt(propositionArticle.Proposition.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s: %s", propositionArticle.Proposition.Id,
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
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterPropositions []dto.Proposition
	err = postgresConnection.Select(&newsletterPropositions, queries.NewsletterProposition().Select().ByNewsletterId(),
		newsletterId)
	if err != nil {
		log.Errorf("Error retrieving the proposition data from newsletter %s from the database: %s",
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
			log.Errorf("Error validating data for proposition %s of newsletter %s: %s", propositionData.Id,
				newsletterId, err.Error())
			continue
		}

		propositions = append(propositions, *propositionDomain)
	}

	return propositions, nil
}

func (instance Proposition) GetLatestPropositionCodes() ([]int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionCodes []int
	err = postgresConnection.Select(&propositionCodes, queries.Proposition().Select().LatestPropositionsCodes())
	if err != nil {
		log.Error("Error obtaining the codes of the last propositions inserted into the database: ", err.Error())
		return nil, err
	}

	return propositionCodes, nil
}
