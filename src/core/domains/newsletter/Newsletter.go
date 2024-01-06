package newsletter

import (
	"github.com/google/uuid"
	"reflect"
	"time"
	"vnc-write-api/core/domains/proposition"
)

type Newsletter struct {
	id            uuid.UUID
	title         string
	content       string
	referenceDate time.Time
	propositions  []proposition.Proposition
	active        bool
	createdAt     time.Time
	updatedAt     time.Time
}

func (instance *Newsletter) NewUpdater() *builder {
	return &builder{newsletter: instance}
}

func (instance *Newsletter) Id() uuid.UUID {
	return instance.id
}

func (instance *Newsletter) Title() string {
	return instance.title
}

func (instance *Newsletter) Content() string {
	return instance.content
}

func (instance *Newsletter) ReferenceDate() time.Time {
	return instance.referenceDate
}

func (instance *Newsletter) Propositions() []proposition.Proposition {
	return instance.propositions
}

func (instance *Newsletter) Active() bool {
	return instance.active
}

func (instance *Newsletter) CreatedAt() time.Time {
	return instance.createdAt
}

func (instance *Newsletter) UpdatedAt() time.Time {
	return instance.updatedAt
}

func (instance *Newsletter) IsZero() bool {
	return reflect.DeepEqual(instance, &Newsletter{})
}
