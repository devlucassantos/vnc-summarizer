package proposition

import (
	"github.com/google/uuid"
	"reflect"
	"time"
	"vnc-write-api/core/domains/deputy"
	"vnc-write-api/core/domains/keyword"
	"vnc-write-api/core/domains/organization"
)

type Proposition struct {
	id              uuid.UUID
	code            int
	originalTextUrl string
	title           string
	summary         string
	submittedAt     time.Time
	deputies        []deputy.Deputy
	organizations   []organization.Organization
	keywords        []keyword.Keyword
	active          bool
	createdAt       time.Time
	updatedAt       time.Time
}

func (instance *Proposition) NewUpdater() *builder {
	return &builder{proposition: instance}
}

func (instance *Proposition) Id() uuid.UUID {
	return instance.id
}

func (instance *Proposition) Code() int {
	return instance.code
}

func (instance *Proposition) OriginalTextUrl() string {
	return instance.originalTextUrl
}

func (instance *Proposition) Title() string {
	return instance.title
}

func (instance *Proposition) Summary() string {
	return instance.summary
}

func (instance *Proposition) SubmittedAt() time.Time {
	return instance.submittedAt
}

func (instance *Proposition) Deputies() []deputy.Deputy {
	return instance.deputies
}

func (instance *Proposition) Organizations() []organization.Organization {
	return instance.organizations
}

func (instance *Proposition) Keywords() []keyword.Keyword {
	return instance.keywords
}

func (instance *Proposition) Active() bool {
	return instance.active
}

func (instance *Proposition) CreatedAt() time.Time {
	return instance.createdAt
}

func (instance *Proposition) UpdatedAt() time.Time {
	return instance.updatedAt
}

func (instance *Proposition) IsZero() bool {
	return reflect.DeepEqual(instance, &Proposition{})
}
