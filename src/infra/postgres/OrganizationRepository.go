package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-write-api/core/domains/organization"
	"vnc-write-api/infra/dto"
	"vnc-write-api/infra/postgres/queries"
)

type Organization struct {
	connectionManager ConnectionManagerInterface
}

func NewOrganizationRepository(connectionManager ConnectionManagerInterface) *Organization {
	return &Organization{
		connectionManager: connectionManager,
	}
}

func (instance Organization) CreateOrganization(organization organization.Organization) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var organizationId uuid.UUID
	err = postgresConnection.QueryRow(queries.Organization().Insert(), organization.Code(), organization.Name(),
		organization.Acronym(), organization.Nickname()).Scan(&organizationId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a organização %d: %s", organization.Code(), err.Error())
		return nil, err
	}

	log.Infof("Organização %d registrada com sucesso com o ID %s", organization.Code(), organizationId)
	return &organizationId, nil
}

func (instance Organization) UpdateOrganization(organization organization.Organization) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	_, err = postgresConnection.Exec(queries.Organization().Update(), organization.Name(), organization.Acronym(),
		organization.Nickname(), organization.Code())
	if err != nil {
		log.Errorf("Erro ao atualizar a organização %d: %s", organization.Code(), err.Error())
		return err
	}

	log.Infof("Dados da organização %d atualizados com sucesso", organization.Code())
	return nil
}

func (instance Organization) GetOrganizationByCode(code int) (*organization.Organization, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var organizationData dto.Organization
	err = postgresConnection.Get(&organizationData, queries.Organization().Select().ByCode(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("Organização %d não encontrada no banco de dados", code)
			return nil, nil
		}
		log.Errorf("Erro ao obter os dados da organização %d no banco de dados: %s", code, err.Error())
		return nil, err
	}

	organizationDomain, err := organization.NewBuilder().
		Id(organizationData.Id).
		Code(organizationData.Code).
		Name(organizationData.Name).
		Acronym(organizationData.Acronym).
		Nickname(organizationData.Nickname).
		Active(organizationData.Active).
		CreatedAt(organizationData.CreatedAt).
		UpdatedAt(organizationData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados da organização %s: %s", organizationData.Id, err.Error())
		return nil, err
	}

	return organizationDomain, nil
}
