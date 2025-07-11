package queries

type eventTypeSqlManager struct{}

func EventType() *eventTypeSqlManager {
	return &eventTypeSqlManager{}
}

type eventTypeSelectSqlManager struct{}

func (eventTypeSqlManager) Select() *eventTypeSelectSqlManager {
	return &eventTypeSelectSqlManager{}
}

func (eventTypeSelectSqlManager) ByCode() string {
	return `SELECT id AS event_type_id, description AS event_type_description, codes AS event_type_codes,
				color AS event_type_color
			FROM event_type
			WHERE active = true AND $1 = ANY(string_to_array(codes, ','))`
}

func (eventTypeSelectSqlManager) DefaultOption() string {
	return `SELECT id AS event_type_id, description AS event_type_description, codes AS event_type_codes,
       			color AS event_type_color
			FROM event_type
			WHERE active = true AND 'default_option' = ANY(string_to_array(codes, ','))`
}
