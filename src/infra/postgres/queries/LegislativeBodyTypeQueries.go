package queries

type legislativeBodyTypeSqlManager struct{}

func LegislativeBodyType() *legislativeBodyTypeSqlManager {
	return &legislativeBodyTypeSqlManager{}
}

func (legislativeBodyTypeSqlManager) Insert() string {
	return `INSERT INTO legislative_body_type(code, description)
			VALUES ($1, $2)
			RETURNING id`
}

type legislativeBodyTypeSelectSqlManager struct{}

func (legislativeBodyTypeSqlManager) Select() *legislativeBodyTypeSelectSqlManager {
	return &legislativeBodyTypeSelectSqlManager{}
}

func (legislativeBodyTypeSelectSqlManager) ByCode() string {
	return `SELECT id AS legislative_body_type_id, code AS legislative_body_type_code,
				description AS legislative_body_type_description
			FROM legislative_body_type
			WHERE active = true AND code = $1`
}
