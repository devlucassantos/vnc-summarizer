package diconteiner

import (
	interfaces "vnc-write-api/core/interfaces/services"
	"vnc-write-api/core/services"
)

func GetBackgroundDataService() interfaces.BackgroundData {
	return services.NewBackgroundDataService(GetDeputyPostgresRepository(), GetOrganizationPostgresRepository(),
		GetPartyPostgresRepository(), GetPropositionPostgresRepository(), GetNewsletterPostgresRepository())
}
