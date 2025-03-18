package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/core/services/utils/datetime"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type Event struct {
	connectionManager connectionManagerInterface
}

func NewEventRepository(connectionManager connectionManagerInterface) *Event {
	return &Event{
		connectionManager: connectionManager,
	}
}

func (instance Event) CreateEvent(event event.Event) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to register event %d: %s", event.Code(), err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	referenceDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	var articleId uuid.UUID
	eventArticle := event.Article()
	articleType := eventArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert(), articleType.Id(), referenceDateTime).Scan(&articleId)
	if err != nil {
		log.Errorf("Error registering event %s as article: %s", event.Id(), err.Error())
		return nil, err
	}

	var eventId uuid.UUID
	eventType := event.Type()
	eventSituation := event.Situation()
	err = transaction.QueryRow(queries.Event().Insert(), event.Code(), event.Title(), event.Description(),
		event.StartsAt(), event.EndsAt(), event.Location(), event.IsInternal(), event.VideoUrl(), event.SpecificType(),
		eventType.Id(), event.SpecificSituation(), eventSituation.Id(), articleId).Scan(&eventId)
	if err != nil {
		log.Errorf("Error registering event %d: %s", event.Code(), err.Error())
		return nil, err
	}

	for _, legislativeBody := range event.LegislativeBodies() {
		_, err = transaction.Exec(queries.EventLegislativeBody().Insert(), eventId, legislativeBody.Id())
		if err != nil {
			log.Errorf("Error registering legislative body %s responsible for event %s: %s", legislativeBody.Id(),
				eventId, err.Error())
			return nil, err
		}
		log.Infof("Legislative body %s responsible for event %s successfully registered", legislativeBody.Id(),
			eventId)
	}

	for _, requirement := range event.Requirements() {
		_, err = transaction.Exec(queries.EventRequirement().Insert(), eventId, requirement.Id())
		if err != nil {
			log.Errorf("Error registering requeriment %s as part of event %s: %s", requirement.Id(), eventId,
				err.Error())
			return nil, err
		}
		log.Infof("Requirement %s successfully registered as part of event %s", requirement.Id(), eventId)
	}

	for _, agendaItem := range event.AgendaItems() {
		regime := agendaItem.Regime()
		proposition := agendaItem.Proposition()
		rapporteur := agendaItem.Rapporteur()
		var rapporteurId, rapporteurPartyId *uuid.UUID
		var rapporteurFederatedUnit *string
		if rapporteur.Id() != uuid.Nil {
			deputyId := rapporteur.Id()
			rapporteurId = &deputyId
			rapporteurParty := rapporteur.Party()
			partyId := rapporteurParty.Id()
			rapporteurPartyId = &partyId
			federatedUnit := rapporteur.FederatedUnit()
			rapporteurFederatedUnit = &federatedUnit
		}
		relatedProposition := agendaItem.RelatedProposition()
		var relatedPropositionId *uuid.UUID
		if relatedProposition.Id() != uuid.Nil {
			propositionId := relatedProposition.Id()
			relatedPropositionId = &propositionId
		}
		voting := agendaItem.Voting()
		var votingId *uuid.UUID
		if voting.Id() != uuid.Nil {
			agendaItemVotingId := voting.Id()
			votingId = &agendaItemVotingId
		}
		var agendaItemSituation *string
		if agendaItem.Situation() != "" {
			situation := agendaItem.Situation()
			agendaItemSituation = &situation
		}

		var agendaItemId uuid.UUID
		err = transaction.QueryRow(queries.EventAgendaItem().Insert(), agendaItem.Title(), agendaItem.Topic(),
			agendaItemSituation, regime.Id(), rapporteurId, rapporteurPartyId, rapporteurFederatedUnit, proposition.Id(),
			relatedPropositionId, votingId, eventId).Scan(&agendaItemId)
		if err != nil {
			log.Errorf("Error registering agenda item as part of event %s: %s", eventId, err.Error())
			return nil, err
		}

		log.Infof("Agenda item %s successfully registered as part of event %s", agendaItemId, eventId)
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to register event %d: %s", event.Code(), err.Error())
		return nil, err
	}

	log.Infof("Event %d successfully registered with ID %s (Article ID: %s)", event.Code(), eventId, articleId)
	return &eventId, nil
}

func (instance Event) UpdateEvent(event event.Event) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	eventSituation := event.Situation()
	_, err = postgresConnection.Exec(queries.Event().Update(), event.Description(), event.StartsAt(), event.EndsAt(),
		event.Location(), event.IsInternal(), event.VideoUrl(), event.SpecificSituation(), eventSituation.Id(),
		event.Code())
	if err != nil {
		log.Errorf("Error updating event %d: %s", event.Code(), err.Error())
		return err
	}

	log.Infof("Event %d successfully updated", event.Code())
	return nil
}

func (instance Event) GetEventsByCodes(codes []int) ([]event.Event, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var eventCodes []interface{}
	for _, code := range codes {
		eventCodes = append(eventCodes, code)
	}

	var events []dto.Event
	err = postgresConnection.Select(&events, queries.Event().Select().ByCodes(len(eventCodes)), eventCodes...)
	if err != nil {
		log.Error("Error retrieving the event data by codes from the database: ", err.Error())
		return nil, err
	}

	var eventSlice []event.Event
	for _, eventData := range events {
		eventBuilder := event.NewBuilder()

		if !eventData.EndsAt.IsZero() {
			eventBuilder.EndsAt(eventData.EndsAt)
		}

		if eventData.VideoUrl != "" {
			eventBuilder.VideoUrl(eventData.VideoUrl)
		}

		eventDomain, err := eventBuilder.
			Id(eventData.Id).
			Code(eventData.Code).
			Title(eventData.Title).
			Description(eventData.Description).
			StartsAt(eventData.StartsAt).
			Location(eventData.Location).
			IsInternal(eventData.IsInternal).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event %s: %s", eventData.Id, err.Error())
			return nil, err
		}

		eventSlice = append(eventSlice, *eventDomain)
	}

	return eventSlice, nil
}

func (instance Event) GetEventsOccurringToday() ([]event.Event, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)
	var events []dto.Event
	err = postgresConnection.Select(&events, queries.Event().Select().OccurringToday())
	if err != nil {
		log.Error("Error retrieving events occurring today from the database: ", err.Error())
		return nil, err
	}

	var eventSlice []event.Event
	for _, eventData := range events {
		eventSituation, err := eventsituation.NewBuilder().
			Id(eventData.EventSituation.Id).
			Description(eventData.EventSituation.Description).
			Color(eventData.EventSituation.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event situation %s of event %s: %s",
				eventData.EventSituation.Id, eventData.Id, err.Error())
			return nil, err
		}

		eventBuilder := event.NewBuilder()

		if !eventData.EndsAt.IsZero() {
			eventBuilder.EndsAt(eventData.EndsAt)
		}

		if eventData.VideoUrl != "" {
			eventBuilder.VideoUrl(eventData.VideoUrl)
		}

		eventDomain, err := eventBuilder.
			Id(eventData.Id).
			Code(eventData.Code).
			Title(eventData.Title).
			Description(eventData.Description).
			StartsAt(eventData.StartsAt).
			Location(eventData.Location).
			IsInternal(eventData.IsInternal).
			SpecificType(eventData.SpecificSituation).
			Situation(*eventSituation).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event %s: %s", eventData.Id, err.Error())
			return nil, err
		}

		eventSlice = append(eventSlice, *eventDomain)
	}

	return eventSlice, nil
}

func (instance Event) GetEventsThatStartedInTheLastThreeMonthsAndHaveNotFinished() ([]event.Event, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)
	var events []dto.Event
	err = postgresConnection.Select(&events, queries.Event().Select().StartedInTheLastThreeMonthsAndHaveNotFinished())
	if err != nil {
		log.Error("Error retrieving events that started in the last three months and have not finished from the "+
			"database: ", err.Error())
		return nil, err
	}

	var eventSlice []event.Event
	for _, eventData := range events {
		eventSituation, err := eventsituation.NewBuilder().
			Id(eventData.EventSituation.Id).
			Description(eventData.EventSituation.Description).
			Color(eventData.EventSituation.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event situation %s of event %s: %s",
				eventData.EventSituation.Id, eventData.Id, err.Error())
			return nil, err
		}

		eventBuilder := event.NewBuilder()

		if !eventData.EndsAt.IsZero() {
			eventBuilder.EndsAt(eventData.EndsAt)
		}

		if eventData.VideoUrl != "" {
			eventBuilder.VideoUrl(eventData.VideoUrl)
		}

		eventDomain, err := eventBuilder.
			Id(eventData.Id).
			Code(eventData.Code).
			Title(eventData.Title).
			Description(eventData.Description).
			StartsAt(eventData.StartsAt).
			Location(eventData.Location).
			IsInternal(eventData.IsInternal).
			SpecificType(eventData.SpecificSituation).
			Situation(*eventSituation).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event %s: %s", eventData.Id, err.Error())
			return nil, err
		}

		eventSlice = append(eventSlice, *eventDomain)
	}

	return eventSlice, nil
}
