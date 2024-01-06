package proposition

import (
	"github.com/google/uuid"
	"reflect"
	"time"
	"vnc-write-api/core/domains/deputy"
	"vnc-write-api/core/domains/organization"
)

type Proposition struct {
	id              uuid.UUID
	code            int
	originalTextUrl string
	title           string
	content         string
	submittedAt     time.Time
	imageUrl        string
	deputies        []deputy.Deputy
	organizations   []organization.Organization
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

func (instance *Proposition) Content() string {
	return instance.content
}

func (instance *Proposition) SubmittedAt() time.Time {
	return instance.submittedAt
}

func (instance *Proposition) ImageUrl() string {
	return instance.imageUrl
}

func (instance *Proposition) Deputies() []deputy.Deputy {
	return instance.deputies
}

func (instance *Proposition) Organizations() []organization.Organization {
	return instance.organizations
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
