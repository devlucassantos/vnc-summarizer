package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/core/services/utils/datetime"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type Voting struct {
	connectionManager connectionManagerInterface
}

func NewVotingRepository(connectionManager connectionManagerInterface) *Voting {
	return &Voting{
		connectionManager: connectionManager,
	}
}

func (instance Voting) CreateVoting(voting voting.Voting) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Error starting transaction to register the voting %s: %s", voting.Code(), err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	referenceDateTime, err := datetime.GetCurrentDateTimeInBrazil()
	if err != nil {
		log.Error("datetime.GetCurrentDateTimeInBrazil(): ", err)
		return nil, err
	}

	var articleId uuid.UUID
	votingArticle := voting.Article()
	articleType := votingArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert(), articleType.Id(), referenceDateTime).Scan(&articleId)
	if err != nil {
		log.Errorf("Error registering voting %s as article: %s", voting.Id(), err.Error())
		return nil, err
	}

	legislativeBody := voting.LegislativeBody()
	mainProposition := voting.MainProposition()
	var mainPropositionId *uuid.UUID
	if mainProposition.Id() != uuid.Nil {
		propositionId := mainProposition.Id()
		mainPropositionId = &propositionId
	}

	var votingId uuid.UUID
	err = transaction.QueryRow(queries.Voting().Insert(), voting.Code(), voting.Description(), voting.Result(),
		voting.ResultAnnouncedAt(), voting.IsApproved(), legislativeBody.Id(), mainPropositionId, articleId).
		Scan(&votingId)
	if err != nil {
		log.Errorf("Error registering voting %s: %s", voting.Code(), err.Error())
		return nil, err
	}

	for _, relatedProposition := range voting.RelatedPropositions() {
		_, err = transaction.Exec(queries.PropositionRelatedToVoting().Insert(), relatedProposition.Id(), votingId)
		if err != nil {
			log.Errorf("Error registering proposition %s related to voting %s: %s", relatedProposition.Id(),
				votingId, err.Error())
			return nil, err
		}
		log.Infof("Proposition %s related to voting %s successfully registered", relatedProposition.Id(),
			votingId)
	}

	for _, affectedProposition := range voting.AffectedPropositions() {
		_, err = transaction.Exec(queries.PropositionAffectedByVoting().Insert(), affectedProposition.Id(), votingId)
		if err != nil {
			log.Errorf("Error registering proposition %s affected by voting %s: %s", affectedProposition.Id(),
				votingId, err.Error())
			return nil, err
		}
		log.Infof("Proposition %s affected by voting %s successfully registered", affectedProposition.Id(),
			votingId)
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Error confirming transaction to register voting %s: %s", voting.Code(), err.Error())
		return nil, err
	}

	log.Infof("Voting %s successfully registered with ID %s (Article ID: %s)", voting.Code(), votingId,
		articleId)
	return &votingId, nil
}

func (instance Voting) GetVotingByCode(code string) (*voting.Voting, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var votingData dto.Voting
	err = postgresConnection.Get(&votingData, queries.Voting().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Voting %s not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for voting %s from the database: %s", code, err.Error())
		return nil, err
	}

	votingDomain, err := voting.NewBuilder().
		Id(votingData.Id).
		Code(votingData.Code).
		Description(votingData.Description).
		Result(votingData.Result).
		ResultAnnouncedAt(votingData.ResultAnnouncedAt).
		IsApproved(votingData.IsApproved).
		Build()
	if err != nil {
		log.Errorf("Error validating data for voting %s: %s", votingData.Id, err.Error())
		return nil, err
	}

	return votingDomain, nil
}

func (instance Voting) GetVotesByCodes(codes []string) ([]voting.Voting, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var votingCodes []interface{}
	for _, code := range codes {
		votingCodes = append(votingCodes, code)
	}

	var votes []dto.Voting
	err = postgresConnection.Select(&votes, queries.Voting().Select().ByCodes(len(votingCodes)), votingCodes...)
	if err != nil {
		log.Error("Error retrieving the voting data by codes from the database: ", err.Error())
		return nil, err
	}

	var votingData []voting.Voting
	for _, votingDto := range votes {
		votingDomain, err := voting.NewBuilder().
			Id(votingDto.Id).
			Code(votingDto.Code).
			Description(votingDto.Description).
			Result(votingDto.Result).
			ResultAnnouncedAt(votingDto.ResultAnnouncedAt).
			IsApproved(votingDto.IsApproved).
			Build()
		if err != nil {
			log.Errorf("Error validating data for voting %s: %s", votingDto.Id, err.Error())
			return nil, err
		}
		votingData = append(votingData, *votingDomain)
	}

	return votingData, nil
}
