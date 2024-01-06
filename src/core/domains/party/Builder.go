package party

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type builder struct {
	party         *Party
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{party: &Party{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID do partido é inválido")
		return instance
	}
	instance.party.id = id
	return instance
}

func (instance *builder) Code(code int) *builder {
	if code <= 0 {
		instance.invalidFields = append(instance.invalidFields, "O código do partido é inválido")
		return instance
	}
	instance.party.code = code
	return instance
}

func (instance *builder) Name(name string) *builder {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O nome do partido é inválido")
		return instance
	}
	instance.party.name = name
	return instance
}

func (instance *builder) Acronym(acronym string) *builder {
	instance.party.acronym = acronym
	return instance
}

func (instance *builder) ImageUrl(imageUrl string) *builder {
	imageUrl = strings.Trim(imageUrl, "/")
	if len(imageUrl) == 0 {
		instance.invalidFields = append(instance.invalidFields, "A URL da imagem do partido é inválida")
		return instance
	}
	instance.party.imageUrl = imageUrl
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.party.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro do partido é inválida")
		return instance
	}
	instance.party.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro do partido é inválida")
		return instance
	}
	instance.party.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*Party, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}
	return instance.party, nil
}
