package queries

type newsletterArticleSqlManager struct{}

func NewsletterArticle() *newsletterArticleSqlManager {
	return &newsletterArticleSqlManager{}
}

func (newsletterArticleSqlManager) Insert() string {
	return `INSERT INTO newsletter_article(newsletter_id, article_id)
			VALUES ($1, $2)`
}

type newsletterArticleSelectSqlManager struct{}

func (newsletterArticleSqlManager) Select() *newsletterArticleSelectSqlManager {
	return &newsletterArticleSelectSqlManager{}
}

func (newsletterArticleSelectSqlManager) ByNewsletterId() string {
	return `SELECT article.id AS article_id,
				article_type.id AS article_type_id, article_type.description AS article_type_description,
				article_type.codes AS article_type_codes, article_type.color AS article_type_color,
				COALESCE(proposition.title, '') AS proposition_title,
				COALESCE(proposition.content, '') AS proposition_content,
				COALESCE(voting.code, '') AS voting_code, COALESCE(voting.result, '') AS voting_result,
				COALESCE(event.title, '') AS event_title, COALESCE(event.description, '') AS event_description
			FROM newsletter_article
				INNER JOIN article ON article.id = newsletter_article.article_id
				INNER JOIN article_type ON article_type.id = article.article_type_id
				LEFT JOIN proposition ON proposition.article_id = article.id
				LEFT JOIN voting ON voting.article_id = article.id
				LEFT JOIN event ON event.article_id = article.id
			WHERE newsletter_article.active = true AND article.active = true AND article_type.active = true AND
				proposition.active IS NOT false AND voting.active IS NOT false AND event.active IS NOT false AND
				newsletter_article.newsletter_id = $1
			ORDER BY article.reference_date_time`
}
