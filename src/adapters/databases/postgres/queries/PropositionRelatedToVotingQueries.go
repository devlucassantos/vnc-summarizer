package queries

type propositionRelatedToVotingSqlManager struct{}

func PropositionRelatedToVoting() *propositionRelatedToVotingSqlManager {
	return &propositionRelatedToVotingSqlManager{}
}

func (propositionRelatedToVotingSqlManager) Insert() string {
	return `INSERT INTO proposition_related_to_voting(proposition_id, voting_id)
			VALUES ($1, $2)`
}
