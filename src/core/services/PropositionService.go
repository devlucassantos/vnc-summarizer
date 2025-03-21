package services

import (
	"errors"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"vnc-summarizer/core/interfaces/repositories"
	"vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/datetime"
	"vnc-summarizer/core/services/utils/requesters"
	"vnc-summarizer/core/services/utils/validators"
)

type Proposition struct {
	propositionRepository     repositories.Proposition
	propositionTypeRepository repositories.PropositionType
	articleTypeRepository     repositories.ArticleType
	authorService             services.Author
}

func NewPropositionService(propositionRepository repositories.Proposition,
	propositionTypeRepository repositories.PropositionType, articleTypeRepository repositories.ArticleType,
	author services.Author) *Proposition {
	return &Proposition{
		propositionRepository:     propositionRepository,
		propositionTypeRepository: propositionTypeRepository,
		articleTypeRepository:     articleTypeRepository,
		authorService:             author,
	}
}

func (instance Proposition) RegisterNewPropositions() {
	codesOfTheMostRecentPropositionsReturned, err := getCodesOfTheMostRecentPropositionsRegisteredInTheChamber()
	if err != nil {
		log.Error("getCodesOfTheMostRecentPropositionsRegisteredInTheChamber(): ", err.Error())
		return
	}

	registeredPropositions, err := instance.propositionRepository.
		GetPropositionsByCodes(codesOfTheMostRecentPropositionsReturned)
	if err != nil {
		log.Error("propositionRepository.GetCodesOfTheMostRecentPropositions(): ", err.Error())
		return
	}

	codesOfTheNewPropositions := getCodesOfTheNewPropositions(codesOfTheMostRecentPropositionsReturned,
		registeredPropositions)
	if codesOfTheNewPropositions != nil {
		log.Infof("%d new propositions were identified for registration: %v", len(codesOfTheNewPropositions),
			codesOfTheNewPropositions)
	} else {
		log.Info("No new propositions were identified for registration")
		return
	}

	for _, propositionCode := range codesOfTheNewPropositions {
		_, err = instance.RegisterNewPropositionByCode(propositionCode)
		if err != nil {
			log.Error("RegisterNewPropositionByCode(): ", err.Error())
		}
	}

	return
}

func getCodesOfTheMostRecentPropositionsRegisteredInTheChamber() ([]int, error) {
	log.Info("Starting the search for the most recent propositions")

	currentDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	var mostRecentPropositionsReturned []map[string]interface{}
	for page := 1; ; page++ {
		chunkSize := 100
		urlOfTheMostRecentPropositions := fmt.Sprintf(
			"https://dadosabertos.camara.leg.br/api/v2/proposicoes?&pagina=%d&itens=%d&dataApresentacaoInicio=%s&ordenarPor=id&ordem=asc",
			page, chunkSize, currentDateTime.AddDate(0, 0, -1).Format("2006-01-02"),
		)
		mostRecentPropositions, err := requesters.GetDataSliceFromUrl(urlOfTheMostRecentPropositions)
		if err != nil {
			log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
			return nil, err
		}

		mostRecentPropositionsReturned = append(mostRecentPropositionsReturned, mostRecentPropositions...)

		if len(mostRecentPropositions) < chunkSize {
			break
		}
	}

	propositionCodes, err := extractPropositionCodes(mostRecentPropositionsReturned)
	if err != nil {
		log.Error("extractPropositionCodes(): ", err.Error())
		return nil, err
	}

	log.Info("Successful search for the latest propositions: ", propositionCodes)
	return propositionCodes, nil
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

func getCodesOfTheNewPropositions(returnedPropositionCodes []int, registeredPropositions []proposition.Proposition) []int {
	var propositionCodesToRegister []int
	for _, propositionCode := range returnedPropositionCodes {
		var propositionAlreadyRegistered bool
		for _, propositionData := range registeredPropositions {
			if propositionData.Code() == propositionCode {
				propositionAlreadyRegistered = true
				break
			}
		}
		if !propositionAlreadyRegistered {
			propositionCodesToRegister = append(propositionCodesToRegister, propositionCode)
		}
	}

	propositionCodesToRegister = converters.IntSliceToUniqueIntSlice(propositionCodesToRegister)

	return propositionCodesToRegister
}

func (instance Proposition) RegisterNewPropositionByCode(code int) (*uuid.UUID, error) {
	propositionData, err := instance.getProposition(code)
	if err != nil {
		if !strings.Contains(err.Error(), "no content") {
			log.Error("getProposition(): ", err.Error())
		}
		return nil, err
	}

	propositionId, err := instance.propositionRepository.CreateProposition(*propositionData)
	if err != nil {
		log.Error("propositionRepository.CreateProposition(): ", err.Error())
		return nil, err
	}

	return propositionId, err
}

func (instance Proposition) getProposition(code int) (*proposition.Proposition, error) {
	log.Info("Starting summary of proposition ", code)
	propositionData, err := instance.getPropositionDataToRegister(code)
	if err != nil && !strings.Contains(err.Error(), "no content") {
		for attempt := 1; attempt <= 3; attempt++ {
			waitingTimeInSeconds := int(math.Pow(4, float64(attempt)))
			log.Warnf("It was not possible to register proposition %d on the %dth attempt, trying again in %d seconds",
				code, attempt, waitingTimeInSeconds)
			time.Sleep(time.Duration(waitingTimeInSeconds) * time.Second)
			propositionData, err = instance.getPropositionDataToRegister(code)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		if !strings.Contains(err.Error(), "no content") {
			log.Errorf("Error summarizing proposition %d: %s", code, err.Error())
		}
		return nil, err
	}

	log.Infof("Proposition %d successfully summarized", code)
	return propositionData, nil
}

func (instance Proposition) getPropositionDataToRegister(propositionCode int) (*proposition.Proposition, error) {
	propositionUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/proposicoes/", propositionCode)
	propositionData, err := requesters.GetDataObjectFromUrl(propositionUrl)
	if err != nil {
		log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	originalTextUrl := fmt.Sprint(propositionData["urlInteiroTeor"])
	if originalTextUrl == "<nil>" {
		errorMessage := fmt.Sprintf("Proposition %d can not be registered because it has no content",
			propositionCode)
		log.Warn(errorMessage)
		err = errors.New(errorMessage)
		return nil, err
	}

	authorsUrl := fmt.Sprint(propositionData["uriAutores"])
	deputies, externalAuthors, err := instance.authorService.GetAuthorsFromAuthorsUrl(authorsUrl)
	if err != nil {
		log.Error("authorService.GetAuthorsFromAuthorsUrl(): ", err.Error())
		return nil, err
	}

	propositionTypeCode := fmt.Sprint(propositionData["codTipo"])
	propositionType, err := instance.propositionTypeRepository.GetPropositionTypeByCodeOrDefaultType(propositionTypeCode)
	if err != nil {
		log.Error("propositionTypeRepository.GetPropositionTypeByCodeOrDefaultType(): ", err.Error())
		return nil, err
	}

	specificType, err := getArticleSpecificType(propositionTypeCode)
	if err != nil {
		log.Error("getArticleSpecificType(): ", err.Error())
		return nil, err
	}

	var propositionContentSummary, propositionText, originalTextMimeType string
	if validators.IsUrlValid(originalTextUrl) {
		propositionText, originalTextMimeType, err = getPropositionContent(propositionCode, originalTextUrl)
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

	chatGptCommand := "Gere um título chamativo utilizando uma linguagem simples e direta para a seguinte matéria para " +
		"uma revista sobre uma proposição política: "
	purpose := fmt.Sprint("Generating the title of proposition ", propositionCode)
	propositionTitle, err := requestToChatGpt(chatGptCommand, propositionContentSummary, purpose)
	if err != nil {
		log.Error("requestToChatGpt(): ", err.Error())
		return nil, err
	}
	propositionTitle = strings.Trim(propositionTitle, "*\"")

	submittedAt, err := time.Parse("2006-01-02T15:04", fmt.Sprint(propositionData["dataApresentacao"]))
	if err != nil {
		log.Errorf("Error converting submission date and time of proposition %d: %s", propositionCode, err.Error())
		return nil, err
	}

	economyModeActive, err := strconv.ParseBool(os.Getenv("ECONOMY_MODE_ACTIVE"))
	if err != nil {
		log.Error("Error converting environment variable ECONOMY_MODE_ACTIVE to boolean, this setting is disabled: ",
			err.Error())
	}

	var propositionImageUrl, propositionImageDescription string
	if !economyModeActive || !strings.Contains(propositionType.Codes(), "default_option") {
		propositionImageUrl, propositionImageDescription, err = getPropositionImage(propositionCode,
			propositionContentSummary)
		if err != nil {
			log.Error("getPropositionImage(): ", err.Error())
			return nil, err
		}
	} else {
		log.Infof("Active economy mode: Image generation for proposition %d was skipped", propositionCode)
	}

	referenceDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	if referenceDateTime.Sub(submittedAt).Hours() > 24 {
		referenceDateTime = &submittedAt
	}

	articleTypeCode := "proposition"
	articleType, err := instance.articleTypeRepository.GetArticleTypeByCode(articleTypeCode)
	if err != nil {
		log.Error("articleTypeRepository.GetArticleTypeByCode(): ", err.Error())
		return nil, err
	}

	articleData, err := article.NewBuilder().Type(*articleType).ReferenceDateTime(*referenceDateTime).Build()
	if err != nil {
		log.Errorf("Error validating article data for proposition %d: %s", propositionCode,
			err.Error())
		return nil, err
	}

	propositionBuilder := proposition.NewBuilder().
		Code(propositionCode).
		OriginalTextUrl(originalTextUrl).
		OriginalTextMimeType(originalTextMimeType).
		Title(propositionTitle).
		Content(propositionContentSummary).
		SubmittedAt(submittedAt).
		SpecificType(specificType).
		Type(*propositionType).
		Deputies(deputies).
		ExternalAuthors(externalAuthors).
		Article(*articleData)

	if propositionImageUrl != "" {
		propositionBuilder.ImageUrl(propositionImageUrl).ImageDescription(propositionImageDescription)
	}

	propositionDataToRegister, err := propositionBuilder.Build()
	if err != nil {
		log.Errorf("Error validating data for proposition %d: %s", propositionCode, err.Error())
		return nil, err
	}

	return propositionDataToRegister, err
}

func getPropositionContent(propositionCode int, propositionUrl string) (string, string, error) {
	log.Info("Extracting content from proposition ", propositionCode)

	propositionContent, err := requestToVncPdfContentExtractorApi(propositionUrl)
	if err != nil && !strings.Contains(err.Error(), "Is this really a PDF?") {
		log.Error("requestToVncPdfContentExtractorApi(): ", err.Error())
		return "", "", err
	} else if err == nil && propositionContent == "" {
		errorMessage := fmt.Sprintf("The file of the original text of the proposition %d is not supported: %s",
			propositionCode, propositionUrl)
		log.Error(errorMessage)
		return "", "", errors.New(errorMessage)
	} else if propositionContent != "" {
		propositionMimeType := "application/pdf"
		return propositionContent, propositionMimeType, nil
	}

	propositionContent, err = getPropositionContentDirectly(propositionUrl)
	if err != nil {
		log.Error("getPropositionContentDirectly(): ", err.Error())
		return "", "", err
	}
	propositionMimeType := "text/html"

	return propositionContent, propositionMimeType, nil
}

func getPropositionContentDirectly(propositionUrl string) (string, error) {
	parsedPropositionUrl, err := url.Parse(propositionUrl)
	if err != nil {
		log.Errorf("Error parsing proposition URL %s: %s", propositionUrl, err.Error())
		return "", err
	}
	queryParams := parsedPropositionUrl.Query()
	propositionContentCode := queryParams.Get("codteor")

	propositionUrl = fmt.Sprintf("https://www.camara.leg.br/internet/ordemdodia/integras/%s.htm",
		propositionContentCode)
	response, err := requesters.GetRequest(propositionUrl)
	if err != nil {
		log.Error("requests.GetRequest(): ", err.Error())
		return "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error interpreting response to request %s: %s", propositionUrl, err.Error())
		return "", err
	}
	responseBodyAsString := string(responseBody)

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error making request to get original proposition content directly: "+
			"[Status: %s; Body: %s]", response.Status, responseBodyAsString)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	return responseBodyAsString, err
}

func getArticleSpecificType(articleTypeCode string) (string, error) {
	articleTypesUrl := "https://dadosabertos.camara.leg.br/api/v2/referencias/proposicoes/siglaTipo"
	articleTypes, err := requesters.GetDataSliceFromUrl(articleTypesUrl)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return "", err
	}

	var articleSpecificType string
	for _, articleType := range articleTypes {
		if fmt.Sprint(articleType["cod"]) == articleTypeCode {
			articleSpecificType = fmt.Sprintf("%s (%s)", articleType["nome"], articleType["cod"])
			break
		}
	}

	return articleSpecificType, nil
}

func getPropositionImage(propositionCode int, propositionContent string) (string, string, error) {
	chatGptCommand := "Gere um prompt para o DALL·E gerar uma imagem para um site jornalistico sobre a seguinte " +
		"proposição política brasileira. É importante que o prompt esteja de acordo com as políticas do DALL·E e que " +
		"seja especificado a necessidade de evitar usar textos nessas imagens: "
	purpose := fmt.Sprint("Generating the prompt for the image of proposition ", propositionCode)
	prompt, err := requestToChatGpt(chatGptCommand, propositionContent, purpose)
	if err != nil {
		log.Error("requestToChatGpt(): ", err.Error())
		return "", "", err
	}

	purpose = fmt.Sprint("Generating the image of proposition ", propositionCode)
	dallEImageUrl, err := requestToDallE(prompt, purpose)
	if err != nil {
		log.Error("requestToDallE(): ", err.Error())
		return "", "", err
	}

	response, err := requesters.GetRequest(dallEImageUrl)
	if err != nil {
		log.Error("requests.GetRequest(): ", err.Error())
		return "", "", err
	}

	image, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error interpreting the image of proposition %d: %s", propositionCode, err.Error())
		return "", "", err
	}

	imageUrl, err := savePropositionImageInAwsS3(propositionCode, image)
	if err != nil {
		log.Error("savePropositionImageInAwsS3(): ", err.Error())
		return "", "", err
	}

	imageDescription, err := requestToChatGptVision(imageUrl)
	if err != nil {
		log.Error("requestToChatGptVision(): ", err.Error())
		return "", "", err
	}

	return imageUrl, imageDescription, nil
}

func (instance Proposition) GetPropositionsByCodes(codes []int) ([]proposition.Proposition, error) {
	propositions, err := instance.propositionRepository.GetPropositionsByCodes(codes)
	if err != nil {
		log.Error("propositionRepository.GetPropositionsByCodes(): ", err.Error())
		return nil, err
	}

	return propositions, err
}
