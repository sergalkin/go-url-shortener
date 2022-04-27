package migrations

import (
	"github.com/sergalkin/go-url-shortener.git/internal/app/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up() (bool, error) {
	m, err := migrate.New(
		"file://../../internal/app/migrations", config.DatabaseDSN())
	if err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return false, err
		}
	}

	return true, nil
}
