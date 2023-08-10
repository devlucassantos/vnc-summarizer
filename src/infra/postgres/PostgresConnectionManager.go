package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type ConnectionManagerInterface interface {
	createConnection() (*sqlx.DB, error)
	endConnection(*sqlx.DB)
	rollbackTransaction(*sqlx.Tx)
}

type ConnectionManager struct{}

func NewPostgresConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

func (ConnectionManager) createConnection() (*sqlx.DB, error) {
	uri, err := getPostgresConnectionURI()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		log.Error("Erro ao criar conexão com o banco de dados: ", err.Error())
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10000)
	db.SetMaxIdleConns(10000)

	return db, nil
}

func (ConnectionManager) endConnection(connection *sqlx.DB) {
	err := connection.Close()
	if err != nil {
		log.Error("Erro ao encerrar conexão com o banco de dados: ", err.Error())
	}
}

func (ConnectionManager) rollbackTransaction(transaction *sqlx.Tx) {
	err := transaction.Rollback()
	if err != nil {
		log.Warn("Erro ao cancelar transação: ", err.Error())
	}
}

func getPostgresConnectionURI() (string, error) {
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	databaseName := os.Getenv("POSTGRESQL_DB")

	connectionData := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, databaseName)
	return connectionData, nil
}
