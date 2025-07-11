package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type Party struct {
	connectionManager connectionManagerInterface
}

func NewPartyRepository(connectionManager connectionManagerInterface) *Party {
	return &Party{
		connectionManager: connectionManager,
	}
}

func (instance Party) CreateParty(party party.Party) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var partyId uuid.UUID
	err = postgresConnection.QueryRow(queries.Party().Insert(), party.Code(), party.Name(), party.Acronym(),
		party.ImageUrl()).Scan(&partyId)
	if err != nil {
		log.Errorf("Error registering party %d: %s", party.Code(), err.Error())
		return nil, err
	}

	log.Infof("Party %d successfully registered with ID %s", party.Code(), partyId)
	return &partyId, nil
}

func (instance Party) UpdateParty(party party.Party) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	_, err = postgresConnection.Exec(queries.Party().Update(), party.Name(), party.Acronym(), party.ImageUrl(),
		party.Code())
	if err != nil {
		log.Errorf("Error updating party %d: %s", party.Code(), err.Error())
		return err
	}

	log.Infof("Party %d successfully updated", party.Code())
	return nil
}

func (instance Party) GetPartyByCode(code int) (*party.Party, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var partyData dto.Party
	err = postgresConnection.Get(&partyData, queries.Party().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Party %d not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for party %d from the database: %s", code, err.Error())
		return nil, err
	}

	partyDomain, err := party.NewBuilder().
		Id(partyData.Id).
		Code(partyData.Code).
		Name(partyData.Name).
		Acronym(partyData.Acronym).
		ImageUrl(partyData.ImageUrl).
		Build()
	if err != nil {
		log.Errorf("Error validating data for party %s: %s", partyData.Id, err.Error())
		return nil, err
	}

	return partyDomain, nil
}
