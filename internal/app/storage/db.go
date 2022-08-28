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
	"github.com/sergalkin/go-url-shortener.git/pkg/sequence"
)

var _ DB = (*db)(nil)

// db - representation of *pgx.Conn and *zap.Logger
type db struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

// linkRow - a representation of link in DB.
type linkRow struct {
	CreatedAt     time.Time
	URLHash       string
	URL           string
	CorrelationID string
	IsDeleted     bool
	ID            int64
	UID           uuid.UUID
}

// BatchRequest - a representation of mass assignment URL request.
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchLink - a representation of returned values of mass assignment URL request.
type BatchLink struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchDelete struct {
	UID string
	Arr []string
}

type DB interface {
	// Ping - checks for connection. If no error returned Ping is considered successful.
	Ping(ctx context.Context) error
	// Close - closes connection.
	Close(ctx context.Context) error
	// Store - stores given url into database
	Store(key *string, url string, uid string)
	// Get - trying to retrieve a URL from database by provided key.
	// Get - returns URL, bool as status of retrieval, bool as status was URL deleted or is it still present.
	Get(key string) (string, bool, bool)
	// LinksByUUID - trying to retrieve slice of UserURLs. On successful retrieval returns true as bool value
	// and false of failure.
	LinksByUUID(uuid string) ([]UserURLs, bool)
	// BatchInsert - mass insert provided links into database
	BatchInsert([]BatchRequest) ([]BatchLink, error)
	// SoftDeleteUserURLs - marks provided links as deleted. uuid - is user unique id, ids - is slice of links that
	// needs to be marked as soft deleted.
	SoftDeleteUserURLs(uuid string, ids []string) error
	DeleteThroughCh(channels ...chan BatchDelete)
	Stats() (int, int, error)

	HasNotNilConn() bool
}

const (
	getURLHash  = `select url_hash from links where url = $1`
	insertLinks = `insert into links (url_hash, url, uid) values ($1,$2,$3) ON CONFLICT ON CONSTRAINT links_url_key DO NOTHING`
	stats       = `select count(id) as links,  count(Distinct uid) as url from links where is_deleted = false`
)

// NewDBConnection - creates new database connection and attempts to run migrations.
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

// Ping - ping database to check for its availability.
func (d *db) Ping(ctx context.Context) error {
	if d.conn == nil {
		return errors.New("error in connection to db")
	}
	return d.conn.Ping(ctx)
}

// Close - closes connection with database.
func (d *db) Close(ctx context.Context) error {
	if d.conn != nil {
		err := d.conn.Close(ctx)
		if err != nil {
			return err
		}
		d.conn = nil
	}

	return nil
}

// Store - stores provided url by key in database
func (d *db) Store(key *string, url string, uid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

// Get - attempt to get url from database by its key
// returns url, bool representation of was url found, bool representation of was url soft deleted.
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

// LinksByUUID - get slice of UserURLs from database by provided uuid and bool representation of was url found or not.
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
		if scanErr := rows.Scan(&r.URLHash, &r.URL); scanErr != nil {
			d.logger.Error(scanErr.Error(), zap.Error(scanErr))
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

// BatchInsert - batch insert links to database with CorrelationID
// additionally adds uuid to uid column in database gotten form uid cookie.
func (d *db) BatchInsert(br []BatchRequest) ([]BatchLink, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback(ctx)

	seqGenerator := sequence.NewSequence()
	batchLinks := make([]BatchLink, 0)

	var uid string
	err = utils.Decode(middleware.GetUUID(), &uid)
	if err != nil {
		return []BatchLink{}, err
	}

	q := "insert into links(url_hash, url, uid, correlation_id) values ($1, $2, $3, $4)"
	for _, val := range br {
		urlHash, errGen := seqGenerator.Generate(5)
		if errGen != nil {
			return []BatchLink{}, errGen
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

// DeleteThroughCh - Soft deletes URL using channels ang go routine.
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

// SoftDeleteUserURLs - marks URL as deleted in database.
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

func (d *db) HasNotNilConn() bool {
	return d.conn != nil
}

// Stats - returns count of urls and users stored in DB. Counts only non-soft deleted records.
func (d *db) Stats() (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var url, users int
	if err := d.conn.QueryRow(ctx, stats).Scan(&url, &users); err != nil {
		return 0, 0, err
	}

	return url, users, nil
}
