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
	return `SELECT id AS article_type_id, description AS article_type_description,
       			color AS article_type_color, sort_order AS article_type_sort_order,
       			created_at AS article_type_created_at, updated_at AS article_type_updated_at
			FROM article_type
			WHERE active = true AND $1 = ANY(string_to_array(codes, ','))`
}

func (articleTypeSelectSqlManager) DefaultOption() string {
	return `SELECT id AS article_type_id, description AS article_type_description,
       			color AS article_type_color, sort_order AS article_type_sort_order,
       			created_at AS article_type_created_at, updated_at AS article_type_updated_at
			FROM article_type
			WHERE active = true AND 'default_option' = ANY(string_to_array(codes, ','))`
}
