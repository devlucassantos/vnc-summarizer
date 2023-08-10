package queries

type propositionAuthorSqlManager struct{}

func PropositionAuthor() *propositionAuthorSqlManager {
	return &propositionAuthorSqlManager{}
}

func (propositionAuthorSqlManager) InsertDeputy() string {
	return `INSERT INTO proposition_author(proposition_id, deputy_id, party_id) VALUES ($1, $2, $3) RETURNING id`
}

func (propositionAuthorSqlManager) InsertOrganization() string {
	return `INSERT INTO proposition_author(proposition_id, organization_id) VALUES ($1, $2) RETURNING id`
}
