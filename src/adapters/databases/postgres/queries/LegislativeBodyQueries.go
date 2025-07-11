package queries

import (
	"fmt"
	"strings"
)

type legislativeBodySqlManager struct{}

func LegislativeBody() *legislativeBodySqlManager {
	return &legislativeBodySqlManager{}
}

func (legislativeBodySqlManager) Insert() string {
	return `INSERT INTO legislative_body(code, name, acronym, legislative_body_type_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
}

type legislativeBodySelectSqlManager struct{}

func (legislativeBodySqlManager) Select() *legislativeBodySelectSqlManager {
	return &legislativeBodySelectSqlManager{}
}

func (legislativeBodySelectSqlManager) ByCode() string {
	return `SELECT legislative_body.id AS legislative_body_id,
				legislative_body.code AS legislative_body_code,
       			legislative_body.name AS legislative_body_name,
       			legislative_body.acronym AS legislative_body_acronym,
				legislative_body_type.id AS legislative_body_type_id,
				legislative_body_type.code AS legislative_body_type_code,
				legislative_body_type.description AS legislative_body_type_description
			FROM legislative_body
				INNER JOIN legislative_body_type ON legislative_body_type.id = legislative_body.legislative_body_type_id
			WHERE legislative_body.active = true AND legislative_body_type.active = true AND
				legislative_body.code = $1`
}

func (legislativeBodySelectSqlManager) ByCodes(numberOfLegislativeBodies int) string {
	var parameters []string
	for i := 1; i <= numberOfLegislativeBodies; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT legislative_body.id AS legislative_body_id,
				legislative_body.code AS legislative_body_code,
				legislative_body.name AS legislative_body_name,
				legislative_body.acronym AS legislative_body_acronym,
				legislative_body_type.id AS legislative_body_type_id,
				legislative_body_type.code AS legislative_body_type_code,
				legislative_body_type.description AS legislative_body_type_description
			FROM legislative_body
				INNER JOIN legislative_body_type ON legislative_body_type.id = legislative_body.legislative_body_type_id
			WHERE legislative_body.active = true AND legislative_body_type.active = true AND
				legislative_body.code IN (%s)`, strings.Join(parameters, ","))
}
