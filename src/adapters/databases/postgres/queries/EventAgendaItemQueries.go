package queries

type eventAgendaItemSqlManager struct{}

func EventAgendaItem() *eventAgendaItemSqlManager {
	return &eventAgendaItemSqlManager{}
}

func (eventAgendaItemSqlManager) Insert() string {
	return `INSERT INTO event_agenda_item(title, topic, situation, agenda_item_regime_id, rapporteur_id,
				rapporteur_party_id, rapporteur_federated_unit, proposition_id, related_proposition_id, voting_id,
                event_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id`
}
