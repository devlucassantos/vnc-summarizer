package queries

type propositionKeywordSqlManager struct{}

func PropositionKeyword() *propositionKeywordSqlManager {
	return &propositionKeywordSqlManager{}
}

func (propositionKeywordSqlManager) Insert() string {
	return `INSERT INTO proposition_keyword(proposition_id, keyword_id) VALUES ($1, $2)`
}
