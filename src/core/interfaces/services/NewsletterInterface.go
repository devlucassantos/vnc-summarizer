package services

import "time"

type Newsletter interface {
	RegisterNewNewsletter(referenceDate time.Time)
}
