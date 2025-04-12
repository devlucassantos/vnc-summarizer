package queries

type externalAuthorTypeSqlManager struct{}

func ExternalAuthorType() *externalAuthorTypeSqlManager {
	return &externalAuthorTypeSqlManager{}
}

func (externalAuthorTypeSqlManager) Insert() string {
	return `INSERT INTO external_author_type(code, description)
			VALUES ($1, $2)
			RETURNING id`
}

type externalAuthorTypeSelectSqlManager struct{}

func (externalAuthorTypeSqlManager) Select() *externalAuthorTypeSelectSqlManager {
	return &externalAuthorTypeSelectSqlManager{}
}

func (externalAuthorTypeSelectSqlManager) ByCode() string {
	return `SELECT id AS external_author_type_id, code AS external_author_type_code,
       			description AS external_author_type_description
			FROM external_author_type
			WHERE active = true AND code = $1`
}
