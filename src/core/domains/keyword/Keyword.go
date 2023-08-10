package keyword

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Keyword struct {
	id        uuid.UUID
	keyword   string
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func (instance *Keyword) NewUpdater() *builder {
	return &builder{keyword: instance}
}

func (instance *Keyword) Id() uuid.UUID {
	return instance.id
}

func (instance *Keyword) Keyword() string {
	return instance.keyword
}

func (instance *Keyword) Active() bool {
	return instance.active
}

func (instance *Keyword) CreatedAt() time.Time {
	return instance.createdAt
}

func (instance *Keyword) UpdatedAt() time.Time {
	return instance.updatedAt
}

func (instance *Keyword) IsEqual(keyword Keyword) bool {
	return instance.keyword == keyword.keyword
}

func (instance *Keyword) IsZero() bool {
	return reflect.DeepEqual(instance, &Keyword{})
}
