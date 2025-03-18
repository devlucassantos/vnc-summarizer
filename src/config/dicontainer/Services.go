package dicontainer

import (
	interfaces "vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/core/services"
)

func GetDeputyService() interfaces.Deputy {
	return services.NewDeputyService(GetDeputyPostgresRepository(), GetPartyPostgresRepository())
}

func GetExternalAuthorService() interfaces.ExternalAuthor {
	return services.NewExternalAuthorService(GetExternalAuthorPostgresRepository(),
		GetExternalAuthorTypePostgresRepository())
}

func GetAuthorService() interfaces.Author {
	return services.NewAuthorService(GetDeputyService(), GetExternalAuthorService())
}

func GetPropositionService() interfaces.Proposition {
	return services.NewPropositionService(GetPropositionPostgresRepository(), GetPropositionTypePostgresRepository(),
		GetArticleTypePostgresRepository(), GetAuthorService())
}

func GetLegislativeBodyService() interfaces.LegislativeBody {
	return services.NewLegislativeBodyService(GetLegislativeBodyPostgresRepository(),
		GetLegislativeBodyTypePostgresRepository())
}

func GetVotingService() interfaces.Voting {
	return services.NewVotingService(GetVotingPostgresRepository(), GetArticleTypePostgresRepository(),
		GetLegislativeBodyService(), GetPropositionService())
}

func GetEventService() interfaces.Event {
	return services.NewEventService(GetEventPostgresRepository(), GetArticleTypePostgresRepository(),
		GetEventTypePostgresRepository(), GetEventSituationPostgresRepository(), GetAgendaItemRegimeRepository(),
		GetDeputyService(), GetLegislativeBodyService(), GetPropositionService(), GetVotingService())
}

func GetNewsletterService() interfaces.Newsletter {
	return services.NewNewsletterService(GetNewsletterPostgresRepository(), GetArticleTypePostgresRepository(),
		GetArticlePostgresRepository())
}
