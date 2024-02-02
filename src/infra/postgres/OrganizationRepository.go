package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
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

	var code *int
	organizationCode := organization.Code()
	if organizationCode > 0 {
		code = &organizationCode
	}

	var organizationId uuid.UUID
	err = postgresConnection.QueryRow(queries.Organization().Insert(), code, organization.Name(),
		organization.Acronym(), organization.Nickname(), organization.Type()).Scan(&organizationId)
	if err != nil {
		log.Errorf("Erro ao cadastrar a organização %s: %s", organization.Name(), err.Error())
		return nil, err
	}

	log.Infof("Organização %s registrada com sucesso com o ID %s", organization.Name(), organizationId)
	return &organizationId, nil
}

func (instance Organization) UpdateOrganization(organization organization.Organization) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	if organization.Code() > 0 {
		_, err = postgresConnection.Exec(queries.Organization().UpdateByCode(), organization.Name(),
			organization.Acronym(), organization.Nickname(), organization.Type(), organization.Code())
	} else {
		_, err = postgresConnection.Exec(queries.Organization().UpdateByNameAndType(), organization.Acronym(),
			organization.Nickname(), organization.Name(), organization.Type())
	}
	if err != nil {
		log.Errorf("Erro ao atualizar a organização %s: %s", organization.Name(), err.Error())
		return err
	}

	log.Infof("Dados da organização %s atualizados com sucesso", organization.Name())
	return nil
}

func (instance Organization) GetOrganization(org organization.Organization) (*organization.Organization, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var organizationData dto.Organization
	if org.Code() > 0 {
		err = postgresConnection.Get(&organizationData, queries.Organization().Select().ByCode(), org.Code())
	} else {
		err = postgresConnection.Get(&organizationData, queries.Organization().Select().ByNameAndType(), org.Name(),
			org.Type())
	}
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("Organização %s não encontrada no banco de dados", org.Name)
			return nil, nil
		}
		log.Errorf("Erro ao obter os dados da organização %s no banco de dados: %s", org.Name, err.Error())
		return nil, err
	}

	organizationBuilder := organization.NewBuilder().
		Id(organizationData.Id)

	if organizationData.Code > 0 {
		organizationBuilder.Code(organizationData.Code)
	}

	organizationDomain, err := organizationBuilder.
		Name(organizationData.Name).
		Acronym(organizationData.Acronym).
		Nickname(organizationData.Nickname).
		Active(organizationData.Active).
		CreatedAt(organizationData.CreatedAt).
		UpdatedAt(organizationData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro durante a construção da estrutura de dados da organização %s: %s", organizationData.Id,
			err.Error())
		return nil, err
	}

	return organizationDomain, nil
}
