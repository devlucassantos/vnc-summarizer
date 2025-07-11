package dto

import (
	"github.com/google/uuid"
	"time"
)

type Voting struct {
	Id                uuid.UUID `db:"voting_id"`
	Code              string    `db:"voting_code"`
	Description       string    `db:"voting_description"`
	Result            string    `db:"voting_result"`
	ResultAnnouncedAt time.Time `db:"voting_result_announced_at"`
	IsApproved        *bool     `db:"voting_is_approved"`
}
