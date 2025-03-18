package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type Deputy struct {
	connectionManager connectionManagerInterface
}

func NewDeputyRepository(connectionManager connectionManagerInterface) *Deputy {
	return &Deputy{
		connectionManager: connectionManager,
	}
}

func (instance Deputy) CreateDeputy(deputy deputy.Deputy) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var deputyId uuid.UUID
	deputyParty := deputy.Party()
	err = postgresConnection.QueryRow(queries.Deputy().Insert(), deputy.Code(), deputy.Cpf(), deputy.Name(),
		deputy.ElectoralName(), deputy.ImageUrl(), deputyParty.Id(), deputy.FederatedUnit()).Scan(&deputyId)
	if err != nil {
		log.Errorf("Error registering deputy %d: %s", deputy.Code(), err.Error())
		return nil, err
	}

	log.Infof("Deputy %d successfully registered with ID %s", deputy.Code(), deputyId)
	return &deputyId, nil
}

func (instance Deputy) UpdateDeputy(deputy deputy.Deputy) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	deputyParty := deputy.Party()
	_, err = postgresConnection.Exec(queries.Deputy().Update(), deputy.Name(), deputy.ElectoralName(),
		deputy.ImageUrl(), deputyParty.Id(), deputy.FederatedUnit(), deputy.Code())
	if err != nil {
		log.Errorf("Error updating deputy %d: %s", deputy.Code(), err.Error())
		return err
	}

	log.Infof("Deputy %d successfully updated", deputy.Code())
	return nil
}

func (instance Deputy) GetDeputyByCode(code int) (*deputy.Deputy, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var deputyData dto.Deputy
	err = postgresConnection.Get(&deputyData, queries.Deputy().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Deputy %d not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for deputy %d from the database: %s", code, err.Error())
		return nil, err
	}

	deputyParty, err := party.NewBuilder().
		Id(deputyData.Party.Id).
		Code(deputyData.Party.Code).
		Name(deputyData.Party.Name).
		Acronym(deputyData.Party.Acronym).
		ImageUrl(deputyData.Party.ImageUrl).
		Build()
	if err != nil {
		log.Errorf("Error validating data for party %s of deputy %s: %s", deputyData.Party.Id, deputyData.Id,
			err.Error())
		return nil, err
	}

	deputyDomain, err := deputy.NewBuilder().
		Id(deputyData.Id).
		Code(deputyData.Code).
		Cpf(deputyData.Cpf).
		Name(deputyData.Name).
		ElectoralName(deputyData.ElectoralName).
		ImageUrl(deputyData.ImageUrl).
		Party(*deputyParty).
		FederatedUnit(deputyData.FederatedUnit).
		Build()
	if err != nil {
		log.Errorf("Error validating data for deputy %s: %s", deputyData.Id, err.Error())
		return nil, err
	}

	return deputyDomain, nil
}
