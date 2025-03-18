package queries

type propositionAffectedByVotingSqlManager struct{}

func PropositionAffectedByVoting() *propositionAffectedByVotingSqlManager {
	return &propositionAffectedByVotingSqlManager{}
}

func (propositionAffectedByVotingSqlManager) Insert() string {
	return `INSERT INTO proposition_affected_by_voting(proposition_id, voting_id)
			VALUES ($1, $2)`
}
