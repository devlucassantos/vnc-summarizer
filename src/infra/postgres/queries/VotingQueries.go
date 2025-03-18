package queries

import (
	"fmt"
	"strings"
)

type votingSqlManager struct{}

func Voting() *votingSqlManager {
	return &votingSqlManager{}
}

func (votingSqlManager) Insert() string {
	return `INSERT INTO voting(code, result, result_announced_at, is_approved, legislative_body_id, main_proposition_id,
            	article_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`
}

type votingSelectSqlManager struct{}

func (votingSqlManager) Select() *votingSelectSqlManager {
	return &votingSelectSqlManager{}
}

func (votingSelectSqlManager) ByCode() string {
	return `SELECT id AS voting_id, code AS voting_code, result AS voting_result,
       			result_announced_at AS voting_result_announced_at, is_approved AS voting_is_approved
			FROM voting
			WHERE active = true AND code = $1`
}

func (votingSelectSqlManager) ByCodes(numberOfVotes int) string {
	var parameters []string
	for i := 1; i <= numberOfVotes; i++ {
		parameters = append(parameters, fmt.Sprintf("$%d", i))
	}

	return fmt.Sprintf(`SELECT id AS voting_id, code AS voting_code, result AS voting_result,
				result_announced_at AS voting_result_announced_at, is_approved AS voting_is_approved
			FROM voting
			WHERE active = true AND code IN (%s)`, strings.Join(parameters, ","))
}
