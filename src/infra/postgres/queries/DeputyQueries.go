package queries

type deputySqlManager struct{}

func Deputy() *deputySqlManager {
	return &deputySqlManager{}
}

func (deputySqlManager) Insert() string {
	return `INSERT INTO deputy(code, cpf, name, electoral_name, image_url, party_id)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
}

func (deputySqlManager) Update() string {
	return `UPDATE deputy SET name = COALESCE($1, name), electoral_name = COALESCE($2, electoral_name),
                  image_url = COALESCE($3, image_url), party_id = COALESCE($4, party_id),
                  updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            WHERE active = true AND code = $5`
}

type deputySelectSqlManager struct{}

func (deputySqlManager) Select() *deputySelectSqlManager {
	return &deputySelectSqlManager{}
}

func (deputySelectSqlManager) ByCode() string {
	return `SELECT COALESCE(deputy.id, '00000000-0000-0000-0000-000000000000') AS deputy_id,
       			COALESCE(deputy.code, 0) AS deputy_code,
       			COALESCE(deputy.cpf, '') AS deputy_cpf,
       			COALESCE(deputy.name, '') AS deputy_name,
       			COALESCE(deputy.electoral_name, '') AS deputy_electoral_name,
       			COALESCE(deputy.image_url, '') AS deputy_image_url,
       			COALESCE(deputy.active, true) AS deputy_active,
       			COALESCE(deputy.created_at, '1970-01-01 00:00:00') AS deputy_created_at,
       			COALESCE(deputy.updated_at, '1970-01-01 00:00:00') AS deputy_updated_at,
       			
        		COALESCE(party.id, '00000000-0000-0000-0000-000000000000') AS party_id,
        		COALESCE(party.code, 0) AS party_code,
        		COALESCE(party.name, '') AS party_name,
        		COALESCE(party.acronym, '') AS party_acronym,
        		COALESCE(party.image_url, '') AS party_image_url,
        		COALESCE(party.active, true) AS party_active,
        		COALESCE(party.created_at, '1970-01-01 00:00:00') AS party_created_at,
        		COALESCE(party.updated_at, '1970-01-01 00:00:00') AS party_updated_at
			FROM deputy
				INNER JOIN party on party.id = deputy.party_id
			WHERE deputy.active = true AND deputy.code = $1`
}
