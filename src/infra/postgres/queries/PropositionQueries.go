package queries

type propositionSqlManager struct{}

func Proposition() *propositionSqlManager {
	return &propositionSqlManager{}
}

func (propositionSqlManager) Insert() string {
	return `INSERT INTO proposition(code, original_text_url, title, content, submitted_at, image_url, specific_type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`
}

type propositionSelectSqlManager struct{}

func (propositionSqlManager) Select() *propositionSelectSqlManager {
	return &propositionSelectSqlManager{}
}

func (propositionSelectSqlManager) ByDate() string {
	return `SELECT id AS proposition_id, code AS proposition_code, original_text_url AS proposition_original_text_url,
       			title AS proposition_title, content AS proposition_content, submitted_at AS proposition_submitted_at,
       			created_at AS proposition_created_at, updated_at AS proposition_updated_at
    		FROM proposition
    		WHERE active = true AND DATE(submitted_at) = $1`
}

func (propositionSelectSqlManager) LatestPropositionsCodes() string {
	return `SELECT code
			FROM proposition
			WHERE active = true
			ORDER BY created_at DESC
			LIMIT 50`
}
