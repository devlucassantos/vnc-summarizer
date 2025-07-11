package services

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-summarizer/core/interfaces/chamber"
	"vnc-summarizer/core/interfaces/postgres"
	"vnc-summarizer/utils/converters"
)

type LegislativeBody struct {
	chamberApi                    chamber.Chamber
	legislativeBodyRepository     postgres.LegislativeBody
	legislativeBodyTypeRepository postgres.LegislativeBodyType
}

func NewLegislativeBodyService(chamberApi chamber.Chamber, legislativeBodyRepository postgres.LegislativeBody,
	legislativeBodyTypeRepository postgres.LegislativeBodyType) *LegislativeBody {
	return &LegislativeBody{
		chamberApi:                    chamberApi,
		legislativeBodyRepository:     legislativeBodyRepository,
		legislativeBodyTypeRepository: legislativeBodyTypeRepository,
	}
}

func (instance LegislativeBody) RegisterNewLegislativeBodyByCode(code int) (*uuid.UUID, error) {
	legislativeBodyData, err := instance.chamberApi.GetLegislativeBodyByCode(code)
	if err != nil {
		log.Error("chamberApi.GetLegislativeBodyByCode(): ", err.Error())
		return nil, err
	}

	typeCode, err := converters.ToInt(legislativeBodyData["codTipoOrgao"])
	if err != nil {
		log.Error("converters.ToInt(): ", err.Error())
		return nil, err
	}

	legislativeBodyType, err := instance.getLegislativeBodyTypeDataByCode(typeCode)
	if err != nil {
		log.Error("getLegislativeBodyTypeDataByCode(): ", err.Error())
		return nil, err
	}

	legislativeBody, err := legislativebody.NewBuilder().
		Code(code).
		Name(fmt.Sprint(legislativeBodyData["nome"])).
		Acronym(fmt.Sprint(legislativeBodyData["sigla"])).
		Type(*legislativeBodyType).
		Build()
	if err != nil {
		log.Errorf("Error validating data for legislative body %d: %s", code, err.Error())
		return nil, err
	}

	legislativeBodyId, err := instance.legislativeBodyRepository.CreateLegislativeBody(*legislativeBody)
	if err != nil {
		log.Error("legislativeBodyRepository.CreateLegislativeBody(): ", err.Error())
		return nil, err
	}

	return legislativeBodyId, nil
}

func (instance LegislativeBody) getLegislativeBodyTypeDataByCode(code int) (*legislativebodytype.LegislativeBodyType,
	error) {
	legislativeBodyType, err := instance.legislativeBodyTypeRepository.GetLegislativeBodyTypeByCode(code)
	if err != nil {
		log.Error("legislativeBodyTypeRepository.GetLegislativeBodyTypeByCode(): ", err.Error())
		return nil, err
	}

	if legislativeBodyType == nil {
		legislativeBodyTypes, err := instance.chamberApi.GetLegislativeBodyTypes()
		if err != nil {
			log.Error("chamberApi.GetLegislativeBodyTypes(): ", err.Error())
			return nil, err
		}

		var legislativeBodyTypeData map[string]interface{}
		for _, legislativeBodyTypeMap := range legislativeBodyTypes {
			if fmt.Sprint(legislativeBodyTypeMap["cod"]) == fmt.Sprint(code) {
				legislativeBodyTypeData = legislativeBodyTypeMap
				break
			}
		}

		description := fmt.Sprint(legislativeBodyTypeData["nome"])

		legislativeBodyType, err = legislativebodytype.NewBuilder().
			Code(code).
			Description(description).
			Build()
		if err != nil {
			log.Errorf("Error validating data for legislative body type %d: %s", code, err.Error())
			return nil, err
		}

		legislativeBodyTypeId, err := instance.legislativeBodyTypeRepository.CreateLegislativeBodyType(
			*legislativeBodyType)
		if err != nil {
			log.Error("legislativeBodyTypeRepository.CreateLegislativeBodyType(): ", err.Error())
			return nil, err
		}

		legislativeBodyType, err = legislativeBodyType.NewUpdater().Id(*legislativeBodyTypeId).Build()
		if err != nil {
			log.Errorf("Error updating legislative body type %s: %s", legislativeBodyTypeId, err.Error())
			return nil, err
		}
	}

	return legislativeBodyType, nil
}

func (instance LegislativeBody) GetLegislativeBodyByCode(code int) (*legislativebody.LegislativeBody, error) {
	legislativeBody, err := instance.legislativeBodyRepository.GetLegislativeBodyByCode(code)
	if err != nil {
		log.Error("legislativeBodyRepository.GetLegislativeBodyByCode(): ", err.Error())
		return nil, err
	}

	return legislativeBody, nil
}

func (instance LegislativeBody) GetLegislativeBodiesByCodes(codes []int) ([]legislativebody.LegislativeBody, error) {
	legislativeBodies, err := instance.legislativeBodyRepository.GetLegislativeBodiesByCodes(codes)
	if err != nil {
		log.Error("legislativeBodyRepository.GetLegislativeBodiesByCodes(): ", err.Error())
		return nil, err
	}

	return legislativeBodies, nil
}

func getCodesOfTheNewLegislativeBodies(returnedLegislativeBodyCodes []int,
	registeredLegislativeBodies []legislativebody.LegislativeBody) []int {
	var codesOfTheLegislativeBodiesToRegister []int
	for _, legislativeBodyCode := range returnedLegislativeBodyCodes {
		var legislativeBodyAlreadyRegistered bool
		for _, legislativeBodyData := range registeredLegislativeBodies {
			if legislativeBodyData.Code() == legislativeBodyCode {
				legislativeBodyAlreadyRegistered = true
				break
			}
		}
		if !legislativeBodyAlreadyRegistered {
			codesOfTheLegislativeBodiesToRegister = append(codesOfTheLegislativeBodiesToRegister, legislativeBodyCode)
		}
	}

	codesOfTheLegislativeBodiesToRegister = converters.IntSliceToUniqueIntSlice(codesOfTheLegislativeBodiesToRegister)

	return codesOfTheLegislativeBodiesToRegister
}
