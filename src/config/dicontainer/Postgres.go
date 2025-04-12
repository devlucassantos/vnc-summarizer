package dicontainer

import (
	"vnc-summarizer/adapters/databases/postgres"
	interfaces "vnc-summarizer/core/interfaces/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetExternalAuthorTypePostgresRepository() interfaces.ExternalAuthorType {
	return postgres.NewExternalAuthorTypeRepository(GetPostgresDatabaseManager())
}

func GetExternalAuthorPostgresRepository() interfaces.ExternalAuthor {
	return postgres.NewExternalAuthorRepository(GetPostgresDatabaseManager())
}

func GetDeputyPostgresRepository() interfaces.Deputy {
	return postgres.NewDeputyRepository(GetPostgresDatabaseManager())
}

func GetPartyPostgresRepository() interfaces.Party {
	return postgres.NewPartyRepository(GetPostgresDatabaseManager())
}

func GetArticleTypePostgresRepository() interfaces.ArticleType {
	return postgres.NewArticleTypeRepository(GetPostgresDatabaseManager())
}

func GetArticlePostgresRepository() interfaces.Article {
	return postgres.NewArticleRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() interfaces.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}

func GetPropositionTypePostgresRepository() interfaces.PropositionType {
	return postgres.NewPropositionTypeRepository(GetPostgresDatabaseManager())
}

func GetLegislativeBodyTypePostgresRepository() interfaces.LegislativeBodyType {
	return postgres.NewLegislativeBodyTypeRepository(GetPostgresDatabaseManager())
}

func GetLegislativeBodyPostgresRepository() interfaces.LegislativeBody {
	return postgres.NewLegislativeBodyRepository(GetPostgresDatabaseManager())
}

func GetVotingPostgresRepository() interfaces.Voting {
	return postgres.NewVotingRepository(GetPostgresDatabaseManager())
}

func GetEventTypePostgresRepository() interfaces.EventType {
	return postgres.NewEventTypeRepository(GetPostgresDatabaseManager())
}

func GetEventSituationPostgresRepository() interfaces.EventSituation {
	return postgres.NewEventSituationRepository(GetPostgresDatabaseManager())
}

func GetAgendaItemRegimeRepository() interfaces.AgendaItemRegime {
	return postgres.NewAgendaItemRegimeRepository(GetPostgresDatabaseManager())
}

func GetEventPostgresRepository() interfaces.Event {
	return postgres.NewEventRepository(GetPostgresDatabaseManager())
}

func GetNewsletterPostgresRepository() interfaces.Newsletter {
	return postgres.NewNewsletterRepository(GetPostgresDatabaseManager())
}
