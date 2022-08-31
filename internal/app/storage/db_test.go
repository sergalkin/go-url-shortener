package storage

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
	"github.com/sergalkin/go-url-shortener.git/pkg/sequence"
)

func Test_db_Ping(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Can ping DB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				errPing := conn.Ping(context.Background())

				assert.NoError(t, errPing)
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Can close DB conn without errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				errPing := conn.Close(context.Background())

				assert.NoError(t, errPing)
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func BenchmarkDb_Store(b *testing.B) {
	cfg := config.NewConfig(
		config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
	)

	conn, err := NewDBConnection(zap.NewNop(), true)

	if err == nil {
		seq := sequence.NewSequence()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// do not stop here timer cause store need generated seq of letters, and it too takes time
			key, _ := seq.Generate(4)

			conn.Store(&key, "test.com", "046cf584-df95-43fd-a2fc-f95a85c7bb95")
		}
		b.StopTimer()
	}

	cfg.DatabaseDSN = ""
}

func Test_fanIn(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"can fanIn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make([]chan BatchDelete, 1)
			assert.NotNil(t, fanIn(ch...))
		})
	}
}

func Test_db_HasNotNilConn(t *testing.T) {
	type fields struct {
		conn *pgx.Conn
	}
	tests := []struct {
		fields fields
		name   string
		want   bool
	}{
		{
			name:   "can return is connection nil or not",
			fields: fields{conn: nil},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &db{
				conn: tt.fields.conn,
			}
			assert.False(t, d.HasNotNilConn())
		})
	}
}

func Test_db_SoftDeleteUserURLs(t *testing.T) {
	type args struct {
		uuid string
		ids  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "has no errors on soft deleting URLs",
			args: args{
				uuid: "64fb79de-24cf-475a-a042-0aa582ca05bb",
				ids:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				errDelete := conn.SoftDeleteUserURLs(tt.args.uuid, tt.args.ids)

				assert.NoError(t, errDelete)
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_LinksByUUID(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{
			name: "can get links by uuid from database",
			uuid: "64fb79de-24cf-475a-a042-0aa582ca05bb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				conn.conn.Exec(context.Background(), "insert into links (url_hash, url, uid) values ('test', 'ya.ru', $1)", tt.uuid)

				urls, result := conn.LinksByUUID(tt.uuid)

				assert.NotEmpty(t, urls)
				assert.True(t, result)

				conn.conn.Exec(context.Background(), "delete from links where uid = '"+tt.uuid+"'")
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_Get(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		URLHash string
	}{
		{
			name:    "can get link by key from database",
			uuid:    "64fb79de-24cf-475a-a042-0aa582ca05bb",
			URLHash: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)

			if err == nil {
				conn.conn.Exec(context.Background(), "insert into links (url_hash, url, uid) values ($1, 'ya.ru', $2)", tt.URLHash, tt.uuid)

				url, result, isDeleted := conn.Get(tt.URLHash)

				assert.NotEmpty(t, url)
				assert.True(t, result)
				assert.False(t, isDeleted)

				conn.conn.Exec(context.Background(), "delete from links where uid = '"+tt.uuid+"'")
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_GetWillReturnDefaultValueIfNothingIsFound(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "get will return default values if nothing is found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)

			if err == nil {
				url, result, isDeleted := conn.Get("somestring")

				assert.Equal(t, "", url)
				assert.False(t, result)
				assert.False(t, isDeleted)
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_DeleteThroughCh(t *testing.T) {
	tests := []struct {
		name string
		uid  string
		urls []string
	}{
		{
			name: "can delete through chanel",
			uid:  "64fb79de-24cf-475a-a042-0aa582ca05bb",
			urls: []string{"test"},
		},
	}
	for _, tt := range tests {
		b := BatchDelete{
			UID: tt.uid,
			Arr: tt.urls,
		}
		inputCh := make(chan BatchDelete, 1)
		inputCh <- b

		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)

			if err == nil {
				conn.conn.Exec(context.Background(), "insert into links (url_hash, url, uid) values ('test', 'ya.ru', $1)", tt.uid)

				conn.DeleteThroughCh(inputCh)
				time.Sleep(250 * time.Millisecond)
				r := conn.conn.QueryRow(context.Background(), "select is_deleted from links where uid =$1", tt.uid)
				var isDeleted bool
				r.Scan(&isDeleted)

				assert.True(t, isDeleted)

				conn.conn.Exec(context.Background(), "delete from links where uid = '"+tt.uid+"'")
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_Store(t *testing.T) {
	tests := []struct {
		name string
		url  string
		key  string
		uid  string
	}{
		{
			name: "can store link in db",
			key:  "test",
			url:  "ya.ru",
			uid:  "046cf584-df95-43fd-a2fc-f95a85c7bb95",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)
			if err == nil {
				conn.Store(&tt.key, tt.url, tt.uid)
				assert.NotEqual(t, "", tt.key)

				conn.conn.Exec(context.Background(), "delete from links where uid = '"+tt.uid+"'")
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_StoreModifiesKeyOnDuplicate(t *testing.T) {
	tests := []struct {
		name string
		url  string
		uid  string
	}{
		{
			name: "Store modifies key on duplicate",
			url:  "ya.ru",
			uid:  "046cf584-df95-43fd-a2fc-f95a85c7bb95",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)

			if err == nil {
				conn.conn.Exec(context.Background(), "insert into links (url_hash, url, uid) values ('test', 'ya.ru', $1)", tt.uid)

				var key string
				conn.Store(&key, tt.url, tt.uid)
				assert.Equal(t, "test", key)

				conn.conn.Exec(context.Background(), "delete from links where uid = '"+tt.uid+"'")
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_Stats(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Stats can be called via db manager.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(zap.NewNop(), false)
			if err == nil {
				_, _, errStats := conn.Stats()
				assert.NoError(t, errStats)
			}

			cfg.DatabaseDSN = ""
		})
	}
}

func Test_db_BatchInsert(t *testing.T) {
	type args struct {
		br  []BatchRequest
		uid string
	}
	tests := []struct {
		name string
		args args
		want []BatchLink
	}{
		{
			name: "can batch insert",
			args: args{
				br: []BatchRequest{
					{CorrelationID: "1", OriginalURL: "ya.ru"},
					{CorrelationID: "1", OriginalURL: "test.ru"},
				},
				uid: "66f29390-381c-4a6a-9df9-74a0247ebe72",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewConfig(
				config.WithDatabaseConnection("postgres://root:root@localhost:5432/postgres?sslmode=disable"),
			)

			conn, err := NewDBConnection(&zap.Logger{}, false)

			if err == nil {
				res, errIns := conn.BatchInsert(tt.args.br, tt.args.uid)
				assert.NoError(t, errIns)
				assert.Len(t, res, 2)
				conn.conn.Exec(context.Background(), "delete from links where uid = '"+tt.args.uid+"'")
			}

			cfg.DatabaseDSN = ""
		})
	}
}
