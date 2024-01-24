package queries

type newsletterPropositionSqlManager struct{}

func NewsletterProposition() *newsletterPropositionSqlManager {
	return &newsletterPropositionSqlManager{}
}

func (newsletterPropositionSqlManager) Insert() string {
	return `INSERT INTO newsletter_proposition(newsletter_id, proposition_id) VALUES ($1, $2)`
}

type newsletterPropositionSelectSqlManager struct{}

func (newsletterPropositionSqlManager) Select() *newsletterPropositionSelectSqlManager {
	return &newsletterPropositionSelectSqlManager{}
}

func (newsletterPropositionSelectSqlManager) ByNewsletterId() string {
	return `SELECT COALESCE(proposition.id, '00000000-0000-0000-0000-000000000000') AS proposition_id,
    			COALESCE(proposition.code, 0) AS proposition_code,
    			COALESCE(proposition.original_text_url, '') AS proposition_original_text_url,
       			COALESCE(proposition.title, '') AS proposition_title,
    			COALESCE(proposition.content, '') AS proposition_content,
    			COALESCE(proposition.submitted_at, '1970-01-01 00:00:00') AS proposition_submitted_at,
       			COALESCE(proposition.active, true) AS proposition_active,
    			COALESCE(proposition.created_at, '1970-01-01 00:00:00') AS proposition_created_at,
    			COALESCE(proposition.updated_at, '1970-01-01 00:00:00') AS proposition_updated_at
    		FROM newsletter_proposition
    			INNER JOIN proposition ON proposition.id = newsletter_proposition.proposition_id
    			INNER JOIN news ON news.newsletter_id = newsletter_proposition.newsletter_id
    		WHERE newsletter_proposition.active = true AND proposition.active = true AND news.active = true AND
    		      newsletter_proposition.newsletter_id = $1`
}
