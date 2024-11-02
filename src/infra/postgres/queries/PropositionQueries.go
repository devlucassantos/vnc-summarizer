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
	return `SELECT article.id AS article_id, proposition.id AS proposition_id, proposition.code AS proposition_code,
       			proposition.original_text_url AS proposition_original_text_url, proposition.title AS proposition_title,
       			proposition.content AS proposition_content, proposition.submitted_at AS proposition_submitted_at,
       			proposition.created_at AS proposition_created_at, proposition.updated_at AS proposition_updated_at,
       			article_type.id AS article_type_id, article_type.description AS article_type_description,
       			article_type.codes AS article_type_codes, article_type.color AS article_type_color,
       			article_type.sort_order AS article_type_sort_order, article_type.created_at AS article_type_created_at,
       			article_type.updated_at AS article_type_updated_at
    		FROM proposition
    		    INNER JOIN article ON article.proposition_id = proposition.id
    			INNER JOIN article_type ON article_type.id = article.article_type_id
    		WHERE proposition.active = true AND DATE(proposition.submitted_at) = $1
    		ORDER BY article.reference_date_time`
}

func (propositionSelectSqlManager) LatestPropositionsCodes() string {
	return `SELECT code
			FROM proposition
			WHERE active = true
			ORDER BY created_at DESC
			LIMIT 50`
}
