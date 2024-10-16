package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type connectionManagerInterface interface {
	createConnection() (*sqlx.DB, error)
	closeConnection(*sqlx.DB)
	rollbackTransaction(*sqlx.Tx)
}

type ConnectionManager struct{}

func NewPostgresConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

func (ConnectionManager) createConnection() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", getPostgresConnectionUri())
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10000)
	db.SetMaxIdleConns(10000)

	return db, nil
}

func (ConnectionManager) closeConnection(connection *sqlx.DB) {
	err := connection.Close()
	if err != nil {
		log.Error("Error closing the connection to the Postgres database: ", err.Error())
	}
}

func (ConnectionManager) rollbackTransaction(transaction *sqlx.Tx) {
	err := transaction.Rollback()
	if err != nil {
		log.Warn("Error canceling the transaction in the Postgres database: ", err.Error())
	}
}

func getPostgresConnectionUri() string {
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	databaseName := os.Getenv("POSTGRESQL_DB")

	connectionData := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, databaseName)
	return connectionData
}
