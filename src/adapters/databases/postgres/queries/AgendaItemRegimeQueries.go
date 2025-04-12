package queries

type agendaItemRegimeSqlManager struct{}

func AgendaItemRegime() *agendaItemRegimeSqlManager {
	return &agendaItemRegimeSqlManager{}
}

func (agendaItemRegimeSqlManager) Insert() string {
	return `INSERT INTO agenda_item_regime(code, description)
			VALUES ($1, $2)
			RETURNING id`
}

type agendaItemRegimeSelectSqlManager struct{}

func (agendaItemRegimeSqlManager) Select() *agendaItemRegimeSelectSqlManager {
	return &agendaItemRegimeSelectSqlManager{}
}

func (agendaItemRegimeSelectSqlManager) ByCode() string {
	return `SELECT id AS agenda_item_regime_id, code AS agenda_item_regime_code,
       			description AS agenda_item_regime_description
			FROM agenda_item_regime
			WHERE active = true AND code = $1`
}
