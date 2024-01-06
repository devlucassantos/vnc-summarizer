package newsletter

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
	"vnc-write-api/core/domains/proposition"
)

type builder struct {
	newsletter    *Newsletter
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{newsletter: &Newsletter{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID do boletim é inválido")
		return instance
	}
	instance.newsletter.id = id
	return instance
}

func (instance *builder) Title(title string) *builder {
	if len(title) < 10 {
		instance.newsletter.title = fmt.Sprintf("Boletim Diário")
		return instance
	}
	instance.newsletter.title = title
	return instance
}

func (instance *builder) Content(content string) *builder {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O conteúdo do boletim é inválido")
		return instance
	}
	instance.newsletter.content = content
	return instance
}

func (instance *builder) ReferenceDate(referenceDate time.Time) *builder {
	if referenceDate.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de referência do boletim é inválida")
		return instance
	}
	instance.newsletter.referenceDate = referenceDate
	return instance
}

func (instance *builder) Propositions(propositions []proposition.Proposition) *builder {
	instance.newsletter.propositions = propositions
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.newsletter.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro do boletim é inválida")
		return instance
	}
	instance.newsletter.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro do boletim é inválida")
		return instance
	}
	instance.newsletter.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*Newsletter, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}
	return instance.newsletter, nil
}
