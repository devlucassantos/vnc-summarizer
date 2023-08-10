package diconteiner

import (
	"vnc-write-api/core/interfaces/repositories"
	"vnc-write-api/infra/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetDeputyPostgresRepository() repositories.Deputy {
	return postgres.NewDeputyRepository(GetPostgresDatabaseManager())
}

func GetKeywordPostgresRepository() repositories.Keyword {
	return postgres.NewKeywordRepository(GetPostgresDatabaseManager())
}

func GetOrganizationPostgresRepository() repositories.Organization {
	return postgres.NewOrganizationRepository(GetPostgresDatabaseManager())
}

func GetPartyPostgresRepository() repositories.Party {
	return postgres.NewPartyRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() repositories.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}
