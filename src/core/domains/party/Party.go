package party

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Party struct {
	id        uuid.UUID
	code      int
	name      string
	acronym   string
	imageUrl  string
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func (instance *Party) NewUpdater() *builder {
	return &builder{party: instance}
}

func (instance *Party) Id() uuid.UUID {
	return instance.id
}

func (instance *Party) Code() int {
	return instance.code
}

func (instance *Party) Name() string {
	return instance.name
}

func (instance *Party) Acronym() string {
	return instance.acronym
}

func (instance *Party) ImageUrl() string {
	return instance.imageUrl
}

func (instance *Party) Active() bool {
	return instance.active
}

func (instance *Party) CreatedAt() time.Time {
	return instance.createdAt
}

func (instance *Party) UpdatedAt() time.Time {
	return instance.updatedAt
}

func (instance *Party) IsEqual(party Party) bool {
	return instance.code == party.code &&
		instance.name == party.name &&
		instance.acronym == party.acronym &&
		instance.imageUrl == party.imageUrl
}

func (instance *Party) IsZero() bool {
	return reflect.DeepEqual(instance, &Party{})
}
