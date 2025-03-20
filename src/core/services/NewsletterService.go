package services

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/labstack/gommon/log"
	"math"
	"strings"
	"time"
	"vnc-summarizer/core/interfaces/repositories"
)

type Newsletter struct {
	newsletterRepository  repositories.Newsletter
	articleTypeRepository repositories.ArticleType
	articleRepository     repositories.Article
}

func NewNewsletterService(newsletterRepository repositories.Newsletter, articleTypeRepository repositories.ArticleType,
	articleRepository repositories.Article) *Newsletter {
	return &Newsletter{
		newsletterRepository:  newsletterRepository,
		articleTypeRepository: articleTypeRepository,
		articleRepository:     articleRepository,
	}
}

func (instance Newsletter) RegisterNewNewsletter(referenceDate time.Time) {
	formattedReferenceDate := referenceDate.Format("02/01/2006")

	articles, err := instance.articleRepository.GetArticlesByReferenceDate(referenceDate)
	if err != nil {
		log.Error("articleRepository.GetArticlesByReferenceDate(): ", err.Error())
		return
	} else if articles == nil {
		log.Infof("No articles were found on %s to generate a new newsletter", formattedReferenceDate)
		return
	}

	var articlesOutsideTheNewsletter []article.Article
	registeredNewsletter, err := instance.newsletterRepository.GetNewsletterByReferenceDate(referenceDate)
	if err != nil {
		log.Error("newsletterRepository.GetNewsletterByReferenceDate(): ", err.Error())
		return
	} else if registeredNewsletter != nil {
		newsletterArticles, err := instance.articleRepository.GetNewsletterArticlesByNewsletterId(
			registeredNewsletter.Id())
		if err != nil {
			log.Error("articleRepository.GetNewsletterArticlesByNewsletterId(): ", err.Error())
			return
		}

		for _, articleData := range articles {
			var isInTheNewsletter bool
			for _, newsletterArticleData := range newsletterArticles {
				if newsletterArticleData.Id() == articleData.Id() {
					isInTheNewsletter = true
					break
				}
			}

			if !isInTheNewsletter {
				articlesOutsideTheNewsletter = append(articlesOutsideTheNewsletter, articleData)
			}
		}
	} else {
		articlesOutsideTheNewsletter = articles
	}

	if articlesOutsideTheNewsletter == nil {
		log.Infof("No new articles were found on %s to update newsletter %s", formattedReferenceDate,
			registeredNewsletter.Id())
		return
	}

	if registeredNewsletter == nil {
		log.Info("Starting generation of the newsletter of ", formattedReferenceDate)
	} else {
		log.Infof("Starting update of the newsletter %s of %s", registeredNewsletter.Id(),
			formattedReferenceDate)
	}

	newsletterData, err := instance.generateNewsletter(articles, referenceDate)
	if err != nil {
		for attempt := 1; attempt <= 3; attempt++ {
			waitingTimeInSeconds := int(math.Pow(5, float64(attempt)))
			log.Warnf("It was not possible to register newsletter of %s on the %dth attempt, trying again in "+
				"%d seconds", formattedReferenceDate, attempt, waitingTimeInSeconds)
			time.Sleep(time.Duration(waitingTimeInSeconds) * time.Second)
			newsletterData, err = instance.generateNewsletter(articles, referenceDate)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.Errorf("Error generating the newsletter of %s: %s", formattedReferenceDate, err.Error())
		return
	}

	if registeredNewsletter != nil {
		newsletterData, err = newsletterData.NewUpdater().Id(registeredNewsletter.Id()).Build()
		if err != nil {
			log.Errorf("Error validating data for newsletter of %s: %s", formattedReferenceDate, err.Error())
			return
		}
	}

	if registeredNewsletter == nil {
		_, err = instance.newsletterRepository.CreateNewsletter(*newsletterData)
		if err != nil {
			log.Error("newsletterRepository.CreateNewsletter(): ", err.Error())
		}
	} else {
		err = instance.newsletterRepository.UpdateNewsletter(*newsletterData, articlesOutsideTheNewsletter)
		if err != nil {
			log.Error("newsletterRepository.UpdateNewsletter(): ", err.Error())
		}
	}

	log.Infof("Newsletter of %s successfully generated", formattedReferenceDate)
	return
}

func (instance Newsletter) generateNewsletter(articles []article.Article, referenceDate time.Time) (
	*newsletter.Newsletter, error) {
	formattedReferenceDate := referenceDate.Format("02/01/2006")

	maximumNumberOfRelevantArticles := 10
	mostRelevantArticles := getMostRelevantArticles(articles, maximumNumberOfRelevantArticles)

	var contentOfArticles string
	for count, articleData := range mostRelevantArticles {
		contentOfArticles += fmt.Sprintf("%dª matéria:\nTítulo: %s\n\nConteúdo: %s\n\n", count+1,
			articleData.Title(), articleData.Content())
	}

	chatGptCommand := "Gere uma descrição para ser usada em um boletim sobre o conjunto de matérias políticas abaixo. " +
		"É importante que a descrição seja curta e chamativa, falando sobre o máximo de matérias possíveis, " +
		"correlacionando os temas e utilizando uma linguagem simples e direta. Não deve ter mais do que 500 " +
		"caracteres e a frequência em que o boletim é disponibilizado não precisa ser mencionada. Matérias:\n\n"
	purpose := fmt.Sprint("Generating the newsletter description of ", formattedReferenceDate)
	description, err := requestToChatGpt(chatGptCommand, contentOfArticles, purpose)
	if err != nil {
		log.Error("requestToChatGpt(): ", err.Error())
		return nil, err
	}

	articleTypeCode := "newsletter"
	articleType, err := instance.articleTypeRepository.GetArticleTypeByCode(articleTypeCode)
	if err != nil {
		log.Error("articleTypeRepository.GetArticleTypeByCode(): ", err.Error())
		return nil, err
	}

	articleData, err := article.NewBuilder().Type(*articleType).Build()
	if err != nil {
		log.Errorf("Error validating article data for newsletter of %s: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	newsletterData, err := newsletter.NewBuilder().
		ReferenceDate(referenceDate).
		Description(description).
		Article(*articleData).
		Articles(articles).
		Build()
	if err != nil {
		log.Errorf("Error validating newsletter data of %s: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	return newsletterData, nil
}

func getMostRelevantArticles(articles []article.Article, maximumNumberOfRelevantArticles int) []article.Article {
	var mostRelevantArticles []article.Article
	if len(articles) > maximumNumberOfRelevantArticles {
		for _, articleData := range articles {
			if len(mostRelevantArticles) >= maximumNumberOfRelevantArticles {
				break
			}

			articleSpecificType := articleData.SpecificType()
			if articleSpecificType.IsZero() || !strings.Contains(articleSpecificType.Codes(), "default_option") {
				mostRelevantArticles = append(mostRelevantArticles, articleData)
			}
		}

		if len(mostRelevantArticles) < maximumNumberOfRelevantArticles {
			for _, articleData := range articles {
				if len(mostRelevantArticles) >= maximumNumberOfRelevantArticles {
					break
				}

				articleSpecificType := articleData.SpecificType()
				if !articleSpecificType.IsZero() &&
					strings.Contains(articleSpecificType.Codes(), "default_option") {
					mostRelevantArticles = append(mostRelevantArticles, articleData)
				}
			}
		}
	} else {
		mostRelevantArticles = articles
	}

	return mostRelevantArticles
}
