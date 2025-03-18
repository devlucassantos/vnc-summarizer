package services

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/labstack/gommon/log"
	"sort"
	"vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/requesters"
)

type Author struct {
	deputyService         services.Deputy
	externalAuthorService services.ExternalAuthor
}

func NewAuthorService(deputyService services.Deputy, externalAuthorService services.ExternalAuthor) *Author {
	return &Author{
		deputyService:         deputyService,
		externalAuthorService: externalAuthorService,
	}
}

func (instance Author) GetAuthorsFromAuthorsUrl(authorsUrl string) ([]deputy.Deputy, []externalauthor.ExternalAuthor,
	error) {
	authors, err := requesters.GetDataSliceFromUrl(authorsUrl)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, nil, err
	}

	deputies, externalAuthors, err := instance.convertAuthorsMapToDeputiesAndExternalAuthors(authors)
	if err != nil {
		log.Error("convertAuthorsMapToDeputiesAndExternalAuthors(): ", err.Error())
		return nil, nil, err
	}

	return deputies, externalAuthors, nil
}

func (instance Author) convertAuthorsMapToDeputiesAndExternalAuthors(authors []map[string]interface{}) ([]deputy.Deputy,
	[]externalauthor.ExternalAuthor, error) {
	var deputies []deputy.Deputy
	var externalAuthors []externalauthor.ExternalAuthor

	sort.Slice(authors, func(i, j int) bool {
		return authors[i]["ordemAssinatura"].(float64) < authors[j]["ordemAssinatura"].(float64)
	})

	for authorIndex, author := range authors {
		authorName := fmt.Sprint(author["nome"])
		authorType := fmt.Sprint(author["tipo"])
		authorTypeCode, err := converters.ToInt(author["codTipo"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, nil, err
		}

		log.Infof("Starting the search for the %dth author: %s - %s", authorIndex+1, authorName, authorType)

		deputyTypeCode := 10000
		if authorTypeCode == deputyTypeCode {
			deputyData, err := instance.deputyService.GetDeputyFromDeputyData(author)
			if err != nil {
				log.Error("deputyService.GetDeputyFromDeputyData(): ", err.Error())
				return nil, nil, err
			}
			deputies = append(deputies, *deputyData)
		} else {
			externalAuthorData, err := instance.externalAuthorService.GetExternalAuthorFromAuthorData(authorName,
				authorTypeCode, authorType)
			if err != nil {
				log.Error("externalAuthorService.GetExternalAuthorFromAuthorData(): ", err.Error())
				return nil, nil, err
			}
			externalAuthors = append(externalAuthors, *externalAuthorData)
		}
		log.Infof("Successful search for %dth author", authorIndex+1)
	}

	return deputies, externalAuthors, nil
}
