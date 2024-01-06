package queries

type newsletterSqlManager struct{}

func Newsletter() *newsletterSqlManager {
	return &newsletterSqlManager{}
}

func (newsletterSqlManager) Insert() string {
	return `INSERT INTO newsletter(title, content, reference_date) VALUES ($1, $2, $3) RETURNING id`
}
