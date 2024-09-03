package queries

type newsletterSqlManager struct{}

func Newsletter() *newsletterSqlManager {
	return &newsletterSqlManager{}
}

func (newsletterSqlManager) Insert() string {
	return `INSERT INTO newsletter(reference_date, title, description)
			VALUES ($1, $2, $3)
			RETURNING id`
}

func (newsletterSqlManager) Update() string {
	return `UPDATE newsletter SET description = COALESCE($1, description),
            	updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            FROM article
            WHERE newsletter.active = true AND article.active = true AND newsletter.id = $2`
}

type newsletterSelectSqlManager struct{}

func (newsletterSqlManager) Select() *newsletterSelectSqlManager {
	return &newsletterSelectSqlManager{}
}

func (newsletterSelectSqlManager) ByReferenceDate() string {
	return `SELECT id AS newsletter_id, reference_date AS newsletter_reference_date, title AS newsletter_title,
       			description AS newsletter_description, created_at AS newsletter_created_at, updated_at AS newsletter_updated_at
    		FROM newsletter
    		WHERE active = true AND reference_date = $1`
}
