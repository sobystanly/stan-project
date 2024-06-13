package db

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"stan-project/cmd/config"
)

type (
	pgConn interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
		Close(ctx context.Context) error
	}
	db struct {
		client pgConn
	}
)

func InitDB(ctx context.Context) (*db, error) {
	connConfig := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", config.Global.PostgresUsername, config.Global.PostgresPassword, config.Global.PostgresAddress, config.Global.PostgresDatabase)
	conn, err := pgx.Connect(ctx, connConfig)
	if err != nil {
		return nil, err
	}
	return &db{client: conn}, nil
}

func (db *db) Close(ctx context.Context) error {
	return db.client.Close(ctx)
}

//go:embed sql/create_risk_table.sql
var createRisksTable string

func (db *db) RunMigrations(ctx context.Context) error {
	_, err := db.client.Exec(ctx, createRisksTable)
	if err != nil {
		return err
	}
	return nil
}
