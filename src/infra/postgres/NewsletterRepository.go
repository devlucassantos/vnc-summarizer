package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-summarizer/infra/dto"
	"vnc-summarizer/infra/postgres/queries"
)

type Newsletter struct {
	connectionManager connectionManagerInterface
}

func NewNewsletterRepository(connectionManager connectionManagerInterface) *Newsletter {
	return &Newsletter{
		connectionManager: connectionManager,
	}
}

func (instance Newsletter) CreateNewsletter(newsletter newsletter.Newsletter, propositions []proposition.Proposition) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro do boletim do dia %d: %s", formattedReferenceDate,
			err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var newsletterId uuid.UUID
	err = transaction.QueryRow(queries.Newsletter().Insert(), newsletter.ReferenceDate(),
		newsletter.Title(), newsletter.Description()).Scan(&newsletterId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o boletim do dia %s: %s", formattedReferenceDate, err.Error())
		return err
	}

	for _, propositionData := range propositions {
		_, err = transaction.Exec(queries.NewsletterProposition().Insert(), newsletterId, propositionData.Id())
		if err != nil {
			log.Errorf("Erro ao cadastrar proposição %d como parte integrante do boletim %s do dia %s: %s",
				newsletterId, propositionData.Code(), formattedReferenceDate, err.Error())
			continue
		}

		log.Infof("Proposição %d cadastrada como parte integrante do boletim do dia %s", propositionData.Code(),
			formattedReferenceDate)
	}

	var articleId uuid.UUID
	newsletterArticle := newsletter.Article()
	articleType := newsletterArticle.Type()
	err = transaction.QueryRow(queries.Article().Insert().Newsletter(), newsletterId, articleType.Id()).Scan(&articleId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o boletim do dia %s como matéria: %s", formattedReferenceDate, err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Erro ao confirmar transação para o cadastro do boletim do dia %d: %s", formattedReferenceDate,
			err.Error())
		return err
	}

	log.Infof("Boletim do dia %s registrado com sucesso com o ID %s (ID da Matéria: %s)",
		formattedReferenceDate, newsletterId, articleId)
	return nil
}

func (instance Newsletter) UpdateNewsletter(newsletter newsletter.Newsletter, newPropositions []proposition.Proposition) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para a atualização do boletim %s do dia %d: %s", newsletter.Id(),
			formattedReferenceDate, err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	_, err = transaction.Exec(queries.Newsletter().Update(), newsletter.Description(), newsletter.Id())
	if err != nil {
		log.Errorf("Erro ao atualizar o boletim %s do dia %s: %s", newsletter.Id(), formattedReferenceDate,
			err.Error())
		return err
	}

	for _, propositionData := range newPropositions {
		_, err = transaction.Exec(queries.NewsletterProposition().Insert(), newsletter.Id(), propositionData.Id())
		if err != nil {
			log.Errorf("Erro ao cadastrar proposição %d como parte integrante do boletim %s do dia %s: %s",
				propositionData.Code(), newsletter.Id(), formattedReferenceDate, err.Error())
			continue
		}

		log.Infof("Proposição %d cadastrada como parte integrante do boletim %s do dia %s", propositionData.Code(),
			newsletter.Id(), formattedReferenceDate)
	}

	_, err = transaction.Exec(queries.Article().Update().NewsletterReferenceDateTime(), newsletter.Id())
	if err != nil {
		log.Errorf("Erro ao atualizar data e hora de referência da matéria do boletim %s do dia %s: %s",
			newsletter.Id(), formattedReferenceDate, err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Errorf("Erro ao confirmar transação para a atualização do boletim %s do dia %s: %s", newsletter.Id(),
			formattedReferenceDate, err.Error())
		return err
	}

	log.Infof("Boletim %s do dia %s atualizado com sucesso", newsletter.Id(), formattedReferenceDate)
	return nil
}

func (instance Newsletter) GetNewsletterByReferenceDate(referenceDate time.Time) (*newsletter.Newsletter, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterData dto.Newsletter
	err = postgresConnection.Get(&newsletterData, queries.Newsletter().Select().ByReferenceDate(), referenceDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Errorf("Erro ao obter os dados do boletim do dia %s no banco de dados: %s",
			referenceDate.Format("02/01/2006"), err.Error())
		return nil, err
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsletterData.Id).
		ReferenceDate(newsletterData.ReferenceDate).
		Title(newsletterData.Title).
		Description(newsletterData.Description).
		CreatedAt(newsletterData.CreatedAt).
		UpdatedAt(newsletterData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do boletim %s: %s", newsletterData.Id,
			err.Error())
		return nil, err
	}

	return newsletterDomain, nil
}
