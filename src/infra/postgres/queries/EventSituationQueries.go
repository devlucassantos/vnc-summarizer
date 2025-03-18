package queries

type eventSituationSqlManager struct{}

func EventSituation() *eventSituationSqlManager {
	return &eventSituationSqlManager{}
}

type eventSituationSelectSqlManager struct{}

func (eventSituationSqlManager) Select() *eventSituationSelectSqlManager {
	return &eventSituationSelectSqlManager{}
}

func (eventSituationSelectSqlManager) ByCode() string {
	return `SELECT id AS event_situation_id, description AS event_situation_description, codes AS event_situation_codes,
       			color AS event_situation_color, is_finished AS event_situation_is_finished
			FROM event_situation
			WHERE active = true AND $1 = ANY(string_to_array(codes, ','))`
}

func (eventSituationSelectSqlManager) DefaultOption() string {
	return `SELECT id AS event_situation_id, description AS event_situation_description, codes AS event_situation_codes,
       			color AS event_situation_color, is_finished AS event_situation_is_finished
			FROM event_situation
			WHERE active = true AND 'default_option' = ANY(string_to_array(codes, ','))`
}
