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
	return `SELECT deputy.id AS deputy_id, deputy.code AS deputy_code, deputy.cpf AS deputy_cpf,
       			deputy.name AS deputy_name, deputy.electoral_name AS deputy_electoral_name,
       			deputy.image_url AS deputy_image_url, deputy.created_at AS deputy_created_at, deputy.updated_at AS deputy_updated_at,
        		party.id AS party_id, party.code AS party_code, party.name AS party_name, party.acronym AS party_acronym,
        		party.image_url AS party_image_url, party.created_at AS party_created_at, party.updated_at AS party_updated_at
			FROM deputy
				INNER JOIN party on party.id = deputy.party_id
			WHERE deputy.active = true AND deputy.code = $1`
}
