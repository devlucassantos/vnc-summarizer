package organization

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Organization struct {
	id        uuid.UUID
	code      int
	name      string
	nickname  string
	acronym   string
	_type     string
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func (instance *Organization) NewUpdater() *builder {
	return &builder{organization: instance}
}

func (instance *Organization) Id() uuid.UUID {
	return instance.id
}

func (instance *Organization) Code() int {
	return instance.code
}

func (instance *Organization) Name() string {
	return instance.name
}

func (instance *Organization) Nickname() string {
	return instance.nickname
}

func (instance *Organization) Acronym() string {
	return instance.acronym
}

func (instance *Organization) Type() string {
	return instance._type
}

func (instance *Organization) Active() bool {
	return instance.active
}

func (instance *Organization) CreatedAt() time.Time {
	return instance.createdAt
}

func (instance *Organization) UpdatedAt() time.Time {
	return instance.updatedAt
}

func (instance *Organization) IsEqual(organization Organization) bool {
	return instance.code == organization.code &&
		instance.name == organization.name &&
		instance.acronym == organization.acronym &&
		instance.nickname == organization.nickname
}

func (instance *Organization) IsZero() bool {
	return reflect.DeepEqual(instance, &Organization{})
}
