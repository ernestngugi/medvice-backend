package db

import "database/sql"

type TestDB struct {
	*sql.Tx
	valid bool
}

func (db *TestDB) Begin() (*sql.Tx, error) {
	return db.Tx, nil
}

func (db *TestDB) Close() error {
	return nil
}

func (db *TestDB) Ping() error {
	return nil
}

func (db *TestDB) Valid() bool {
	return db.valid
}

func NewTestDB(tx *sql.Tx) *TestDB {
	return &TestDB{
		Tx:    tx,
		valid: tx != nil,
	}
}
