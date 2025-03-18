package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/labstack/gommon/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"vnc-summarizer/core/interfaces/repositories"
)

type ExternalAuthor struct {
	externalAuthorRepository     repositories.ExternalAuthor
	externalAuthorTypeRepository repositories.ExternalAuthorType
}

func NewExternalAuthorService(externalAuthorRepository repositories.ExternalAuthor,
	externalAuthorTypeRepository repositories.ExternalAuthorType) *ExternalAuthor {
	return &ExternalAuthor{
		externalAuthorRepository:     externalAuthorRepository,
		externalAuthorTypeRepository: externalAuthorTypeRepository,
	}
}

func (instance ExternalAuthor) GetExternalAuthorFromAuthorData(authorName string, authorTypeCode int,
	authorType string) (*externalauthor.ExternalAuthor, error) {
	caser := cases.Title(language.BrazilianPortuguese)
	externalAuthorTypeData, err := externalauthortype.NewBuilder().
		Code(authorTypeCode).
		Description(caser.String(strings.ToLower(authorType))).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author type %d: %s", authorTypeCode, err.Error())
		return nil, err
	}

	externalAuthorData, err := externalauthor.NewBuilder().
		Name(caser.String(strings.ToLower(authorName))).
		Type(*externalAuthorTypeData).
		Build()
	if err != nil {
		log.Errorf("Error validating data for external author %s - %s: %s", authorName, authorType, err.Error())
		return nil, err
	}

	updatedExternalAuthor, err := instance.getExternalAuthorFromDatabase(externalAuthorData)

	return updatedExternalAuthor, nil
}

func (instance ExternalAuthor) getExternalAuthorFromDatabase(
	externalAuthorDomain *externalauthor.ExternalAuthor) (*externalauthor.ExternalAuthor, error) {
	externalAuthorType := externalAuthorDomain.Type()
	registeredExternalAuthor, err := instance.externalAuthorRepository.GetExternalAuthorByNameAndTypeCode(
		externalAuthorDomain.Name(), externalAuthorType.Code())
	if err != nil {
		log.Error("externalAuthorRepository.GetExternalAuthorByNameAndTypeCode(): ", err.Error())
		return nil, err
	} else if registeredExternalAuthor != nil {
		return registeredExternalAuthor, nil
	}

	registeredExternalAuthorType, err := instance.externalAuthorTypeRepository.GetExternalAuthorTypeByCode(
		externalAuthorType.Code())
	if err != nil {
		log.Error("externalAuthorTypeRepository.GetExternalAuthorTypeByCode(): ", err.Error())
		return nil, err
	}

	var updatedExternalAuthor *externalauthor.ExternalAuthor
	if registeredExternalAuthorType != nil {
		updatedExternalAuthor, err = externalAuthorDomain.NewUpdater().Type(*registeredExternalAuthorType).Build()
	} else {
		externalAuthorTypeId, err := instance.externalAuthorTypeRepository.CreateExternalAuthorType(
			externalAuthorType)
		if err != nil {
			log.Error("externalAuthorTypeRepository.CreateExternalAuthorType(): ", err.Error())
			return nil, err
		}

		updatedExternalAuthorType, err := externalAuthorType.NewUpdater().Id(*externalAuthorTypeId).Build()
		if err != nil {
			log.Errorf("Error updating external author type %d: %s", externalAuthorType.Code(),
				err.Error())
			return nil, err
		}

		updatedExternalAuthor, err = externalAuthorDomain.NewUpdater().Type(*updatedExternalAuthorType).Build()
	}
	if err != nil {
		log.Errorf("Error updating external author %s: %s", externalAuthorDomain.Name(), err.Error())
		return nil, err
	}

	externalAuthorId, err := instance.externalAuthorRepository.CreateExternalAuthor(*externalAuthorDomain)
	if err != nil {
		log.Error("externalAuthorRepository.CreateExternalAuthor(): ", err.Error())
		return nil, err
	}

	updatedExternalAuthor, err = externalAuthorDomain.NewUpdater().Id(*externalAuthorId).Build()
	if err != nil {
		log.Errorf("Error updating external author %s: %s", externalAuthorDomain.Name(), err.Error())
		return nil, err
	}

	return updatedExternalAuthor, nil
}
