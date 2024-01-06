package news

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type builder struct {
	news          *News
	invalidFields []string
}

func NewBuilder() *builder {
	return &builder{news: &News{}}
}

func (instance *builder) Id(id uuid.UUID) *builder {
	if id.ID() == 0 {
		instance.invalidFields = append(instance.invalidFields, "O ID do boletim é inválido")
		return instance
	}
	instance.news.id = id
	return instance
}

func (instance *builder) Title(title string) *builder {
	if len(title) < 10 {
		instance.invalidFields = append(instance.invalidFields, "O título da matéria é inválido")
		return instance
	}
	instance.news.title = title
	return instance
}

func (instance *builder) Content(content string) *builder {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O conteúdo da matéria é inválido")
		return instance
	}
	instance.news.content = content
	return instance
}

func (instance *builder) Views(views int) *builder {
	if views < 0 {
		instance.invalidFields = append(instance.invalidFields, "O número de visualizações da matéria é inválido")
		return instance
	}
	instance.news.views = views
	return instance
}

func (instance *builder) Type(_type string) *builder {
	_type = strings.TrimSpace(_type)
	if len(_type) == 0 {
		instance.invalidFields = append(instance.invalidFields, "O tipo da matéria é inválido")
		return instance
	}
	instance.news._type = _type
	return instance
}

func (instance *builder) Active(active bool) *builder {
	instance.news.active = active
	return instance
}

func (instance *builder) CreatedAt(createdAt time.Time) *builder {
	if createdAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de criação do registro do boletim é inválida")
		return instance
	}
	instance.news.createdAt = createdAt
	return instance
}

func (instance *builder) UpdatedAt(updatedAt time.Time) *builder {
	if updatedAt.IsZero() {
		instance.invalidFields = append(instance.invalidFields, "A data de atualização do registro do boletim é inválida")
		return instance
	}
	instance.news.updatedAt = updatedAt
	return instance
}

func (instance *builder) Build() (*News, error) {
	if len(instance.invalidFields) > 0 {
		return nil, errors.New(strings.Join(instance.invalidFields, ";"))
	}
	return instance.news, nil
}
