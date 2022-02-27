package db

import (
	"log"
	"database/sql"
    _ "github.com/lib/pq"
)

// DB Schema Link: https://dbdiagram.io/d/61cc73453205b45b73d0fdfe

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