package proposition

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
	"vnc-write-api/core/domains/deputy"
	"vnc-write-api/core/domains/keyword"
	"vnc-write-api/core/domains/organization"
)

type builder struct {
	proposition   *Proposition
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{proposition: &Proposition{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID da proposição é inválido")
		return instance
	}
	instance.proposition.id = id
	return instance
}

func (instance *builder) Code(code int) *builder {
	if code <= 0 {
		instance.invalidFields = append(instance.invalidFields, "O código da proposição é inválido")
		return instance
	}
	instance.proposition.code = code
	return instance
}

func (instance *builder) Title(title string) *builder {
	if len(title) < 10 {
		instance.proposition.title = fmt.Sprintf("Nova proposição de %s", time.Now())
		return instance
	}
	instance.proposition.title = title
	return instance
}

func (instance *builder) OriginalTextUrl(originalTextUrl string) *builder {
	originalTextUrl = strings.Trim(originalTextUrl, "/")
	if len(originalTextUrl) == 0 {
		instance.invalidFields = append(instance.invalidFields, "A URL do texto original da proposição é inválida")
		return instance
	}
	instance.proposition.originalTextUrl = originalTextUrl
	return instance
}

func (instance *builder) Summary(summary string) *builder {
	summary = strings.TrimSpace(summary)
	if len(summary) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O resumo da proposição é inválido")
		return instance
	}
	instance.proposition.summary = summary
	return instance
}

func (instance *builder) SubmittedAt(submittedAt time.Time) *builder {
	if submittedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de apresentação da proposição é inválida")
		return instance
	}
	instance.proposition.submittedAt = submittedAt
	return instance
}

func (instance *builder) Deputies(deputies []deputy.Deputy) *builder {
	instance.proposition.deputies = deputies
	return instance
}

func (instance *builder) Organizations(organizations []organization.Organization) *builder {
	instance.proposition.organizations = organizations
	return instance
}

func (instance *builder) Keywords(keywords []keyword.Keyword) *builder {
	instance.proposition.keywords = keywords
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.proposition.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro da proposição é inválida")
		return instance
	}
	instance.proposition.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro da proposição é inválida")
		return instance
	}
	instance.proposition.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*Proposition, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}

	return instance.proposition, nil
}
