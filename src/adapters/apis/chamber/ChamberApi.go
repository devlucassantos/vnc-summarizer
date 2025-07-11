package chamber

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"net/url"
	"strings"
	"vnc-summarizer/utils/datetime"
	"vnc-summarizer/utils/requesters"
)

type Chamber struct{}

func NewChamberApi() *Chamber {
	return &Chamber{}
}

func (instance Chamber) GetMostRecentPropositions() ([]map[string]interface{}, error) {
	currentDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	var mostRecentPropositionsReturned []map[string]interface{}
	for page := 1; ; page++ {
		chunkSize := 100
		urlOfTheMostRecentPropositions := fmt.Sprintf(
			"https://dadosabertos.camara.leg.br/api/v2/proposicoes?pagina=%d&itens=%d&dataApresentacaoInicio=%s&ordenarPor=id&ordem=asc",
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

	return mostRecentPropositionsReturned, nil
}

func (instance Chamber) GetPropositionByCode(code int) (map[string]interface{}, error) {
	propositionUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/proposicoes/", code)
	proposition, err := requesters.GetDataObjectFromUrl(propositionUrl)
	if err != nil {
		log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	return proposition, nil
}

func (instance Chamber) GetPropositionContentDirectly(propositionUrl string) (string, string, error) {
	parsedPropositionUrl, err := url.Parse(propositionUrl)
	if err != nil {
		log.Errorf("Error parsing proposition URL %s: %s", propositionUrl, err.Error())
		return "", "", err
	}
	queryParams := parsedPropositionUrl.Query()
	propositionContentCode := queryParams.Get("codteor")

	propositionUrl = fmt.Sprintf("https://www.camara.leg.br/internet/ordemdodia/integras/%s.htm",
		propositionContentCode)
	response, err := requesters.GetRequest(propositionUrl)
	if err != nil {
		log.Error("requests.GetRequest(): ", err.Error())
		return "", "", err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error interpreting response to request %s: %s", propositionUrl, err.Error())
		return "", "", err
	}
	responseBodyAsString := string(responseBody)

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error making request to get original proposition content directly: "+
			"[Status: %s; Body: %s]", response.Status, responseBodyAsString)
		log.Error(errorMessage)
		return "", "", errors.New(errorMessage)
	}

	return propositionUrl, responseBodyAsString, err
}

func (instance Chamber) GetPropositionTypes() ([]map[string]interface{}, error) {
	propositionTypesUrl := "https://dadosabertos.camara.leg.br/api/v2/referencias/proposicoes/siglaTipo"
	propositionTypes, err := requesters.GetDataSliceFromUrl(propositionTypesUrl)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	return propositionTypes, nil
}

func (instance Chamber) GetPartyByAcronym(acronym string) (map[string]interface{}, error) {
	partyUrlByAcronym := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/partidos?sigla=", acronym)
	parties, err := requesters.GetDataSliceFromUrl(partyUrlByAcronym)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	return parties[0], nil
}

func (instance Chamber) GetLegislativeBodyByCode(code int) (map[string]interface{}, error) {
	legislativeBodyUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/orgaos/", code)
	legislativeBody, err := requesters.GetDataObjectFromUrl(legislativeBodyUrl)
	if err != nil {
		log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	return legislativeBody, nil
}

func (instance Chamber) GetLegislativeBodyTypes() ([]map[string]interface{}, error) {
	urlOfLegislativeBodyTypes := "https://dadosabertos.camara.leg.br/api/v2/referencias/tiposOrgao"
	legislativeBodyTypes, err := requesters.GetDataSliceFromUrl(urlOfLegislativeBodyTypes)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	return legislativeBodyTypes, nil
}

func (instance Chamber) GetMostRecentVotes() ([]map[string]interface{}, error) {
	currentDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	var mostRecentVotesReturned []map[string]interface{}
	for page := 1; ; page++ {
		chunkSize := 100
		urlOfTheMostRecentVotes := fmt.Sprintf(
			"https://dadosabertos.camara.leg.br/api/v2/votacoes?&pagina=%d&itens=%d&dataInicio=%s&ordenarPor=id&ordem=asc",
			page, chunkSize, currentDateTime.AddDate(0, 0, -1).Format("2006-01-02"),
		)
		mostRecentVotes, err := requesters.GetDataSliceFromUrl(urlOfTheMostRecentVotes)
		if err != nil {
			log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
			return nil, err
		}

		mostRecentVotesReturned = append(mostRecentVotesReturned, mostRecentVotes...)

		if len(mostRecentVotes) < chunkSize {
			break
		}
	}

	return mostRecentVotesReturned, nil
}

func (instance Chamber) GetVotingByCode(code string) (map[string]interface{}, error) {
	votingUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/votacoes/", code)
	voting, err := requesters.GetDataObjectFromUrl(votingUrl)
	if err != nil {
		log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	return voting, nil
}

func (instance Chamber) GetMostRecentEvents() ([]map[string]interface{}, error) {
	currentDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDatetimeInBrazil(): ", err)
		return nil, err
	}

	var mostRecentEventsReturned []map[string]interface{}
	for page := 1; ; page++ {
		chunkSize := 100
		urlOfTheMostRecentEvents := fmt.Sprintf(
			"https://dadosabertos.camara.leg.br/api/v2/eventos?pagina=%d&itens=%d&dataInicio=%s&ordenarPor=id&ordem=asc",
			page, chunkSize, currentDateTime.AddDate(0, 0, -1).Format("2006-01-02"),
		)
		mostRecentEvents, err := requesters.GetDataSliceFromUrl(urlOfTheMostRecentEvents)
		if err != nil {
			log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
			return nil, err
		}

		mostRecentEventsReturned = append(mostRecentEventsReturned, mostRecentEvents...)

		if len(mostRecentEvents) < chunkSize {
			break
		}
	}

	return mostRecentEventsReturned, nil
}

func (instance Chamber) GetEventByCode(code int) (map[string]interface{}, error) {
	eventUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/eventos/", code)
	event, err := requesters.GetDataObjectFromUrl(eventUrl)
	if err != nil {
		log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	return event, nil
}

func (instance Chamber) GetEventsByCodes(eventCodes []string) ([]map[string]interface{}, error) {
	eventsUrl := fmt.Sprintf("https://dadosabertos.camara.leg.br/api/v2/eventos?id=%s&itens=%d",
		strings.Join(eventCodes, ","), len(eventCodes))
	events, err := requesters.GetDataSliceFromUrl(eventsUrl)
	if err != nil {
		log.Error("requesters.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	return events, nil
}

func (instance Chamber) GetEventTypes() ([]map[string]interface{}, error) {
	eventTypesUrl := "https://dadosabertos.camara.leg.br/api/v2/referencias/eventos/codTipoEvento"
	eventTypes, err := requesters.GetDataSliceFromUrl(eventTypesUrl)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	return eventTypes, nil
}

func (instance Chamber) GetEventSituations() ([]map[string]interface{}, error) {
	eventSituationsUrl := "https://dadosabertos.camara.leg.br/api/v2/referencias/situacoesEvento"
	eventSituations, err := requesters.GetDataSliceFromUrl(eventSituationsUrl)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	return eventSituations, nil
}
