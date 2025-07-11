package queries

type eventRequirementSqlManager struct{}

func EventRequirement() *eventRequirementSqlManager {
	return &eventRequirementSqlManager{}
}

func (eventRequirementSqlManager) Insert() string {
	return `INSERT INTO event_requirement(event_id, proposition_id)
			VALUES ($1, $2)`
}
