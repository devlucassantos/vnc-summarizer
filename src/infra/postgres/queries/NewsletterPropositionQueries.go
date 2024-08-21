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
	return `SELECT proposition.id AS proposition_id, proposition.code AS proposition_code,
    			proposition.original_text_url AS proposition_original_text_url, proposition.title AS proposition_title,
    			proposition.content AS proposition_content, proposition.submitted_at AS proposition_submitted_at,
       			proposition.created_at AS proposition_created_at, proposition.updated_at AS proposition_updated_at
    		FROM newsletter_proposition
    			INNER JOIN proposition ON proposition.id = newsletter_proposition.proposition_id
    			INNER JOIN article ON article.newsletter_id = newsletter_proposition.newsletter_id
    		WHERE newsletter_proposition.active = true AND proposition.active = true AND article.active = true AND
    		      newsletter_proposition.newsletter_id = $1`
}
