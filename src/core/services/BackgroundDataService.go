package services

import (
	"encoding/json"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"vnc-write-api/core/interfaces/repositories"
)

type BackgroundData struct {
	deputyRepository       repositories.Deputy
	organizationRepository repositories.Organization
	partyRepository        repositories.Party
	propositionRepository  repositories.Proposition
	newsletterRepository   repositories.Newsletter
}

func NewBackgroundDataService(deputyRepository repositories.Deputy, organizationRepository repositories.Organization,
	partyRepository repositories.Party, propositionRepository repositories.Proposition,
	newsletterRepository repositories.Newsletter) *BackgroundData {
	return &BackgroundData{
		deputyRepository:       deputyRepository,
		organizationRepository: organizationRepository,
		partyRepository:        partyRepository,
		propositionRepository:  propositionRepository,
		newsletterRepository:   newsletterRepository,
	}
}

func (instance BackgroundData) RegisterNewPropositions() {
	propositionCodes, err := getLatestPropositionsRegisteredAtCamara()
	if err != nil {
		return
	}

	latestPropositionCodes, err := instance.propositionRepository.GetLatestPropositionCodes()
	if err != nil {
		return
	}

	newPropositionCodes := findNewPropositionsByCodes(latestPropositionCodes, propositionCodes)
	if newPropositionCodes != nil {
		log.Infof("Foram identificadas %d novas proposições para registro", len(newPropositionCodes))
	} else {
		log.Info("Não foi identificada nenhuma nova proposição para registro")
		return
	}

	for _, propositionCode := range newPropositionCodes {
		propositionData, err := instance.getProposition(propositionCode)
		if err != nil {
			continue
		}

		var deputies []deputy.Deputy
		for _, deputyData := range propositionData.Deputies() {
			deputyParty := deputyData.Party()
			registeredParty, err := instance.partyRepository.GetPartyByCode(deputyParty.Code())
			if err != nil {
				continue
			}

			var partyId *uuid.UUID
			if registeredParty == nil {
				partyId, err = instance.partyRepository.CreateParty(deputyParty)
			} else if !registeredParty.IsEqual(deputyParty) {
				err = instance.partyRepository.UpdateParty(deputyParty)
			}
			if err != nil {
				continue
			}

			var updatedParty *party.Party
			if partyId == nil {
				updatedParty, err = deputyParty.NewUpdater().Id(registeredParty.Id()).Build()
			} else {
				updatedParty, err = deputyParty.NewUpdater().Id(*partyId).Build()
			}
			if err != nil {
				log.Errorf("Erro ao atualizar partido %s: %s", partyId, err.Error())
				continue
			}

			updatedDeputy, err := deputyData.NewUpdater().Party(*updatedParty).Build()
			if err != nil {
				log.Error("Erro ao atualizar partido %s do deputado(a) %d: %s", partyId, deputyData.Code(), err.Error())
				continue
			}
			registeredDeputy, err := instance.deputyRepository.GetDeputyByCode(updatedDeputy.Code())
			if err != nil {
				continue
			}

			var deputyId *uuid.UUID
			if registeredDeputy == nil {
				deputyId, err = instance.deputyRepository.CreateDeputy(*updatedDeputy)
			} else if !registeredParty.IsEqual(deputyParty) {
				err = instance.deputyRepository.UpdateDeputy(*updatedDeputy)
			}
			if err != nil {
				continue
			}

			if deputyId == nil {
				updatedDeputy, err = updatedDeputy.NewUpdater().Id(registeredDeputy.Id()).Build()
			} else {
				updatedDeputy, err = updatedDeputy.NewUpdater().Id(*deputyId).Build()
			}
			if err != nil {
				log.Errorf("Erro ao atualizar deputado(a) %d: %s", deputyData.Code(), err.Error())
				continue
			}

			deputies = append(deputies, *updatedDeputy)
		}

		var organizations []organization.Organization
		for _, organizationData := range propositionData.Organizations() {
			registeredOrganization, err := instance.organizationRepository.GetOrganization(organizationData)
			if err != nil {
				continue
			}

			var organizationId *uuid.UUID
			if registeredOrganization == nil {
				organizationId, err = instance.organizationRepository.CreateOrganization(organizationData)
			} else if !registeredOrganization.IsEqual(organizationData) {
				err = instance.organizationRepository.UpdateOrganization(organizationData)
			}
			if err != nil {
				continue
			}

			var updatedOrganization *organization.Organization
			if organizationId == nil {
				updatedOrganization, err = organizationData.NewUpdater().Id(registeredOrganization.Id()).Build()
			} else {
				updatedOrganization, err = organizationData.NewUpdater().Id(*organizationId).Build()
			}
			if err != nil {
				log.Errorf("Erro ao atualizar organização %d: %s", organizationData.Code(), err.Error())
				continue
			}

			organizations = append(organizations, *updatedOrganization)
		}

		var updatedProposition *proposition.Proposition
		updatedProposition, err = propositionData.NewUpdater().
			Deputies(deputies).
			Organizations(organizations).
			Build()
		if err != nil {
			log.Errorf("Erro ao atualizar proposição %d: %s", propositionData.Code(), err.Error())
			continue
		}

		_, err = instance.propositionRepository.CreateProposition(*updatedProposition)
		if err != nil {
			continue
		}
	}

	return
}

func getLatestPropositionsRegisteredAtCamara() ([]int, error) {
	log.Info("Iniciando busca das últimas proposições")

	latestPropositionsUrl := "https://dadosabertos.camara.leg.br/api/v2/proposicoes?ordenarPor=id&ordem=desc&itens=25"
	response, err := getRequest(latestPropositionsUrl)
	if err != nil {
		return nil, err
	}

	body, err := readResponseBody(response)
	if err != nil {
		return nil, err
	}

	newPropositions, err := getDataFromRequestMap(body)
	if err != nil {
		return nil, err
	}

	propositionCodes, err := extractPropositionCodes(newPropositions)
	if err != nil {
		return nil, err
	}

	return propositionCodes, nil
}

func getRequest(url string) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Errorf("Erro ao realizar a requisição %s: %s", url, err.Error())
		return nil, err
	}

	return response, err
}

func readResponseBody(response *http.Response) (map[string]interface{}, error) {
	content := make(map[string]interface{})
	err := json.NewDecoder(response.Body).Decode(&content)
	if err != nil {
		log.Error("Erro ao ler o corpo da resposta: ", err.Error())
		return nil, err
	}

	return content, nil
}

func getDataFromRequestMap(body map[string]interface{}) ([]map[string]interface{}, error) {
	jsonData, err := extractDataAsJson(body["dados"])
	if err != nil {
		return nil, err
	}

	var resultMap []map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		log.Error("Erro ao converter JSON para []map[string]interface{}: ", err.Error())
		return nil, err
	}

	return resultMap, nil
}

func extractDataAsJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("Erro ao converter extrair os dados e converter para JSON: ", err.Error())
		return nil, err
	}

	return jsonData, nil
}

func extractPropositionCodes(propositions []map[string]interface{}) ([]int, error) {
	log.Info("Iniciando extração dos códigos das últimas proposições")

	var propositionCodes []int
	for _, propositionData := range propositions {
		code, err := convertInterfaceToInt(propositionData["id"])
		if err != nil {
			return nil, err
		}
		propositionCodes = append(propositionCodes, code)
	}

	return propositionCodes, nil
}

func convertInterfaceToInt(data interface{}) (int, error) {
	number, err := strconv.ParseFloat(fmt.Sprint(data), 64)
	if err != nil {
		log.Error("Erro ao converter interface{} em inteiro: %s", err.Error())
		return 0, err
	}

	return int(number), nil
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
	log.Info("Iniciando síntese da proposição ", propositionCode)
	propositionData, err := instance.getPropositionDataToRegister(propositionCode)
	if err != nil {
		for attempt := 1; attempt <= 3; attempt++ {
			waitingTimeInSeconds := int(math.Pow(5, float64(attempt)))
			log.Warnf("Não foi possível registrar a proposição %d na %dª tentativa, tentando novamente em %d segundos",
				propositionCode, attempt, waitingTimeInSeconds)
			time.Sleep(time.Duration(waitingTimeInSeconds) * time.Second)
			propositionData, err = instance.getPropositionDataToRegister(propositionCode)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.Error("Erro ao sintetizar a proposição ", propositionCode)
		return nil, err
	}

	log.Infof("Proposição %d sintetizada com sucesso", propositionCode)
	return propositionData, nil
}

func (instance BackgroundData) getPropositionDataToRegister(propositionCode int) (*proposition.Proposition, error) {
	propositionSummaryUrl := fmt.Sprintf("https://dadosabertos.camara.leg.br/api/v2/proposicoes/%d", propositionCode)
	propositionData, err := getDataFromUrl(propositionSummaryUrl)
	if err != nil {
		return nil, err
	}

	originalTextUrl := fmt.Sprint(propositionData["urlInteiroTeor"])
	propositionText, err := getPropositionContent(originalTextUrl)
	if err != nil {
		return nil, err
	}

	chatGptCommand := "Resuma a seguinte proposição política de forma simples e direta, como se estivesse escrevendo para uma " +
		"revista. O texto produzido deve conter no máximo três parágrafos:"
	purpose := fmt.Sprint("Síntese do conteúdo da proposição ", propositionCode)
	propositionContentSummary, err := requestToChatGpt(chatGptCommand, propositionText, purpose)
	if err != nil {
		return nil, err
	}

	deputies, organizations, err := getAuthorsOfTheProposition(fmt.Sprint(propositionData["uriAutores"]))
	if err != nil {
		return nil, err
	}

	chatGptCommand = "Gere um título chamativo para a seguinte matéria para uma revista sobre uma proposição política:"
	purpose = fmt.Sprint("Geração do título da proposição ", propositionCode)
	propositionTitle, err := requestToChatGpt(chatGptCommand, propositionContentSummary, purpose)
	if err != nil {
		return nil, err
	}

	submittedAt, err := time.Parse("2006-01-02T15:04", fmt.Sprint(propositionData["dataApresentacao"]))
	if err != nil {
		log.Errorf("Erro ao converter data de apresentação da proposição %d: %s", propositionCode, err.Error())
		return nil, err
	}

	propositionDataToRegister, err := proposition.NewBuilder().
		Code(propositionCode).
		OriginalTextUrl(originalTextUrl).
		Title(strings.Trim(propositionTitle, "\"")).
		Content(propositionContentSummary).
		SubmittedAt(submittedAt).
		Deputies(deputies).
		Organizations(organizations).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados da proposição %d: %s", propositionCode, err.Error())
		return nil, err
	}

	return propositionDataToRegister, err
}

func getDataFromUrl(url string) (map[string]interface{}, error) {
	response, err := getRequest(url)
	if err != nil {
		return nil, err
	}

	body, err := readResponseBody(response)
	if err != nil {
		return nil, err
	}

	resultMap, err := extractDataAsMap(body["dados"])
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

func extractDataAsMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := extractDataAsJson(data)
	if err != nil {
		return nil, err
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		log.Error("Erro ao converter JSON para map[string]interface{}: ", err.Error())
		return nil, err
	}

	return resultMap, err
}

func getPropositionContent(url string) (string, error) {
	log.Info("Iniciando leitura do conteúdo da proposição")

	response, err := getRequest(url)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("./", "temp-pdf-*.pdf")
	if err != nil {
		log.Error("Erro ao criar o arquivo temporário: ", err.Error())
		return "", err
	}

	defer removeTempFile(tempFile)

	_, err = io.Copy(tempFile, response.Body)
	if err != nil {
		log.Error("Erro ao salvar o conteúdo da proposição no arquivo temporário: ", err.Error())
		return "", err
	}

	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		log.Error("Erro ao acessar o conteúdo da proposição no arquivo temporário: ", err.Error())
		return "", err
	}

	pdfReader, err := model.NewPdfReader(tempFile)
	if err != nil {
		log.Error("Erro ao acessar o conteúdo da proposição no arquivo temporário: ", err.Error())
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Error("Erro ao verificar o número de páginas do arquivo temporário: ", err.Error())
		return "", err
	}

	var fullText string

	for pageNumber := 1; pageNumber <= numPages; pageNumber++ {
		page, err := pdfReader.GetPage(pageNumber)
		if err != nil {
			log.Errorf("Erro ao buscar a página %d do arquivo temporário: %s", pageNumber, err.Error())
			return "", err
		}

		contentExtractor, err := extractor.New(page)
		if err != nil {
			log.Errorf("Erro ao criar o extrator de conteúdo da página %d do arquivo temporário: %s", pageNumber, err.Error())
			return "", err
		}

		text, err := contentExtractor.ExtractText()
		if err != nil {
			log.Errorf("Erro ao extrair o conteúdo da página %d do arquivo temporário: %s", pageNumber, err.Error())
			return "", err
		}

		fullText += text
	}

	log.Info("Extração do conteúdo da proposição finalizada com sucesso")
	return fullText, nil
}

func removeTempFile(file *os.File) {
	err := os.Remove(file.Name())
	if err != nil {
		log.Errorf("Erro ao apagar arquivo temporário %s: %s", file.Name(), err.Error())
	}

	err = file.Close()
	if err != nil {
		log.Errorf("Erro ao fechar o arquivo %s: %s", file.Name(), err.Error())
	}
}

func closeResponseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Info("Erro ao encerrar resposta da requisição realizada ao ChatGPT")
	}
}

func getAuthorsOfTheProposition(url string) ([]deputy.Deputy, []organization.Organization, error) {
	response, err := getRequest(url)
	if err != nil {
		return nil, nil, err
	}

	body, err := readResponseBody(response)
	if err != nil {
		return nil, nil, err
	}

	authors, err := getDataFromRequestMap(body)
	if err != nil {
		return nil, nil, err
	}

	deputies, organizations, err := convertAuthorsMapToDeputiesAndOrganizations(authors)
	if err != nil {
		return nil, nil, err
	}

	return deputies, organizations, nil
}

func convertAuthorsMapToDeputiesAndOrganizations(authors []map[string]interface{}) ([]deputy.Deputy,
	[]organization.Organization, error) {
	var deputies []deputy.Deputy
	var organizations []organization.Organization

	for authorIndex, author := range authors {
		log.Infof("Iniciando busca dos dados do %d° autor: %s", authorIndex+1, author["nome"])
		authorData, err := getDataFromUrl(fmt.Sprint(author["uri"]))
		if err != nil {
			return nil, nil, err
		}

		if fmt.Sprint(author["tipo"]) == "Deputado" {
			deputyCode, err := convertInterfaceToInt(authorData["id"])
			if err != nil {
				return nil, nil, err
			}

			authorLastStatus, err := extractDataAsMap(authorData["ultimoStatus"])
			if err != nil {
				return nil, nil, err
			}

			partyMap, err := getDataFromUrl(fmt.Sprint(authorLastStatus["uriPartido"]))
			if err != nil {
				return nil, nil, err
			}

			partyCode, err := convertInterfaceToInt(partyMap["id"])
			if err != nil {
				return nil, nil, err
			}

			partyData, err := party.NewBuilder().
				Code(partyCode).
				Name(fmt.Sprint(partyMap["nome"])).
				Acronym(fmt.Sprint(partyMap["sigla"])).
				ImageUrl(fmt.Sprint(partyMap["urlLogo"])).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do partido %d do(a) deputado(a) %d: %s",
					partyCode, deputyCode, err.Error())
				return nil, nil, err
			}

			deputyData, err := deputy.NewBuilder().
				Code(deputyCode).
				Cpf(fmt.Sprint(authorData["cpf"])).
				Name(fmt.Sprint(authorData["nomeCivil"])).
				ElectoralName(fmt.Sprint(authorLastStatus["nomeEleitoral"])).
				ImageUrl(fmt.Sprint(authorLastStatus["urlFoto"])).
				Party(*partyData).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do(a) deputado(a) %d: %s", deputyCode, err.Error())
				return nil, nil, err
			}

			deputies = append(deputies, *deputyData)
		} else {
			var organizationData *organization.Organization
			if authorData != nil {
				code, err := convertInterfaceToInt(authorData["id"])
				if err != nil {
					return nil, nil, err
				}

				organizationData, err = organization.NewBuilder().
					Code(code).
					Name(fmt.Sprint(authorData["nome"])).
					Acronym(fmt.Sprint(authorData["sigla"])).
					Nickname(fmt.Sprint(authorData["apelido"])).
					Type(fmt.Sprint(author["tipo"])).
					Build()
			} else {
				organizationData, err = organization.NewBuilder().
					Name(fmt.Sprint(author["nome"])).
					Type(fmt.Sprint(author["tipo"])).
					Build()
			}
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da organização %s: %s", fmt.Sprint(authorData["nome"]), err.Error())
				return nil, nil, err
			}

			organizations = append(organizations, *organizationData)
		}
		log.Info("Busca dos dados do autor finalizada com sucesso")
	}

	return deputies, organizations, nil
}

func (instance BackgroundData) RegisterNewNewsletter() {
	referenceDate := time.Now().AddDate(0, 0, -1)
	formattedReferenceDate := referenceDate.Format("02/01/2006")

	log.Info("Iniciando geração do boletim do dia ", formattedReferenceDate)

	propositions, err := instance.propositionRepository.GetPropositionsByDate(referenceDate)
	if err != nil {
		return
	} else if propositions == nil {
		log.Infof("Não foram encontradas proposições no dia %s para gerar um novo boletim", formattedReferenceDate)
		return
	}

	newsletterData, err := generateNewsletter(propositions, referenceDate)
	if err != nil {
		for attempt := 1; attempt <= 3; attempt++ {
			waitingTimeInSeconds := int(math.Pow(5, float64(attempt)))
			log.Warnf("Não foi possível criar o boletim do dia %s na %dª tentativa, tentando novamente em %d segundos",
				formattedReferenceDate, attempt, waitingTimeInSeconds)
			time.Sleep(time.Duration(waitingTimeInSeconds) * time.Second)
			newsletterData, err = generateNewsletter(propositions, referenceDate)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.Error("Erro ao gerar o do boletim do dia ", formattedReferenceDate)
		return
	}

	log.Infof("Boletim do dia %s gerado com sucesso", formattedReferenceDate)

	_, err = instance.newsletterRepository.CreateNewsletter(*newsletterData)
	if err != nil {
		return
	}

	return
}

func generateNewsletter(propositions []proposition.Proposition, referenceDate time.Time) (*newsletter.Newsletter, error) {
	formattedReferenceDate := referenceDate.Format("02/01/2006")

	var contentOfPropositions string
	for count, propositionData := range propositions {
		contentOfPropositions += fmt.Sprintf("Título da %dª matéria: %s\nConteúdo: %s\n\n", count+1, propositionData.Title(),
			propositionData.Content())
	}

	chatGptCommand := "Crie um boletim de notícias em formato de texto corrido a partir das seguintes matérias sobre " +
		"algumas proposições políticas:"
	purpose := fmt.Sprint("Geração do boletim do dia ", formattedReferenceDate)
	newsletterContent, err := requestToChatGpt(chatGptCommand, contentOfPropositions, purpose)
	if err != nil {
		return nil, err
	}

	chatGptCommand = "Gere um único título chamativo para a seguinte matéria para uma revista sobre um conjunto de " +
		"proposições políticas:"
	purpose = fmt.Sprint("Geração do título do boletim do dia ", formattedReferenceDate)
	newsletterTitle, err := requestToChatGpt(chatGptCommand, newsletterContent, purpose)
	if err != nil {
		return nil, err
	}

	newsletterData, err := newsletter.NewBuilder().
		Title(strings.Trim(newsletterTitle, "\"")).
		Content(newsletterContent).
		ReferenceDate(referenceDate).
		Propositions(propositions).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados do boletim do dia %s: %s", formattedReferenceDate, err.Error())
		return nil, err
	}

	return newsletterData, nil
}
