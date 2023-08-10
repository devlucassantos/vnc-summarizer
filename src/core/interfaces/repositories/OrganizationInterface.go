package repositories

import (
	"github.com/google/uuid"
	"vnc-write-api/core/domains/organization"
)

type Organization interface {
	CreateOrganization(organization organization.Organization) (*uuid.UUID, error)
	UpdateOrganization(organization organization.Organization) error
	GetOrganizationByCode(code int) (*organization.Organization, error)
}
