package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type SQLOperations interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type DB interface {
	SQLOperations
	Close() error
	Ping() error
}

type RowScanner interface {
	Scan(dest ...any) error
}

type AppDB struct {
	*sql.DB
}

func InitDB() DB {
	return initDBWithURL(os.Getenv("DATABASE_URL"))
}

func initDBWithURL(databaseURL string) DB {

	if databaseURL == "" {
		log.Fatal("database url is empty")
	}

	appDB, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("sql open error %v", err)
	}

	db := &AppDB{
		DB: appDB,
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("database ping error %v", err)
	}

	return db
}
