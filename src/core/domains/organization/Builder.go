package organization

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type builder struct {
	organization  *Organization
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{organization: &Organization{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID da organização é inválido")
		return instance
	}
	instance.organization.id = id
	return instance
}

func (instance *builder) Code(code int) *builder {
	if code <= 0 {
		instance.invalidFields = append(instance.invalidFields, "O código da organização é inválido")
		return instance
	}
	instance.organization.code = code
	return instance
}

func (instance *builder) Name(name string) *builder {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O nome da organização é inválido")
		return instance
	}
	instance.organization.name = name
	return instance
}

func (instance *builder) Nickname(nickname string) *builder {
	instance.organization.nickname = nickname
	return instance
}

func (instance *builder) Acronym(acronym string) *builder {
	instance.organization.acronym = acronym
	return instance
}

func (instance *builder) Type(_type string) *builder {
	_type = strings.TrimSpace(_type)
	if len(_type) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O tipo da organização é inválido")
		return instance
	}
	instance.organization._type = _type
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.organization.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro da organização é inválida")
		return instance
	}
	instance.organization.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro da organização é inválida")
		return instance
	}
	instance.organization.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*Organization, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}
	return instance.organization, nil
}
