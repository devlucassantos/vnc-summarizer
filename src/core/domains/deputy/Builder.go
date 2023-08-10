package deputy

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
	"vnc-write-api/core/domains/party"
)

type builder struct {
	deputy        *Deputy
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{deputy: &Deputy{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID do(a) deputado(a) é inválido")
		return instance
	}
	instance.deputy.id = id
	return instance
}

func (instance *builder) Code(code int) *builder {
	if code <= 0 {
		instance.invalidFields = append(instance.invalidFields, "O código do(a) deputado(a) é inválido")
		return instance
	}
	instance.deputy.code = code
	return instance
}

func (instance *builder) Cpf(cpf string) *builder {
	instance.deputy.cpf = cpf
	return instance
}

func (instance *builder) Name(name string) *builder {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O nome do(a) deputado(a) é inválido")
		return instance
	}
	instance.deputy.name = name
	return instance
}

func (instance *builder) ElectoralName(electoralName string) *builder {
	instance.deputy.electoralName = electoralName
	return instance
}

func (instance *builder) ImageUrl(imageUrl string) *builder {
	imageUrl = strings.Trim(imageUrl, "/")
	if len(imageUrl) == 0 {
		instance.invalidFields = append(instance.invalidFields, "A URL da imagem do(a) deputado(a) é inválida")
		return instance
	}
	instance.deputy.imageUrl = imageUrl
	return instance
}

func (instance *builder) CurrentParty(party party.Party) *builder {
	instance.deputy.currentParty = party
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.deputy.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro do(a) deputado(a) é inválida")
		return instance
	}
	instance.deputy.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro do(a) deputado(a) é inválida")
		return instance
	}
	instance.deputy.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*Deputy, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}

	return instance.deputy, nil
}
