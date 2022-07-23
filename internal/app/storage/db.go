package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/internal/app/middleware"
	"github.com/sergalkin/go-url-shortener.git/internal/app/migrations"
	"github.com/sergalkin/go-url-shortener.git/internal/app/utils"
)

var _ DB = (*db)(nil)

type db struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

type linkRow struct {
	ID            int64
	URLHash       string
	URL           string
	UID           uuid.UUID
	CreatedAt     time.Time
	CorrelationID string
	IsDeleted     bool
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchLink struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchDelete struct {
	UID string
	Arr []string
}

type DB interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
	Store(key *string, url string)
	Get(key string) (string, bool, bool)
	LinksByUUID(uuid string) ([]UserURLs, bool)
	BatchInsert([]BatchRequest) ([]BatchLink, error)
	SoftDeleteUserURLs(uuid string, ids []string) error
	DeleteThroughCh(channels ...chan BatchDelete)
}

const (
	getURLHash  = `select url_hash from links where url = $1`
	insertLinks = `insert into links (url_hash, url, uid) values ($1,$2,$3) ON CONFLICT ON CONSTRAINT links_url_key DO NOTHING`
)

func NewDBConnection(l *zap.Logger, isNeedToRunMigrations bool) (*db, error) {
	var database = &db{conn: nil, logger: l}

	if len(config.DatabaseDSN()) > 0 {
		conn, err := pgx.Connect(context.Background(), config.DatabaseDSN())
		if err != nil {
			return nil, err
		}

		database.conn = conn
	}

	if database.conn != nil && isNeedToRunMigrations {
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
		d.logger.Error(err.Error(), zap.Error(err))
	}

	r, err := d.conn.Exec(ctx, insertLinks, *key, url, uid)
	if err != nil {
		d.logger.Error(err.Error(), zap.Error(err))
	}

	if r.RowsAffected() == 0 {
		row := d.conn.QueryRow(ctx, getURLHash, url)

		var tempKey *string
		err = row.Scan(&tempKey)
		if err != nil {
			d.logger.Error(err.Error(), zap.Error(err))
		}

		*key = *tempKey
	}
}

func (d *db) Get(key string) (string, bool, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var url string
	var isDeleted bool

	q := fmt.Sprintf("select url,is_deleted from links where url_hash = '%s'", key)
	if err := d.conn.QueryRow(ctx, q).Scan(&url, &isDeleted); err != nil {
		return "", false, false
	}

	return url, true, isDeleted
}

func (d *db) LinksByUUID(uuid string) ([]UserURLs, bool) {
	var userUrls []UserURLs

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := fmt.Sprintf("select url_hash, url from links where uid = '%s'", uuid)
	rows, err := d.conn.Query(ctx, q)
	if err != nil {
		d.logger.Error(err.Error(), zap.Error(err))
		return userUrls, false
	}
	defer rows.Close()

	for rows.Next() {
		var r linkRow
		if err := rows.Scan(&r.URLHash, &r.URL); err != nil {
			d.logger.Error(err.Error(), zap.Error(err))
			return userUrls, false
		}

		userUrls = append(userUrls, UserURLs{
			ShortURL:    r.URLHash,
			OriginalURL: r.URL,
		})
	}

	err = rows.Err()
	if err != nil {
		d.logger.Error(err.Error(), zap.Error(err))
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

func (d *db) DeleteThroughCh(channels ...chan BatchDelete) {
	out := fanIn(channels...)

	for c := range out {
		go func() {
			err := d.SoftDeleteUserURLs(c.UID, c.Arr)
			if err != nil {
				d.logger.Error(err.Error(), zap.Error(err))
			}
		}()
		close(out)
	}
}

func fanIn(channels ...chan BatchDelete) chan BatchDelete {
	outCh := make(chan BatchDelete)

	go func() {
		wg := &sync.WaitGroup{}

		for _, ch := range channels {
			wg.Add(1)

			go func(ch chan BatchDelete) {
				defer wg.Done()
				for i := range ch {
					outCh <- i
				}
			}(ch)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func (d *db) SoftDeleteUserURLs(uuid string, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		d.logger.Error(err.Error(), zap.Error(err))
	}
	defer tx.Rollback(ctx)

	stmt := fmt.Sprintf("update links set is_deleted = true where uid = '%s' and url_hash in ('%s')",
		uuid,
		strings.Join(ids, "','"))

	_, err = tx.Exec(ctx, stmt)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
