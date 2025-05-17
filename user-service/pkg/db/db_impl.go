package db

import (
	"database/sql"
	"fmt"
	"user-service/internal/config"

	_ "github.com/lib/pq"
)

type dbImpl struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	sslmode  string
}

func (db dbImpl) OpenConnect() (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.host,
		db.port,
		db.user,
		db.password,
		db.dbname,
		db.sslmode,
	)
	conn, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func New(dbConfig config.DatabaseConfig) (DB, error) {
	return &dbImpl{
			host:     dbConfig.GetHost(),
			port:     dbConfig.GetPort(),
			user:     dbConfig.GetUser(),
			password: dbConfig.GetPassword(),
			dbname:   dbConfig.GetDBName(),
			sslmode:  dbConfig.GetSSLMode(),
		},
		nil
}
