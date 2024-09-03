package diconteiner

import (
	interfaces "vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/core/services"
)

func GetBackgroundDataService() interfaces.BackgroundData {
	return services.NewBackgroundDataService(GetDeputyPostgresRepository(), GetExternalAuthorPostgresRepository(),
		GetPartyPostgresRepository(), GetPropositionPostgresRepository(), GetNewsletterPostgresRepository(),
		GetArticleTypePostgresRepository())
}
