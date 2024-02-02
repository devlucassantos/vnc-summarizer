package queries

type newsSqlManager struct{}

func News() *newsSqlManager {
	return &newsSqlManager{}
}

func (newsSqlManager) InsertProposition() string {
	return `INSERT INTO news(proposition_id) VALUES ($1) RETURNING id`
}

func (newsSqlManager) InsertNewsletter() string {
	return `INSERT INTO news(newsletter_id) VALUES ($1) RETURNING id`
}

func (newsSqlManager) UpdateNewsletterReferenceDateTime() string {
	return `UPDATE news SET reference_date_time = TIMEZONE('America/Sao_Paulo'::TEXT, NOW()),
                updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
            FROM newsletter WHERE news.active = true AND newsletter.active = true AND news.newsletter_id = $1`
}
