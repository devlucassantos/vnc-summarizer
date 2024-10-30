package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"vnc-summarizer/core/interfaces/repositories"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/requests"
	"vnc-summarizer/core/services/utils/validators"
)

type BackgroundData struct {
	deputyRepository         repositories.Deputy
	externalAuthorRepository repositories.ExternalAuthor
	partyRepository          repositories.Party
	propositionRepository    repositories.Proposition
	newsletterRepository     repositories.Newsletter
	articleTypeRepository    repositories.ArticleType
}

func NewBackgroundDataService(deputyRepository repositories.Deputy, externalAuthorRepository repositories.ExternalAuthor,
	partyRepository repositories.Party, propositionRepository repositories.Proposition,
	newsletterRepository repositories.Newsletter, articleTypeRepository repositories.ArticleType) *BackgroundData {
	return &BackgroundData{
		deputyRepository:         deputyRepository,
		externalAuthorRepository: externalAuthorRepository,
		partyRepository:          partyRepository,
		propositionRepository:    propositionRepository,
		newsletterRepository:     newsletterRepository,
		articleTypeRepository:    articleTypeRepository,
	}
}

func (instance BackgroundData) RegisterNewPropositions() {
	propositionCodes, err := getLatestPropositionsRegisteredAtCamara()
	if err != nil {
		log.Error("getLatestPropositionsRegisteredAtCamara(): ", err.Error())
		return
	}

	latestPropositionCodes, err := instance.propositionRepository.GetLatestPropositionCodes()
	if err != nil {
		log.Error("propositionRepository.GetLatestPropositionCodes()(): ", err.Error())
		return
	}

	newPropositionCodes := findNewPropositionsByCodes(latestPropositionCodes, propositionCodes)
	if newPropositionCodes != nil {
		log.Infof("%d new propositions were identified for registration", len(newPropositionCodes))
	} else {
		log.Info("No new propositions were identified for registration")
		return
	}

	for _, propositionCode := range newPropositionCodes {
		propositionData, err := instance.getProposition(propositionCode)
		if err != nil {
			log.Warn("getProposition(): ", err.Error())
			continue
		}

		var deputies []deputy.Deputy
		for _, deputyData := range propositionData.Deputies() {
			deputyParty := deputyData.Party()
			registeredParty, err := instance.partyRepository.GetPartyByCode(deputyParty.Code())
			if err != nil {
				log.Warn("partyRepository.GetPartyByCode(): ", err.Error())
				continue
			}

			var partyId *uuid.UUID
			if registeredParty == nil {
				partyId, err = instance.partyRepository.CreateParty(deputyParty)
				if err != nil {
					log.Warn("partyRepository.CreateParty(): ", err.Error())
					continue
				}
			} else if !registeredParty.IsEqual(deputyParty) {
				err = instance.partyRepository.UpdateParty(deputyParty)
				if err != nil {
					log.Warn("partyRepository.UpdateParty(): ", err.Error())
					continue
				}
			}

			var updatedParty *party.Party
			if partyId == nil {
				updatedParty, err = deputyParty.NewUpdater().Id(registeredParty.Id()).Build()
				if err != nil {
					log.Warnf("Error updating party %s: %s", registeredParty.Id(), err.Error())
					continue
				}
			} else {
				updatedParty, err = deputyParty.NewUpdater().Id(*partyId).Build()
				if err != nil {
					log.Warnf("Error updating party %s: %s", partyId, err.Error())
					continue
				}
			}

			updatedDeputy, err := deputyData.NewUpdater().Party(*updatedParty).Build()
			if err != nil {
				log.Warnf("Error updating party %s of deputy %d: %s", partyId, deputyData.Code(), err.Error())
				continue
			}

			registeredDeputy, err := instance.deputyRepository.GetDeputyByCode(updatedDeputy.Code())
			if err != nil {
				log.Warn("deputyRepository.GetDeputyByCode(): ", err.Error())
				continue
			}

			var deputyId *uuid.UUID
			if registeredDeputy == nil {
				deputyId, err = instance.deputyRepository.CreateDeputy(*updatedDeputy)
				if err != nil {
					log.Warn("deputyRepository.CreateDeputy(): ", err.Error())
					continue
				}
			} else if !registeredDeputy.IsEqual(*updatedDeputy) {
				err = instance.deputyRepository.UpdateDeputy(*updatedDeputy)
				if err != nil {
					log.Warn("deputyRepository.UpdateDeputy(): ", err.Error())
					continue
				}
			}

			if deputyId == nil {
				updatedDeputy, err = updatedDeputy.NewUpdater().Id(registeredDeputy.Id()).Build()
				if err != nil {
					log.Warnf("Error updating deputy %d: %s", registeredDeputy.Id(), err.Error())
					continue
				}
			} else {
				updatedDeputy, err = updatedDeputy.NewUpdater().Id(*deputyId).Build()
				if err != nil {
					log.Warnf("Error updating deputy %d: %s", deputyId, err.Error())
					continue
				}
			}

			deputies = append(deputies, *updatedDeputy)
		}

		var externalAuthors []external.ExternalAuthor
		for _, externalAuthorData := range propositionData.ExternalAuthors() {
			registeredExternalAuthor, err := instance.externalAuthorRepository.GetExternalAuthorByNameAndType(
				externalAuthorData.Name(), externalAuthorData.Type())
			if err != nil {
				log.Warn("externalAuthorRepository.GetExternalAuthorByNameAndType(): ", err.Error())
				continue
			}

			var externalAuthorId *uuid.UUID
			if registeredExternalAuthor == nil {
				externalAuthorId, err = instance.externalAuthorRepository.CreateExternalAuthor(externalAuthorData)
				if err != nil {
					log.Warn("externalAuthorRepository.CreateExternalAuthor(): ", err.Error())
					continue
				}
			}

			if externalAuthorId == nil {
				externalAuthors = append(externalAuthors, *registeredExternalAuthor)
				continue
			}

			updatedExternalAuthor, err := externalAuthorData.NewUpdater().Id(*externalAuthorId).Build()
			if err != nil {
				log.Warnf("Error updating external author %s: %s", externalAuthorData.Name(), err.Error())
				continue
			}

			externalAuthors = append(externalAuthors, *updatedExternalAuthor)
		}

		var updatedProposition *proposition.Proposition
		updatedProposition, err = propositionData.NewUpdater().
			Deputies(deputies).
			ExternalAuthors(externalAuthors).
			Build()
		if err != nil {
			log.Warnf("Error updating proposition %d: %s", propositionData.Code(), err.Error())
			continue
		}

		err = instance.propositionRepository.CreateProposition(*updatedProposition)
		if err != nil {
			log.Warn("propositionRepository.CreateProposition(): ", err.Error())
			continue
		}
	}

	return
}

func getLatestPropositionsRegisteredAtCamara() ([]int, error) {
	log.Info("Starting the search for the latest propositions")

	latestPropositionsUrl := "https://dadosabertos.camara.leg.br/api/v2/proposicoes?ordenarPor=id&ordem=desc&itens=25"
	newPropositions, err := getDataSliceFromUrl(latestPropositionsUrl)
	if err != nil {
		log.Error("getDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	propositionCodes, err := extractPropositionCodes(newPropositions)
	if err != nil {
		log.Error("extractPropositionCodes(): ", err.Error())
		return nil, err
	}

	sort.Ints(propositionCodes)

	log.Info("Successful search for the latest propositions: ", propositionCodes)
	return propositionCodes, nil
}

func getDataSliceFromUrl(url string) ([]map[string]interface{}, error) {
	response, err := requests.GetRequest(url)
	if err != nil {
		log.Error("requests.GetRequest(): ", err.Error())
		return nil, err
	}

	responseBody, err := requests.DecodeResponseBody(response)
	if err != nil {
		log.Error("requests.DecodeResponseBody(): ", err.Error())
		return nil, err
	}

	resultMapSlice, err := converters.ToMapSlice(responseBody["dados"])
	if err != nil {
		log.Error("converters.ToMapSlice(): ", err.Error())
		return nil, err
	}

	return resultMapSlice, nil
}

func extractPropositionCodes(propositions []map[string]interface{}) ([]int, error) {
	var propositionCodes []int
	for _, propositionData := range propositions {
		code, err := converters.ToInt(propositionData["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		propositionCodes = append(propositionCodes, code)
	}

	return propositionCodes, nil
}

func findNewPropositionsByCodes(latestPropositionCodesRegistered, latestPropositionCodesReturned []int) []int {
	latestPropositionCodes := make(map[int]bool)
	var newPropositions []int

	for _, code := range latestPropositionCodesRegistered {
		latestPropositionCodes[code] = true
	}

	for _, code := range latestPropositionCodesReturned {
		if !latestPropositionCodes[code] {
			newPropositions = append(newPropositions, code)
		}
	}

	return newPropositions
}

func (instance BackgroundData) getProposition(propositionCode int) (*proposition.Proposition, error) {
	log.Info("Starting summary of proposition ", propositionCode)
	propositionData, err := instance.getPropositionDataToRegister(propositionCode)
	if err != nil {
		for attempt := 1; attempt <= 3; attempt++ {
			waitingTimeInSeconds := int(math.Pow(5, float64(attempt)))
			log.Warnf("It was not possible to register proposition %d on the %dth attempt, trying again in %d seconds",
				propositionCode, attempt, waitingTimeInSeconds)
			time.Sleep(time.Duration(waitingTimeInSeconds) * time.Second)
			propositionData, err = instance.getPropositionDataToRegister(propositionCode)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.Error("Error summarizing proposition ", propositionCode)
		return nil, err
	}

	log.Infof("Proposition %d successfully summarized", propositionCode)
	return propositionData, nil
}

func (instance BackgroundData) getPropositionDataToRegister(propositionCode int) (*proposition.Proposition, error) {
	propositionUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/proposicoes/", propositionCode)
	propositionData, err := getDataObjectFromUrl(propositionUrl)
	if err != nil {
		log.Error("getDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	deputies, externalAuthors, err := getAuthorsOfTheProposition(fmt.Sprint(propositionData["uriAutores"]))
	if err != nil {
		log.Error("getAuthorsOfTheProposition(): ", err.Error())
		return nil, err
	}

	articleTypeCode := fmt.Sprint(propositionData["codTipo"])
	articleType, err := instance.articleTypeRepository.GetArticleTypeByCodeOrDefaultType(articleTypeCode)
	if err != nil {
		log.Error("articleTypeRepository.GetArticleTypeByCodeOrDefaultType(): ", err.Error())
		return nil, err
	}

	specificType, err := getArticleTypeDescription(articleTypeCode)
	if err != nil {
		log.Error("getArticleTypeDescription(): ", err.Error())
		return nil, err
	}

	originalTextUrl := fmt.Sprint(propositionData["urlInteiroTeor"])
	var propositionContentSummary string
	if validators.IsUrlValid(originalTextUrl) {
		propositionText, err := getPropositionContent(propositionCode, originalTextUrl)
		if err != nil {
			log.Error("getPropositionContent(): ", err.Error())
			return nil, err
		}

		chatGptCommand := "Resuma a seguinte proposição política de forma simples e direta, como se estivesse escrevendo " +
			"para uma revista. O texto produzido deve conter no máximo três parágrafos:"
		purpose := fmt.Sprint("Summary of the content of proposition ", propositionCode)
		propositionContentSummary, err = requestToChatGpt(chatGptCommand, propositionText, purpose)
		if err != nil {
			log.Error("requestToChatGpt(): ", err.Error())
			return nil, err
		}
	} else {
		propositionContentSummary = fmt.Sprint(propositionData["ementa"])
	}

	chatGptCommand := "Gere um título chamativo para a seguinte matéria para uma revista sobre uma proposição política:"
	purpose := fmt.Sprint("Generating the title of proposition ", propositionCode)
	propositionTitle, err := requestToChatGpt(chatGptCommand, propositionContentSummary, purpose)
	if err != nil {
		log.Error("requestToChatGpt(): ", err.Error())
		return nil, err
	}

	submittedAt, err := time.Parse("2006-01-02T15:04", fmt.Sprint(propositionData["dataApresentacao"]))
	if err != nil {
		log.Errorf("Error converting submission date of proposition %d: %s", propositionCode, err.Error())
		return nil, err
	}

	activeEconomyMode, err := strconv.ParseBool(os.Getenv("ACTIVE_ECONOMY_MODE"))
	if err != nil {
		log.Error("Error converting environment variable ACTIVE_ECONOMY_MODE to boolean, this setting is disabled: ",
			err.Error())
	}

	var propositionImageUrl string
	if !strings.Contains(articleType.Codes(), "default_option") || !activeEconomyMode {
		propositionImageUrl, err = getPropositionImage(propositionCode, propositionContentSummary)
		if err != nil {
			log.Error("getPropositionImage(): ", err.Error())
			return nil, err
		}
	} else {
		log.Infof("Active economy mode: Image generation for proposition %d was skipped", propositionCode)
	}

	articleData, err := article.NewBuilder().Type(*articleType).Build()
	if err != nil {
		log.Errorf("Error validating article data for proposition %d: %s", propositionCode,
			err.Error())
		return nil, err
	}

	propositionBuilder := proposition.NewBuilder().
		Code(propositionCode).
		OriginalTextUrl(originalTextUrl).
		Title(strings.Trim(propositionTitle, "\"")).
		Content(propositionContentSummary).
		SubmittedAt(submittedAt).
		SpecificType(specificType).
		Deputies(deputies).
		ExternalAuthors(externalAuthors).
		Article(*articleData)

	if propositionImageUrl != "" {
		propositionBuilder.ImageUrl(propositionImageUrl)
	}

	propositionDataToRegister, err := propositionBuilder.Build()
	if err != nil {
		log.Errorf("Error validating data for proposition %d: %s", propositionCode, err.Error())
		return nil, err
	}

	return propositionDataToRegister, err
}

func getDataObjectFromUrl(url string) (map[string]interface{}, error) {
	response, err := requests.GetRequest(url)
	if err != nil {
		log.Error("requests.GetRequest(): ", err.Error())
		return nil, err
	}

	responseBody, err := requests.DecodeResponseBody(response)
	if err != nil {
		log.Error("requests.DecodeResponseBody(): ", err.Error())
		return nil, err
	}

	resultMap, err := converters.ToMap(responseBody["dados"])
	if err != nil {
		log.Error("converters.ToMap(): ", err.Error())
		return nil, err
	}

	return resultMap, nil
}

func getPropositionContent(propositionCode int, propositionPdfUrl string) (string, error) {
	log.Info("Extracting content from proposition ", propositionCode)

	pdfContentExtractorAddress := fmt.Sprintf("%s/api/v1/extract-content", os.Getenv("VNC_PDF_CONTENT_EXTRACTOR_API_ADDRESS"))
	body := map[string]string{
		"pdf_url": propositionPdfUrl,
	}
	requestBody, err := converters.ToJson(body)
	if err != nil {
		log.Error("converters.ToJson(): ", err.Error())
		return "", err
	}

	request, err := http.NewRequest("POST", pdfContentExtractorAddress, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Error building the request for communication with VNC PDF Content Extractor API: ", err.Error())
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error("Error making request to VNC PDF Content Extractor API: ", err.Error())
		return "", err
	}
	defer requests.CloseResponseBody(request, response)

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error making request to VNC PDF Content Extractor API: [Status code: %s]",
			response.Status)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	responseBody, err := requests.DecodeResponseBody(response)
	if err != nil {
		log.Error("requests.DecodeResponseBody(): ", err.Error())
		return "", err
	}

	propositionContent := fmt.Sprint(responseBody["content"])

	log.Info("Successful extraction of content from proposition ", propositionCode)
	return propositionContent, nil
}

func getAuthorsOfTheProposition(authorsUrl string) ([]deputy.Deputy, []external.ExternalAuthor, error) {
	authors, err := getDataSliceFromUrl(authorsUrl)
	if err != nil {
		log.Error("getDataSliceFromUrl(): ", err.Error())
		return nil, nil, err
	}

	deputies, externalAuthors, err := convertAuthorsMapToDeputiesAndExternalAuthors(authors)
	if err != nil {
		log.Error("convertAuthorsMapToDeputiesAndExternalAuthors(): ", err.Error())
		return nil, nil, err
	}

	return deputies, externalAuthors, nil
}

func convertAuthorsMapToDeputiesAndExternalAuthors(authors []map[string]interface{}) ([]deputy.Deputy,
	[]external.ExternalAuthor, error) {
	var deputies []deputy.Deputy
	var externalAuthors []external.ExternalAuthor

	sort.Slice(authors, func(i, j int) bool {
		return authors[i]["ordemAssinatura"].(float64) < authors[j]["ordemAssinatura"].(float64)
	})

	for authorIndex, author := range authors {
		authorName := fmt.Sprint(author["nome"])
		authorType := fmt.Sprint(author["tipo"])
		authorUrl := fmt.Sprint(author["uri"])

		log.Infof("Starting the search for the %dth author: %s - %s", authorIndex+1, authorName, authorType)

		var authorData map[string]interface{}
		var err error
		if authorUrl != "" {
			authorData, err = getDataObjectFromUrl(authorUrl)
			if err != nil {
				log.Warnf("Error searching data for author %s - %s: %s", authorName, authorType, err.Error())
				return nil, nil, err
			}
		}

		if authorType == "Deputado(a)" || authorType == "Deputado" {
			deputyCode, err := converters.ToInt(authorData["id"])
			if err != nil {
				log.Error("converters.ToInt(): ", err.Error())
				return nil, nil, err
			}

			authorLastStatus, err := converters.ToMap(authorData["ultimoStatus"])
			if err != nil {
				log.Error("converters.ToMap(): ", err.Error())
				return nil, nil, err
			}

			urlOfPartiesWithAcronym := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/partidos?sigla=",
				authorLastStatus["siglaPartido"])
			partiesWithTheAcronym, err := getDataSliceFromUrl(urlOfPartiesWithAcronym)
			if err != nil {
				log.Error("getDataSliceFromUrl(): ", err.Error())
				return nil, nil, err
			}

			partyUrl := fmt.Sprint(partiesWithTheAcronym[0]["uri"])
			partyMap, err := getDataObjectFromUrl(partyUrl)
			if err != nil {
				log.Error("getDataObjectFromUrl(): ", err.Error())
				return nil, nil, err
			}

			partyCode, err := converters.ToInt(partyMap["id"])
			if err != nil {
				log.Error("converters.ToInt(): ", err.Error())
				return nil, nil, err
			}

			partyData, err := party.NewBuilder().
				Code(partyCode).
				Name(fmt.Sprint(partyMap["nome"])).
				Acronym(fmt.Sprint(partyMap["sigla"])).
				ImageUrl(fmt.Sprint(partyMap["urlLogo"])).
				Build()
			if err != nil {
				log.Errorf("Error validating data for party %d of deputy %d: %s", partyCode, deputyCode,
					err.Error())
				return nil, nil, err
			}

			deputyData, err := deputy.NewBuilder().
				Code(deputyCode).
				Cpf(fmt.Sprint(authorData["cpf"])).
				Name(cases.Title(language.BrazilianPortuguese).String(strings.ToLower(fmt.Sprint(authorData["nomeCivil"])))).
				ElectoralName(fmt.Sprint(authorLastStatus["nomeEleitoral"])).
				ImageUrl(fmt.Sprint(authorLastStatus["urlFoto"])).
				Party(*partyData).
				Build()
			if err != nil {
				log.Errorf("Error validating data for deputy %d: %s", deputyCode, err.Error())
				return nil, nil, err
			}

			deputies = append(deputies, *deputyData)
		} else {

			caser := cases.Title(language.BrazilianPortuguese)
			externalAuthorData, err := external.NewBuilder().
				Name(caser.String(strings.ToLower(authorName))).
				Type(caser.String(strings.ToLower(authorType))).
				Build()
			if err != nil {
				log.Errorf("Error validating data for external author %s - %s: %s", authorName, authorType, err.Error())
				return nil, nil, err
			}

			externalAuthors = append(externalAuthors, *externalAuthorData)
		}
		log.Infof("Successful search for %dth author", authorIndex+1)
	}

	return deputies, externalAuthors, nil
}

func getArticleTypeDescription(articleTypeCode string) (string, error) {
	articleTypesUrl := "https://dadosabertos.camara.leg.br/api/v2/referencias/proposicoes/siglaTipo"
	articleTypes, err := getDataSliceFromUrl(articleTypesUrl)
	if err != nil {
		log.Error("getDataSliceFromUrl(): ", err.Error())
		return "", err
	}

	for _, articleType := range articleTypes {
		if fmt.Sprint(articleType["cod"]) == articleTypeCode {
			return fmt.Sprintf("%s (%s)", articleType["nome"], articleType["cod"]), nil
		}
	}

	return "", nil
}

func getPropositionImage(propositionCode int, propositionContent string) (string, error) {
	prompt, err := requestToChatGpt("Gere um prompt para o DALL·E gerar uma imagem para um site jornalistico "+
		"sobre a seguinte proposição política brasileira. É importante que o prompt esteja de acordo com as políticas "+
		"do DALL·E e que seja especificado a necessidade de evitar usar textos nessas imagens: ",
		propositionContent,
		fmt.Sprint("Geração do prompt da imagem da proposição ", propositionCode))
	if err != nil {
		log.Error("requestToChatGpt(): ", err.Error())
		return "", err
	}

	purpose := fmt.Sprint("Generating the image of proposition ", propositionCode)
	dallEImageUrl, err := requestToDallE(prompt, purpose)
	if err != nil {
		log.Error("requestToDallE(): ", err.Error())
		return "", err
	}

	response, err := requests.GetRequest(dallEImageUrl)
	if err != nil {
		log.Error("requests.GetRequest(): ", err.Error())
		return "", err
	}

	image, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Error interpreting the image of proposition %d: %s", propositionCode, err.Error())
		return "", err
	}

	imageUrl, err := savePropositionImageInAwsS3(propositionCode, image)
	if err != nil {
		log.Error("savePropositionImageInAwsS3(): ", err.Error())
		return "", err
	}

	return imageUrl, nil
}

func (instance BackgroundData) RegisterNewNewsletter(referenceDate time.Time) {
	formattedReferenceDate := referenceDate.Format("02/01/2006")

	propositions, err := instance.propositionRepository.GetPropositionsByDate(referenceDate)
	if err != nil {
		log.Error("propositionRepository.GetPropositionsByDate(): ", err.Error())
		return
	} else if propositions == nil {
		log.Infof("No propositions were found on %s to generate a new newsletter", formattedReferenceDate)
		return
	}

	var propositionOutsideTheNewsletter []proposition.Proposition
	registeredNewsletter, err := instance.newsletterRepository.GetNewsletterByReferenceDate(referenceDate)
	if err != nil {
		log.Error("newsletterRepository.GetNewsletterByReferenceDate(): ", err.Error())
		return
	} else if registeredNewsletter != nil {
		newsletterPropositions, err := instance.propositionRepository.GetPropositionsByNewsletterId(registeredNewsletter.Id())
		if err != nil {
			log.Error("propositionRepository.GetPropositionsByNewsletterId(): ", err.Error())
			return
		}

		for _, propositionData := range propositions {
			var isInTheNewsletter bool
			for _, newsletterPropositionData := range newsletterPropositions {
				if newsletterPropositionData.Id() == propositionData.Id() {
					isInTheNewsletter = true
					break
				}
			}

			if !isInTheNewsletter {
				propositionOutsideTheNewsletter = append(propositionOutsideTheNewsletter, propositionData)
			}
		}
	} else {
		propositionOutsideTheNewsletter = propositions
	}

	if propositionOutsideTheNewsletter == nil {
		log.Infof("No new propositions were found on %s to update newsletter %s", formattedReferenceDate,
			registeredNewsletter.Id())
		return
	}

	if registeredNewsletter == nil {
		log.Info("Starting generation of the newsletter of ", formattedReferenceDate)
	} else {
		log.Infof("Starting update of the newsletter %s of %s", registeredNewsletter.Id(),
			formattedReferenceDate)
	}

	newsletterData, err := instance.generateNewsletter(propositions, referenceDate)
	if err != nil {
		for attempt := 1; attempt <= 3; attempt++ {
			waitingTimeInSeconds := int(math.Pow(5, float64(attempt)))
			log.Warnf("It was not possible to register newsletter of %s on the %dth attempt, trying again in %d seconds",
				formattedReferenceDate, attempt, waitingTimeInSeconds)
			time.Sleep(time.Duration(waitingTimeInSeconds) * time.Second)
			newsletterData, err = instance.generateNewsletter(propositions, referenceDate)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.Error("Error generating the newsletter of ", formattedReferenceDate)
		return
	}

	log.Infof("Newsletter of %s successfully generated ", formattedReferenceDate)

	if registeredNewsletter != nil {
		newsletterData, err = newsletterData.NewUpdater().Id(registeredNewsletter.Id()).Build()
		if err != nil {
			log.Errorf("Error updating the newsletter data structure of %s: %s", formattedReferenceDate, err.Error())
			return
		}
	}

	if registeredNewsletter == nil {
		err = instance.newsletterRepository.CreateNewsletter(*newsletterData, propositions)
		if err != nil {
			log.Error("Error registering the newsletter of ", formattedReferenceDate)
		}
	} else {
		err = instance.newsletterRepository.UpdateNewsletter(*newsletterData, propositionOutsideTheNewsletter)
		if err != nil {
			log.Error("Error updating the newsletter of ", formattedReferenceDate)
		}
	}

	return
}

func (instance BackgroundData) generateNewsletter(propositions []proposition.Proposition, referenceDate time.Time) (*newsletter.Newsletter, error) {
	formattedReferenceDate := referenceDate.Format("02/01/2006")

	var contentOfPropositions string
	for count, propositionData := range propositions {
		contentOfPropositions += fmt.Sprintf("Título da %dª matéria: %s\nConteúdo: %s\n\n", count+1,
			propositionData.Title(), propositionData.Content())
	}

	chatGptCommand := "Gere uma descrição para ser usada no boletim abaixo sobre o conjunto de proposições  políticas " +
		"abaixo. É importante que a descrição seja curta e chamativa, falando sobre o máximo de proposições possíveis e " +
		"correlacionando os temas. Não deve ter mais do que 500 caracteres. Proposições:"
	purpose := fmt.Sprint("Generating the newsletter description of ", formattedReferenceDate)
	newsletterDescription, err := requestToChatGpt(chatGptCommand, contentOfPropositions, purpose)
	if err != nil {
		log.Error("requestToChatGpt(): ", err.Error())
		return nil, err
	}

	articleType, err := instance.articleTypeRepository.GetArticleTypeByCodeOrDefaultType("newsletter")
	if err != nil {
		log.Error("articleTypeRepository.GetArticleTypeByCodeOrDefaultType(): ", err.Error())
		return nil, err
	}

	articleData, err := article.NewBuilder().Type(*articleType).Build()
	if err != nil {
		log.Errorf("Error validating article data for newsletter of %s: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	newsletterData, err := newsletter.NewBuilder().
		ReferenceDate(referenceDate).
		Title(fmt.Sprint("Boletim do dia ", formattedReferenceDate)).
		Description(newsletterDescription).
		Article(*articleData).
		Build()
	if err != nil {
		log.Errorf("Error validating newsletter data of %s: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	return newsletterData, nil
}
