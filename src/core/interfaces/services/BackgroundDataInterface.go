package services

import "time"

type BackgroundData interface {
	RegisterNewNewsletter(referenceDate time.Time)
	RegisterNewPropositions()
}
