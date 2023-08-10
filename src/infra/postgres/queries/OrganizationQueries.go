package queries

type organizationSqlManager struct{}

func Organization() *organizationSqlManager {
	return &organizationSqlManager{}
}

func (organizationSqlManager) Insert() string {
	return `INSERT INTO organization(code, name, acronym, nickname) VALUES ($1, $2, $3, $4) RETURNING id`
}

func (organizationSqlManager) Update() string {
	return `UPDATE organization SET name = COALESCE($1, name), acronym = COALESCE($2, acronym),
            	nickname = COALESCE($3, nickname), updated_at = NOW()
            WHERE active = true AND code = $4`
}

type organizationSelectSqlManager struct{}

func (organizationSqlManager) Select() *organizationSelectSqlManager {
	return &organizationSelectSqlManager{}
}

func (organizationSelectSqlManager) ByCode() string {
	return `SELECT COALESCE(organization.id, '00000000-0000-0000-0000-000000000000') AS organization_id,
        		COALESCE(organization.code, 0) AS organization_code,
        		COALESCE(organization.name, '') AS organization_name,
        		COALESCE(organization.nickname, '') AS organization_nickname,
        		COALESCE(organization.acronym, '') AS organization_acronym,
        		COALESCE(organization.active, true) AS organization_active,
        		COALESCE(organization.created_at, '1970-01-01 00:00:00') AS organization_created_at,
        		COALESCE(organization.updated_at, '1970-01-01 00:00:00') AS organization_updated_at
			FROM organization WHERE active = true AND code = $1`
}
