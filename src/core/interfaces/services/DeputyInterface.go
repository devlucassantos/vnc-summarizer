package services

import "github.com/devlucassantos/vnc-domains/src/domains/deputy"

type Deputy interface {
	GetDeputyFromDeputyData(deputyData map[string]interface{}) (*deputy.Deputy, error)
}
