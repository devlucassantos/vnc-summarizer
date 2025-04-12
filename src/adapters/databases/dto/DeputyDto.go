package dto

import (
	"github.com/google/uuid"
)

type Deputy struct {
	Id            uuid.UUID `db:"deputy_id"`
	Code          int       `db:"deputy_code"`
	Cpf           string    `db:"deputy_cpf"`
	Name          string    `db:"deputy_name"`
	ElectoralName string    `db:"deputy_electoral_name"`
	ImageUrl      string    `db:"deputy_image_url"`
	FederatedUnit string    `db:"deputy_federated_unit"`
	*Party
}
