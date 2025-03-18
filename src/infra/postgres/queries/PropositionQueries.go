package queries

import (
	"fmt"
	"strings"
)

type propositionSqlManager struct{}

func Proposition() *propositionSqlManager {
	return &propositionSqlManager{}
}

func (propositionSqlManager) Insert() string {
	return `INSERT INTO proposition(code, original_text_url, original_text_mime_type, title, content, submitted_at,
                        image_url, image_description, specific_type, proposition_type_id, article_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id`
}

type propositionSelectSqlManager struct{}

func (propositionSqlManager) Select() *propositionSelectSqlManager {
	return &propositionSelectSqlManager{}
}

func (propositionSelectSqlManager) ByCodes(numberOfPropositions int) string {
	var parameters []string
	for i := 1; i <= numberOfPropositions; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS proposition_id, code AS proposition_code,
				original_text_url AS proposition_original_text_url,
				original_text_mime_type AS proposition_original_text_mime_type, title AS proposition_title,
				content AS proposition_content, submitted_at AS proposition_submitted_at
			FROM proposition
			WHERE active = true AND code IN (%s)`, strings.Join(parameters, ","))
}
