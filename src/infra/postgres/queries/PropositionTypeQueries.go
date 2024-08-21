package queries

type propositionTypeSqlManager struct{}

func PropositionType() *propositionTypeSqlManager {
	return &propositionTypeSqlManager{}
}

type propositionTypeSelectSqlManager struct{}

func (propositionTypeSqlManager) Select() *propositionTypeSelectSqlManager {
	return &propositionTypeSelectSqlManager{}
}

func (propositionTypeSelectSqlManager) ByCode() string {
	return `SELECT id AS proposition_type_id, description AS proposition_type_description,
       			codes AS proposition_type_codes, created_at AS proposition_type_created_at,
				updated_at AS proposition_type_updated_at
			FROM proposition_type WHERE active = true AND $1 = ANY(string_to_array(codes, ','))`
}

func (propositionTypeSelectSqlManager) OtherOption() string {
	return `SELECT id AS proposition_type_id, description AS proposition_type_description,
       			codes AS proposition_type_codes, created_at AS proposition_type_created_at,
				updated_at AS proposition_type_updated_at
			FROM proposition_type WHERE active = true AND description = 'Outras Proposições'`
}
