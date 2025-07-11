package services

import "github.com/devlucassantos/vnc-domains/src/domains/externalauthor"

type ExternalAuthor interface {
	GetExternalAuthorFromAuthorData(authorName string, authorTypeCode int, authorType string) (
		*externalauthor.ExternalAuthor, error)
}
