package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
	"vnc-write-api/infra/dto"
	"vnc-write-api/infra/postgres/queries"
)

type Newsletter struct {
	connectionManager ConnectionManagerInterface
}

func NewNewsletterRepository(connectionManager ConnectionManagerInterface) *Newsletter {
	return &Newsletter{
		connectionManager: connectionManager,
	}
}

func (instance Newsletter) CreateNewsletter(newsletter newsletter.Newsletter) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro do boletim do dia %d: %s", formattedReferenceDate,
			err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var newsletterId uuid.UUID
	err = transaction.QueryRow(queries.Newsletter().Insert(), newsletter.Title(), newsletter.Content(),
		newsletter.ReferenceDate()).Scan(&newsletterId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o boletim do dia %s: %s", formattedReferenceDate, err.Error())
		return err
	}

	for _, propositionData := range newsletter.Propositions() {
		_, err = transaction.Exec(queries.NewsletterProposition().Insert(), newsletterId, propositionData.Id())
		if err != nil {
			log.Errorf("Erro ao cadastrar proposição %d como parte integrante do boletim do dia  %s: %s",
				propositionData.Code(), formattedReferenceDate, err.Error())
			continue
		}

		log.Infof("Proposição %d cadastrada como parte integrante do boletim do dia %s", propositionData.Code(),
			formattedReferenceDate)
	}

	var newsId uuid.UUID
	err = transaction.QueryRow(queries.News().InsertNewsletter(), newsletterId).Scan(&newsId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o boletim do dia %s como matéria: %s", formattedReferenceDate, err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Error("Erro ao confirmar transação para o cadastro do boletim do dia %d: %s", formattedReferenceDate,
			err.Error())
		return err
	}

	log.Infof("Boletim do dia %s registrado com sucesso com o ID %s (ID da Matéria: %s)",
		formattedReferenceDate, newsletterId, newsId)
	return nil
}

func (instance Newsletter) UpdateNewsletter(newsletter newsletter.Newsletter, newPropositions []proposition.Proposition) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para a atualização do boletim %s do dia %d: %s", newsletter.Id(),
			formattedReferenceDate, err.Error())
		return err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	_, err = transaction.Exec(queries.Newsletter().Update(), newsletter.Title(), newsletter.Content(), newsletter.Id())
	if err != nil {
		log.Errorf("Erro ao atualizar o boletim %s do dia %s: %s", newsletter.Id(), formattedReferenceDate,
			err.Error())
		return err
	}

	for _, propositionData := range newPropositions {
		_, err = transaction.Exec(queries.NewsletterProposition().Insert(), newsletter.Id(), propositionData.Id())
		if err != nil {
			log.Errorf("Erro ao cadastrar proposição %d como parte integrante do boletim %s do dia  %s: %s",
				newsletter.Id(), propositionData.Code(), formattedReferenceDate, err.Error())
			continue
		}

		log.Infof("Proposição %d cadastrada como parte integrante do boletim %s do dia %s", propositionData.Code(),
			newsletter.Id(), formattedReferenceDate)
	}

	_, err = transaction.Exec(queries.News().UpdateNewsletterReferenceDateTime(), newsletter.Id())
	if err != nil {
		log.Errorf("Erro ao atualizar data e hora de referência da matéria do boletim %s do dia %s: %s",
			newsletter.Id(), formattedReferenceDate, err.Error())
		return err
	}

	err = transaction.Commit()
	if err != nil {
		log.Error("Erro ao confirmar transação para a atualização do boletim %s do dia %s: %s", newsletter.Id(),
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
	defer instance.connectionManager.endConnection(postgresConnection)

	var newsData dto.News
	err = postgresConnection.Get(&newsData, queries.Newsletter().Select().ByReferenceDate(), referenceDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Errorf("Erro ao obter os dados do boletim do dia %s no banco de dados: %s",
			referenceDate.Format("02/01/2006"), err.Error())
		return nil, err
	}

	var newsletterPropositions []dto.Proposition
	err = postgresConnection.Select(&newsletterPropositions, queries.NewsletterProposition().Select().ByNewsletterId(),
		newsData.Newsletter.Id)
	if err != nil {
		log.Errorf("Erro ao obter os dados das proposições do boletim %s no banco de dados: %s",
			newsData.Newsletter.Id, err.Error())
		return nil, err
	}

	var propositions []proposition.Proposition
	for _, propositionData := range newsletterPropositions {
		propositionDomain, err := proposition.NewBuilder().
			Id(propositionData.Id).
			Code(propositionData.Code).
			OriginalTextUrl(propositionData.OriginalTextUrl).
			Title(propositionData.Title).
			Content(propositionData.Content).
			SubmittedAt(propositionData.SubmittedAt).
			Active(propositionData.Active).
			CreatedAt(propositionData.CreatedAt).
			UpdatedAt(propositionData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro durante a construção da estrutura de dados da proposição %s do boletim %s: %s",
				propositionData.Id, newsData.Newsletter.Id, err.Error())
			continue
		}

		propositions = append(propositions, *propositionDomain)
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsData.Newsletter.Id).
		Title(newsData.Newsletter.Title).
		Content(newsData.Newsletter.Content).
		ReferenceDate(newsData.Newsletter.ReferenceDate).
		Propositions(propositions).
		Active(newsData.Newsletter.Active).
		CreatedAt(newsData.Newsletter.CreatedAt).
		UpdatedAt(newsData.Newsletter.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro durante a construção da estrutura de dados do boletim %s: %s", newsData.Newsletter.Id,
			err.Error())
		return nil, err
	}

	return newsletterDomain, nil
}
