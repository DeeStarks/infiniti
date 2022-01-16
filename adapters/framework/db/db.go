package db

import (
	"log"
	"database/sql"
)

type DBAdapter struct {
	db		*sql.DB
}

func NewDBAdapter(driver, source string) (*DBAdapter, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Printf("Connected to %s\n", source)
	return &DBAdapter{db: db}, nil
}

func (adapter *DBAdapter) CloseDBConnection() error {
	return adapter.db.Close()
}