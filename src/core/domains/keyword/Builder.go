package keyword

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type builder struct {
	keyword       *Keyword
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{keyword: &Keyword{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID da palavra-chave é inválido")
		return instance
	}
	instance.keyword.id = id
	return instance
}

func (instance *builder) Keyword(keyword string) *builder {
	instance.keyword.keyword = keyword
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.keyword.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro da palavra-chave é inválida")
		return instance
	}
	instance.keyword.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro da palavra-chave é inválida")
		return instance
	}
	instance.keyword.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*Keyword, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}

	return instance.keyword, nil
}
