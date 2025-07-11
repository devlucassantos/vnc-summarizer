package queries

type partySqlManager struct{}

func Party() *partySqlManager {
	return &partySqlManager{}
}

func (partySqlManager) Insert() string {
	return `INSERT INTO party(code, name, acronym, image_url)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
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
	return `SELECT id AS party_id, code AS party_code, name AS party_name, acronym AS party_acronym,
       			image_url AS party_image_url
			FROM party
			WHERE active = true AND code = $1`
}
