package queries

type propositionAuthorSqlManager struct{}

func PropositionAuthor() *propositionAuthorSqlManager {
	return &propositionAuthorSqlManager{}
}

type propositionAuthorInsertSqlManager struct{}

func (propositionAuthorSqlManager) Insert() *propositionAuthorInsertSqlManager {
	return &propositionAuthorInsertSqlManager{}
}

func (propositionAuthorInsertSqlManager) Deputy() string {
	return `INSERT INTO proposition_author(proposition_id, deputy_id, party_id, federated_unit)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
}

func (propositionAuthorInsertSqlManager) ExternalAuthor() string {
	return `INSERT INTO proposition_author(proposition_id, external_author_id)
			VALUES ($1, $2)
			RETURNING id`
}
