package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
	"vnc-summarizer/utils/datetime"
)

type Newsletter struct {
	connectionManager connectionManagerInterface
}

func NewNewsletterRepository(connectionManager connectionManagerInterface) *Newsletter {
	return &Newsletter{
		connectionManager: connectionManager,
	}
}

func (instance Newsletter) CreateNewsletter(newsletter newsletter.Newsletter) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to register the newsletter of %s: %s", formattedReferenceDate,
			err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	referenceDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	var articleId uuid.UUID
	newsletterArticle := newsletter.Article()
	articleType := newsletterArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert(), articleType.Id(), referenceDateTime).Scan(&articleId)
	if err != nil {
		log.Errorf("Error registering newsletter of %s as article: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	var newsletterId uuid.UUID
	err = transaction.QueryRow(queries.Newsletter().Insert(), newsletter.ReferenceDate(), newsletter.Description(),
		articleId).Scan(&newsletterId)
	if err != nil {
		log.Errorf("Error registering the newsletter of %s: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	for _, articleData := range newsletter.Articles() {
		_, err = transaction.Exec(queries.NewsletterArticle().Insert(), newsletterId, articleData.Id())
		if err != nil {
			log.Errorf("Error registering article %s as part of newsletter %s of %s: %s", articleData.Id(),
				newsletterId, formattedReferenceDate, err.Error())
			return nil, err
		}

		log.Infof("Article %s registered as part of newsletter of %s", articleData.Id(),
			formattedReferenceDate)
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to register newsletter of %s: %s", formattedReferenceDate,
			err.Error())
		return nil, err
	}

	log.Infof("Newsletter of %s successfully registered with ID %s (Article ID: %s))", formattedReferenceDate,
		newsletterId, articleId)
	return &newsletterId, nil
}

func (instance Newsletter) UpdateNewsletter(newsletter newsletter.Newsletter, newArticles []article.Article) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to update newsletter %s of %s: %s", newsletter.Id(),
			formattedReferenceDate, err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	_, err = transaction.Exec(queries.Newsletter().Update(), newsletter.Description(), newsletter.Id())
	if err != nil {
		log.Errorf("Error updating newsletter %s of %s: %s", newsletter.Id(), formattedReferenceDate,
			err.Error())
		return err
	}

	for _, articleData := range newArticles {
		_, err = transaction.Exec(queries.NewsletterArticle().Insert(), newsletter.Id(), articleData.Id())
		if err != nil {
			log.Errorf("Error registering article %s as part of newsletter %s of %s: %s",
				articleData.Id(), newsletter.Id(), formattedReferenceDate, err.Error())
			return err
		}

		log.Infof("Article %s registered as part of newsletter %s of %s", articleData.Id(), newsletter.Id(),
			formattedReferenceDate)
	}

	_, err = transaction.Exec(queries.Article().Update().NewsletterReferenceDateTime(), newsletter.Id())
	if err != nil {
		log.Errorf("Error updating article reference date and time for newsletter %s of %s: %s",
			newsletter.Id(), formattedReferenceDate, err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to update newsletter %s of %s: %s", newsletter.Id(),
			formattedReferenceDate, err.Error())
		return err
	}

	log.Infof("Newsletter %s of %s successfully updated", newsletter.Id(), formattedReferenceDate)
	return nil
}

func (instance Newsletter) GetNewsletterByReferenceDate(referenceDate time.Time) (*newsletter.Newsletter, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterData dto.Newsletter
	err = postgresConnection.Get(&newsletterData, queries.Newsletter().Select().ByReferenceDate(), referenceDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		log.Errorf("Error retrieving data for newsletter of %s from the database: %s",
			referenceDate.Format("02/01/2006"), err.Error())
		return nil, err
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsletterData.Id).
		ReferenceDate(newsletterData.ReferenceDate).
		Description(newsletterData.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for newsletter %s: %s", newsletterData.Id, err.Error())
		return nil, err
	}

	return newsletterDomain, nil
}
