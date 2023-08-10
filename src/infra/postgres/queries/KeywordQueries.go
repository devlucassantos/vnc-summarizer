package queries

type keywordSqlManager struct{}

func Keyword() *keywordSqlManager {
	return &keywordSqlManager{}
}

func (keywordSqlManager) Insert() string {
	return `INSERT INTO keyword(keyword) VALUES ($1) RETURNING id`
}

type keywordSelectSqlManager struct{}

func (keywordSqlManager) Select() *keywordSelectSqlManager {
	return &keywordSelectSqlManager{}
}

func (keywordSelectSqlManager) ByKeyword() string {
	return `SELECT COALESCE(keyword.id, '00000000-0000-0000-0000-000000000000') AS keyword_id,
       			COALESCE(keyword.keyword, '') AS keyword_keyword,
       			COALESCE(keyword.active, true) AS keyword_active,
       			COALESCE(keyword.created_at, '1970-01-01 00:00:00') AS keyword_created_at,
       			COALESCE(keyword.updated_at, '1970-01-01 00:00:00') AS keyword_updated_at
			FROM keyword WHERE keyword.keyword = $1`
}
