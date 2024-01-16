package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
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

func (instance Newsletter) CreateNewsletter(newsletter newsletter.Newsletter) (*uuid.UUID, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	formattedReferenceDate := newsletter.ReferenceDate().Format("02/01/2006")

	transaction, err := postgresConnection.Beginx()
	if err != nil {
		log.Errorf("Erro ao iniciar transação para o cadastro do boletim do dia %d: %s", formattedReferenceDate,
			err.Error())
		return nil, err
	}
	defer instance.connectionManager.rollbackTransaction(transaction)

	var newsletterId uuid.UUID
	err = transaction.QueryRow(queries.Newsletter().Insert(), newsletter.Title(), newsletter.Content(),
		newsletter.ReferenceDate()).Scan(&newsletterId)
	if err != nil {
		log.Errorf("Erro ao cadastrar o boletim do dia %s: %s", formattedReferenceDate, err.Error())
		return nil, err
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
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		log.Error("Erro ao confirmar transação para o cadastro do boletim do dia %d: %s", formattedReferenceDate,
			err.Error())
		return nil, err
	}

	log.Infof("Boletim do dia %s registrado com sucesso com o ID %s (ID da Matéria: %s)",
		formattedReferenceDate, newsletterId, newsId)
	return &newsletterId, nil
}
