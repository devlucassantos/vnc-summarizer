package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/google/uuid"
)

type Organization interface {
	CreateOrganization(organization organization.Organization) (*uuid.UUID, error)
	UpdateOrganization(organization organization.Organization) error
	GetOrganization(organization organization.Organization) (*organization.Organization, error)
}
