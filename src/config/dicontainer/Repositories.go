package dicontainer

import (
	"vnc-summarizer/core/interfaces/repositories"
	"vnc-summarizer/infra/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetExternalAuthorTypePostgresRepository() repositories.ExternalAuthorType {
	return postgres.NewExternalAuthorTypeRepository(GetPostgresDatabaseManager())
}

func GetExternalAuthorPostgresRepository() repositories.ExternalAuthor {
	return postgres.NewExternalAuthorRepository(GetPostgresDatabaseManager())
}

func GetDeputyPostgresRepository() repositories.Deputy {
	return postgres.NewDeputyRepository(GetPostgresDatabaseManager())
}

func GetPartyPostgresRepository() repositories.Party {
	return postgres.NewPartyRepository(GetPostgresDatabaseManager())
}

func GetArticleTypePostgresRepository() repositories.ArticleType {
	return postgres.NewArticleTypeRepository(GetPostgresDatabaseManager())
}

func GetArticlePostgresRepository() repositories.Article {
	return postgres.NewArticleRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() repositories.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}

func GetPropositionTypePostgresRepository() repositories.PropositionType {
	return postgres.NewPropositionTypeRepository(GetPostgresDatabaseManager())
}

func GetLegislativeBodyTypePostgresRepository() repositories.LegislativeBodyType {
	return postgres.NewLegislativeBodyTypeRepository(GetPostgresDatabaseManager())
}

func GetLegislativeBodyPostgresRepository() repositories.LegislativeBody {
	return postgres.NewLegislativeBodyRepository(GetPostgresDatabaseManager())
}

func GetVotingPostgresRepository() repositories.Voting {
	return postgres.NewVotingRepository(GetPostgresDatabaseManager())
}

func GetEventTypePostgresRepository() repositories.EventType {
	return postgres.NewEventTypeRepository(GetPostgresDatabaseManager())
}

func GetEventSituationPostgresRepository() repositories.EventSituation {
	return postgres.NewEventSituationRepository(GetPostgresDatabaseManager())
}

func GetAgendaItemRegimeRepository() repositories.AgendaItemRegime {
	return postgres.NewAgendaItemRegimeRepository(GetPostgresDatabaseManager())
}

func GetEventPostgresRepository() repositories.Event {
	return postgres.NewEventRepository(GetPostgresDatabaseManager())
}

func GetNewsletterPostgresRepository() repositories.Newsletter {
	return postgres.NewNewsletterRepository(GetPostgresDatabaseManager())
}
