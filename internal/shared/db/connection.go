package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Db struct {
    DB *sql.DB
}

func NewDatabase(connStr string) (*Db, error) {
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        return nil, fmt.Errorf("error opening database connection: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    return &Db{
        DB: db,
    }, nil
}

func (db *Db) CloseConnection() {
    db.DB.Close()
}