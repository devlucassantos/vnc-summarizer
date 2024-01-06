package queries

type propositionSqlManager struct{}

func (propositionSqlManager) Insert() string {
	return `INSERT INTO proposition(code, original_text_url, title, content, submitted_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
}

func Proposition() *propositionSqlManager {
	return &propositionSqlManager{}
}

type propositionSelectSqlManager struct{}

func (propositionSqlManager) Select() *propositionSelectSqlManager {
	return &propositionSelectSqlManager{}
}

func (propositionSelectSqlManager) ByDate() string {
	return `SELECT COALESCE(id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
    			COALESCE(code, 0) AS proposition_code,
    			COALESCE(original_text_url, '') AS proposition_original_text_url,
       			COALESCE(title, '') AS proposition_title,
    			COALESCE(content, '') AS proposition_content,
    			COALESCE(submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
       			COALESCE(active, true) AS proposition_active,
    			COALESCE(created_at, '1970-01-01 00:00:00') AS proposition_created_at,
    			COALESCE(updated_at, '1970-01-01 00:00:00') AS proposition_updated_at
    		FROM proposition WHERE active = true AND DATE(submitted_at) = $1`
}

func (propositionSelectSqlManager) LatestPropositionsCodes() string {
	return `SELECT code FROM proposition WHERE active = true ORDER BY created_at DESC LIMIT 50`
}
