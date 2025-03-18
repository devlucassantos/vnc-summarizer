package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
)

type Author interface {
	GetAuthorsFromAuthorsUrl(authorsUrl string) ([]deputy.Deputy, []externalauthor.ExternalAuthor, error)
}
