package queries

type organizationSqlManager struct{}

func Organization() *organizationSqlManager {
	return &organizationSqlManager{}
}

func (organizationSqlManager) Insert() string {
	return `INSERT INTO organization(code, name, acronym, nickname, type) VALUES ($1, $2, $3, $4, $5) RETURNING id`
}

func (organizationSqlManager) UpdateByCode() string {
	return `UPDATE organization SET name = COALESCE($1, name), acronym = COALESCE($2, acronym),
            	nickname = COALESCE($3, nickname), type = COALESCE($4, type),
				updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            WHERE active = true AND code = $4`
}

func (organizationSqlManager) UpdateByNameAndType() string {
	return `UPDATE organization SET acronym = COALESCE($1, acronym), nickname = COALESCE($2, nickname),
				updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            WHERE active = true AND name = $3 AND type = $4`
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

func (organizationSelectSqlManager) ByNameAndType() string {
	return `SELECT COALESCE(organization.id, '00000000-0000-0000-0000-000000000000') AS organization_id,
        		COALESCE(organization.code, 0) AS organization_code,
        		COALESCE(organization.name, '') AS organization_name,
        		COALESCE(organization.nickname, '') AS organization_nickname,
        		COALESCE(organization.acronym, '') AS organization_acronym,
        		COALESCE(organization.active, true) AS organization_active,
        		COALESCE(organization.created_at, '1970-01-01 00:00:00') AS organization_created_at,
        		COALESCE(organization.updated_at, '1970-01-01 00:00:00') AS organization_updated_at
			FROM organization WHERE active = true AND name = $1 AND type = $2`
}
