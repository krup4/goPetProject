package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB(ctx context.Context) error {
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) NOT NULL,
		password BYTEA NOT NULL,
		name VARCHAR(255)
	);
	`

	_, err := Pool.Exec(ctx, createTablesSQL)
	return err
}

func ConnectToDB(connString string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	log.Println("Connected to PostgreSQL!")
	return nil
}

func Close() {
	Pool.Close()
}
