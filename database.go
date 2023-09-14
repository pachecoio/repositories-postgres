package main

import "github.com/jmoiron/sqlx"

type PostgresDB struct {
	DB *sqlx.DB
}

func NewPostgresDB(db *sqlx.DB) (*PostgresDB, error) {
	return &PostgresDB{DB: db}, nil
}

func (db *PostgresDB) Disconnect() error {
	return db.DB.Close()
}
