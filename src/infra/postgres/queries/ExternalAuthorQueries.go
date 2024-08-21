package queries

type externalAuthorSqlManager struct{}

func ExternalAuthor() *externalAuthorSqlManager {
	return &externalAuthorSqlManager{}
}

func (externalAuthorSqlManager) Insert() string {
	return `INSERT INTO external_author(name, type) VALUES ($1, $2) RETURNING id`
}

type externalAuthorSelectSqlManager struct{}

func (externalAuthorSqlManager) Select() *externalAuthorSelectSqlManager {
	return &externalAuthorSelectSqlManager{}
}

func (externalAuthorSelectSqlManager) ByNameAndType() string {
	return `SELECT id AS external_author_id, name AS external_author_name, type AS external_author_type,
        		created_at AS external_author_created_at, updated_at AS external_author_updated_at
			FROM external_author
			WHERE active = true AND name = $1 AND type = $2`
}
