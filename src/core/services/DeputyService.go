package services

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"vnc-summarizer/core/interfaces/repositories"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/replacers"
	"vnc-summarizer/core/services/utils/requesters"
)

type Deputy struct {
	deputyRepository repositories.Deputy
	partyRepository  repositories.Party
}

func NewDeputyService(deputyRepository repositories.Deputy, partyRepository repositories.Party) *Deputy {
	return &Deputy{
		deputyRepository: deputyRepository,
		partyRepository:  partyRepository,
	}
}

func (instance Deputy) GetDeputyFromDeputyData(deputyData map[string]interface{}) (*deputy.Deputy, error) {
	deputyDataUrl := fmt.Sprint(deputyData["uri"])
	deputyDetails, err := requesters.GetDataObjectFromUrl(deputyDataUrl)
	if err != nil {
		log.Errorf("Error searching data for deputy %v: %s", deputyData["nome"], err.Error())
		return nil, err
	}

	deputyCode, err := converters.ToInt(deputyDetails["id"])
	if err != nil {
		log.Error("converters.ToInt(): ", err.Error())
		return nil, err
	}

	deputyLastStatus, err := converters.ToMap(deputyDetails["ultimoStatus"])
	if err != nil {
		log.Error("converters.ToMap(): ", err.Error())
		return nil, err
	}

	partyAcronym := fmt.Sprint(deputyLastStatus["siglaPartido"])
	partyAcronym = strings.ToUpper(strings.Trim(partyAcronym, "*"))
	partyAcronym = replacers.RemoveSpellingAccents(partyAcronym)

	urlOfExternalPartyData := fmt.Sprint("https://dadosabertos.camara.leg.br/api/v2/partidos?sigla=",
		partyAcronym)
	parties, err := requesters.GetDataSliceFromUrl(urlOfExternalPartyData)
	if err != nil {
		log.Error("requests.GetDataSliceFromUrl(): ", err.Error())
		return nil, err
	}

	partyUrl := fmt.Sprint(parties[0]["uri"])
	partyData, err := requesters.GetDataObjectFromUrl(partyUrl)
	if err != nil {
		log.Error("getDataObjectFromUrl(): ", err.Error())
		return nil, err
	}

	partyCode, err := converters.ToInt(partyData["id"])
	if err != nil {
		log.Error("converters.ToInt(): ", err.Error())
		return nil, err
	}

	partyDomain, err := party.NewBuilder().
		Code(partyCode).
		Name(fmt.Sprint(partyData["nome"])).
		Acronym(fmt.Sprint(partyData["sigla"])).
		ImageUrl(fmt.Sprint(partyData["urlLogo"])).
		Build()
	if err != nil {
		log.Errorf("Error validating data for party %d of deputy %d: %s", partyCode, deputyCode, err.Error())
		return nil, err
	}

	deputyDomain, err := deputy.NewBuilder().
		Code(deputyCode).
		Cpf(fmt.Sprint(deputyDetails["cpf"])).
		Name(cases.Title(language.BrazilianPortuguese).String(fmt.Sprint(deputyDetails["nomeCivil"]))).
		ElectoralName(cases.Title(language.BrazilianPortuguese).String(fmt.Sprint(deputyLastStatus["nomeEleitoral"]))).
		ImageUrl(fmt.Sprint(deputyLastStatus["urlFoto"])).
		Party(*partyDomain).
		FederatedUnit(fmt.Sprint(deputyLastStatus["siglaUf"])).
		Build()
	if err != nil {
		log.Errorf("Error validating data for deputy %d: %s", deputyCode, err.Error())
		return nil, err
	}

	updatedDeputy, err := instance.getDeputyFromDatabase(deputyDomain)
	if err != nil {
		log.Error("getDeputyFromDatabase(): ", err.Error())
		return nil, err
	}

	return updatedDeputy, nil
}

func (instance Deputy) getDeputyFromDatabase(deputyDomain *deputy.Deputy) (*deputy.Deputy, error) {
	deputyParty := deputyDomain.Party()
	registeredParty, err := instance.partyRepository.GetPartyByCode(deputyParty.Code())
	if err != nil {
		log.Error("partyRepository.GetPartyByCode(): ", err.Error())
		return nil, err
	}

	var partyId *uuid.UUID
	if registeredParty == nil {
		partyId, err = instance.partyRepository.CreateParty(deputyParty)
		if err != nil {
			log.Error("partyRepository.CreateParty(): ", err.Error())
			return nil, err
		}
	} else if !registeredParty.IsEqual(deputyParty) {
		err = instance.partyRepository.UpdateParty(deputyParty)
		if err != nil {
			log.Error("partyRepository.UpdateParty(): ", err.Error())
			return nil, err
		}
	}

	var updatedParty *party.Party
	if partyId == nil {
		updatedParty, err = deputyParty.NewUpdater().Id(registeredParty.Id()).Build()
		if err != nil {
			log.Errorf("Error updating party %s: %s", registeredParty.Id(), err.Error())
			return nil, err
		}
	} else {
		updatedParty, err = deputyParty.NewUpdater().Id(*partyId).Build()
		if err != nil {
			log.Errorf("Error updating party %s: %s", partyId, err.Error())
			return nil, err
		}
	}

	updatedDeputy, err := deputyDomain.NewUpdater().Party(*updatedParty).Build()
	if err != nil {
		log.Errorf("Error updating party %s of deputy %d: %s", partyId, deputyDomain.Code(), err.Error())
		return nil, err
	}

	registeredDeputy, err := instance.deputyRepository.GetDeputyByCode(updatedDeputy.Code())
	if err != nil {
		log.Error("deputyRepository.GetDeputyByCode(): ", err.Error())
		return nil, err
	}

	var deputyId *uuid.UUID
	if registeredDeputy == nil {
		deputyId, err = instance.deputyRepository.CreateDeputy(*updatedDeputy)
		if err != nil {
			log.Error("deputyRepository.CreateDeputy(): ", err.Error())
			return nil, err
		}
	} else if !registeredDeputy.IsEqual(*updatedDeputy) {
		err = instance.deputyRepository.UpdateDeputy(*updatedDeputy)
		if err != nil {
			log.Error("deputyRepository.UpdateDeputy(): ", err.Error())
			return nil, err
		}
	}

	if deputyId == nil {
		updatedDeputy, err = updatedDeputy.NewUpdater().Id(registeredDeputy.Id()).Build()
		if err != nil {
			log.Errorf("Error updating deputy %s: %s", registeredDeputy.Id(), err.Error())
			return nil, err
		}
	} else {
		updatedDeputy, err = updatedDeputy.NewUpdater().Id(*deputyId).Build()
		if err != nil {
			log.Errorf("Error updating deputy %s: %s", deputyId, err.Error())
			return nil, err
		}
	}

	return updatedDeputy, nil
}
