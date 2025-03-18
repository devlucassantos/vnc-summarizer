package services

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"strconv"
	"strings"
	"time"
	"vnc-summarizer/core/interfaces/repositories"
	"vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/datetime"
	"vnc-summarizer/core/services/utils/requesters"
)

type Voting struct {
	votingRepository       repositories.Voting
	articleTypeRepository  repositories.ArticleType
	legislativeBodyService services.LegislativeBody
	propositionService     services.Proposition
}

func NewVotingService(votingRepository repositories.Voting, articleTypeRepository repositories.ArticleType,
	legislativeBodyService services.LegislativeBody, propositionService services.Proposition) *Voting {
	return &Voting{
		votingRepository:       votingRepository,
		articleTypeRepository:  articleTypeRepository,
		legislativeBodyService: legislativeBodyService,
		propositionService:     propositionService,
	}
}

func (instance Voting) RegisterNewVotes() {
	codesOfTheMostRecentVotesReturned, err := getCodesOfTheMostRecentVotesRegisteredInTheChamber()
	if err != nil {
		log.Error("getCodesOfTheMostRecentVotesRegisteredInTheChamber(): ", err.Error())
		return
	}

	votesRegistered, err := instance.votingRepository.GetVotesByCodes(codesOfTheMostRecentVotesReturned)
	if err != nil {
		log.Error("votingRepository.GetVotesByCodes(): ", err.Error())
		return
	}

	codesOfTheNewVotes := getCodesOfTheNewVotes(codesOfTheMostRecentVotesReturned, votesRegistered)
	if codesOfTheNewVotes != nil {
		log.Infof("%d new votes were identified for registration: %v", len(codesOfTheNewVotes),
			codesOfTheNewVotes)
	} else {
		log.Info("No new votes were identified for registration")
		return
	}

	for _, votingCode := range codesOfTheNewVotes {
		_, err = instance.RegisterNewVotingByCode(votingCode)
		if err != nil {
			log.Error("RegisterNewVotingByCode(): ", err.Error())
		}
	}

	return
}

func getCodesOfTheMostRecentVotesRegisteredInTheChamber() ([]string, error) {
	log.Info("Starting the search for the most recent votes")

	currentDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	urlOfTheMostRecentVotes := fmt.Sprintf(
		"https://dadosabertos.camara.leg.br/api/v2/votacoes?itens=100&dataInicio=%s&dataFim=%s&ordem=asc&ordenarPor=id",
		currentDateTime.AddDate(0, 0, -1).Format("2006-01-02"),
		currentDateTime.Format("2006-01-02"),
	)
	mostRecentVotesReturned, err := requesters.GetDataSliceFromUrl(urlOfTheMostRecentVotes)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	votingCodes, err := extractVotingCodes(mostRecentVotesReturned)
	if err != nil {
		log.Error("extractVotingCodes(): ", err.Error())
		return nil, err
	}

	log.Info("Successful search for the latest votes: ", votingCodes)
	return votingCodes, nil
}

func extractVotingCodes(votes []map[string]interface{}) ([]string, error) {
	var votingCodes []string
	for _, votingData := range votes {
		votingCodes = append(votingCodes, fmt.Sprint(votingData["id"]))
	}

	return votingCodes, nil
}

func getCodesOfTheNewVotes(returnedVotingCodes []string, registeredVotes []voting.Voting) []string {
	var votingCodesToRegister []string
	for _, votingCode := range returnedVotingCodes {
		var votingAlreadyRegistered bool
		for _, votingData := range registeredVotes {
			if votingData.Code() == votingCode {
				votingAlreadyRegistered = true
				break
			}
		}
		if !votingAlreadyRegistered {
			votingCodesToRegister = append(votingCodesToRegister, votingCode)
		}
	}

	votingCodesToRegister = converters.StringSliceToUniqueStringSlice(votingCodesToRegister)

	return votingCodesToRegister
}

func (instance Voting) RegisterNewVotingByCode(code string) (*uuid.UUID, error) {
	votingData, err := instance.getVotingDataToRegister(code)
	if err != nil {
		log.Errorf("Error retrieving data for voting %s: %s", code, err.Error())
		return nil, err
	}

	votingId, err := instance.votingRepository.CreateVoting(*votingData)
	if err != nil {
		log.Error("votingRepository.CreateVoting(): ", err.Error())
		return nil, err
	}

	return votingId, nil
}

func (instance Voting) getVotingDataToRegister(code string) (*voting.Voting, error) {
	log.Info("Starting data search for voting ", code)

	votingUrl := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/votacoes/", code)
	votingData, err := requesters.GetDataObjectFromUrl(votingUrl)
	if err != nil {
		log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	resultAnnouncedAtAsString := fmt.Sprint(votingData["dataHoraRegistro"])
	if resultAnnouncedAtAsString == "<nil>" {
		mainPropositionData, err := getMainPropositionData(votingData["ultimaApresentacaoProposicao"])
		if err != nil {
			log.Error("getMainPropositionData(): ", err.Error())
			return nil, err
		}

		resultAnnouncedAtAsString = fmt.Sprint(mainPropositionData["dataHoraRegistro"])
	}

	resultAnnouncedAt, err := time.Parse("2006-01-02T15:04:05", resultAnnouncedAtAsString)
	if err != nil {
		log.Errorf("Error converting date and time of result announcement of voting %s: %s", code,
			err.Error())
		return nil, err
	}

	isApproved, err := strconv.ParseBool(fmt.Sprint(votingData["aprovacao"]))
	if err != nil {
		log.Errorf("Error converting result of voting %s to boolean: %s", code, err.Error())
	}

	legislativeBodyCode, err := converters.ToInt(votingData["idOrgao"])
	if err != nil {
		log.Error("converters.ToInt(): ", err.Error())
		return nil, err
	}

	legislativeBody, err := instance.legislativeBodyService.GetLegislativeBodyByCode(legislativeBodyCode)
	if err != nil {
		log.Error("legislativeBodyService.GetLegislativeBodyDataByCode(): ", err.Error())
		return nil, err
	}

	if legislativeBody == nil {
		legislativeBodyId, err := instance.legislativeBodyService.RegisterNewLegislativeBodyByCode(legislativeBodyCode)
		if err != nil {
			log.Error("legislativeBodyService.RegisterNewLegislativeBodyByCode(): ", err.Error())
			return nil, err
		}

		legislativeBodyRegistered, err := legislativebody.NewBuilder().Id(*legislativeBodyId).Build()
		if err != nil {
			log.Errorf("Error updating legislative body %s: %s", legislativeBodyId, err.Error())
			return nil, err
		}
		legislativeBody = legislativeBodyRegistered
	}

	mainProposition, relatedPropositions, affectedPropositions, err := instance.getVotingRelatedPropositions(votingData)
	if err != nil {
		log.Error("getVotingRelatedPropositions(): ", err.Error())
		return nil, err
	}

	articleTypeCode := "voting"
	articleType, err := instance.articleTypeRepository.GetArticleTypeByCode(articleTypeCode)
	if err != nil {
		log.Error("articleTypeRepository.GetArticleTypeByCode(): ", err.Error())
		return nil, err
	}

	articleData, err := article.NewBuilder().Type(*articleType).Build()
	if err != nil {
		log.Errorf("Error validating article data for voting %s: %s", code, err.Error())
		return nil, err
	}

	votingBuilder := voting.NewBuilder().
		Code(code).
		Result(fmt.Sprint(votingData["descricao"])).
		ResultAnnouncedAt(resultAnnouncedAt).
		IsApproved(isApproved).
		LegislativeBody(*legislativeBody).
		RelatedPropositions(relatedPropositions).
		AffectedPropositions(affectedPropositions).
		Article(*articleData)

	if mainProposition != nil {
		votingBuilder.MainProposition(*mainProposition)
	}

	votingDomain, err := votingBuilder.Build()
	if err != nil {
		log.Errorf("Error validating data for voting %s: %s", code, err.Error())
		return nil, err
	}

	log.Infof("Data search for voting %s successful", code)
	return votingDomain, err
}

func (instance Voting) getVotingRelatedPropositions(votingData map[string]interface{}) (*proposition.Proposition,
	[]proposition.Proposition, []proposition.Proposition, error) {
	var mainPropositionData map[string]interface{}
	var err error
	mainPropositionData, err = getMainPropositionData(votingData["ultimaApresentacaoProposicao"])
	if err != nil {
		log.Error("getMainPropositionData(): ", err.Error())
		return nil, nil, nil, err
	}

	relatedPropositionsMapSlice, err := converters.ToMapSlice(votingData["objetosPossiveis"])
	if err != nil {
		log.Error("converters.ToMapSlice(): ", err.Error())
		return nil, nil, nil, err
	}

	affectedPropositionsMapSlice, err := converters.ToMapSlice(votingData["proposicoesAfetadas"])
	if err != nil {
		log.Error("converters.ToMapSlice(): ", err.Error())
		return nil, nil, nil, err
	}

	votingRelatedPropositionCodes, err := getVotingRelatedPropositionCodes(mainPropositionData,
		relatedPropositionsMapSlice, affectedPropositionsMapSlice)
	if err != nil {
		log.Error("getVotingRelatedPropositionCodes(): ", err.Error())
		return nil, nil, nil, err
	}

	relatedPropositionsAlreadyRegistered, err := instance.propositionService.GetPropositionsByCodes(
		votingRelatedPropositionCodes)
	if err != nil {
		log.Error("propositionRepository.GetPropositionsByCodes(): ", err.Error())
		return nil, nil, nil, err
	}

	codesOfThePropositionsToRegister := getCodesOfTheNewPropositions(votingRelatedPropositionCodes,
		relatedPropositionsAlreadyRegistered)

	registeredPropositions, err := instance.registerNewVotingRelatedPropositions(codesOfThePropositionsToRegister)
	if err != nil {
		log.Error("registerNewVotingRelatedPropositions(): ", err.Error())
		return nil, nil, nil, err
	}

	votingRelatedPropositions := append(relatedPropositionsAlreadyRegistered, registeredPropositions...)

	mainProposition, relatedPropositions, affectedPropositions := extractVotingRelatedPropositions(mainPropositionData,
		relatedPropositionsMapSlice, affectedPropositionsMapSlice, votingRelatedPropositions)

	return &mainProposition, relatedPropositions, affectedPropositions, nil
}

func getMainPropositionData(mainProposition interface{}) (map[string]interface{}, error) {
	mainPropositionMap, err := converters.ToMap(mainProposition)
	if err != nil {
		log.Error("converters.ToMap(): ", err.Error())
		return nil, err
	}

	var mainPropositionData map[string]interface{}
	UrlOfTheMainProposition := fmt.Sprint(mainPropositionMap["uriProposicaoCitada"])
	if UrlOfTheMainProposition != "<nil>" {
		mainPropositionData, err = requesters.GetDataObjectFromUrl(UrlOfTheMainProposition)
		if err != nil {
			log.Error("requests.GetDataObjectFromUrl(): ", err.Error())
			return nil, err
		}
	}

	return mainPropositionData, nil
}

func getVotingRelatedPropositionCodes(mainProposition map[string]interface{},
	relatedPropositions, affectedPropositions []map[string]interface{}) ([]int, error) {
	var votingRelatedPropositionCodes []int

	if mainProposition != nil {
		mainPropositionCode, err := converters.ToInt(mainProposition["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		votingRelatedPropositionCodes = append(votingRelatedPropositionCodes, mainPropositionCode)
	}

	for _, relatedProposition := range relatedPropositions {
		relatedPropositionCode, err := converters.ToInt(relatedProposition["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		votingRelatedPropositionCodes = append(votingRelatedPropositionCodes, relatedPropositionCode)
	}

	for _, affectedProposition := range affectedPropositions {
		affectedPropositionCode, err := converters.ToInt(affectedProposition["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		votingRelatedPropositionCodes = append(votingRelatedPropositionCodes, affectedPropositionCode)
	}

	return votingRelatedPropositionCodes, nil
}

func (instance Voting) registerNewVotingRelatedPropositions(codesOfThePropositionsToRegister []int) (
	[]proposition.Proposition, error) {
	var registeredPropositions []proposition.Proposition
	for _, propositionCode := range codesOfThePropositionsToRegister {
		propositionId, err := instance.propositionService.RegisterNewPropositionByCode(propositionCode)
		if err != nil {
			if strings.Contains(err.Error(), "no content") {
				continue
			}
			log.Error("propositionService.RegisterNewPropositionByCode(): ", err.Error())
			return nil, err
		}

		propositionDomain, err := proposition.NewBuilder().Id(*propositionId).Code(propositionCode).Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s: %s", propositionId, err.Error())
			return nil, err
		}
		registeredPropositions = append(registeredPropositions, *propositionDomain)
	}

	return registeredPropositions, nil
}

func extractVotingRelatedPropositions(mainPropositionMap map[string]interface{},
	relatedPropositionsMapSlice, affectedPropositionsMapSlice []map[string]interface{},
	propositionsAlreadyRegistered []proposition.Proposition) (proposition.Proposition, []proposition.Proposition,
	[]proposition.Proposition) {
	var mainProposition proposition.Proposition
	var relatedPropositions, affectedPropositions []proposition.Proposition
	for _, propositionData := range propositionsAlreadyRegistered {
		if mainPropositionMap != nil {
			if propositionData.Code() == int(mainPropositionMap["id"].(float64)) {
				mainProposition = propositionData
			}
		}

		for _, relatedPropositionData := range relatedPropositionsMapSlice {
			if propositionData.Code() == int(relatedPropositionData["id"].(float64)) {
				relatedPropositions = append(relatedPropositions, propositionData)
			}
		}

		for _, affectedPropositionData := range affectedPropositionsMapSlice {
			if propositionData.Code() == int(affectedPropositionData["id"].(float64)) {
				affectedPropositions = append(affectedPropositions, propositionData)
			}
		}
	}

	return mainProposition, relatedPropositions, affectedPropositions
}

func (instance Voting) GetVotesByCodes(codes []string) ([]voting.Voting, error) {
	votes, err := instance.votingRepository.GetVotesByCodes(codes)
	if err != nil {
		log.Error("votingRepository.GetVotesByCodes(): ", err.Error())
		return nil, err
	}

	return votes, err
}
