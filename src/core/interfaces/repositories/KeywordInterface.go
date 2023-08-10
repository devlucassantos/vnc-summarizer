package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/keyword"
)

type Keyword interface {
	CreateKeyword(keyword keyword.Keyword) (*uuid.UUID, error)
	GetKeywordByKeyword(keyword string) (*keyword.Keyword, error)
}
