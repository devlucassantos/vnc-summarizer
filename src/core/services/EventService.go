package services

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/agendaitemregime"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/devlucassantos/vnc-domains/src/domains/eventagendaitem"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"path"
	"strings"
	"time"
	"vnc-summarizer/core/interfaces/chamber"
	"vnc-summarizer/core/interfaces/chatgpt"
	"vnc-summarizer/core/interfaces/postgres"
	"vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/utils/converters"
	"vnc-summarizer/utils/requesters"
	"vnc-summarizer/utils/splitters"
	"vnc-summarizer/utils/validators"
)

type Event struct {
	deputyService              services.Deputy
	legislativeBodyService     services.LegislativeBody
	propositionService         services.Proposition
	votingService              services.Voting
	chamberApi                 chamber.Chamber
	chatGptApi                 chatgpt.ChatGpt
	eventRepository            postgres.Event
	articleTypeRepository      postgres.ArticleType
	eventTypeRepository        postgres.EventType
	eventSituationRepository   postgres.EventSituation
	agendaItemRegimeRepository postgres.AgendaItemRegime
}

func NewEventService(deputyService services.Deputy, legislativeBodyService services.LegislativeBody,
	propositionService services.Proposition, votingService services.Voting, chamberApi chamber.Chamber,
	chatGptApi chatgpt.ChatGpt, eventRepository postgres.Event, articleTypeRepository postgres.ArticleType,
	eventTypeRepository postgres.EventType, eventSituationRepository postgres.EventSituation,
	agendaItemRegimeRepository postgres.AgendaItemRegime) *Event {
	return &Event{
		deputyService:              deputyService,
		legislativeBodyService:     legislativeBodyService,
		propositionService:         propositionService,
		votingService:              votingService,
		chamberApi:                 chamberApi,
		chatGptApi:                 chatGptApi,
		eventRepository:            eventRepository,
		articleTypeRepository:      articleTypeRepository,
		eventTypeRepository:        eventTypeRepository,
		eventSituationRepository:   eventSituationRepository,
		agendaItemRegimeRepository: agendaItemRegimeRepository,
	}
}

func (instance Event) RegisterNewEvents() {
	codesOfTheMostRecentEventsReturned, err := instance.getCodesOfTheMostRecentEventsRegisteredInTheChamber()
	if err != nil {
		log.Error("getCodesOfTheMostRecentEventsRegisteredInTheChamber(): ", err.Error())
		return
	} else if codesOfTheMostRecentEventsReturned == nil {
		log.Info("No new events were identified for registration")
		return
	}

	eventsRegistered, err := instance.eventRepository.GetEventsByCodes(codesOfTheMostRecentEventsReturned)
	if err != nil {
		log.Error("eventRepository.GetEventsByCodes(): ", err.Error())
		return
	}

	codesOfTheNewEvents := getCodesOfTheNewEvents(codesOfTheMostRecentEventsReturned, eventsRegistered)
	if codesOfTheNewEvents != nil {
		log.Infof("%d new events were identified for registration: %v", len(codesOfTheNewEvents),
			codesOfTheNewEvents)
	} else {
		log.Info("No new events were identified for registration")
		return
	}

	for _, eventCode := range codesOfTheNewEvents {
		_, err = instance.RegisterNewEventByCode(eventCode)
		if err != nil {
			log.Error("RegisterNewEventByCode(): ", err.Error())
		}
	}

	return
}

func (instance Event) getCodesOfTheMostRecentEventsRegisteredInTheChamber() ([]int, error) {
	log.Info("Starting the search for the most recent events")

	mostRecentEvents, err := instance.chamberApi.GetMostRecentEvents()
	if err != nil {
		log.Error("chamberApi.GetMostRecentEvents(): ", err.Error())
		return nil, err
	}

	eventCodes, err := extractCodesFromEvents(mostRecentEvents)
	if err != nil {
		log.Error("extractCodesFromEvents(): ", err.Error())
		return nil, err
	}

	log.Info("Successful search for the most recent events: ", eventCodes)
	return eventCodes, nil
}

func extractCodesFromEvents(events []map[string]interface{}) ([]int, error) {
	var eventCodes []int
	for _, eventData := range events {
		eventCode, err := converters.ToInt(eventData["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		eventCodes = append(eventCodes, eventCode)
	}

	return eventCodes, nil
}

func getCodesOfTheNewEvents(returnedEventCodes []int, registeredEvents []event.Event) []int {
	var eventCodesToRegister []int
	for _, eventCode := range returnedEventCodes {
		var eventAlreadyRegistered bool
		for _, eventData := range registeredEvents {
			if eventData.Code() == eventCode {
				eventAlreadyRegistered = true
				break
			}
		}
		if !eventAlreadyRegistered {
			eventCodesToRegister = append(eventCodesToRegister, eventCode)
		}
	}

	eventCodesToRegister = converters.IntSliceToUniqueIntSlice(eventCodesToRegister)

	return eventCodesToRegister
}

func (instance Event) RegisterNewEventByCode(code int) (*uuid.UUID, error) {
	eventData, err := instance.getEventDataToRegister(code)
	if err != nil {
		log.Errorf("Error retrieving data for event %d: %s", code, err.Error())
		return nil, err
	} else if eventData == nil {
		return nil, nil
	}

	eventId, err := instance.eventRepository.CreateEvent(*eventData)
	if err != nil {
		log.Error("eventRepository.CreateEvent(): ", err.Error())
		return nil, err
	}

	return eventId, nil
}

func (instance Event) getEventDataToRegister(code int) (*event.Event, error) {
	log.Info("Starting data search for event ", code)

	eventData, err := instance.chamberApi.GetEventByCode(code)
	if err != nil {
		log.Error("chamberApi.GetEventByCode(): ", err.Error())
		return nil, err
	}

	startsAt, err := time.Parse("2006-01-02T15:04", fmt.Sprint(eventData["dataHoraInicio"]))
	if err != nil {
		log.Errorf("Error converting date and time of start of event %d: %s", code, err.Error())
		return nil, err
	}

	var endsAt time.Time
	endsAtAsString := fmt.Sprint(eventData["dataHoraFim"])
	if endsAtAsString != "<nil>" {
		endsAt, err = time.Parse("2006-01-02T15:04", endsAtAsString)
		if err != nil {
			log.Errorf("Error converting date and time of end of event %d: %s", code, err.Error())
			return nil, err
		}
	}

	location, isInternal, err := getEventLocation(eventData)
	if err != nil {
		log.Error("getEventLocation(): ", err.Error())
		return nil, err
	}

	eventTypeDescription := fmt.Sprint(eventData["descricaoTipo"])
	eventType, specificType, err := instance.getEventTypeByDescription(eventTypeDescription)
	if err != nil {
		log.Error("getEventTypeByDescription(): ", err.Error())
		return nil, err
	}

	eventSituationDescription := fmt.Sprint(eventData["situacao"])
	eventSituation, specificSituation, err := instance.getEventSituationByDescription(eventSituationDescription)
	if err != nil {
		log.Error("getEventSituationByDescription(): ", err.Error())
		return nil, err
	}

	legislativeBodyData, err := converters.ToMapSlice(eventData["orgaos"])
	if err != nil {
		log.Error("converters.ToMapSlice(): ", err.Error())
		return nil, err
	}

	legislativeBodies, err := instance.getLegislativeBodies(legislativeBodyData)
	if err != nil {
		log.Error("getLegislativeBodies(): ", err.Error())
		return nil, err
	}

	requirementData, err := converters.ToMapSlice(eventData["requerimentos"])
	if err != nil {
		log.Error("converters.ToMapSlice(): ", err.Error())
		return nil, err
	}

	requirements, err := instance.getEventRequirements(requirementData)
	if err != nil {
		log.Error("getEventRequirements(): ", err.Error())
		return nil, err
	}

	agendaItemUrl := fmt.Sprint(eventData["urlDocumentoPauta"])
	agendaItems, err := instance.getEventAgendaItems(agendaItemUrl)
	if err != nil {
		log.Error("getEventAgendaItems(): ", err.Error())
		return nil, err
	}

	eventTopics, err := instance.getEventTopics(requirements, agendaItems)
	if err != nil {
		log.Error("getEventTopics(): ", err.Error())
		return nil, err
	}

	if eventTopics == "" {
		log.Warnf("Event %d not registered as there were no propositions related to it", code)
		return nil, nil
	}

	chatGptCommand := "Gere um título que usando uma linguagem simples e direta seja chamativo para uma matéria " +
		"jornalistica sobre um evento que tratou dos seguintes temas:"
	purpose := fmt.Sprint("Generating the title of event ", code)
	title, err := instance.chatGptApi.MakeRequest(chatGptCommand, eventTopics, purpose)
	if err != nil {
		log.Error("chatGptApi.MakeRequest(): ", err.Error())
		return nil, err
	}
	title = strings.Trim(title, "*\"")

	articleTypeCode := "event"
	articleType, err := instance.articleTypeRepository.GetArticleTypeByCode(articleTypeCode)
	if err != nil {
		log.Error("articleTypeRepository.GetArticleTypeByCode(): ", err.Error())
		return nil, err
	}

	articleData, err := article.NewBuilder().Type(*articleType).ReferenceDateTime(startsAt).Build()
	if err != nil {
		log.Errorf("Error validating article data for event %d: %s", code, err.Error())
		return nil, err
	}

	eventBuilder := event.NewBuilder().
		Code(code).
		Title(title).
		Description(fmt.Sprint(eventData["descricao"])).
		StartsAt(startsAt).
		Location(location).
		IsInternal(isInternal).
		SpecificType(specificType).
		Type(*eventType).
		SpecificSituation(specificSituation).
		Situation(*eventSituation).
		LegislativeBodies(legislativeBodies).
		Requirements(requirements).
		AgendaItems(agendaItems).
		Article(*articleData)

	if !endsAt.IsZero() {
		eventBuilder.EndsAt(endsAt)
	}

	videoUrl := fmt.Sprint(eventData["urlRegistro"])
	if videoUrl != "<nil>" {
		if validators.IsUrlValid(videoUrl) {
			eventBuilder.VideoUrl(videoUrl)
		} else {
			log.Infof("Event %d video URL is invalid", code)
		}
	}

	eventDomain, err := eventBuilder.Build()
	if err != nil {
		log.Errorf("Error validating data for event %d: %s", code, err.Error())
		return nil, err
	}

	log.Infof("Data search for event %d successful", code)
	return eventDomain, err
}

func getEventLocation(eventData map[string]interface{}) (string, bool, error) {
	locationInTheChamber, err := converters.ToMap(eventData["localCamara"])
	if err != nil {
		log.Error("converters.ToMap(): ", err.Error())
		return "", false, err
	}

	location := fmt.Sprint(locationInTheChamber["nome"])
	isInternal := true
	if location == "<nil>" {
		location = fmt.Sprint(eventData["localExterno"])
		isInternal = false
	}

	return location, isInternal, nil
}

func (instance Event) getEventTypeByDescription(eventTypeDescription string) (*eventtype.EventType, string, error) {
	eventTypes, err := instance.chamberApi.GetEventTypes()
	if err != nil {
		log.Error("chamberApi.GetEventTypes(): ", err.Error())
		return nil, "", err
	}

	var eventTypeCode string
	for _, eventTypeMap := range eventTypes {
		if fmt.Sprint(eventTypeMap["nome"]) == eventTypeDescription {
			eventTypeCode = fmt.Sprint(eventTypeMap["cod"])
			break
		}
	}

	eventType, err := instance.eventTypeRepository.GetEventTypeByCodeOrDefaultType(eventTypeCode)
	if err != nil {
		log.Error("eventTypeRepository.GetEventTypeByCodeOrDefaultType(): ", err.Error())
		return nil, "", err
	}

	eventSpecificType := fmt.Sprintf("%s (%s)", eventTypeDescription, eventTypeCode)

	return eventType, eventSpecificType, nil
}

func (instance Event) getEventSituationByDescription(eventSituationDescription string) (*eventsituation.EventSituation,
	string, error) {
	eventSituations, err := instance.chamberApi.GetEventSituations()
	if err != nil {
		log.Error("chamberApi.GetEventSituations(): ", err.Error())
		return nil, "", err
	}

	var eventSituationCode string
	for _, eventSituationMap := range eventSituations {
		if strings.TrimSpace(fmt.Sprint(eventSituationMap["nome"])) == eventSituationDescription {
			eventSituationCode = fmt.Sprint(eventSituationMap["cod"])
			break
		}
	}

	eventSituation, err := instance.eventSituationRepository.GetEventSituationByCodeOrDefaultSituation(eventSituationCode)
	if err != nil {
		log.Error("eventSituationRepository.GetEventSituationByCodeOrDefaultSituation(): ", err.Error())
		return nil, "", err
	}

	eventSpecificSituation := fmt.Sprintf("%s (%s)", eventSituationDescription, eventSituationCode)

	return eventSituation, eventSpecificSituation, nil
}

func (instance Event) getLegislativeBodies(legislativeBodyData []map[string]interface{}) (
	[]legislativebody.LegislativeBody, error) {
	codesOfTheReturnedLegislativeBodies, err := extractCodesFromLegislativeBodies(legislativeBodyData)
	if err != nil {
		log.Error("extractCodesFromLegislativeBodies(): ", err.Error())
		return nil, err
	}

	var registeredLegislativeBodies []legislativebody.LegislativeBody
	if codesOfTheReturnedLegislativeBodies != nil {
		registeredLegislativeBodies, err = instance.legislativeBodyService.GetLegislativeBodiesByCodes(
			codesOfTheReturnedLegislativeBodies)
		if err != nil {
			log.Error("legislativeBodyService.GetLegislativeBodiesByCodes(): ", err.Error())
			return nil, err
		}

		codesOfTheLegislativeBodiesToRegister := getCodesOfTheNewLegislativeBodies(codesOfTheReturnedLegislativeBodies,
			registeredLegislativeBodies)
		for _, legislativeBodyCode := range codesOfTheLegislativeBodiesToRegister {
			legislativeBodyId, err := instance.legislativeBodyService.RegisterNewLegislativeBodyByCode(legislativeBodyCode)
			if err != nil {
				log.Error("legislativeBodyService.RegisterNewLegislativeBodyByCode(): ", err.Error())
				return nil, err
			}

			legislativeBodyDomain, err := legislativebody.NewBuilder().
				Id(*legislativeBodyId).
				Code(legislativeBodyCode).
				Build()
			if err != nil {
				log.Errorf("Error validating data for legislative body %s: %s", legislativeBodyId, err.Error())
				return nil, err
			}
			registeredLegislativeBodies = append(registeredLegislativeBodies, *legislativeBodyDomain)
		}
	}

	return registeredLegislativeBodies, nil
}

func extractCodesFromLegislativeBodies(legislativeBodyData []map[string]interface{}) ([]int, error) {
	var legislativeBodyCodes []int
	for _, legislativeBody := range legislativeBodyData {
		legislativeBodyCode, err := converters.ToInt(legislativeBody["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		legislativeBodyCodes = append(legislativeBodyCodes, legislativeBodyCode)
	}

	return legislativeBodyCodes, nil
}

func (instance Event) getEventRequirements(requirements []map[string]interface{}) ([]proposition.Proposition, error) {
	propositionCodes, err := getPropositionCodesFromEventRequirements(requirements)
	if err != nil {
		log.Error("getPropositionCodesFromEventRequirements(): ", err.Error())
		return nil, err
	}

	var propositions []proposition.Proposition
	if propositionCodes != nil {
		propositions, err = instance.propositionService.GetPropositionsByCodes(propositionCodes)
		if err != nil {
			log.Error("propositionService.GetPropositionsByCodes(): ", err.Error())
			return nil, err
		}

		codesOfThePropositionsToRegister := getCodesOfTheNewPropositions(propositionCodes, propositions)
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
			propositions = append(propositions, *propositionDomain)
		}
	}

	return propositions, nil
}

func getPropositionCodesFromEventRequirements(requirements []map[string]interface{}) ([]int, error) {
	var propositionCodes []int
	for _, requirement := range requirements {
		propositionCode := path.Base(fmt.Sprint(requirement["uri"]))
		code, err := converters.ToInt(propositionCode)
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		propositionCodes = append(propositionCodes, code)
	}

	return propositionCodes, nil
}

func (instance Event) getEventAgendaItems(agendaItemUrl string) ([]eventagendaitem.EventAgendaItem,
	error) {
	agendaItemData, err := requesters.GetDataSliceFromUrl(agendaItemUrl)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	} else if fmt.Sprint(agendaItemData) == "[]" {
		return nil, nil
	}

	propositionsRelatedToTheEventAgendaItems, err := instance.getPropositionsRelatedToTheEventAgendaItems(
		agendaItemData)
	if err != nil {
		log.Error("getPropositionsRelatedToTheEventAgendaItems(): ", err.Error())
		return nil, err
	}

	votesRelatedToTheEventAgendaItems, err := instance.getVotingRelatedToTheEventAgendaItem(agendaItemData)
	if err != nil {
		log.Error("instance.getVotingRelatedToTheEventAgendaItem(): ", err.Error())
		return nil, err
	}

	var agendaItems []eventagendaitem.EventAgendaItem
	for _, agendaItem := range agendaItemData {
		agendaItemRegimeCode, err := converters.ToInt(agendaItem["codRegime"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}

		agendaItemRegime, err := instance.getAgendaItemRegime(agendaItemRegimeCode, fmt.Sprint(agendaItem["regime"]))
		if err != nil {
			log.Error("instance.getAgendaItemRegime(): ", err.Error())
			return nil, err
		}

		var rapporteur deputy.Deputy
		if fmt.Sprint(agendaItem["relator"]) != "<nil>" {
			rapporteurData, err := converters.ToMap(agendaItem["relator"])
			if err != nil {
				log.Error("converters.ToMap(): ", err.Error())
				return nil, err
			}

			agendaItemRapporteur, err := instance.deputyService.GetDeputyFromDeputyData(rapporteurData)
			if err != nil {
				log.Error("deputyService.GetDeputyFromDeputyData(): ", err.Error())
				return nil, err
			}

			rapporteur = *agendaItemRapporteur
		}

		propositionData, err := converters.ToMap(agendaItem["proposicao_"])
		if err != nil {
			log.Error("converters.ToMap(): ", err.Error())
			return nil, err
		}

		relatedPropositionData, err := converters.ToMap(agendaItem["proposicaoRelacionada_"])
		if err != nil {
			log.Error("converters.ToMap(): ", err.Error())
			return nil, err
		}

		var agendaItemVoting voting.Voting
		if votesRelatedToTheEventAgendaItems != nil {
			agendaItemVoting = votesRelatedToTheEventAgendaItems[path.Base(fmt.Sprint(agendaItem["uriVotacao"]))]
		}

		mainPropositionId, err := converters.ToInt(propositionData["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}

		var relatedPropositionId int
		if fmt.Sprint(relatedPropositionData["id"]) != "<nil>" {
			relatedPropositionId, err = converters.ToInt(relatedPropositionData["id"])
			if err != nil {
				log.Error("converters.ToInt(): ", err.Error())
				return nil, err
			}
		}

		agendaItemBuilder := eventagendaitem.NewBuilder().
			Title(fmt.Sprint(agendaItem["titulo"])).
			Topic(fmt.Sprint(agendaItem["topico"])).
			Regime(*agendaItemRegime).
			Rapporteur(rapporteur).
			Proposition(propositionsRelatedToTheEventAgendaItems[mainPropositionId]).
			RelatedProposition(propositionsRelatedToTheEventAgendaItems[relatedPropositionId]).
			Voting(agendaItemVoting)

		situation := fmt.Sprint(agendaItem["situacaoItem"])
		if situation != "<nil>" {
			agendaItemBuilder.Situation(situation)
		}

		agendaItemDomain, err := agendaItemBuilder.Build()
		if err != nil {
			log.Errorf("Error validating data for agenda item %s: %s", fmt.Sprint(agendaItem["title"]),
				err.Error())
			return nil, err
		}
		agendaItems = append(agendaItems, *agendaItemDomain)
	}

	return agendaItems, nil
}

func (instance Event) getAgendaItemRegime(code int, description string) (*agendaitemregime.AgendaItemRegime, error) {
	agendaItemRegime, err := instance.agendaItemRegimeRepository.GetAgendaItemRegimeByCode(code)
	if err != nil {
		log.Error("agendaItemRegimeRepository.GetAgendaItemRegimeByCode(): ", err.Error())
		return nil, err
	}

	if agendaItemRegime == nil {
		agendaItemRegime, err = agendaitemregime.NewBuilder().Code(code).Description(description).Build()
		if err != nil {
			log.Errorf("Error validating data for agenda item regime %d: %s", code, err.Error())
			return nil, err
		}

		agendaItemRegimeId, err := instance.agendaItemRegimeRepository.CreateAgendaItemRegime(*agendaItemRegime)
		if err != nil {
			log.Error("agendaItemRegimeRepository.CreateAgendaItemRegime(): ", err.Error())
			return nil, err
		}

		agendaItemRegime, err = agendaItemRegime.NewUpdater().Id(*agendaItemRegimeId).Build()
		if err != nil {
			log.Errorf("Error updating legislative body type %s: %s", agendaItemRegimeId, err.Error())
			return nil, err
		}
	}

	return agendaItemRegime, nil
}

func (instance Event) getPropositionsRelatedToTheEventAgendaItems(agendaItems []map[string]interface{}) (
	map[int]proposition.Proposition, error) {
	var propositionCodes []int
	for _, agendaItemData := range agendaItems {
		propositionData, err := converters.ToMap(agendaItemData["proposicao_"])
		if err != nil {
			log.Error("converters.ToMap(): ", err.Error())
			return nil, err
		}

		propositionCode, err := converters.ToInt(propositionData["id"])
		if err != nil {
			log.Error("converters.ToInt(): ", err.Error())
			return nil, err
		}
		propositionCodes = append(propositionCodes, propositionCode)

		if fmt.Sprint(agendaItemData["proposicaoRelacionada_"]) != "<nil>" {
			relatedPropositionData, err := converters.ToMap(agendaItemData["proposicaoRelacionada_"])
			if err != nil {
				log.Error("converters.ToMap(): ", err.Error())
				return nil, err
			}

			propositionCode, err = converters.ToInt(relatedPropositionData["id"])
			if err != nil {
				log.Error("converters.ToInt(): ", err.Error())
				return nil, err
			}
			propositionCodes = append(propositionCodes, propositionCode)
		}
	}

	propositions, err := instance.propositionService.GetPropositionsByCodes(propositionCodes)
	if err != nil {
		log.Error("propositionService.GetPropositionsByCodes(): ", err.Error())
		return nil, err
	}

	codesOfThePropositionsToRegister := getCodesOfTheNewPropositions(propositionCodes, propositions)
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
		propositions = append(propositions, *propositionDomain)
	}

	propositionsRelatedToTheEventAgendaItems := map[int]proposition.Proposition{}
	for _, propositionData := range propositions {
		propositionsRelatedToTheEventAgendaItems[propositionData.Code()] = propositionData
	}

	return propositionsRelatedToTheEventAgendaItems, nil
}

func (instance Event) getVotingRelatedToTheEventAgendaItem(agendaItems []map[string]interface{}) (
	map[string]voting.Voting, error) {
	var returnedVotingCodes []string
	for _, agendaItemData := range agendaItems {
		votingCode := path.Base(fmt.Sprint(agendaItemData["uriVotacao"]))
		if votingCode == "<nil>" {
			continue
		}
		returnedVotingCodes = append(returnedVotingCodes, votingCode)
	}

	if returnedVotingCodes == nil {
		return nil, nil
	}

	registeredVotes, err := instance.votingService.GetVotesByCodes(returnedVotingCodes)
	if err != nil {
		log.Error("votingService.GetVotesByCodes(): ", err.Error())
		return nil, err
	}

	votesToRegister := getCodesOfTheNewVotes(returnedVotingCodes, registeredVotes)
	for _, votingCode := range votesToRegister {
		votingId, err := instance.votingService.RegisterNewVotingByCode(votingCode)
		if err != nil {
			log.Error("votingService.RegisterNewVotingByCode(): ", err.Error())
			return nil, err
		}

		votingDomain, err := voting.NewBuilder().Id(*votingId).Code(votingCode).Build()
		if err != nil {
			log.Errorf("Error validating data for voting %s: %s", votingId, err.Error())
			return nil, err
		}
		registeredVotes = append(registeredVotes, *votingDomain)
	}

	votesRelatedToTheEventAgendaItems := map[string]voting.Voting{}
	for _, votingData := range registeredVotes {
		votesRelatedToTheEventAgendaItems[votingData.Code()] = votingData
	}

	return votesRelatedToTheEventAgendaItems, nil
}

func (instance Event) getEventTopics(requirements []proposition.Proposition,
	agendaItems []eventagendaitem.EventAgendaItem) (string, error) {
	var propositionCodes []int
	for _, requirement := range requirements {
		propositionCodes = append(propositionCodes, requirement.Code())
	}
	for _, agendaItem := range agendaItems {
		agendaItemProposition := agendaItem.Proposition()
		propositionCodes = append(propositionCodes, agendaItemProposition.Code())
	}

	var eventTopics string
	if propositionCodes != nil {
		propositions, err := instance.propositionService.GetPropositionsByCodes(propositionCodes)
		if err != nil {
			log.Error("propositionService.GetPropositionsByCodes(): ", err.Error())
			return "", err
		}

		for _, propositionData := range propositions {
			eventTopics += fmt.Sprint(propositionData.Title(), "\n\n")
		}
	}

	return eventTopics, nil
}

func (instance Event) UpdateEventsOccurringToday() {
	events, err := instance.eventRepository.GetEventsOccurringToday()
	if err != nil {
		log.Error("eventRepository.GetEventsOccurringToday(): ", err.Error())
		return
	}

	instance.updateEvents(events)
}

func (instance Event) UpdateEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished() {
	events, err := instance.eventRepository.GetEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished()
	if err != nil {
		log.Error("eventRepository.GetEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished(): ", err.Error())
		return
	}

	instance.updateEvents(events)
}

func (instance Event) updateEvents(events []event.Event) {
	if len(events) == 0 {
		log.Info("No events found for update")
		return
	}

	eventsToUpdate, err := instance.getEventsToUpdate(events)
	if err != nil {
		log.Error("getEventsToUpdate(): ", err.Error())
		return
	}

	if len(eventsToUpdate) == 0 {
		log.Info("No outdated events found")
		return
	}

	for _, eventData := range eventsToUpdate {
		err = instance.eventRepository.UpdateEvent(eventData)
		if err != nil {
			log.Error("eventRepository.UpdateEvent(): ", err.Error())
		}
	}
}

func (instance Event) getEventsToUpdate(events []event.Event) ([]event.Event, error) {
	eventCodesAsString := getEventCodesAsString(events)

	log.Info("Starting data search to update events: ", eventCodesAsString)

	chunkSize := 100
	eventCodeChunks := splitters.StringSlice(eventCodesAsString, chunkSize)

	var updatedEvents []event.Event
	for _, eventCodeChunk := range eventCodeChunks {
		eventsToUpdate, err := instance.chamberApi.GetEventsByCodes(eventCodeChunk)
		if err != nil {
			log.Error("chamberApi.GetEventsByCodes(): ", err.Error())
			return nil, nil
		}

		for _, eventData := range eventsToUpdate {
			eventDomain, err := instance.getEventDataToUpdate(eventData)
			if err != nil {
				log.Error("getEventDataToUpdate(): ", err.Error())
				return nil, err
			}

			updatedEvents = append(updatedEvents, *eventDomain)
		}
	}

	updatedEvents = removeEventsThatHaveNotBeenUpdated(updatedEvents, events)

	log.Info("Successful data search for event update: ", eventCodesAsString)
	return updatedEvents, nil
}

func getEventCodesAsString(events []event.Event) []string {
	var eventCodesAsString []string
	for _, eventData := range events {
		eventCodesAsString = append(eventCodesAsString, fmt.Sprint(eventData.Code()))
	}

	return eventCodesAsString
}

func removeEventsThatHaveNotBeenUpdated(updatedEvents []event.Event, registeredEvents []event.Event) []event.Event {
	var eventsWithUpdates []event.Event
	for _, updatedEvent := range updatedEvents {
		for _, registeredEvent := range registeredEvents {
			if registeredEvent.Code() == updatedEvent.Code() {
				if !registeredEvent.IsEqual(updatedEvent) {
					eventsWithUpdates = append(eventsWithUpdates, updatedEvent)
				}
				break
			}
		}
	}

	if len(eventsWithUpdates) > 0 {
		eventCodesAsString := getEventCodesAsString(eventsWithUpdates)
		log.Info("Events that need to be updated: ", eventCodesAsString)
	}

	return eventsWithUpdates
}

func (instance Event) getEventDataToUpdate(eventData map[string]interface{}) (*event.Event, error) {
	code, err := converters.ToInt(eventData["id"])
	if err != nil {
		log.Error("converters.ToInt(): ", err.Error())
		return nil, err
	}

	startsAt, err := time.Parse("2006-01-02T15:04", fmt.Sprint(eventData["dataHoraInicio"]))
	if err != nil {
		log.Errorf("Error converting date and time of start of event %d: %s", code, err.Error())
		return nil, err
	}

	var endsAt time.Time
	endsAtAsString := fmt.Sprint(eventData["dataHoraFim"])
	if endsAtAsString != "<nil>" {
		endsAt, err = time.Parse("2006-01-02T15:04", endsAtAsString)
		if err != nil {
			log.Errorf("Error converting date and time of end of event %d: %s", code, err.Error())
			return nil, err
		}
	}

	location, isInternal, err := getEventLocation(eventData)
	if err != nil {
		log.Error("getEventLocation(): ", err.Error())
		return nil, err
	}

	eventSituationDescription := fmt.Sprint(eventData["situacao"])
	eventSituation, specificSituation, err := instance.getEventSituationByDescription(eventSituationDescription)
	if err != nil {
		log.Error("getEventSituationByDescription(): ", err.Error())
		return nil, err
	}

	eventBuilder := event.NewBuilder().
		Code(code).
		Description(fmt.Sprint(eventData["descricao"])).
		StartsAt(startsAt).
		Location(location).
		IsInternal(isInternal).
		SpecificSituation(specificSituation).
		Situation(*eventSituation)

	if !endsAt.IsZero() {
		eventBuilder.EndsAt(endsAt)
	}

	videoUrl := fmt.Sprint(eventData["urlRegistro"])
	if videoUrl != "<nil>" {
		if validators.IsUrlValid(videoUrl) {
			eventBuilder.VideoUrl(videoUrl)
		} else {
			log.Infof("Event %d video URL is invalid", code)
		}
	}

	eventDomain, err := eventBuilder.Build()
	if err != nil {
		log.Errorf("Error validating data for event %d: %s", code, err.Error())
		return nil, err
	}

	return eventDomain, nil
}
