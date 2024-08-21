package diconteiner

import (
	"vnc-summarizer/core/interfaces/repositories"
	"vnc-summarizer/infra/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetDeputyPostgresRepository() repositories.Deputy {
	return postgres.NewDeputyRepository(GetPostgresDatabaseManager())
}

func GetExternalAuthorPostgresRepository() repositories.ExternalAuthor {
	return postgres.NewExternalAuthorRepository(GetPostgresDatabaseManager())
}

func GetPartyPostgresRepository() repositories.Party {
	return postgres.NewPartyRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() repositories.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}

func GetPropositionTypePostgresRepository() repositories.PropositionType {
	return postgres.NewPropositionTypeRepository(GetPostgresDatabaseManager())
}

func GetNewsletterPostgresRepository() repositories.Newsletter {
	return postgres.NewNewsletterRepository(GetPostgresDatabaseManager())
}
