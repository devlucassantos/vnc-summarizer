package chamber

type Chamber interface {
	GetMostRecentPropositions() ([]map[string]interface{}, error)
	GetPropositionByCode(code int) (map[string]interface{}, error)
	GetPropositionContentDirectly(propositionUrl string) (string, string, error)
	GetPropositionTypes() ([]map[string]interface{}, error)
	GetPartyByAcronym(acronym string) (map[string]interface{}, error)
	GetLegislativeBodyByCode(code int) (map[string]interface{}, error)
	GetLegislativeBodyTypes() ([]map[string]interface{}, error)
	GetMostRecentVotes() ([]map[string]interface{}, error)
	GetVotingByCode(code string) (map[string]interface{}, error)
	GetMostRecentEvents() ([]map[string]interface{}, error)
	GetEventByCode(code int) (map[string]interface{}, error)
	GetEventsByCodes(eventCodes []string) ([]map[string]interface{}, error)
	GetEventTypes() ([]map[string]interface{}, error)
	GetEventSituations() ([]map[string]interface{}, error)
}
