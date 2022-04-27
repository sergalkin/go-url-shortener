package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
)

var _ DB = (*db)(nil)

type db struct {
	conn *pgx.Conn
}

type DB interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

func NewDBConnection() (*db, error) {
	var database = &db{conn: nil}

	if len(config.DatabaseDSN()) > 0 {
		conn, err := pgx.Connect(context.Background(), config.DatabaseDSN())
		if err != nil {
			return nil, err
		}

		database.conn = conn
	}

	return database, nil
}

func (d *db) Ping(ctx context.Context) error {
	if d.conn == nil {
		return errors.New("error in connection to db")
	}
	return d.conn.Ping(ctx)
}

func (d *db) Close(ctx context.Context) error {
	if d.conn == nil {
		return errors.New("error in connection to db")
	}
	return d.conn.Close(ctx)
}
