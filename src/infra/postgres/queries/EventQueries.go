package queries

import (
	"fmt"
	"strings"
)

type eventSqlManager struct{}

func Event() *eventSqlManager {
	return &eventSqlManager{}
}

func (eventSqlManager) Insert() string {
	return `INSERT INTO event(code, title, description, starts_at, ends_at, location, is_internal, video_url,
				specific_type, event_type_id, specific_situation, event_situation_id, article_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING id`
}

func (eventSqlManager) Update() string {
	return `UPDATE event SET description = $1, starts_at = $2, ends_at = $3, location = $4, is_internal = $5,
                 video_url = $6, specific_situation = $7, event_situation_id = $8,
                 updated_at = TIMEZONE('America/Sao_Paulo'::TEXT, NOW())
			WHERE event.active = true AND event.code = $9`
}

type eventSelectSqlManager struct{}

func (eventSqlManager) Select() *eventSelectSqlManager {
	return &eventSelectSqlManager{}
}

func (eventSelectSqlManager) ByCodes(numberOfEvents int) string {
	var parameters []string
	for i := 1; i <= numberOfEvents; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS event_id, code AS event_code, title AS event_title,
				description AS event_description, starts_at AS event_starts_at,
				COALESCE(ends_at, '1970-01-01 00:00:00') AS event_ends_at, location AS event_location,
				is_internal AS event_is_internal, COALESCE(video_url, '') AS event_video_url
			FROM event
			WHERE active = true AND code IN (%s)`, strings.Join(parameters, ","))
}

func (eventSelectSqlManager) OccurringToday() string {
	return `SELECT event.id AS event_id, event.code AS event_code, event.title AS event_title,
				event.description AS event_description, event.starts_at AS event_starts_at,
				COALESCE(event.ends_at, '1970-01-01 00:00:00') AS event_ends_at, event.location AS event_location,
				event.is_internal AS event_is_internal, COALESCE(event.video_url, '') AS event_video_url,
				event.specific_situation AS event_specific_situation,
				event_situation.id AS event_situation_id, event_situation.description AS event_situation_description,
				event_situation.codes AS event_situation_codes, event_situation.color AS event_situation_color,
				event_situation.is_finished AS event_situation_is_finished
			FROM event
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
			WHERE event.active = true AND event_situation.active = true AND event_situation.is_finished = false AND
				DATE_TRUNC('day', event.starts_at) <= CURRENT_DATE AND DATE_TRUNC('day', event.ends_at) >= CURRENT_DATE
			ORDER BY event.starts_at, event.ends_at`
}

func (eventSelectSqlManager) StartedInTheLastThreeMonthsAndHaveNotFinished() string {
	return `SELECT event.id AS event_id, event.code AS event_code, event.title AS event_title,
				event.description AS event_description, event.starts_at AS event_starts_at,
				COALESCE(event.ends_at, '1970-01-01 00:00:00') AS event_ends_at, event.location AS event_location,
				event.is_internal AS event_is_internal, COALESCE(event.video_url, '') AS event_video_url,
				event.specific_situation AS event_specific_situation,
				event_situation.id AS event_situation_id, event_situation.description AS event_situation_description,
				event_situation.codes AS event_situation_codes, event_situation.color AS event_situation_color,
				event_situation.is_finished AS event_situation_is_finished
			FROM event
				INNER JOIN event_situation ON event_situation.id = event.event_situation_id
			WHERE event.active = true AND event_situation.active = true AND event_situation.is_finished = false AND
			      DATE_TRUNC('day', event.starts_at) >= CURRENT_DATE - INTERVAL '3 months'
			ORDER BY event.starts_at, event.ends_at`
}
