package queries

type newsletterPropositionSqlManager struct{}

func NewsletterProposition() *newsletterPropositionSqlManager {
	return &newsletterPropositionSqlManager{}
}

func (newsletterPropositionSqlManager) Insert() string {
	return `INSERT INTO newsletter_proposition(newsletter_id, proposition_id) VALUES ($1, $2)`
}
