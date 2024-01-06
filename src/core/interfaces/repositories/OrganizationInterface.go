package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/organization"
)

type Organization interface {
	CreateOrganization(organization organization.Organization) (*uuid.UUID, error)
	UpdateOrganization(organization organization.Organization) error
	GetOrganization(organization organization.Organization) (*organization.Organization, error)
}
