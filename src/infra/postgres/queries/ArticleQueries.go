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
				article_type.codes AS article_type_codes,
				COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(proposition_type.id, '00000000-0000-0000-0000-000000000000') AS proposition_type_id,
				COALESCE(proposition_type.description, '') AS proposition_type_description,
				COALESCE(proposition_type.codes, '') AS proposition_type_codes,
				COALESCE(voting.id, '00000000-0000-0000-0000-000000000000') AS voting_id,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.description, '') AS voting_description,
				COALESCE(event.id, '00000000-0000-0000-0000-000000000000') AS event_id,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description,
				COALESCE(event_type.id, '00000000-0000-0000-0000-000000000000') AS event_type_id,
				COALESCE(event_type.description, '') AS event_type_description,
				COALESCE(event_type.codes, '') AS event_type_codes
			FROM article
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN proposition_type ON proposition_type.id = proposition.proposition_type_id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
				LEFT JOIN event_type ON event_type.id = event.event_type_id
			WHERE article.active = true AND article_type.active = true AND proposition.active IS NOT false AND
				proposition_type.active IS NOT false AND voting.active IS NOT false AND
				event.active IS NOT false AND event_type.active IS NOT false AND
				DATE_TRUNC('day', article.reference_date_time) = DATE_TRUNC('day', $1::timestamp)
			ORDER BY article.reference_date_time`
}
