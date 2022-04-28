package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/migrations"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

var _ DB = (*db)(nil)

type db struct {
	conn *pgx.Conn
}

type linkRow struct {
	ID            int64
	URLHash       string
	URL           string
	UID           uuid.UUID
	CreatedAt     time.Time
	CorrelationID string
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchLink struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type DB interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Store(key *string, url string)
	Get(key string) (string, bool)
	LinksByUUID(uuid string) ([]UserURLs, bool)
	BatchInsert([]BatchRequest) ([]BatchLink, error)
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

func (d *db) Store(key *string, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uid string
	err := utils.Decode(middleware.GetUUID(), &uid)
	if err != nil {
		fmt.Println(err)
	}

	q := fmt.Sprintf(
		"insert into links (url_hash, url, uid) values ('%s', '%s', '%s') "+
			"ON CONFLICT ON CONSTRAINT links_url_key DO NOTHING", *key, url, uid,
	)

	r, err := d.conn.Exec(ctx, q)
	if err != nil {
		fmt.Println(err)
	}

	if r.RowsAffected() == 0 {
		q = "select url_hash from links where url = $1"
		row := d.conn.QueryRow(ctx, q, url)

		var tempKey *string
		err = row.Scan(&tempKey)

		*key = *tempKey
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

	q := fmt.Sprintf("select url_hash, url from links where uid = '%s'", uuid)
	rows, err := d.conn.Query(ctx, q)
	if err != nil {
		fmt.Println(err)
		return userUrls, false
	}
	defer rows.Close()

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

func (d *db) BatchInsert(br []BatchRequest) ([]BatchLink, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback(ctx)

	seqGenerator := utils.NewSequence()
	batchLinks := make([]BatchLink, 0)

	var uid string
	err = utils.Decode(middleware.GetUUID(), &uid)
	if err != nil {
		return []BatchLink{}, err
	}

	q := "insert into links(url_hash, url, uid, correlation_id) values ($1, $2, $3, $4)"
	for _, val := range br {
		urlHash, err := seqGenerator.Generate(5)
		if err != nil {
			return []BatchLink{}, err
		}

		_, err = tx.Exec(ctx, q, urlHash, val.OriginalURL, uid, val.CorrelationID)
		if err != nil {
			return []BatchLink{}, err
		}
		batchLinks = append(batchLinks, BatchLink{
			CorrelationID: val.CorrelationID,
			ShortURL:      config.BaseURL() + "/" + urlHash,
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []BatchLink{}, err
	}

	return batchLinks, nil
}
