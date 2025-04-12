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
	"os"
	"strconv"
	"strings"
	"time"
	"vnc-summarizer/core/interfaces/chamber"
	"vnc-summarizer/core/interfaces/chatgpt"
	"vnc-summarizer/core/interfaces/dalle"
	"vnc-summarizer/core/interfaces/pdfcontentextractor"
	"vnc-summarizer/core/interfaces/postgres"
	"vnc-summarizer/core/interfaces/s3"
	"vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/utils/converters"
	"vnc-summarizer/utils/datetime"
	"vnc-summarizer/utils/requesters"
)

type Proposition struct {
	authorService             services.Author
	chamberApi                chamber.Chamber
	chatGptApi                chatgpt.ChatGpt
	dallEApi                  dalle.DallE
	vncPdfContentExtractor    pdfcontentextractor.VncPdfContentExtractor
	awsS3Api                  s3.AwsS3
	propositionRepository     postgres.Proposition
	propositionTypeRepository postgres.PropositionType
	articleTypeRepository     postgres.ArticleType
}

func NewPropositionService(authorService services.Author, chamberApi chamber.Chamber,
	chatGptApi chatgpt.ChatGpt, dallEApi dalle.DallE, vncPdfContentExtractor pdfcontentextractor.VncPdfContentExtractor,
	awsS3Api s3.AwsS3, propositionRepository postgres.Proposition, propositionTypeRepository postgres.PropositionType,
	articleTypeRepository postgres.ArticleType) *Proposition {
	return &Proposition{
		authorService:             authorService,
		chamberApi:                chamberApi,
		chatGptApi:                chatGptApi,
		dallEApi:                  dallEApi,
		vncPdfContentExtractor:    vncPdfContentExtractor,
		awsS3Api:                  awsS3Api,
		propositionRepository:     propositionRepository,
		propositionTypeRepository: propositionTypeRepository,
		articleTypeRepository:     articleTypeRepository,
	}
}

func (instance Proposition) RegisterNewPropositions() {
	codesOfTheMostRecentPropositionsReturned, err := instance.getCodesOfTheMostRecentPropositionsRegisteredInTheChamber()
	if err != nil {
		log.Error("getCodesOfTheMostRecentPropositionsRegisteredInTheChamber(): ", err.Error())
		return
	} else if codesOfTheMostRecentPropositionsReturned == nil {
		log.Info("No new propositions were identified for registration")
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

func (instance Proposition) getCodesOfTheMostRecentPropositionsRegisteredInTheChamber() ([]int, error) {
	log.Info("Starting the search for the most recent propositions")

	mostRecentPropositions, err := instance.chamberApi.GetMostRecentPropositions()
	if err != nil {
		log.Error("chamberApi.GetMostRecentPropositions(): ", err.Error())
		return nil, err
	}

	propositionCodes, err := extractPropositionCodes(mostRecentPropositions)
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
	propositionData, err := instance.chamberApi.GetPropositionByCode(propositionCode)
	if err != nil {
		log.Error("chamberApi.GetPropositionByCode(): ", err.Error())
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

	specificType, err := instance.getPropositionSpecificType(propositionTypeCode)
	if err != nil {
		log.Error("getPropositionSpecificType(): ", err.Error())
		return nil, err
	}

	originalTextUrl, propositionText, originalTextMimeType, err := instance.getPropositionContent(propositionCode,
		originalTextUrl)
	if err != nil {
		log.Error("getPropositionContent(): ", err.Error())
		return nil, err
	}

	chatGptCommand := "Resuma a seguinte proposição política de forma simples e direta, como se estivesse escrevendo " +
		"para uma revista. O texto produzido deve conter no máximo três parágrafos:"
	purpose := fmt.Sprint("Summary of the content of proposition ", propositionCode)
	propositionContentSummary, err := instance.chatGptApi.MakeRequest(chatGptCommand, propositionText, purpose)
	if err != nil {
		log.Error("chatGptApi.MakeRequest(): ", err.Error())
		return nil, err
	}

	chatGptCommand = "Gere um título chamativo utilizando uma linguagem simples e direta para a seguinte matéria para " +
		"uma revista sobre uma proposição política: "
	purpose = fmt.Sprint("Generating the title of proposition ", propositionCode)
	propositionTitle, err := instance.chatGptApi.MakeRequest(chatGptCommand, propositionContentSummary, purpose)
	if err != nil {
		log.Error("chatGptApi.MakeRequest(): ", err.Error())
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
		propositionImageUrl, propositionImageDescription, err = instance.getPropositionImage(propositionCode,
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

func (instance Proposition) getPropositionContent(propositionCode int, propositionUrl string) (string, string, string, error) {
	log.Info("Extracting content from proposition ", propositionCode)

	propositionContent, err := instance.vncPdfContentExtractor.MakeRequest(propositionUrl)
	if err != nil && !strings.Contains(err.Error(), "Is this really a PDF?") {
		log.Error("vncPdfContentExtractor.MakeRequest(): ", err.Error())
		return "", "", "", err
	} else if err == nil && propositionContent == "" {
		errorMessage := fmt.Sprintf("The file of the original text of the proposition %d is not supported: %s",
			propositionCode, propositionUrl)
		log.Error(errorMessage)
		return "", "", "", errors.New(errorMessage)
	} else if propositionContent != "" {
		propositionMimeType := "application/pdf"
		return propositionUrl, propositionContent, propositionMimeType, nil
	}

	propositionUrl, propositionContent, err = instance.chamberApi.GetPropositionContentDirectly(propositionUrl)
	if err != nil {
		log.Error("chamberApi.GetPropositionContentDirectly(): ", err.Error())
		return "", "", "", err
	}
	propositionMimeType := "text/html"

	return propositionUrl, propositionContent, propositionMimeType, nil
}

func (instance Proposition) getPropositionSpecificType(propositionTypeCode string) (string, error) {
	propositionTypes, err := instance.chamberApi.GetPropositionTypes()
	if err != nil {
		log.Error("chamberApi.GetPropositionTypes(): ", err.Error())
		return "", err
	}

	var propositionSpecificType string
	for _, propositionType := range propositionTypes {
		if fmt.Sprint(propositionType["cod"]) == propositionTypeCode {
			propositionSpecificType = fmt.Sprintf("%s (%s)", propositionType["nome"], propositionType["cod"])
			break
		}
	}

	return propositionSpecificType, nil
}

func (instance Proposition) getPropositionImage(propositionCode int, propositionContent string) (string, string, error) {
	chatGptCommand := "Gere um prompt para o DALL·E gerar uma imagem para um site jornalistico sobre a seguinte " +
		"proposição política brasileira. É importante que o prompt esteja de acordo com as políticas do DALL·E e que " +
		"seja especificado a necessidade de evitar usar textos nessas imagens: "
	purpose := fmt.Sprint("Generating the prompt for the image of proposition ", propositionCode)
	prompt, err := instance.chatGptApi.MakeRequest(chatGptCommand, propositionContent, purpose)
	if err != nil {
		log.Error("chatGptApi.MakeRequest(): ", err.Error())
		return "", "", err
	}

	purpose = fmt.Sprint("Generating the image of proposition ", propositionCode)
	dallEImageUrl, err := instance.dallEApi.MakeRequest(prompt, purpose)
	if err != nil {
		log.Error("dallEApi.MakeRequest(): ", err.Error())
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

	imageUrl, err := instance.awsS3Api.SavePropositionImage(propositionCode, image)
	if err != nil {
		log.Error("awsS3Api.SavePropositionImage(): ", err.Error())
		return "", "", err
	}

	imageDescription, err := instance.chatGptApi.MakeRequestToVision(imageUrl)
	if err != nil {
		log.Error("chatGptApi.MakeRequestToVision(): ", err.Error())
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
