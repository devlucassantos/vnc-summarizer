package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-write-api/infra/dto"
	"vnc-write-api/infra/postgres/queries"
)

type Deputy struct {
	connectionManager ConnectionManagerInterface
}

func NewDeputyRepository(connectionManager ConnectionManagerInterface) *Deputy {
	return &Deputy{
		connectionManager: connectionManager,
	}
}

func (instance Deputy) CreateDeputy(deputy deputy.Deputy) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var deputyId uuid.UUID
	deputyParty := deputy.Party()
	err = postgresConnection.QueryRow(queries.Deputy().Insert(), deputy.Code(), deputy.Cpf(), deputy.Name(),
		deputy.ElectoralName(), deputy.ImageUrl(), deputyParty.Id()).Scan(&deputyId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o(a) deputado(a) %d: %s", deputy.Code(), err.Error())
		return nil, err
	}

	log.Infof("Deputado(a) %d registrado(a) com sucesso com o ID %s", deputy.Code(), deputyId)
	return &deputyId, nil
}

func (instance Deputy) UpdateDeputy(deputy deputy.Deputy) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	deputyParty := deputy.Party()
	_, err = postgresConnection.Exec(queries.Deputy().Update(), deputy.Name(), deputy.ElectoralName(),
		deputy.ImageUrl(), deputyParty.Id(), deputy.Code())
	if err != nil {
		log.Errorf("Erro ao atualizar o(a) deputado(a) %d: %s", deputy.Code(), err.Error())
		return err
	}

	log.Infof("Dados do(a) deputado(a) %d atualizados com sucesso", deputy.Code())
	return nil
}

func (instance Deputy) GetDeputyByCode(code int) (*deputy.Deputy, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var deputyData dto.Deputy
	err = postgresConnection.Get(&deputyData, queries.Deputy().Select().ByCode(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("Deputado(a) %d não encontrado(a) no banco de dados", code)
			return nil, nil
		}
		log.Errorf("Erro ao obter os dados do(a) deputado(a) %d no banco de dados: %s", code, err.Error())
		return nil, err
	}

	deputyParty, err := party.NewBuilder().
		Id(deputyData.Party.Id).
		Code(deputyData.Party.Code).
		Name(deputyData.Party.Name).
		Acronym(deputyData.Party.Acronym).
		ImageUrl(deputyData.Party.ImageUrl).
		Active(deputyData.Party.Active).
		CreatedAt(deputyData.Party.CreatedAt).
		UpdatedAt(deputyData.Party.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados do partido %s do(a) deputado(a) %s: %s",
			deputyData.Party.Id, deputyData.Id, err.Error())
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
		Active(deputyData.Active).
		CreatedAt(deputyData.CreatedAt).
		UpdatedAt(deputyData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados do(a) deputado(a) %s: %s", deputyData.Id, err.Error())
		return nil, err
	}

	return deputyDomain, nil
}
