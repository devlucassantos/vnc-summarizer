package queries

type eventLegislativeBodySqlManager struct{}

func EventLegislativeBody() *eventLegislativeBodySqlManager {
	return &eventLegislativeBodySqlManager{}
}

func (eventLegislativeBodySqlManager) Insert() string {
	return `INSERT INTO event_legislative_body(event_id, legislative_body_id)
			VALUES ($1, $2)`
}
