package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/adapters/databases/dto"
	"vnc-summarizer/adapters/databases/postgres/queries"
)

type LegislativeBody struct {
	connectionManager connectionManagerInterface
}

func NewLegislativeBodyRepository(connectionManager connectionManagerInterface) *LegislativeBody {
	return &LegislativeBody{
		connectionManager: connectionManager,
	}
}

func (instance LegislativeBody) CreateLegislativeBody(legislativeBody legislativebody.LegislativeBody) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var legislativeBodyId uuid.UUID
	legislativeBodyType := legislativeBody.Type()
	err = postgresConnection.QueryRow(queries.LegislativeBody().Insert(), legislativeBody.Code(), legislativeBody.Name(),
		legislativeBody.Acronym(), legislativeBodyType.Id()).Scan(&legislativeBodyId)
	if err != nil {
		log.Errorf("Error registering legislative body %d: %s", legislativeBody.Code(), err.Error())
		return nil, err
	}

	log.Infof("Legislative body %d successfully registered with ID %s", legislativeBody.Code(), legislativeBodyId)
	return &legislativeBodyId, nil
}

func (instance LegislativeBody) GetLegislativeBodyByCode(code int) (*legislativebody.LegislativeBody, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var legislativeBody dto.LegislativeBody
	err = postgresConnection.Get(&legislativeBody, queries.LegislativeBody().Select().ByCode(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("Legislative body %d not found in database", code)
			return nil, nil
		}
		log.Errorf("Error retrieving data for external author %d from the database: %s", code, err.Error())
		return nil, err
	}

	legislativeBodyType, err := legislativebodytype.NewBuilder().
		Id(legislativeBody.LegislativeBodyType.Id).
		Code(legislativeBody.LegislativeBodyType.Code).
		Description(legislativeBody.LegislativeBodyType.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for legislative body type %s: %s",
			legislativeBody.LegislativeBodyType.Id, err.Error())
		return nil, err
	}

	legislativeBodyDomain, err := legislativebody.NewBuilder().
		Id(legislativeBody.Id).
		Name(legislativeBody.Name).
		Type(*legislativeBodyType).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author %d: %s", legislativeBody.Id, err.Error())
		return nil, err
	}

	return legislativeBodyDomain, nil
}

func (instance LegislativeBody) GetLegislativeBodiesByCodes(codes []int) ([]legislativebody.LegislativeBody, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("connectionManager.createConnection(): ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var legislativeBodyCodes []interface{}
	for _, code := range codes {
		legislativeBodyCodes = append(legislativeBodyCodes, code)
	}

	var legislativeBodyData []dto.LegislativeBody
	err = postgresConnection.Select(&legislativeBodyData, queries.LegislativeBody().Select().ByCodes(
		len(legislativeBodyCodes)), legislativeBodyCodes...)
	if err != nil {
		log.Error("Error retrieving the legislative body data by codes from the database: ", err.Error())
		return nil, err
	}

	var legislativeBodies []legislativebody.LegislativeBody
	for _, legislativeBody := range legislativeBodyData {
		legislativeBodyType, err := legislativebodytype.NewBuilder().
			Id(legislativeBody.LegislativeBodyType.Id).
			Code(legislativeBody.LegislativeBodyType.Code).
			Description(legislativeBody.LegislativeBodyType.Description).
			Build()
		if err != nil {
			log.Errorf("Error validating data for legislative body type %s: %s",
				legislativeBody.LegislativeBodyType.Id, err.Error())
			return nil, err
		}

		legislativeBodyDomain, err := legislativebody.NewBuilder().
			Id(legislativeBody.Id).
			Code(legislativeBody.Code).
			Name(legislativeBody.Name).
			Acronym(legislativeBody.Acronym).
			Type(*legislativeBodyType).
			Build()
		if err != nil {
			log.Errorf("Error validating data for external author %d: %s", legislativeBody.Id, err.Error())
			return nil, err
		}
		legislativeBodies = append(legislativeBodies, *legislativeBodyDomain)
	}

	return legislativeBodies, nil
}
