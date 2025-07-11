package dicontainer

import (
	interfaces "vnc-summarizer/core/interfaces/services"
	"vnc-summarizer/core/services"
)

func GetDeputyService() interfaces.Deputy {
	return services.NewDeputyService(GetChamberApi(), GetDeputyPostgresRepository(), GetPartyPostgresRepository())
}

func GetExternalAuthorService() interfaces.ExternalAuthor {
	return services.NewExternalAuthorService(GetExternalAuthorPostgresRepository(),
		GetExternalAuthorTypePostgresRepository())
}

func GetAuthorService() interfaces.Author {
	return services.NewAuthorService(GetDeputyService(), GetExternalAuthorService())
}

func GetPropositionService() interfaces.Proposition {
	return services.NewPropositionService(GetAuthorService(), GetChamberApi(), GetChatGptApi(), GetDallEApi(),
		GetVncPdfContentExtractorApi(), GetAwsS3(), GetPropositionPostgresRepository(),
		GetPropositionTypePostgresRepository(), GetArticleTypePostgresRepository())
}

func GetLegislativeBodyService() interfaces.LegislativeBody {
	return services.NewLegislativeBodyService(GetChamberApi(), GetLegislativeBodyPostgresRepository(),
		GetLegislativeBodyTypePostgresRepository())
}

func GetVotingService() interfaces.Voting {
	return services.NewVotingService(GetChamberApi(), GetChatGptApi(), GetVotingPostgresRepository(),
		GetArticleTypePostgresRepository(), GetLegislativeBodyService(), GetPropositionService())
}

func GetEventService() interfaces.Event {
	return services.NewEventService(GetDeputyService(), GetLegislativeBodyService(), GetPropositionService(),
		GetVotingService(), GetChamberApi(), GetChatGptApi(), GetEventPostgresRepository(),
		GetArticleTypePostgresRepository(), GetEventTypePostgresRepository(), GetEventSituationPostgresRepository(),
		GetAgendaItemRegimeRepository())
}

func GetNewsletterService() interfaces.Newsletter {
	return services.NewNewsletterService(GetChatGptApi(), GetNewsletterPostgresRepository(),
		GetArticleTypePostgresRepository(), GetArticlePostgresRepository())
}
