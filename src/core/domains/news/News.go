package news

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type News struct {
	id        uuid.UUID
	title     string
	content   string
	views     int
	_type     string
	active    bool
	createdAt time.Time
	updatedAt time.Time
}

func (instance *News) NewUpdater() *builder {
	return &builder{news: instance}
}

func (instance *News) Id() uuid.UUID {
	return instance.id
}

func (instance *News) Title() string {
	return instance.title
}

func (instance *News) Content() string {
	return instance.content
}

func (instance *News) Views() int {
	return instance.views
}

func (instance *News) Type() string {
	return instance._type
}

func (instance *News) Active() bool {
	return instance.active
}

func (instance *News) CreatedAt() time.Time {
	return instance.createdAt
}

func (instance *News) UpdatedAt() time.Time {
	return instance.updatedAt
}

func (instance *News) IsZero() bool {
	return reflect.DeepEqual(instance, &News{})
}
