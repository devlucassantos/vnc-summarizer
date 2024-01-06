package queries

type newsSqlManager struct{}

func News() *newsSqlManager {
	return &newsSqlManager{}
}

func (newsSqlManager) InsertProposition() string {
	return `INSERT INTO news(proposition_id) VALUES ($1) RETURNING id`
}

func (newsSqlManager) InsertNewsletter() string {
	return `INSERT INTO news(newsletter_id) VALUES ($1) RETURNING id`
}
