package queries

type articleSqlManager struct{}

func Article() *articleSqlManager {
	return &articleSqlManager{}
}

type articleInsertSqlManager struct{}

func (articleSqlManager) Insert() *articleInsertSqlManager {
	return &articleInsertSqlManager{}
}

func (articleInsertSqlManager) Proposition() string {
	return `INSERT INTO article(proposition_id, article_type_id)
			VALUES ($1, $2)
			RETURNING id`
}

func (articleInsertSqlManager) Newsletter() string {
	return `INSERT INTO article(newsletter_id, article_type_id)
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
			WHERE article.active = true AND newsletter.active = true AND article.newsletter_id = $1`
}
