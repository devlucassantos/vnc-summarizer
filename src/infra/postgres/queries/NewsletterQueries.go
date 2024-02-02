package queries

type newsletterSqlManager struct{}

func Newsletter() *newsletterSqlManager {
	return &newsletterSqlManager{}
}

func (newsletterSqlManager) Insert() string {
	return `INSERT INTO newsletter(title, content, image_url, reference_date) VALUES ($1, $2, '', $3) RETURNING id`
}

func (newsletterSqlManager) Update() string {
	return `UPDATE newsletter SET title = COALESCE($1, title), content = COALESCE($2, content),
                      image_url = COALESCE('', image_url), updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            FROM news WHERE newsletter.active = true AND news.active = true AND newsletter.id = $3`
}

type newsletterSelectSqlManager struct{}

func (newsletterSqlManager) Select() *newsletterSelectSqlManager {
	return &newsletterSelectSqlManager{}
}

func (newsletterSelectSqlManager) ByReferenceDate() string {
	return `SELECT COALESCE(newsletter.id, '00000000-0000-0000-0000-000000000000') AS newsletter_id,
       			COALESCE(newsletter.title, '') AS newsletter_title,
    			COALESCE(newsletter.content, '') AS newsletter_content,
    			COALESCE(newsletter.reference_date, '1970-01-01 00:00:00') AS newsletter_reference_date,
       			COALESCE(newsletter.active, true) AS newsletter_active,
    			COALESCE(newsletter.created_at, '1970-01-01 00:00:00') AS newsletter_created_at,
    			COALESCE(newsletter.updated_at, '1970-01-01 00:00:00') AS newsletter_updated_at,
    			
				COALESCE(news.id, '00000000-0000-0000-0000-000000000000') AS news_id
    		FROM newsletter
    			INNER JOIN news ON news.newsletter_id = newsletter.id
    		WHERE newsletter.active = true AND news.active = true AND newsletter.reference_date = $1`
}
