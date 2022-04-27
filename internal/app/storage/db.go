package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/migrations"
)

var _ DB = (*db)(nil)

type db struct {
	conn *pgx.Conn
}

type linkRow struct {
	ID        int64
	URLHash   string
	URL       string
	UUIDHASH  string
	CREATEDAT time.Time
}

type DB interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Store(key string, url string)
	Get(key string) (string, bool)
	LinksByUUID(uuid string) ([]UserURLs, bool)
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

	if database.conn != nil {
		_, err := migrations.Up()
		if err != nil {
			return database, err
		}
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

func (d *db) Store(key string, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := fmt.Sprintf("insert into links (url_hash, url, uuid_hash) values ('%s', '%s', '%s');", key, url, middleware.GetUUID())

	if _, err := d.conn.Exec(ctx, q); err != nil {
		fmt.Println(err)
	}
}

func (d *db) Get(key string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var url string

	q := fmt.Sprintf("select url from links where url_hash = '%s'", key)
	if err := d.conn.QueryRow(ctx, q).Scan(&url); err != nil {
		return "", false
	}

	return url, true
}

func (d *db) LinksByUUID(uuid string) ([]UserURLs, bool) {
	var userUrls []UserURLs

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := fmt.Sprintf("select url_hash, url from links where uuid_hash = '%s'", uuid)
	rows, err := d.conn.Query(ctx, q)
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return userUrls, false
	}

	for rows.Next() {
		var r linkRow
		if err := rows.Scan(&r.URLHash, &r.URL); err != nil {
			fmt.Println(err)
			return userUrls, false
		}

		userUrls = append(userUrls, UserURLs{
			ShortURL:    r.URLHash,
			OriginalURL: r.URL,
		})
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return userUrls, false
	}

	if len(userUrls) == 0 {
		return userUrls, false
	}

	return userUrls, true
}
