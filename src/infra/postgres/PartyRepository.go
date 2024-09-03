package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
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
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var partyId uuid.UUID
	err = postgresConnection.QueryRow(queries.Party().Insert(), party.Code(), party.Name(), party.Acronym(),
		party.ImageUrl()).Scan(&partyId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o partido %d: %s", party.Code(), err.Error())
		return nil, err
	}

	log.Infof("Partido %d registrada com sucesso com o ID %s", party.Code(), partyId)
	return &partyId, nil
}

func (instance Party) UpdateParty(party party.Party) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	_, err = postgresConnection.Exec(queries.Party().Update(), party.Name(), party.Acronym(), party.ImageUrl(),
		party.Code())
	if err != nil {
		log.Errorf("Erro ao atualizar o partido %d: %s", party.Code(), err.Error())
		return err
	}

	log.Infof("Dados do partido %d atualizados com sucesso", party.Code())
	return nil
}

func (instance Party) GetPartyByCode(code int) (*party.Party, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var partyData dto.Party
	err = postgresConnection.Get(&partyData, queries.Party().Select().ByCode(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("Partido %d não encontrado no banco de dados", code)
			return nil, nil
		}
		log.Errorf("Erro ao obter os dados do partido %s no banco de dados: %s", code, err.Error())
		return nil, err
	}

	partyDomain, err := party.NewBuilder().
		Id(partyData.Id).
		Code(partyData.Code).
		Name(partyData.Name).
		Acronym(partyData.Acronym).
		ImageUrl(partyData.ImageUrl).
		CreatedAt(partyData.CreatedAt).
		UpdatedAt(partyData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do partido %s: %s", partyData.Id, err.Error())
		return nil, err
	}

	return partyDomain, nil
}
