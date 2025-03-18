package queries

type articleSqlManager struct{}

func Article() *articleSqlManager {
	return &articleSqlManager{}
}

func (articleSqlManager) Insert() string {
	return `INSERT INTO article(article_type_id, reference_date_time)
			VALUES ($1, $2)
			RETURNING id`
}

type articleUpdateSqlManager struct{}

func (articleSqlManager) Update() *articleUpdateSqlManager {
	return &articleUpdateSqlManager{}
}

func (articleUpdateSqlManager) NewsletterReferenceDateTime() string {
	return `UPDATE article SET reference_date_time = TIMEZONE('America/Sao_Paulo'::TEXT, NOW()),
                updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            FROM newsletter
			WHERE article.id = newsletter.article_id AND newsletter.active = true AND article.active = true AND
				newsletter.id = $1`
}

type articleSelectSqlManager struct{}

func (articleSqlManager) Select() *articleSelectSqlManager {
	return &articleSelectSqlManager{}
}

func (articleSelectSqlManager) ByReferenceDate() string {
	return `SELECT article.id AS article_id,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.result, '') AS voting_result,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				voting.active IS NOT false AND event.active IS NOT false AND
				DATE_TRUNC('day', article.reference_date_time) = DATE_TRUNC('day', $1::timestamp)
			ORDER BY article.reference_date_time`
}
