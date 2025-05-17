package db

import "database/sql"

type DB interface {
	OpenConnect() (*sql.DB, error)
}
