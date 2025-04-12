package queries

type articleTypeSqlManager struct{}

func ArticleType() *articleTypeSqlManager {
	return &articleTypeSqlManager{}
}

type articleTypeSelectSqlManager struct{}

func (articleTypeSqlManager) Select() *articleTypeSelectSqlManager {
	return &articleTypeSelectSqlManager{}
}

func (articleTypeSelectSqlManager) ByCode() string {
	return `SELECT id AS article_type_id, description AS article_type_description, codes AS article_type_codes,
       			color AS article_type_color
			FROM article_type
			WHERE active = true AND $1 = ANY(string_to_array(codes, ','))`
}
