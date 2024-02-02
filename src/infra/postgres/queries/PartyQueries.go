package queries

type partySqlManager struct{}

func Party() *partySqlManager {
	return &partySqlManager{}
}

func (partySqlManager) Insert() string {
	return `INSERT INTO party(code, name, acronym, image_url) VALUES ($1, $2, $3, $4) RETURNING id`
}

func (partySqlManager) Update() string {
	return `UPDATE party SET name = COALESCE($1, name), acronym = COALESCE($2, acronym),
                 image_url = COALESCE($3, image_url), updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            WHERE active = true AND code = $4`
}

type partySelectSqlManager struct{}

func (partySqlManager) Select() *partySelectSqlManager {
	return &partySelectSqlManager{}
}

func (partySelectSqlManager) ByCode() string {
	return `SELECT COALESCE(id, '00000000-0000-0000-0000-000000000000') AS party_id,
        		COALESCE(code, 0) AS party_code,
        		COALESCE(name, '') AS party_name,
        		COALESCE(acronym, '') AS party_acronym,
        		COALESCE(image_url, '') AS party_image_url,
        		COALESCE(active, true) AS party_active,
        		COALESCE(created_at, '1970-01-01 00:00:00') AS party_created_at,
        		COALESCE(updated_at, '1970-01-01 00:00:00') AS party_updated_at
			FROM party WHERE active = true AND code = $1`
}
