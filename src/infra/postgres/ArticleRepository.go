package postgres

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type Article struct {
	connectionManager connectionManagerInterface
}

func NewArticleRepository(connectionManager connectionManagerInterface) *Article {
	return &Article{
		connectionManager: connectionManager,
	}
}

func (instance Article) GetArticlesByReferenceDate(referenceDate time.Time) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	err = postgresConnection.Select(&articles, queries.Article().Select().ByReferenceDate(), referenceDate)
	if err != nil {
		log.Errorf("Error retrieving articles by reference date %s from the database: %s", referenceDate, err.Error())
		return nil, err
	}

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			Type(*articleType)

		var articleDomain *article.Article
		var articleErr error
		if articleData.Proposition != nil && articleData.Proposition.Id != uuid.Nil {
			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Proposition.PropositionType.Id).
				Description(articleData.Proposition.PropositionType.Description).
				Codes(articleData.Proposition.PropositionType.Codes).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
					articleData.Proposition.PropositionType.Id, articleData.Proposition.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				SpecificType(*articleSpecificType).
				Build()
		} else if articleData.Voting != nil && articleData.Voting.Id != uuid.Nil {
			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Description).
				Build()
		} else {
			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Event.EventType.Id).
				Description(articleData.Event.EventType.Description).
				Codes(articleData.Event.EventType.Codes).
				Build()
			if err != nil {
				log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
					articleData.Event.EventType.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(articleData.Event.Title).
				Content(articleData.Event.Description).
				SpecificType(*articleSpecificType).
				Build()
		}
		if articleErr != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, articleErr.Error())
			return nil, articleErr
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	return articleSlice, nil
}

func (instance Article) GetNewsletterArticlesByNewsletterId(newsletterId uuid.UUID) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	err = postgresConnection.Select(&articles, queries.NewsletterArticle().Select().ByNewsletterId(), newsletterId)
	if err != nil {
		log.Errorf("Error retrieving articles related to newsletter %s from the database: %s", newsletterId,
			err.Error())
		return nil, err
	}

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Codes(articleData.ArticleType.Codes).
			Description(articleData.ArticleType.Description).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s",
				articleData.ArticleType.Id, articleData.Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder()

		if articleData.Proposition.Id != uuid.Nil {
			articleBuilder.Title(articleData.Proposition.Title).Content(articleData.Proposition.Content)
		} else if articleData.Voting.Id != uuid.Nil {
			articleBuilder.Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Description)
		} else if articleData.Event.Id != uuid.Nil {
			articleBuilder.Title(articleData.Event.Title).Content(articleData.Event.Description)
		}

		articleDomain, err := articleBuilder.
			Id(articleData.Id).
			Type(*articleType).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, err.Error())
			return nil, err
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	return articleSlice, nil
}
