package queries

type propositionSqlManager struct{}

func (propositionSqlManager) Insert() string {
	return `INSERT INTO proposition(code, original_text_url, title, summary, submitted_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
}

func Proposition() *propositionSqlManager {
	return &propositionSqlManager{}
}

type propositionSelectSqlManager struct{}

func (propositionSqlManager) Select() *propositionSelectSqlManager {
	return &propositionSelectSqlManager{}
}

func (propositionSelectSqlManager) LatestPropositionsCodes() string {
	return `SELECT code FROM proposition WHERE active = true ORDER BY created_at DESC LIMIT 50`
}
